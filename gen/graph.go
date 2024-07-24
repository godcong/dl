// Copyright (c) 2024 GodCong. All rights reserved.

// Package gen for Default Loader
package gen

const (
	defaultTagName  = "default"
	defaultFuncName = "Default"
)

// Graph represents the graph structure.
type Graph struct {
	Package string
	Imports []string
	Structs []*Struct
}

// Struct represents the structure.
type Struct struct {
	Name            string
	DefaultFuncName string
	Fields          []*Field
}

// IsValid checks if the struct is valid.
func (s Struct) IsValid() bool {
	return len(s.Fields) > 0
}

// Field represents a field in the struct.
type Field struct {
	IsBasic bool
	Name    string
	Type    string
	Value   string
}

// IsValid checks if the field is valid.
func (f Field) IsValid() bool {
	return f.Name != "" && f.Value != ""
}
