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
	if len(field.Names) == 0 {
		return nil
	}

	if v, ok := field.Type.(*ast.StructType); ok {
		sub := &Struct{
			Name:            field.Names[0].String(),
			DefaultFuncName: defaultFuncName,
		}
		parseStructTags(sub, v)

		// TODO: now used reflect to set the unsupported type by `dl.Load`
		return &Field{
			Name:    sub.Name,
			IsBasic: false,
		}
	}

	if field.Tag == nil {
		return nil
	}

	fieldName := field.Names[0].String()
	tags := NewStructTag(field.Tag.Value)
	val, ok := tags.Lookup(tagName)
	if !(ok && validateTag(val)) {
		return nil
	}
	fieldType := parseType(field.Type)
	debugPrint("field tag:", fmt.Sprintf("tagName: %s, fieldName: %s, fieldType: %s, tagVal: %s",
		tagName, fieldName, fieldType, val))
	return &Field{
		IsBasic: true,
		Name:    fieldName,
		Type:    fieldType,
		Value:   val,
	}
}

func validateTag(val string) bool {
	return val != "" && val != "-"
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
	default:
		panic(fmt.Sprintf("unknown type %T", x))
	}
}

func parseStructTags(gs *Struct, x *ast.StructType) {
	for _, field := range x.Fields.List {
		debugPrint("struct tags:", fmt.Sprintf("Type(%T)", field.Type), fmt.Sprintf("Value(%+v) ", field))
		// switch field.Type.(type) {
		// case *ast.StructType:
		// 	sub := &Struct{
		// 		Name:            field.Names[0].String(),
		// 		DefaultFuncName: defaultFuncName,
		// 	}
		// 	gs.Fields = append(gs.Fields, &Field{
		// 		Name: sub.Name,
		// 		// Type:    sub.Name,
		// 		// Value:   formatValue(sub.Name, field.Tag.Value),
		// 		IsBasic: false,
		// 	})
		// }

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
