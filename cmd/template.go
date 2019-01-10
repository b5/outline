package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/b5/outline/lib"
	"github.com/spf13/cobra"
)

// backticks don't work with golang string literals, so use "'" as a stand-in & strings.Replace
var mdIndex = strings.Replace(`{{- define "mdFn" }}
#### '{{ .Signature }}'
{{- if ne .Description "" }}
{{ .Description }}
{{- end -}}
{{- if gt (len .Params) 0 }}
**parameters:**
| name | type | description |
|------|------|-------------|
{{ range .Params -}}
| '{{ .Name }}' | '{{ .Type }}' | {{ .Description }} |
{{ end -}}
{{- end -}}
{{- end -}}

{{- range . -}}
# {{ .Name }}
{{ if ne .Description "" }}{{ .Description }}{{ end }}
{{- if gt (len .Functions) 0 }}

## Functions
{{ range .Functions -}}
{{ template "mdFn" . }}
{{ end -}}
{{- end }}
{{ if gt (len .Types) 0 }}
## Types

{{ range .Types -}}
### '{{ .Name }}'
{{ if ne .Description "" }}{{ .Description }}{{ end -}}
{{ if gt (len .Fields) 0 }}
**Fields**
| name | type | description |
|------|------|-------------|
{{ range .Fields -}}
| {{ .Name }} | {{ .Type }} | {{ .Description }} |
{{ end -}}
{{ end -}}
{{ if gt (len .Methods) 0 }}
**Methods**
{{- range .Methods -}}
{{ template "mdFn" . }}
{{ end -}}
{{- if gt (len .Operators) 0 }}
**Operators**
| operator | description |
|----------|-------------|
{{ range .Operators -}}
	| {{ .Opr }} | {{ .Description }} |
{{ end }}
{{ end }}
{{ end }}
{{- end -}}
{{- end -}}
{{ end }}`, "'", "`", -1)

// TemplateCmd parses outline documents & executes them against a template
var TemplateCmd = &cobra.Command{
	Use:     "template",
	Aliases: []string{"md"},
	Short:   "execute outline documents against a template",
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

		var docs lib.Docs
		for _, fp := range args {
			f, err := os.Open(fp)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			read, err := lib.Parse(f)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			docs = append(docs, read...)
		}

		noSort, err := cmd.Flags().GetBool("no-sort")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if !noSort {
			docs.Sort()
		}

		if err := t.Execute(os.Stdout, docs); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	TemplateCmd.Flags().StringP("template", "t", "", "template file to load. overrides preset")
	TemplateCmd.Flags().Bool("no-sort", false, "don't alpha-sort fields & outline documents")
}
