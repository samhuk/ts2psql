package main

type TypeDeclaration struct {
	Name                     string
	MetaData                 TypeDeclarationMetaData
	TypePropertyDeclarations []TypePropertyDeclaration
}

type TypeDeclarationMetaData struct {
	TableName string
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
	Unique     bool
	MaxLength  int
	NumberType string
	ColumnName string
	Fk         ForeignKeyMetaData
}

type ForeignKeyMetaData struct {
	TypeName string `json:"type"`
	Property string
}

type Config struct {
	Include []string
	File    string
	OutFile string
	OutDir  string // TODO
	Verbose bool
}
