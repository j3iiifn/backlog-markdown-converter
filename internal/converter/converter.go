package converter

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Convert はMarkdownテキストをBacklog記法に変換します
func Convert(markdown string) (string, error) {
	if markdown == "" {
		return "", nil
	}

	// goldmarkでMarkdownをパース
	md := goldmark.New()
	reader := text.NewReader([]byte(markdown))
	document := md.Parser().Parse(reader)

	// ASTをウォークしてBacklog記法に変換
	var buffer bytes.Buffer
	source := reader.Source()

	err := ast.Walk(document, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		switch node := n.(type) {
		case *ast.Heading:
			if entering {
				writeHeading(&buffer, node, source)
				return ast.WalkSkipChildren, nil
			}

		case *ast.Text:
			if entering && !isChildOfHeading(node) {
				writeText(&buffer, node, source)
			}

		case *ast.Paragraph:
			if !entering && node.NextSibling() != nil {
				buffer.WriteString("\n")
			}
		}

		return ast.WalkContinue, nil
	})

	if err != nil {
		return "", err
	}

	result := buffer.String()
	// 末尾の不要な改行を除去
	result = strings.TrimSuffix(result, "\n")

	return result, nil
}

// writeHeading は見出しノードをBacklog記法で出力します
func writeHeading(buffer *bytes.Buffer, heading *ast.Heading, source []byte) {
	level := heading.Level
	prefix := strings.Repeat("*", level)
	buffer.WriteString(prefix + " ")

	// 見出しのテキスト内容を取得
	for child := heading.FirstChild(); child != nil; child = child.NextSibling() {
		if textNode, ok := child.(*ast.Text); ok {
			buffer.Write(textNode.Segment.Value(source))
		}
	}
	buffer.WriteString("\n")
}

// writeText はテキストノードを出力します（改行も含めて）
func writeText(buffer *bytes.Buffer, textNode *ast.Text, source []byte) {
	segment := textNode.Segment
	buffer.Write(segment.Value(source))

	// セグメント後に改行があるかチェック
	if segment.Stop < len(source) && source[segment.Stop] == '\n' {
		buffer.WriteString("\n")
	}
}

// isChildOfHeading はノードが見出しの子要素かどうかを判定します
func isChildOfHeading(node ast.Node) bool {
	parent := node.Parent()
	for parent != nil {
		if _, ok := parent.(*ast.Heading); ok {
			return true
		}
		parent = parent.Parent()
	}
	return false
}
