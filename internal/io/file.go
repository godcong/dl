// Copyright (c) 2024 GodCong. All rights reserved.

// Package io for Default Loader
package io

import (
	"bytes"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	goversion "github.com/caarlos0/go-version"

	"github.com/godcong/dl/gen"
	"github.com/godcong/dl/internal/shell"
	"github.com/godcong/dl/internal/tpl"
)

const (
	goFileSuffix      = ".go"
	defaultFileSuffix = "_default.go"
	goTestFileSuffix  = "_test.go"
)

func isGoFile(file os.DirEntry) bool {
	name := file.Name()
	return !file.IsDir() && strings.HasSuffix(name, goFileSuffix) && !strings.HasSuffix(name,
		goTestFileSuffix) && !strings.HasSuffix(name, defaultFileSuffix)
}

// ReadDir returns a list of Go files in the specified directory excluding "_test.go" and "_default.go".
// It takes a string parameter `file` representing the directory path and returns a slice of strings and an error.
func ReadDir(dir string) ([]string, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, f := range dirEntries {
		if isGoFile(f) {
			files = append(files, filepath.Join(dir, f.Name()))
		}
	}

	return files, nil
}

// WriteGraph write file to fileName with graph
func WriteGraph(fileName string, info goversion.Info, graph *gen.Graph) error {
	if graph.Structs == nil {
		return nil
	}

	fileName = strings.Replace(fileName, ".go", defaultFileSuffix, 1)

	temple, err := template.New("defaultLoader").Parse(tpl.StructTemplate)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(nil)
	if err := temple.ExecuteTemplate(buf, "header", &info); err != nil {
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

	if err := os.WriteFile(fileName, formatted, 0600); err != nil {
		return err
	}

	if err := shell.ExecGoImports("", "", fileName); err != nil {
		return err
	}

	return nil
}
