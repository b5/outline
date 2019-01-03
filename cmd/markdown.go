package cmd

import (
	"fmt"
	"os"
	"text/template"

	"github.com/b5/outline/lib"
	"github.com/spf13/cobra"
)

const mdTmpl = `{{ range . }}
# Package {{ .Name }}

{{ .Description }}

{{ if gt (len .Functions) 0}}
## Functions
{{ range .Functions }}
### {{ .Signature }}
{{ .Description }}
{{ end }}
{{ end }}


{{ end }}`

// MarkdownCmd converts an outline to a markdown document
var MarkdownCmd = &cobra.Command{
	Use:     "markdown",
	Aliases: []string{"md"},
	Short:   "Convert docs to markdown syntax",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		t := template.Must(template.New("markdown").Parse(mdTmpl))
		var docs []*lib.Doc
		for _, fp := range args {
			f, err := os.Open(fp)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			doc, err := lib.Parse(f)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			docs = append(docs, doc)
		}

		if err := t.Execute(os.Stdout, docs); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	// MarkdownCmd.Flags().StringP("export", "e", "config.json", "path to configuration json file")
}
