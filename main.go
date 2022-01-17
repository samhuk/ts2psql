/* This converts Typescript type declarations to PostgreSql CREATE TABLE statements.
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var typeDeclarationStartRegex = regexp.MustCompile(`\/\*\s*ts2psql\s*({?.*}?)\s*\*\/\s*export\s+type\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\s*=\s*{`)
var typeDeclarationEndRegex = regexp.MustCompile(`\/\*\s*ts2psql\s*end\s*\*\/`)
var typePropertyRegex = regexp.MustCompile(`\/\*\s*ts2psql\s*({?.*}?)\s*\*\/\s*([a-zA-Z_$][a-zA-Z0-9_$]*)\??:\s*([a-zA-Z_$][a-zA-Z0-9_$]*)`)

func main() {
	// Read in config from file, determine input file paths
	config := readConfigFile()
	if config.Verbose {
		fmt.Println("* Using config file:", CONFIG_FILE_PATH)
	}

	paths := determineIncludedFilePaths(config)

	if config.Verbose {
		fmt.Printf("* Parsing %v files:\n", len(paths))
		for _, path := range paths {
			fmt.Println("  ", path)
		}
	}

	// Parse files, creating instances of TypeDeclaration
	var typeDeclarations, err = parseFiles(paths)
	if err != nil {
		fmt.Println("Error occured parsing files:", err)
	}

	// Create the CREATE TABLE statements
	createTableStatements := make([]string, len(typeDeclarations))
	for i := 0; i < len(typeDeclarations); i++ {
		createTableStatements[i] = createCreateTableStatement(typeDeclarations[i], typeDeclarations)
	}

	// Create directories for OutFile config if they don't already exist
	outFileDir := path.Dir(config.OutFile)
	if !determineIfPathExists(outFileDir) {
		if config.Verbose {
			fmt.Printf("* OutFile directory (%v) does not exist. Creating it.\n", outFileDir)
		}
		os.MkdirAll(outFileDir, 0666)
	}

	// Write CREATE TABLE statements to output file
	if config.Verbose {
		fmt.Println("* Writing output to", config.OutFile)
	}
	err = os.WriteFile(config.OutFile, []byte(strings.Join(createTableStatements, "\n\n")), 0666)
	if err != nil {
		fmt.Printf("Error occured writing to OutFile (%s): %v", config.OutFile, err)
	}
}

func determineIncludedFilePaths(config Config) []string {
	// If single file specified
	if len(config.File) > 0 {
		paths := make([]string, 1)
		paths[0] = config.File
		return paths
	}

	// If list of glob patterns specified
	if len(config.Include) > 0 {
		// Iterate through glob patterns, adding up all the found paths
		paths := make([]string, 0)
		for i := 0; i < len(config.Include); i++ {
			foundPaths, _ := filepath.Glob(config.Include[i])
			paths = append(paths, foundPaths...)
		}
		// Remove duplicate paths
		return removeDuplicateStr(paths)
	}

	return make([]string, 0)
}

func parseFiles(paths []string) ([]TypeDeclaration, error) {
	typeDeclarations := make([]TypeDeclaration, 0)
	for i := 0; i < len(paths); i++ {
		newTypeDeclarations, _ := parseFile(paths[i])
		typeDeclarations = append(typeDeclarations, newTypeDeclarations...)
	}
	return typeDeclarations, nil
}

func parseFile(path string) ([]TypeDeclaration, error) {
	fileText := readFile(path)
	i := 0

	typeDeclarations := make([]TypeDeclaration, 0)

	for i < len(fileText) {
		fileTextSlice := fileText[i:]
		data := typeDeclarationStartRegex.FindStringSubmatchIndex(fileTextSlice)
		if len(data) == 0 {
			break
		}
		end := typeDeclarationEndRegex.FindStringIndex(fileTextSlice)[1]
		typeDeclarationMetaDataJsonText := fileTextSlice[data[2]:data[3]]
		typeName := fileTextSlice[data[4]:data[5]]
		metaData := TypeDeclarationMetaData{}
		_ = json.Unmarshal([]byte(typeDeclarationMetaDataJsonText), &metaData)
		if len(metaData.TableName) == 0 {
			metaData.TableName = toSnakeCase(typeName)
		}
		typePropertyDeclarations, _ := parseTypeDeclarationText(fileTextSlice[data[5]:end])
		typeDeclaration := TypeDeclaration{typeName, metaData, typePropertyDeclarations}
		typeDeclarations = append(typeDeclarations, typeDeclaration)
		i += end
	}

	return typeDeclarations, nil
}

func parseTypeDeclarationText(typeDeclarationText string) ([]TypePropertyDeclaration, error) {
	i := 0

	typePropertyDeclarations := make([]TypePropertyDeclaration, 0)

	for i < len(typeDeclarationText) {
		textSlice := typeDeclarationText[i:]
		data := typePropertyRegex.FindStringSubmatchIndex(textSlice)
		if len(data) == 0 {
			break
		}
		TypePropertyDeclarationMetaDataJsonText := textSlice[data[2]:data[3]]
		name := textSlice[data[4]:data[5]]
		typeName := textSlice[data[6]:data[7]]
		metaData := TypePropertyDeclarationMetaData{}
		_ = json.Unmarshal([]byte(TypePropertyDeclarationMetaDataJsonText), &metaData)
		if len(metaData.ColumnName) == 0 {
			metaData.ColumnName = toSnakeCase(name)
		}
		typePropertyDeclaration := TypePropertyDeclaration{name, typeName, false, metaData}
		typePropertyDeclarations = append(typePropertyDeclarations, typePropertyDeclaration)
		i += data[7]
	}

	return typePropertyDeclarations, nil
}

func createCreateTableStatement(typeDeclaration TypeDeclaration, typeDeclarations []TypeDeclaration) string {
	str := ""
	str += "CREATE TABLE "

	str += typeDeclaration.MetaData.TableName
	str += " ( \n  "
	str += strings.Join(createColumnDeclarationStatements(typeDeclaration.TypePropertyDeclarations, typeDeclarations), ",\n  ")
	str += "\n);"
	return str
}

func createColumnDeclarationStatements(typePropertyDeclarations []TypePropertyDeclaration, typeDeclarations []TypeDeclaration) []string {
	statements := make([]string, len(typePropertyDeclarations))
	for i := 0; i < len(typePropertyDeclarations); i++ {
		statements[i] = createColumnDeclarationStatement(typePropertyDeclarations[i], typeDeclarations)
	}
	return statements
}

// Creates the column declaration string, e.g. "id INTEGER serial PRIMARY KEY NOT NULL"
func createColumnDeclarationStatement(typePropertyDeclaration TypePropertyDeclaration, typeDeclarations []TypeDeclaration) string {
	str := ""
	// Column name
	str += typePropertyDeclaration.MetaData.ColumnName + " "

	// Field type name
	str += createSqlTypeName(typePropertyDeclaration)
	str += " "

	// Optionally add "serial" property
	if typePropertyDeclaration.MetaData.Serial {
		str += "serial "
	}
	// Optionally add  "PRIMARY KEY" property
	if typePropertyDeclaration.MetaData.PrimaryKey {
		str += "PRIMARY KEY "
	}
	// Optionally add  "UNIQUE" property
	if typePropertyDeclaration.MetaData.Unique {
		str += "UNIQUE "
	}
	// Optionally add  "NOT NULL" property
	if !typePropertyDeclaration.Optional {
		str += "NOT NULL "
	}

	str = strings.TrimSpace(str)

	// Add foreign key statement
	foreignTypeName := typePropertyDeclaration.MetaData.Fk.TypeName // Alias
	if len(foreignTypeName) > 0 {
		foundForeignTypeDeclarationTableName := findTypeDeclarationTableNameByTypeName(foreignTypeName, typeDeclarations)
		str += ",\n  FOREIGN KEY (" + typePropertyDeclaration.MetaData.ColumnName + ")\n    REFERENCES " + foundForeignTypeDeclarationTableName + " (" + typePropertyDeclaration.MetaData.Fk.Property + ")"
	}

	return strings.TrimSpace(str)
}

func findTypeDeclarationTableNameByTypeName(name string, typeDeclarations []TypeDeclaration) string {
	for i := 0; i < len(typeDeclarations); i++ {
		if typeDeclarations[i].Name == name {
			return typeDeclarations[i].MetaData.TableName
		}
	}
	return "[ERROR: \"" + name + "\" is not a valid foreign key type. It does not match any of the type definition names]"
}

func createSqlTypeName(typePropertyDeclaration TypePropertyDeclaration) string {
	switch typePropertyDeclaration.TypeName {
	// String
	case "string":
		str := ""
		str += "VARCHAR("
		if typePropertyDeclaration.MetaData.MaxLength == 0 {
			str += "50"
		} else {
			str += fmt.Sprint(typePropertyDeclaration.MetaData.MaxLength)
		}
		str += ")"
		return str
	// Number
	case "number":
		if len(typePropertyDeclaration.MetaData.NumberType) > 0 {
			return fmt.Sprint(typePropertyDeclaration.MetaData.NumberType)
		} else {
			return "INTEGER"
		}
	// Boolean
	case "boolean":
		return "BOOLEAN"
	// Date
	case "Date":
		return "TIMESTAMP"
	// Invalid TS type
	default:
		return "[ERROR: \"" + typePropertyDeclaration.TypeName + "\" is not a valid type name]"
	}
}
