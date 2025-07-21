package converter

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	gast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
)

// Convert はMarkdownテキストをBacklog記法に変換します
func Convert(markdown string) (string, error) {
	if markdown == "" {
		return "", nil
	}

	// goldmarkでMarkdownをパース（GFM拡張を有効化）
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))
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

		case *ast.Emphasis:
			if entering {
				writeEmphasis(&buffer, node, source)
				return ast.WalkSkipChildren, nil
			}

		case *gast.Strikethrough:
			if entering {
				writeStrikethrough(&buffer, node, source)
				return ast.WalkSkipChildren, nil
			}

		case *ast.List:
			// リストは子要素（ListItem）の処理に任せる
			// 何もしない

		case *ast.ListItem:
			if entering {
				writeListItem(&buffer, node, source)
				// ネストリストを含む可能性があるので、子要素も処理
				return ast.WalkContinue, nil
			}

		case *ast.Link:
			if entering {
				writeLink(&buffer, node, source)
				return ast.WalkSkipChildren, nil
			}

		case *ast.CodeSpan:
			if entering {
				writeCodeSpan(&buffer, node, source)
				return ast.WalkSkipChildren, nil
			}

		case *ast.FencedCodeBlock:
			if entering {
				writeFencedCodeBlock(&buffer, node, source)
				return ast.WalkSkipChildren, nil
			}

		case *ast.Blockquote:
			if entering {
				writeBlockquote(&buffer, node, source)
				return ast.WalkSkipChildren, nil
			}

		case *gast.Table:
			if entering {
				writeTable(&buffer, node, source)
				return ast.WalkSkipChildren, nil
			}

		case *ast.Text:
			if entering && !isChildOfHeading(node) && !isChildOfEmphasis(node) && !isChildOfStrikethrough(node) && !isChildOfListItem(node) && !isChildOfLink(node) && !isChildOfCodeSpan(node) && !isChildOfFencedCodeBlock(node) && !isChildOfBlockquote(node) && !isChildOfTable(node) {
				writeText(&buffer, node, source)
			}

		case *ast.Paragraph:
			if !entering && node.NextSibling() != nil && !isNextSiblingList(node) && !isChildOfBlockquote(node) && !isChildOfListItem(node) {
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

// writeEmphasis は太字・斜体ノードをBacklog記法で出力します
func writeEmphasis(buffer *bytes.Buffer, emphasis *ast.Emphasis, source []byte) {
	switch emphasis.Level {
	case 2:
		// 太字の場合
		buffer.WriteString("''")
		// 太字内のテキスト内容を取得
		for child := emphasis.FirstChild(); child != nil; child = child.NextSibling() {
			if textNode, ok := child.(*ast.Text); ok {
				buffer.Write(textNode.Segment.Value(source))
			}
		}
		buffer.WriteString("''")
	case 1:
		// 斜体の場合
		buffer.WriteString("'''")
		// 斜体内のテキスト内容を取得
		for child := emphasis.FirstChild(); child != nil; child = child.NextSibling() {
			if textNode, ok := child.(*ast.Text); ok {
				buffer.Write(textNode.Segment.Value(source))
			}
		}
		buffer.WriteString("'''")
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

// isChildOfEmphasis はノードが太字/斜体の子要素かどうかを判定します
func isChildOfEmphasis(node ast.Node) bool {
	parent := node.Parent()
	for parent != nil {
		if _, ok := parent.(*ast.Emphasis); ok {
			return true
		}
		parent = parent.Parent()
	}
	return false
}

// writeStrikethrough は打ち消し線ノードをBacklog記法で出力します
func writeStrikethrough(buffer *bytes.Buffer, strikethrough ast.Node, source []byte) {
	buffer.WriteString("%%")
	// 打ち消し線内のテキスト内容を取得
	for child := strikethrough.FirstChild(); child != nil; child = child.NextSibling() {
		if textNode, ok := child.(*ast.Text); ok {
			buffer.Write(textNode.Segment.Value(source))
		}
	}
	buffer.WriteString("%%")
}

// isChildOfStrikethrough はノードが打ち消し線の子要素かどうかを判定します
func isChildOfStrikethrough(node ast.Node) bool {
	parent := node.Parent()
	for parent != nil {
		if _, ok := parent.(*gast.Strikethrough); ok {
			return true
		}
		parent = parent.Parent()
	}
	return false
}

// writeListItem はリストアイテムノードをBacklog記法で出力します
func writeListItem(buffer *bytes.Buffer, listItem *ast.ListItem, source []byte) {
	// ネストレベルを計算
	nestLevel := calculateListNestLevel(listItem)

	// 親リストが番号付きリストかどうかを判定
	isOrderedList := isInOrderedList(listItem)

	// プレフィックスを生成
	var prefix string
	if isOrderedList {
		// 番号付きリストの場合は「+」を使用（ネスト関係なく平坦化）
		prefix = "+"
	} else {
		// 通常のリストの場合はネストレベルに応じて「-」を繰り返し
		prefix = strings.Repeat("-", nestLevel)
	}
	buffer.WriteString(prefix + " ")

	// リストアイテムの最初のテキスト内容のみを取得（ネストは別処理）
	for child := listItem.FirstChild(); child != nil; child = child.NextSibling() {
		switch childNode := child.(type) {
		case *ast.Paragraph:
			// Paragraphの子要素（Text）を処理
			for grandChild := childNode.FirstChild(); grandChild != nil; grandChild = grandChild.NextSibling() {
				if textNode, ok := grandChild.(*ast.Text); ok {
					buffer.Write(textNode.Segment.Value(source))
				}
			}
			// 最初のParagraphのみ処理して終了
			buffer.WriteString("\n")
			return
		case *ast.TextBlock:
			// TextBlockの子要素（Text）を処理
			for grandChild := childNode.FirstChild(); grandChild != nil; grandChild = grandChild.NextSibling() {
				if textNode, ok := grandChild.(*ast.Text); ok {
					buffer.Write(textNode.Segment.Value(source))
				}
			}
			// 最初のTextBlockのみ処理して終了
			buffer.WriteString("\n")
			return
		}
	}
	buffer.WriteString("\n")
}

// calculateListNestLevel はリストアイテムのネストレベルを計算します
func calculateListNestLevel(listItem *ast.ListItem) int {
	level := 1
	parent := listItem.Parent()

	for parent != nil {
		// 親がListの場合、さらにその親のListItemを探す
		if _, ok := parent.(*ast.List); ok {
			grandParent := parent.Parent()
			if _, ok := grandParent.(*ast.ListItem); ok {
				level++
				parent = grandParent.Parent()
			} else {
				break
			}
		} else {
			parent = parent.Parent()
		}
	}

	return level
}

// isChildOfListItem はノードがリストアイテムの子要素かどうかを判定します
func isChildOfListItem(node ast.Node) bool {
	parent := node.Parent()
	for parent != nil {
		if _, ok := parent.(*ast.ListItem); ok {
			return true
		}
		parent = parent.Parent()
	}
	return false
}

// isNextSiblingList は次の兄弟ノードがListかどうかを判定します
func isNextSiblingList(node ast.Node) bool {
	next := node.NextSibling()
	_, ok := next.(*ast.List)
	return ok
}

// isInOrderedList はリストアイテムが番号付きリストの中にあるかどうかを判定します
func isInOrderedList(listItem *ast.ListItem) bool {
	parent := listItem.Parent()
	for parent != nil {
		if list, ok := parent.(*ast.List); ok {
			return list.IsOrdered()
		}
		parent = parent.Parent()
	}
	return false
}

// writeLink はリンクノードをBacklog記法で出力します
func writeLink(buffer *bytes.Buffer, link *ast.Link, source []byte) {
	buffer.WriteString("[[")

	// リンクテキストを取得
	for child := link.FirstChild(); child != nil; child = child.NextSibling() {
		if textNode, ok := child.(*ast.Text); ok {
			buffer.Write(textNode.Segment.Value(source))
		}
	}

	buffer.WriteString(":")
	buffer.Write(link.Destination)
	buffer.WriteString("]]")
}

// isChildOfLink はノードがリンクの子要素かどうかを判定します
func isChildOfLink(node ast.Node) bool {
	parent := node.Parent()
	for parent != nil {
		if _, ok := parent.(*ast.Link); ok {
			return true
		}
		parent = parent.Parent()
	}
	return false
}

// writeCodeSpan はインラインコードノードをBacklog記法で出力します
func writeCodeSpan(buffer *bytes.Buffer, codeSpan *ast.CodeSpan, source []byte) {
	buffer.WriteString("{code}")
	// インラインコード内のテキスト内容を取得
	for child := codeSpan.FirstChild(); child != nil; child = child.NextSibling() {
		if textNode, ok := child.(*ast.Text); ok {
			buffer.Write(textNode.Segment.Value(source))
		}
	}
	buffer.WriteString("{/code}")
}

// isChildOfCodeSpan はノードがインラインコードの子要素かどうかを判定します
func isChildOfCodeSpan(node ast.Node) bool {
	parent := node.Parent()
	for parent != nil {
		if _, ok := parent.(*ast.CodeSpan); ok {
			return true
		}
		parent = parent.Parent()
	}
	return false
}

// writeFencedCodeBlock はコードブロックノードをBacklog記法で出力します
func writeFencedCodeBlock(buffer *bytes.Buffer, codeBlock *ast.FencedCodeBlock, source []byte) {
	// 開始タグ
	buffer.WriteString(">{code")

	// 言語指定がある場合
	if codeBlock.Language(source) != nil {
		lang := string(codeBlock.Language(source))
		if lang != "" {
			buffer.WriteString(":" + lang)
		}
	}
	buffer.WriteString("}\n")

	// コードブロックの内容を出力
	for i := 0; i < codeBlock.Lines().Len(); i++ {
		line := codeBlock.Lines().At(i)
		buffer.Write(line.Value(source))
	}

	// 終了タグ
	buffer.WriteString("{/code}<")

	// 次の兄弟ノードがある場合は改行を追加
	if codeBlock.NextSibling() != nil {
		buffer.WriteString("\n")
	}
}

// isChildOfFencedCodeBlock はノードがコードブロックの子要素かどうかを判定します
func isChildOfFencedCodeBlock(node ast.Node) bool {
	parent := node.Parent()
	for parent != nil {
		if _, ok := parent.(*ast.FencedCodeBlock); ok {
			return true
		}
		parent = parent.Parent()
	}
	return false
}

// writeBlockquote は引用ノードをBacklog記法で出力します
func writeBlockquote(buffer *bytes.Buffer, blockquote *ast.Blockquote, source []byte) {
	for child := blockquote.FirstChild(); child != nil; child = child.NextSibling() {
		if paragraph, ok := child.(*ast.Paragraph); ok {
			// パラグラフ内の各Textノードを別々の行として処理
			isFirstText := true
			for textChild := paragraph.FirstChild(); textChild != nil; textChild = textChild.NextSibling() {
				if textNode, ok := textChild.(*ast.Text); ok {
					if !isFirstText {
						buffer.WriteString("\n")
					}

					text := string(textNode.Segment.Value(source))
					buffer.WriteString("> " + text)

					isFirstText = false
				}
			}
		} else if nestedBlockquote, ok := child.(*ast.Blockquote); ok {
			// ネストした引用を処理（改行を追加）
			buffer.WriteString("\n")
			processNestedBlockquote(buffer, nestedBlockquote, source, 2)
		}
	}

	// 引用ブロックの後に続く要素がある場合は改行を追加
	if blockquote.NextSibling() != nil {
		buffer.WriteString("\n\n")
	}
}

// processNestedBlockquote はネストした引用を処理します
func processNestedBlockquote(buffer *bytes.Buffer, blockquote *ast.Blockquote, source []byte, level int) {
	for child := blockquote.FirstChild(); child != nil; child = child.NextSibling() {
		if paragraph, ok := child.(*ast.Paragraph); ok {
			// パラグラフ内のテキストを取得
			var textContent strings.Builder
			for textChild := paragraph.FirstChild(); textChild != nil; textChild = textChild.NextSibling() {
				if textNode, ok := textChild.(*ast.Text); ok {
					textContent.Write(textNode.Segment.Value(source))
				}
			}

			// 引用プレフィックスを生成
			prefix := strings.Repeat("> ", level)
			buffer.WriteString(prefix + textContent.String())

			// 次の子要素がある場合は改行を追加
			if child.NextSibling() != nil {
				buffer.WriteString("\n")
			}
		}
	}
}

// isChildOfBlockquote はノードが引用の子要素かどうかを判定します
func isChildOfBlockquote(node ast.Node) bool {
	parent := node.Parent()
	for parent != nil {
		if _, ok := parent.(*ast.Blockquote); ok {
			return true
		}
		parent = parent.Parent()
	}
	return false
}

// writeTable はテーブルノードをBacklog記法で出力します
func writeTable(buffer *bytes.Buffer, table ast.Node, source []byte) {
	gfmTable := table.(*gast.Table)

	// テーブルの子要素を処理
	for child := gfmTable.FirstChild(); child != nil; child = child.NextSibling() {
		switch childNode := child.(type) {
		case *gast.TableHeader:
			// ヘッダー行を出力
			writeTableHeader(buffer, childNode, source)
		case *gast.TableRow:
			// データ行を出力
			writeTableRow(buffer, childNode, source)
		}
	}

	// テーブルの後に続く要素がある場合は改行を追加
	if table.NextSibling() != nil {
		buffer.WriteString("\n")
	}
}

// writeTableHeader はテーブルヘッダーを出力します
func writeTableHeader(buffer *bytes.Buffer, tableHeader *gast.TableHeader, source []byte) {
	buffer.WriteString("|")

	// ヘッダーセルを直接処理
	for cell := tableHeader.FirstChild(); cell != nil; cell = cell.NextSibling() {
		if tableCell, ok := cell.(*gast.TableCell); ok {
			buffer.WriteString("*")

			// セル内のテキストを取得
			for cellChild := tableCell.FirstChild(); cellChild != nil; cellChild = cellChild.NextSibling() {
				if textNode, ok := cellChild.(*ast.Text); ok {
					text := strings.TrimSpace(string(textNode.Segment.Value(source)))
					buffer.WriteString(text)
				}
			}

			buffer.WriteString("|")
		}
	}
	buffer.WriteString("\n")
}

// writeTableRow はテーブル行を出力します
func writeTableRow(buffer *bytes.Buffer, tableRow *gast.TableRow, source []byte) {
	buffer.WriteString("|")

	// セル内容を出力
	for cell := tableRow.FirstChild(); cell != nil; cell = cell.NextSibling() {
		if tableCell, ok := cell.(*gast.TableCell); ok {
			// セル内のテキストを取得
			for cellChild := tableCell.FirstChild(); cellChild != nil; cellChild = cellChild.NextSibling() {
				if textNode, ok := cellChild.(*ast.Text); ok {
					text := strings.TrimSpace(string(textNode.Segment.Value(source)))
					buffer.WriteString(text)
				}
			}

			buffer.WriteString("|")
		}
	}
	buffer.WriteString("\n")
}

// isChildOfTable はノードがテーブルの子要素かどうかを判定します
func isChildOfTable(node ast.Node) bool {
	parent := node.Parent()
	for parent != nil {
		if _, ok := parent.(*gast.Table); ok {
			return true
		}
		if _, ok := parent.(*gast.TableRow); ok {
			return true
		}
		if _, ok := parent.(*gast.TableCell); ok {
			return true
		}
		parent = parent.Parent()
	}
	return false
}
