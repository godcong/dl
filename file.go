package dl

import (
	"bytes"
	"go/format"
	"os"
	"strings"

	"github.com/godcong/dl/gen"
)

func GetFileList(file string) ([]string, error) {
	files, err := os.ReadDir(file)
	if err != nil {
		return nil, err
	}
	var filelist []string
	for _, f := range files {
		if !f.IsDir() &&
			strings.HasSuffix(f.Name(), ".go") &&
			!strings.HasSuffix(f.Name(), "_test.go") &&
			!strings.HasSuffix(f.Name(), "_default.go") {
			filelist = append(filelist, file+"/"+f.Name())
		}
	}
	return filelist, nil
}

// WriteFile write file to fileName with graph
func WriteFile(fileName string, graph *gen.Graph) error {
	temple, err := graph.Temple.Parse(loadTemplate)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(nil)
	if err := temple.ExecuteTemplate(buf, "header", &graph.Header); err != nil {
		return err
	}
	if err := temple.Execute(buf, &graph); err != nil {
		return err
	}
	if err := temple.ExecuteTemplate(buf, "structs", &graph); err != nil {
		return err
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	if err := os.WriteFile(fileName, formatted, 0644); err != nil {
		return err
	}

	if err := ExecGoImports("", "", fileName); err != nil {
		return err
	}

	return nil
}
