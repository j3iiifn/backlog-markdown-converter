package converter

import (
	"testing"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "空文字列",
			input:    "",
			expected: "",
			hasError: false,
		},
		{
			name:     "通常のテキスト",
			input:    "Hello, World!",
			expected: "Hello, World!",
			hasError: false,
		},
		{
			name:     "複数行テキスト",
			input:    "Line 1\nLine 2\nLine 3",
			expected: "Line 1\nLine 2\nLine 3",
			hasError: false,
		},
		{
			name:     "H1見出し変換",
			input:    "# 見出し1",
			expected: "* 見出し1",
			hasError: false,
		},
		{
			name:     "H2見出し変換",
			input:    "## 見出し2",
			expected: "** 見出し2",
			hasError: false,
		},
		{
			name:     "H3見出し変換",
			input:    "### 見出し3",
			expected: "*** 見出し3",
			hasError: false,
		},
		{
			name:     "複数見出し変換",
			input:    "# 見出し1\n## 見出し2\n### 見出し3",
			expected: "* 見出し1\n** 見出し2\n*** 見出し3",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Convert(tt.input)

			if tt.hasError && err == nil {
				t.Errorf("期待されたエラーが発生しませんでした")
				return
			}

			if !tt.hasError && err != nil {
				t.Errorf("予期しないエラーが発生しました: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("期待値: %q, 実際の値: %q", tt.expected, result)
			}
		})
	}
}
