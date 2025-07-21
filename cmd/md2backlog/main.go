package main

import (
	"fmt"
	"io"
	"os"

	"md2backlog/internal/converter"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

var (
	inputFile  string
	outputFile string
)

var rootCmd = &cobra.Command{
	Use:     "md2backlog",
	Short:   "Convert Markdown to Backlog notation",
	Long:    "A CLI tool to convert Markdown text to Backlog notation using AST parsing",
	Version: version,
	Run:     runConvert,
}

func runConvert(cmd *cobra.Command, args []string) {
	var input []byte
	var err error

	// 入力の読み取り
	if inputFile != "" {
		input, err = os.ReadFile(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input file: %v\n", err)
			os.Exit(1)
		}
	} else {
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
	}

	// Markdownをバックログ記法に変換
	result, err := converter.Convert(string(input))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting: %v\n", err)
		os.Exit(1)
	}

	// 出力の書き込み
	if outputFile != "" {
		err = os.WriteFile(outputFile, []byte(result), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Print(result)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input Markdown file (default: stdin)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: stdout)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
