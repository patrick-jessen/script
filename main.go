package main

import (
	"fmt"
	"os"

	"github.com/patrick-jessen/script/compiler/analyzer"
	"github.com/patrick-jessen/script/compiler/generator"
	"github.com/patrick-jessen/script/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Version:           "0.1",
	Use:               "script",
	Short:             "script is a compiler",
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

func init() {
	buildCmd := &cobra.Command{
		Use:   "build [path]",
		Short: "Build to WebAssembly",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			dir := args[0]
			// format, _ := cmd.Flags().GetString("format")
			// output, _ := cmd.Flags().GetString("output")

			debugTokens, _ := cmd.Flags().GetBool("tokens")
			debugAST, _ := cmd.Flags().GetBool("ast")

			config.DebugTokens = debugTokens
			config.DebugAST = debugAST

			analyzer := analyzer.New(dir)
			analyzer.Run()

			generator := generator.New(analyzer)
			generator.Run()

			return nil
		},
	}
	buildCmd.Flags().StringP("format", "f", "wat", "Output format")
	buildCmd.Flags().StringP("output", "o", "", "Output file (default \"stdout\")")
	buildCmd.Flags().Bool("tokens", false, "Debug tokens")
	buildCmd.Flags().Bool("ast", false, "Debug AST")

	rootCmd.AddCommand(buildCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
