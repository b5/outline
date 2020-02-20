package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// logger
var log = logrus.New()

// RootCmd is the walk command
var RootCmd = &cobra.Command{
	Short: "outline is a tool for outlining software packages",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug, err := cmd.Flags().GetBool("debug"); err == nil && debug {
			logrus.SetLevel(logrus.DebugLevel)
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().Bool("debug", false, "show debug output")
	RootCmd.AddCommand(
		FmtCmd,
		TemplateCmd,
		PackageCmd,
	)
}
