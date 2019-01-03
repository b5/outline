package cmd

import (
	"fmt"
	"os"

	"github.com/b5/outline/lib"
	"github.com/spf13/cobra"
)

// FmtCmd prints present configuration info
var FmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "format input",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
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

			data, err := doc.MarshalIndent(0, "  ")
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			fmt.Print(string(data))
		}
	},
}

func init() {
	// FmtCmd.Flags().StringP("export", "e", "config.json", "path to configuration json file")
}
