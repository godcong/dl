// Copyright (c) 2024 GodCong. All rights reserved.

// Package gen for Default Loader
package gen

import (
	_ "embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// ParseFromTags parse default tags from struct.
func ParseFromTags(fileName string) (*Graph, error) {
	graph := Graph{
		Package: "",
		Imports: nil,
		Structs: nil,
	}

	// positions are relative to fset
	fset := token.NewFileSet()
	// Parse the file given in arguments
	f, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	graph.Package = f.Name.Name
	// range over the objects in the scope of this generated AST and check for StructType. Then range over fields
	// contained in that struct.
	ast.Inspect(f, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.File:
			for _, imp := range t.Imports {
				graph.Imports = append(graph.Imports, imp.Path.Value)
			}
		case *ast.TypeSpec:
			if v, ok := t.Type.(*ast.StructType); ok {
				s := &Struct{
					Name:            t.Name.Name,
					DefaultFuncName: defaultFuncName,
				}
				parseStructTags(s, v)
				if s.IsValid() {
					graph.Structs = append(graph.Structs, s)
				}
			}
		}
		return true
	})
	return &graph, nil
}

func parseFieldTag(field *ast.Field, tagName string) *Field {
	fieldName := field.Names[0].String()

	tagValues := strings.Fields(strings.Trim(field.Tag.Value, "`"))
	val := ""
	for i := range tagValues {
		if strings.HasPrefix(tagValues[i], tagName) {
			val = strings.TrimPrefix(tagValues[i], fmt.Sprintf("%s:", tagName))
			break
		}
	}
	if val == "" {
		return nil
	}

	fieldType := parseType(field.Type)
	val = strings.TrimPrefix(val, "\"")
	val = strings.TrimSuffix(val, "\"")
	return &Field{
		Name:  fieldName,
		Type:  fieldType,
		Value: val,
	}
}

func parseType(x ast.Expr) string {
	switch v := x.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.StarExpr:
		return fmt.Sprintf("*%s", parseType(v.X))
	case *ast.ArrayType:
		return fmt.Sprintf("[]%s", parseType(v.Elt))
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", parseType(v.Key), parseType(v.Value))
	}

	return fmt.Sprintf("%v", x)
}

func parseStructTags(gs *Struct, x *ast.StructType) {
	for _, field := range x.Fields.List {
		if len(field.Names) == 0 {
			continue
		}
		if field.Tag == nil {
			continue
		}

		tagValue := parseFieldTag(field, defaultTagName)
		if tagValue != nil {
			gs.Fields = append(gs.Fields, formatField(tagValue))
		}
	}
}

func formatField(value *Field) *Field {
	value.Value = formatValue(value.Type, value.Value)
	return value
}

func formatDefaultValue(typo string) (string, bool) {
	basicType := true
	switch typo {
	case "bool":
		fallthrough
	case "int", "int8", "int16", "int32", "int64":
		fallthrough
	case "uint", "uint8", "uint16", "uint32", "uint64":
		fallthrough
	case "float32", "float64":
		fallthrough
	case "string":
		fallthrough
	case "[]byte":
		basicType = true
	default:
		basicType = false
	}
	return typo, basicType
}

func formatValue(typo string, value string) string {
	switch {
	case strings.HasPrefix(typo, "*"):
		innerType := typo[1:]
		value = formatValue(innerType, value)
		return fmt.Sprintf("dl.Pointer(%s)", value)
	case strings.HasPrefix(typo, "string"):
		value = fmt.Sprintf("\"%s\"", value)
	case strings.HasPrefix(typo, "[]byte"):
		value = fmt.Sprintf("[]byte(\"%s\")", value)
	case strings.HasPrefix(typo, "map["):
		innerType := typo[3:]
		values := toArray(value, ",", 1)
		keyType, valueType := mapKeyValueTypes(innerType)
		for i := range values {
			kv := strings.Split(values[i], ":")
			if len(kv) == 2 {
				kv[0] = formatValue(keyType, kv[0])
				kv[1] = formatValue(valueType, kv[1])
				values[i] = fmt.Sprintf("%s:%s", kv[0], kv[1])
			}
		}
		value = fmt.Sprintf("%s{%s}", typo, strings.Join(values, ","))
	case strings.HasPrefix(typo, "[]"):
		innerType := typo[2:]
		values := toArray(value, ",", 1)
		for i := range values {
			values[i] = formatValue(innerType, values[i])
		}
		value = fmt.Sprintf("%s{%s}", typo, strings.Join(values, ","))
	}
	return value
}

func toArray(value string, split string, bit int) []string {
	if bit >= len(value) {
		return nil
	}
	return strings.Split(value[bit:len(value)-bit], split)
}

func mapKeyIndex(value string) (int, int) {
	sta, end := -1, -1
	count := 0
	for i := range value {
		switch value[i] {
		case '[':
			if sta == -1 {
				sta = i
			}
			count++
		case ']':
			count--
			if count == 0 {
				return sta, i
			}
		}
	}
	return sta, end
}

func mapKeyValueTypes(typo string) (string, string) {
	sta, end := mapKeyIndex(typo)
	if sta == -1 || end == -1 || end < sta {
		return "", ""
	}

	key := typo[sta+1 : end]
	var value string
	if end+1 < len(typo) {
		value = typo[end+1:]
	}
	return key, value
}
