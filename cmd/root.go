// Copyright (c) 2024 GodCong. All rights reserved.

// Package cmd
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	goversion "github.com/caarlos0/go-version"
	"github.com/spf13/cobra"

	"github.com/godcong/dl/gen"
	"github.com/godcong/dl/internal/io"
)

const helpExample = `
### Step 1: Define Struct with Default Tags

Add the ` + "`default`" + ` tag to your struct fields to specify their default values:

// example: demo.go
type Demo struct {
	Name string ` + "`default:" + `"demo"
}

### Step 2: Generate Default Value Loading Method

Run DL to generate the necessary loading method for your struct:

$> dl -f ./demo.go
`

const asciiArt = "\n   ___      ___          ____  __               __       \n  / _ \\___ / _/__ ___ __/ / /_/ / ___  ___ ____/ /__ ____\n / // / -_) _/ _ `/ // / / __/ /_/ _ \\/ _ `/ _  / -_) __/\n/____/\\__/_/ \\_,_/\\_,_/_/\\__/____|___/\\_,_/\\_,_/\\__/_/   \n                                                         \n"
const website = "https://github.com/godcong/dl"

// build tool goreleaser tags
//nolint:gochecknoglobals
var (
	version   = ""
	commit    = ""
	treeState = ""
	date      = ""
	builtBy   = ""
	debug     = false
)

var helpCmd = &cobra.Command{
	Use:     "help",
	Example: helpExample,
	Run: func(cmd *cobra.Command, args []string) {
		gv := buildVersion(version, commit, date, builtBy, treeState)
		fmt.Println(gv)
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
		if debug {
			gen.Debug()
		}
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
			filelist, err = io.ReadDir(file)
			if err != nil {
				return err
			}
		} else {
			filelist = append(filelist, file)
		}

		head := buildVersion(version, commit, date, builtBy, treeState)
		for _, s := range filelist {
			graph, err := gen.ParseFromFile(s)
			if err != nil {
				return err
			}
			if err := io.WriteGraph(s, head, graph, true); err != nil {
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

	rootCmd.SetHelpCommand(helpCmd)
	rootCmd.Version = buildVersion(version, commit, date, builtBy, treeState).String()
	rootCmd.Flags().StringP("file", "f", ".", "load go files or directories")
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "debug mode")
}

func buildVersion(version, commit, date, builtBy, treeState string) goversion.Info {
	return goversion.GetVersionInfo(
		goversion.WithAppDetails("Default Loader", "A default value generate tool for go structs", website),
		goversion.WithASCIIName(asciiArt),
		func(i *goversion.Info) {
			if commit != "" {
				i.GitCommit = commit
			}
			if version != "" {
				i.GitVersion = version
			}
			if treeState != "" {
				i.GitTreeState = treeState
			}
			if date != "" {
				i.BuildDate = date
			}
			if builtBy != "" {
				i.BuiltBy = builtBy
			}
		},
	)
}
