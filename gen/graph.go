package gen

import (
	"text/template"
)

type Header struct {
	Version   string
	BuildDate string
	BuiltBy   string
}

type Graph struct {
	Header  Header `template:"header"`
	Package string
	Imports []string
	Structs []*Struct
	Temple  *template.Template
}

type Struct struct {
	Name            string
	DefaultFuncName string
	Fields          []*Field
}

type Field struct {
	IsStruct bool
	Name     string
	Type     string
	Value    string
}
