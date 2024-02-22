package cmd

import (
	"github.com/Beigelman/ludaapi/scripts/cmd/createusers"
	"github.com/Beigelman/ludaapi/scripts/cmd/importincomes"
	"github.com/Beigelman/ludaapi/scripts/cmd/importsplit"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nossas-despesas-scripts",
	Short: "Scripts para automatizar tarefas de importação de arquivos",
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
	rootCmd.AddCommand(createusers.Cmd())
	rootCmd.AddCommand(importsplit.Cmd())
	rootCmd.AddCommand(importincomes.Cmd())
}
