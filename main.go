/* This converts Typescript type declarations to PostgreSql CREATE TABLE statements.
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

type TypeDeclaration struct {
	Name                     string
	MetaData                 TypeDeclarationMetaData
	TypePropertyDeclarations []TypePropertyDeclaration
}

type TypeDeclarationMetaData struct {
	TableName   string
	ToSnakeCase bool
}

type TypePropertyDeclaration struct {
	Name     string
	TypeName string
	Optional bool
	MetaData TypePropertyDeclarationMetaData
}

type TypePropertyDeclarationMetaData struct {
	PrimaryKey bool
	Serial     bool
}

var PATHS = [1]string{"./test.ts"}

var typeDeclarationStartRegex = regexp.MustCompile(`\/\*\s*ts2psql\s*({.*})\s*\*\/\s*export\s+type\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\s*=\s*{`)
var typeDeclarationEndRegex = regexp.MustCompile(`\/\*\s*ts2psql\s*end\s*\*\/`)
var typePropertyRegex = regexp.MustCompile(`\/\*\s*ts2psql\s*({?.*}?)\s*\*\/\s*([a-zA-Z_$][a-zA-Z0-9_$]*)\??:\s*([a-zA-Z_$][a-zA-Z0-9_$]*)`)

func main() {
	// Gzip the large files
	for i := 0; i < len(PATHS); i++ {
		var _, err = parseFiles(PATHS[:])
		if err != nil {
			fmt.Printf("Error occured converting file %v:\n%v\n", PATHS[i], err)
		}
	}
}

func parseFiles(paths []string) ([]TypeDeclaration, error) {
	typeDeclarations := make([]TypeDeclaration, 0)
	for i := 0; i < len(PATHS); i++ {
		typeDeclaration, _ := parseFile(PATHS[i])
		typeDeclarations = append(typeDeclarations, typeDeclaration...)
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
		end := typeDeclarationEndRegex.FindStringIndex(fileTextSlice)[1]
		typeDeclarationMetaDataJsonText := fileTextSlice[data[2]:data[3]]
		typeName := fileTextSlice[data[4]:data[5]]
		metaData := TypeDeclarationMetaData{}
		_ = json.Unmarshal([]byte(typeDeclarationMetaDataJsonText), &metaData)
		typePropertyDeclarations, _ := parseTypeDeclarationText(fileTextSlice[data[5]:end])

		typeDeclaration := TypeDeclaration{typeName, metaData, typePropertyDeclarations}
		typeDeclarations = append(typeDeclarations, typeDeclaration)
		i += end
		fmt.Println(typeDeclaration)
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
		typePropertyDeclaration := TypePropertyDeclaration{name, typeName, false, metaData}
		typePropertyDeclarations = append(typePropertyDeclarations, typePropertyDeclaration)
		i += data[7]
	}

	return typePropertyDeclarations, nil
}

func convertTypeDeclarationToCreateTableStatement(typeDeclaration TypeDeclaration) string {
	return ""
}

func readFile(path string) string {
	fileText, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error occured reading in file %v:\n%v\n", path, err)
	}
	return string(fileText)
}
