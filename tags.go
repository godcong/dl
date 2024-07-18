package dl

import (
	_ "embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"text/template"

	"github.com/godcong/dl/gen"
)

const (
	defaultTagName    = "default"
	defaultFileSuffix = "_default.go"
)

//go:embed template.go.tmpl
var loadTemplate string

// GenerateFromTags generates snake case json tags so that you won't need to write them. Can be also extended to xml or sql tags
func GenerateFromTags(header gen.Header, fileName string) error {
	graph := gen.Graph{
		Header:  header,
		Package: "",
		Imports: nil,
		Structs: nil,
		Temple:  template.New("defaultLoader"),
	}

	// positions are relative to fset
	fset := token.NewFileSet()
	// Parse the file given in arguments
	f, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		return err
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
			switch sv := t.Type.(type) {
			case *ast.StructType:
				s := &gen.Struct{
					Name:            t.Name.Name,
					DefaultFuncName: "Default",
				}
				processStructTags(s, sv)
				graph.Structs = append(graph.Structs, s)
			}
		}
		return true
	})

	// replace *.go to *_default.go
	fileName = strings.Replace(fileName, ".go", defaultFileSuffix, 1)

	if err := WriteFile(fileName, &graph); err != nil {
		return err
	}

	return nil
}

func parseFieldTag(field *ast.Field, tagName string) *gen.Field {
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
	fieldType := fmt.Sprintf("%s", field.Type)
	if strings.Compare(fieldType, "string") != 0 {
		val = strings.TrimPrefix(val, "\"")
		val = strings.TrimSuffix(val, "\"")
	}
	return &gen.Field{
		Name:  fieldName,
		Type:  fmt.Sprintf("%s", field.Type),
		Value: val,
	}
}

func processStructTags(gs *gen.Struct, x *ast.StructType) {
	for _, field := range x.Fields.List {
		if len(field.Names) == 0 {
			continue
		}
		if field.Tag == nil {
			return
		}

		tagDefaults := parseFieldTag(field, defaultTagName)
		if tagDefaults != nil {
			gs.Fields = append(gs.Fields, tagDefaults)
		}
	}
}
