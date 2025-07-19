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
		{
			name:     "太字変換",
			input:    "**太字のテキスト**",
			expected: "''太字のテキスト''",
			hasError: false,
		},
		{
			name:     "複数太字変換",
			input:    "**太字1** と **太字2**",
			expected: "''太字1'' と ''太字2''",
			hasError: false,
		},
		{
			name:     "太字と通常テキスト混合",
			input:    "通常のテキストと**太字**の混合",
			expected: "通常のテキストと''太字''の混合",
			hasError: false,
		},
		{
			name:     "斜体変換",
			input:    "*斜体のテキスト*",
			expected: "'''斜体のテキスト'''",
			hasError: false,
		},
		{
			name:     "複数斜体変換",
			input:    "*斜体1* と *斜体2*",
			expected: "'''斜体1''' と '''斜体2'''",
			hasError: false,
		},
		{
			name:     "斜体と通常テキスト混合",
			input:    "通常のテキストと*斜体*の混合",
			expected: "通常のテキストと'''斜体'''の混合",
			hasError: false,
		},
		{
			name:     "打ち消し線変換",
			input:    "~~打ち消し線のテキスト~~",
			expected: "%%打ち消し線のテキスト%%",
			hasError: false,
		},
		{
			name:     "複数打ち消し線変換",
			input:    "~~打ち消し1~~ と ~~打ち消し2~~",
			expected: "%%打ち消し1%% と %%打ち消し2%%",
			hasError: false,
		},
		{
			name:     "打ち消し線と通常テキスト混合",
			input:    "通常のテキストと~~打ち消し線~~の混合",
			expected: "通常のテキストと%%打ち消し線%%の混合",
			hasError: false,
		},
		{
			name:     "単一階層箇条書きリスト（ハイフン）",
			input:    "- アイテム1\n- アイテム2\n- アイテム3",
			expected: "- アイテム1\n- アイテム2\n- アイテム3",
			hasError: false,
		},
		{
			name:     "単一階層箇条書きリスト（アスタリスク）",
			input:    "* アイテム1\n* アイテム2\n* アイテム3",
			expected: "- アイテム1\n- アイテム2\n- アイテム3",
			hasError: false,
		},
		{
			name:     "単一階層箇条書きリスト混合記号",
			input:    "- アイテム1\n* アイテム2\n- アイテム3",
			expected: "- アイテム1\n- アイテム2\n- アイテム3",
			hasError: false,
		},
		{
			name:     "箇条書きリストと通常テキスト混合",
			input:    "通常のテキスト\n- アイテム1\n- アイテム2\n\n続きのテキスト",
			expected: "通常のテキスト\n- アイテム1\n- アイテム2\n続きのテキスト",
			hasError: false,
		},
		{
			name:     "ネスト箇条書きリスト（2階層）",
			input:    "- レベル1\n  - レベル2",
			expected: "- レベル1\n-- レベル2",
			hasError: false,
		},
		{
			name:     "ネスト箇条書きリスト（3階層）",
			input:    "- レベル1\n  - レベル2\n    - レベル3",
			expected: "- レベル1\n-- レベル2\n--- レベル3",
			hasError: false,
		},
		{
			name:     "ネスト箇条書きリスト（混合記号）",
			input:    "- レベル1\n  * レベル2\n    - レベル3",
			expected: "- レベル1\n-- レベル2\n--- レベル3",
			hasError: false,
		},
		{
			name:     "ネスト箇条書きリスト（複数項目）",
			input:    "- レベル1-1\n  - レベル2-1\n  - レベル2-2\n- レベル1-2\n  - レベル2-3",
			expected: "- レベル1-1\n-- レベル2-1\n-- レベル2-2\n- レベル1-2\n-- レベル2-3",
			hasError: false,
		},
		{
			name:     "番号付きリスト（単一階層）",
			input:    "1. アイテム1\n2. アイテム2\n3. アイテム3",
			expected: "+ アイテム1\n+ アイテム2\n+ アイテム3",
			hasError: false,
		},
		{
			name:     "番号付きリスト（ネスト）",
			input:    "1. レベル1\n   1. レベル2\n   2. レベル2-2\n2. レベル1-2",
			expected: "+ レベル1\n+ レベル2\n+ レベル2-2\n+ レベル1-2",
			hasError: false,
		},
		{
			name:     "番号付きリストと通常テキスト混合",
			input:    "通常のテキスト\n1. アイテム1\n2. アイテム2\n\n続きのテキスト",
			expected: "通常のテキスト\n+ アイテム1\n+ アイテム2\n続きのテキスト",
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
