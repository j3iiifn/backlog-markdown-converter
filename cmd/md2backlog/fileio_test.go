package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileInputOutput(t *testing.T) {
	// テスト前にフラグをリセット
	inputFile = ""
	outputFile = ""
	defer func() {
		inputFile = ""
		outputFile = ""
	}()

	tests := []struct {
		name         string
		inputContent string
		expected     string
	}{
		{
			name:         "simple markdown conversion",
			inputContent: "# Heading\n\n**Bold** text.",
			expected:     "* Heading\n''Bold'' text.",
		},
		{
			name:         "complex markdown",
			inputContent: "# Heading\n\n- List item\n- Another item\n\n`inline code`",
			expected:     "* Heading\n- List item\n- Another item\n{code}inline code{/code}",
		},
		{
			name:         "empty file",
			inputContent: "",
			expected:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 一時ディレクトリを作成
			tmpDir, err := os.MkdirTemp("", "md2backlog-test")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// 入力ファイルを作成
			inputFile := filepath.Join(tmpDir, "input.md")
			if err := os.WriteFile(inputFile, []byte(tt.inputContent), 0644); err != nil {
				t.Fatalf("Failed to write input file: %v", err)
			}

			// 出力ファイルパス
			outputFile := filepath.Join(tmpDir, "output.txt")

			// コマンド実行
			rootCmd.SetArgs([]string{"-i", inputFile, "-o", outputFile})
			err = rootCmd.Execute()
			if err != nil {
				t.Fatalf("Command execution failed: %v", err)
			}

			// 出力ファイルを読み取り
			outputBytes, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read output file: %v", err)
			}
			output := strings.TrimSpace(string(outputBytes))

			if output != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, output)
			}
		})
	}
}

func TestInputFlagOnly(t *testing.T) {
	// テスト前にフラグをリセット
	inputFile = ""
	outputFile = ""
	defer func() {
		inputFile = ""
		outputFile = ""
	}()
	// 一時ディレクトリを作成
	tmpDir, err := os.MkdirTemp("", "md2backlog-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 入力ファイルを作成
	inputContent := "# Test Heading\n\n**Bold** text."
	expectedOutput := "* Test Heading\n''Bold'' text."
	inputFile := filepath.Join(tmpDir, "input.md")
	if err := os.WriteFile(inputFile, []byte(inputContent), 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	// 標準出力をキャプチャ
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// コマンド実行
	rootCmd.SetArgs([]string{"-i", inputFile})
	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("Command execution failed: %v", err)
	}

	// 標準出力を復元
	w.Close()
	os.Stdout = oldStdout

	// 出力を読み取り
	var buf strings.Builder
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}
	output := strings.TrimSpace(buf.String())

	if output != expectedOutput {
		t.Errorf("Expected %q, got %q", expectedOutput, output)
	}
}

func TestOutputFlagOnly(t *testing.T) {
	// テスト前にフラグをリセット
	inputFile = ""
	outputFile = ""
	defer func() {
		inputFile = ""
		outputFile = ""
	}()
	// 一時ディレクトリを作成
	tmpDir, err := os.MkdirTemp("", "md2backlog-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 出力ファイルパス
	outputFile := filepath.Join(tmpDir, "output.txt")
	inputContent := "# Test Output\n\n*italic* text."
	expectedOutput := "* Test Output\n'''italic''' text."

	// 標準入力をモック
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r

	// 入力データを書き込み
	go func() {
		defer w.Close()
		_, _ = io.WriteString(w, inputContent)
	}()

	// コマンド実行
	rootCmd.SetArgs([]string{"-o", outputFile})
	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("Command execution failed: %v", err)
	}

	// 標準入力を復元
	os.Stdin = oldStdin

	// 出力ファイルを読み取り
	outputBytes, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	output := strings.TrimSpace(string(outputBytes))

	if output != expectedOutput {
		t.Errorf("Expected %q, got %q", expectedOutput, output)
	}
}
