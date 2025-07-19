package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:     "md2backlog",
	Short:   "Convert Markdown to Backlog notation",
	Long:    "A CLI tool to convert Markdown text to Backlog notation using AST parsing",
	Version: version,
	Run:     runConvert,
}

func runConvert(cmd *cobra.Command, args []string) {
	if _, err := io.Copy(os.Stdout, os.Stdin); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
