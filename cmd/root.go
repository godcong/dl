// Copyright (c) 2024 GodCong. All rights reserved.

// Package main
package main

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/godcong/dl"
	"github.com/godcong/dl/gen"
)

const helpExample = `
	# Add default value to struct field with file(e.g. demo.go)
	type Demo struct {
		Name string ` + "`" + `default:"demo"` + "`" + `
	}   
	# Run the tool to generate default value
	$ dl -f demo.go
`

var Version string
var BuiltBy string

var helpCmd = &cobra.Command{
	Use:     "help",
	Example: helpExample,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dl",
	Short: "A default value generate tool from tag ",
	Long: `A default value generate tool from tag.
	you can use tag like: default:"default value" to set default value for the field.
	and this tool will generate the default value for the field.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}
		openFile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer openFile.Close()
		stat, err := openFile.Stat()
		if err != nil {
			return err
		}
		if !stat.IsDir() && !strings.HasSuffix(stat.Name(), ".go") {
			return errors.New("file must be a go file or a directory")
		}

		var filelist []string
		if stat.IsDir() {
			filelist, err = dl.GetFileList(file)
			if err != nil {
				return err
			}
		} else {
			filelist = append(filelist, file)
		}

		header := gen.Header{
			BuildDate: time.Now().Format("2006-01-02 15:04:05"),
			BuiltBy:   BuiltBy,
			Version:   Version,
		}
		for _, s := range filelist {
			err := dl.GenerateFromTags(header, s)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SilenceUsage = true
	rootCmd.AddCommand(helpCmd)
	rootCmd.Flags().StringP("file", "f", ".", "load go files or directories")
}

func main() {
	Execute()
}
