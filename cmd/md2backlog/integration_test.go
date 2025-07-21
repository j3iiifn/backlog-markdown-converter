package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"md2backlog/internal/converter"

	"github.com/spf13/cobra"
)

// TestCompleteMarkdownConversion は実際のMarkdown変換を統合的にテストする
func TestCompleteMarkdownConversion(t *testing.T) {
	defer resetRootCmd()

	simpleMarkdown := `# Main Heading

**Bold text** and *italic text*

- List item 1
- List item 2

[Link text](https://example.com)`

	output := runFileConversionTest(t, simpleMarkdown)

	// 基本的な変換が動作していることを確認
	checkConversions := []struct {
		expected string
		errMsg   string
	}{
		{"* Main Heading", "Heading conversion failed"},
		{"''Bold text''", "Bold conversion failed"},
		{"'''italic text'''", "Italic conversion failed"},
		{"- List item", "List conversion failed"},
		{"[[Link text:https://example.com]]", "Link conversion failed"},
	}

	for _, check := range checkConversions {
		if !strings.Contains(output, check.expected) {
			t.Errorf("%s", check.errMsg)
		}
	}
}

// TestStdinToStdoutIntegration は標準入出力での統合テストを実行する
func TestStdinToStdoutIntegration(t *testing.T) {
	defer resetRootCmd()

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "mixed formatting",
			input:    "# Title\n\n**bold** and *italic*",
			expected: "* Title\n''bold'' and '''italic'''",
		},
		{
			name:     "code elements",
			input:    "`inline` and\n```\nblock code\n```",
			expected: "{code}inline{/code} and\n\n>{code}\nblock code\n{/code}<",
		},
		{
			name:     "lists and links",
			input:    "- Item 1\n- [Link](url)",
			expected: "- Item 1\n- \n[[Link:url]]",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 標準入力をモック
			oldStdin := os.Stdin
			stdinR, stdinW, _ := os.Pipe()
			os.Stdin = stdinR

			// 標準出力をキャプチャ
			oldStdout := os.Stdout
			stdoutR, stdoutW, _ := os.Pipe()
			os.Stdout = stdoutW

			// 入力データを書き込み
			go func() {
				defer stdinW.Close()
				_, _ = io.WriteString(stdinW, tc.input)
			}()

			// コマンド実行
			rootCmd.SetArgs([]string{})
			err := rootCmd.Execute()
			if err != nil {
				t.Fatalf("Command execution failed: %v", err)
			}

			// 標準出力を復元
			stdoutW.Close()
			os.Stdout = oldStdout
			os.Stdin = oldStdin

			// 出力を読み取り
			var buf strings.Builder
			if _, err := io.Copy(&buf, stdoutR); err != nil {
				t.Fatalf("Failed to read output: %v", err)
			}
			output := strings.TrimSpace(buf.String())

			if output != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, output)
			}
		})
	}
}

// TestVersionIntegration はバージョンコマンドの統合テストを実行する
func TestVersionIntegration(t *testing.T) {
	defer resetRootCmd()

	output := captureStdout(t, func() {
		rootCmd.SetArgs([]string{"--version"})
		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("Version command failed: %v", err)
		}
	})

	expectedVersion := "md2backlog version 0.1.0"
	if output != expectedVersion {
		t.Errorf("Expected %q, got %q", expectedVersion, output)
	}
}

// TestBasicIntegration は基本的な統合テストを実行する
func TestBasicIntegration(t *testing.T) {
	defer resetRootCmd()

	t.Run("successful conversion with valid input", func(t *testing.T) {
		// 一時ディレクトリを作成
		tmpDir, err := os.MkdirTemp("", "md2backlog-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		// 入力ファイルを作成
		inputContent := "# Test Heading\n\n**Bold** and *italic* text"
		expectedOutput := "* Test Heading\n''Bold'' and '''italic''' text"
		inputPath := filepath.Join(tmpDir, "test_input.md")
		if err := os.WriteFile(inputPath, []byte(inputContent), 0644); err != nil {
			t.Fatalf("Failed to write input file: %v", err)
		}

		// 出力ファイルパス
		outputPath := filepath.Join(tmpDir, "test_output.txt")

		// 独立したコマンドインスタンスを作成
		var testInputFile, testOutputFile string
		testCmd := &cobra.Command{
			Use:     "md2backlog",
			Short:   "Convert Markdown to Backlog notation",
			Version: version,
			Run: func(cmd *cobra.Command, args []string) {
				runConvertTest(testInputFile, testOutputFile)
			},
		}
		testCmd.Flags().StringVarP(&testInputFile, "input", "i", "", "Input file")
		testCmd.Flags().StringVarP(&testOutputFile, "output", "o", "", "Output file")

		// コマンド実行
		testCmd.SetArgs([]string{"-i", inputPath, "-o", outputPath})
		err = testCmd.Execute()
		if err != nil {
			t.Fatalf("Command execution failed: %v", err)
		}

		// 出力ファイルが作成されているか確認
		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			t.Fatalf("Output file was not created at path: %s", outputPath)
		}

		// 出力ファイルを確認
		outputBytes, err := os.ReadFile(outputPath)
		if err != nil {
			t.Fatalf("Failed to read output file: %v", err)
		}
		output := strings.TrimSpace(string(outputBytes))

		if output != expectedOutput {
			t.Errorf("Expected %q, got %q", expectedOutput, output)
		}
	})
}

// runConvertTest は統合テスト用の変換実行関数
func runConvertTest(inputFile, outputFile string) {
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

// resetFlags はテスト用のフラグリセット関数
func resetFlags() {
	inputFile = ""
	outputFile = ""
}

// resetRootCmd はrootCmdを初期状態にリセットする
func resetRootCmd() {
	resetFlags()
	rootCmd.SetArgs(nil)
	rootCmd.ResetFlags()
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input Markdown file (default: stdin)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: stdout)")
}

// runFileConversionTest はファイル変換の統合テストを実行し、結果を返す
func runFileConversionTest(t *testing.T, inputContent string) string {
	// 一時ディレクトリを作成
	tmpDir, err := os.MkdirTemp("", "md2backlog-integration-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 入力ファイルを作成
	inputPath := filepath.Join(tmpDir, "input.md")
	if err := os.WriteFile(inputPath, []byte(inputContent), 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	// 出力ファイルパス
	outputPath := filepath.Join(tmpDir, "output.txt")

	// コマンド実行
	rootCmd.SetArgs([]string{"-i", inputPath, "-o", outputPath})
	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("Command execution failed: %v", err)
	}

	// 出力ファイルを読み取り
	outputBytes, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	return string(outputBytes)
}

// captureStdout は標準出力をキャプチャして返す
func captureStdout(t *testing.T, fn func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = oldStdout

	var buf strings.Builder
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("Failed to read stdout: %v", err)
	}
	return strings.TrimSpace(buf.String())
}
