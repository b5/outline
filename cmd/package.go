package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/b5/outline/lib"
	"github.com/spf13/cobra"
	parseutil "gopkg.in/src-d/go-parse-utils.v1"
)

// PackageCmd extracts and execute outline documents from a go package against a template",
var PackageCmd = &cobra.Command{
	Use:     "package",
	Aliases: []string{"pkg"},
	Short:   "exctract and execute outline documents from a go package against a template",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		t := template.Must(template.New("mdIndex").Parse(mdIndex))

		str, err := cmd.Flags().GetString("template")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if str != "" {
			t, err = template.ParseFiles(str)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		docs := map[string]*lib.Doc{}
		for _, pkg := range args {
			pkg, err := parseutil.PackageAST(pkg)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			for _, f := range pkg.Files {
				for _, c := range f.Comments {
					buf := strings.NewReader(c.Text())
					read, err := lib.Parse(buf)
					if err != nil {
						fmt.Println(err.Error())
						os.Exit(1)
					}

					for _, doc := range read {
						if found, ok := docs[doc.Name]; ok {
							merge(found, doc)
							continue
						}

						docs[doc.Name] = doc
					}
				}
			}
		}

		noSort, err := cmd.Flags().GetBool("no-sort")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var list lib.Docs
		for _, doc := range docs {
			list = append(list, doc)
		}

		if !noSort {
			list.Sort()
		}

		if err := t.Execute(os.Stdout, list); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

func merge(a, b *lib.Doc) {
	if a.Description == "" {
		a.Description = b.Description
	}

	if a.Path == "" {
		a.Path = b.Path
	}

	a.Types = append(a.Types, b.Types...)
	a.Functions = append(a.Functions, b.Functions...)
}

func init() {
	PackageCmd.Flags().StringP("template", "t", "", "template file to load. overrides preset")
	PackageCmd.Flags().Bool("no-sort", false, "don't alpha-sort fields & outline documents")
}
