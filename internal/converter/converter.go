package converter

import (
	"strings"
)

// Convert はMarkdownテキストをBacklog記法に変換します
func Convert(markdown string) (string, error) {
	if markdown == "" {
		return "", nil
	}

	// 行ごとに処理して変換
	lines := strings.Split(markdown, "\n")
	var result []string

	for _, line := range lines {
		converted := convertLine(line)
		result = append(result, converted)
	}

	return strings.Join(result, "\n"), nil
}

// convertLine は1行のMarkdownをBacklog記法に変換します
func convertLine(line string) string {
	// 見出し変換
	if strings.HasPrefix(line, "#") {
		return convertHeading(line)
	}

	// 他の変換ルールはここに追加（将来の実装）

	// 変換対象でない場合はそのまま返す
	return line
}

// convertHeading は見出し行をBacklog記法に変換します
func convertHeading(line string) string {
	level := 0
	for _, char := range line {
		if char == '#' {
			level++
		} else {
			break
		}
	}

	// # の後のテキストを取得（先頭空白も除去）
	text := strings.TrimSpace(line[level:])

	// Backlog記法に変換: # → *, ## → **, ### → ***
	prefix := strings.Repeat("*", level)
	return prefix + " " + text
}
