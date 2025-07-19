package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestPassthrough(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple text",
			input:    "Hello, World!",
			expected: "Hello, World!",
		},
		{
			name:     "multiline text",
			input:    "Line 1\nLine 2\nLine 3",
			expected: "Line 1\nLine 2\nLine 3",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "markdown text",
			input:    "# Heading\n\n**Bold** text and *italic* text.",
			expected: "# Heading\n\n**Bold** text and *italic* text.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 標準入力をモック
			oldStdin := os.Stdin
			r, w, _ := os.Pipe()
			os.Stdin = r

			// 標準出力をキャプチャ
			oldStdout := os.Stdout
			rOut, wOut, _ := os.Pipe()
			os.Stdout = wOut

			// 入力データを書き込み
			go func() {
				defer w.Close()
				if _, err := io.WriteString(w, tt.input); err != nil {
					t.Errorf("Failed to write to pipe: %v", err)
				}
			}()

			// コマンド実行
			rootCmd.SetArgs([]string{})
			err := rootCmd.Execute()
			if err != nil {
				t.Fatalf("Command execution failed: %v", err)
			}

			// 標準出力を復元
			wOut.Close()
			os.Stdout = oldStdout
			os.Stdin = oldStdin

			// 出力を読み取り
			var buf bytes.Buffer
			if _, err := io.Copy(&buf, rOut); err != nil {
				t.Fatalf("Failed to read output: %v", err)
			}
			output := strings.TrimSpace(buf.String())

			if output != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, output)
			}
		})
	}
}
