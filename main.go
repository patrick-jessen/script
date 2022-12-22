package main

import (
	"fmt"
	"os"

	"github.com/patrick-jessen/script/compiler"
	"github.com/patrick-jessen/script/compiler/config"

	"github.com/patrick-jessen/script/lang/jlang"
	"github.com/patrick-jessen/script/lang/jsonlang"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Version:           "0.1",
	Use:               "script",
	Short:             "script is a compiler",
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

var languages = map[string]compiler.Compiler{
	"json":  &jsonlang.JSONLanguage{},
	"jlang": &jlang.JLang{},
}

func init() {
	buildCmd := &cobra.Command{
		Use:   "build [path]",
		Short: "Build to WebAssembly",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Extract arguments
			path := args[0]
			inputLang, _ := cmd.Flags().GetString("lang")
			outputFile, _ := cmd.Flags().GetString("output")
			outputFormat, _ := cmd.Flags().GetString("format")
			config.NoColor, _ = cmd.Flags().GetBool("no-color")
			config.DebugTokens, _ = cmd.Flags().GetBool("debug-tokens")
			config.DebugAST, _ = cmd.Flags().GetBool("debug-ast")

			lang, ok := languages[inputLang]
			if !ok {
				return fmt.Errorf("unsupported language")
			}
			if outputFormat != "wat" && outputFormat != "wasm" {
				return fmt.Errorf("unsupported format")
			}

			// Run the compiler
			generatedOutput := lang.Compile(path)
			if generatedOutput == nil {
				os.Exit(1)
			}

			// Ouput the result
			if len(outputFile) > 0 {
				err := os.WriteFile(outputFile, generatedOutput, 0666)
				if err != nil {
					return err
				}
			} else {
				if outputFormat == "wat" {
					fmt.Println(string(generatedOutput))
				} else {
					os.Stdout.Write(generatedOutput)
				}
			}
			return nil
		},
	}
	buildCmd.Flags().StringP("lang", "l", "jlang", `Input langugage`)
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
