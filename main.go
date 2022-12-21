package main

import (
	"fmt"
	"os"

	"github.com/patrick-jessen/script/compiler"
	"github.com/patrick-jessen/script/config"
	"github.com/patrick-jessen/script/lang/jlang"
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
			// Extract arguments
			dir := args[0]

			output, _ := cmd.Flags().GetString("output")
			format, _ := cmd.Flags().GetString("format")
			if format != "wat" && format != "wasm" {
				return fmt.Errorf("unsupported format")
			}

			debugTokens, _ := cmd.Flags().GetBool("debug-tokens")
			debugAST, _ := cmd.Flags().GetBool("debug-ast")
			noColor, _ := cmd.Flags().GetBool("no-color")
			config.DebugTokens = debugTokens
			config.DebugAST = debugAST
			config.NoColor = noColor

			// Run the compiler
			analyzer := compiler.NewAnalyzer(dir, &jlang.JLang{})
			err := analyzer.Run()
			if err != nil {
				os.Exit(1)
			}

			generatedOutput := []byte("<generated output>")

			// Ouput the result
			if len(output) > 0 {
				err := os.WriteFile(output, generatedOutput, 0666)
				if err != nil {
					return err
				}
			} else {
				if format == "wat" {
					fmt.Println(string(generatedOutput))
				} else {
					os.Stdout.Write(generatedOutput)
				}
			}
			return nil
		},
	}
	buildCmd.Flags().StringP("format", "f", "wat", `Output format ["wasm", "wat"]`)
	buildCmd.Flags().StringP("output", "o", "", "Output file (if not specified output is written to stdout)")
	buildCmd.Flags().Bool("no-color", false, `Disable color in debug output`)
	buildCmd.Flags().Bool("debug-tokens", false, "Debug tokens")
	buildCmd.Flags().Bool("debug-ast", false, "Debug AST")

	rootCmd.AddCommand(buildCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
