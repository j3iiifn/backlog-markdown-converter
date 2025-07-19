# md2backlog実装タスクリスト

## プロジェクト概要
Go製CLIツール「md2backlog」をTDD（Test Driven Development）で実装する。
MarkdownテキストをAST解析によりBacklog記法に変換するツール。

## 実装方針
- **TDD**: Red-Green-Refactorサイクルを厳守
- **コミット粒度**: 各ステップ完了時または意味のある機能追加時
- **品質保証**: コミット前に必ず`goimports`と`golangci-lint`実行

---

## Phase 0: プロジェクトセットアップ

### ✅ Task 0.1: 基本セットアップ
- [x] `go mod init`でGoモジュール初期化
- [x] `go get github.com/spf13/cobra`でCLIフレームワーク追加
- [x] `go get github.com/yuin/goldmark`でMarkdownパーサー追加
- [x] 基本ディレクトリ構造作成（`cmd/`, `internal/`）

---

## Phase 1: 基本CLI機能（高優先度）

### ✅ Task 1.1: バージョンコマンド実装
**目標**: `md2backlog --version`でバージョン番号表示
- [x] **(Red)** バージョン表示のテスト作成
- [x] **(Green)** Cobraでバージョンコマンド実装
- [x] **(Refactor)** コード整理

### ✅ Task 1.2: 標準入出力パススルー
**目標**: 標準入力の文字列をそのまま標準出力に返す
- [x] **(Red)** パススルーのテスト作成
- [x] **(Green)** 入力読み取り・出力ロジック実装
- [x] **(Refactor)** コード整理
- [x] **(Commit)**

### ⬜ Task 1.3: 変換関数骨格作成
**目標**: `Convert(markdown string) (string, error)`関数の雛形
- [ ] **(Red)** Convert関数の基本テスト作成
- [ ] **(Green)** 文字列をそのまま返すConvert関数実装
- [ ] **(Refactor)** パッケージ構造整理
- [ ] **(Commit)**

---

## Phase 2: 基本変換ルール（中優先度）

### ⬜ Task 2.1: 見出し変換
**変換ルール**: `# H1` → `* H1`, `## H2` → `** H2`
- [ ] **(Red)** 見出し変換テスト作成
- [ ] **(Green)** goldmark AST使用の見出し変換実装
- [ ] **(Refactor)** AST Walker パターン整理
- [ ] **(Commit)**

### ⬜ Task 2.2: 太字変換
**変換ルール**: `**bold**` → `''bold''`
- [ ] **(Red)** 太字変換テスト作成
- [ ] **(Green)** 太字変換ロジック実装
- [ ] **(Refactor)** コード整理
- [ ] **(Commit)**

### ⬜ Task 2.3: 斜体変換
**変換ルール**: `*italic*` → `'''italic'''`
- [ ] **(Red)** 斜体変換テスト作成
- [ ] **(Green)** 斜体変換ロジック実装
- [ ] **(Refactor)** コード整理
- [ ] **(Commit)**

### ⬜ Task 2.4: 打ち消し線変換
**変換ルール**: `~~strike~~` → `%%strike%%`
- [ ] **(Red)** 打ち消し線変換テスト作成
- [ ] **(Green)** 打ち消し線変換ロジック実装
- [ ] **(Refactor)** コード整理
- [ ] **(Commit)**

### ⬜ Task 2.5: 単一階層箇条書きリスト
**変換ルール**: `- item` または `* item` → `- item`
- [ ] **(Red)** 単一リスト変換テスト作成
- [ ] **(Green)** リスト変換ロジック実装
- [ ] **(Refactor)** コード整理
- [ ] **(Commit)**

### ⬜ Task 2.6: ネスト箇条書きリスト
**変換ルール**: `- L1\n  - L2\n    - L3` → `- L1\n-- L2\n--- L3`
- [ ] **(Red)** ネストリスト変換テスト作成
- [ ] **(Green)** ネストレベル検出・変換実装
- [ ] **(Refactor)** コード整理
- [ ] **(Commit)**

### ⬜ Task 2.7: 番号付きリスト
**変換ルール**: `1. item` → `+ item`（ネストは平坦化）
- [ ] **(Red)** 番号付きリスト変換テスト作成
- [ ] **(Green)** 番号付きリスト変換実装
- [ ] **(Refactor)** コード整理
- [ ] **(Commit)**

### ⬜ Task 2.8: リンク変換
**変換ルール**: `[Link](URL)` → `[[Link:URL]]`
- [ ] **(Red)** リンク変換テスト作成
- [ ] **(Green)** リンク変換ロジック実装
- [ ] **(Refactor)** コード整理
- [ ] **(Commit)**

### ⬜ Task 2.9: インラインコード変換
**変換ルール**: `` `inline code` `` → `{code}inline code{/code}`
- [ ] **(Red)** インラインコード変換テスト作成
- [ ] **(Green)** インラインコード変換実装
- [ ] **(Refactor)** コード整理
- [ ] **(Commit)**

### ⬜ Task 2.10: コードブロック変換
**変換ルール**: `\`\`\`lang\ncode\n\`\`\`` → `>{code:lang}\ncode\n{/code}<`
- [ ] **(Red)** コードブロック変換テスト作成
- [ ] **(Green)** コードブロック変換実装
- [ ] **(Refactor)** コード整理
- [ ] **(Commit)**

---

## Phase 3: 追加機能（低優先度）

### ⬜ Task 3.1: 引用処理
**変換ルール**: `> Quote` → `> Quote`（変更なし）
- [ ] **(Red)** 引用処理テスト作成
- [ ] **(Green)** 引用をそのまま通すロジック実装
- [ ] **(Refactor)** コード整理
- [ ] **(Commit)**

### ⬜ Task 3.2: テーブル変換
**変換ルール**: `|th1|...\n|---|...\n|td1|...` → `|*th1|...\n|td1|...`
- [ ] **(Red)** テーブル変換テスト作成
- [ ] **(Green)** テーブルヘッダー変換実装
- [ ] **(Refactor)** コード整理
- [ ] **(Commit)**

### ⬜ Task 3.3: 非対応記法変換
**変換ルール**: 
- `- [x] task` → `- task`（チェックボックス除去）
- HTML タグ除去
- `text[^1]` → `text[1]`（脚注変換）
- [ ] **(Red)** 非対応記法変換テスト作成
- [ ] **(Green)** 各非対応記法の変換実装
- [ ] **(Refactor)** コード整理
- [ ] **(Commit)**

---

## Phase 4: 統合・品質保証（中優先度）

### ⬜ Task 4.1: ファイル入出力統合
**目標**: `-i`, `-o`フラグでファイル入出力対応
- [ ] **(Red)** ファイル入出力テスト作成
- [ ] **(Green)** ファイル読み書きロジック統合
- [ ] **(Refactor)** エラーハンドリング整理
- [ ] **(Commit)**

### ⬜ Task 4.2: 統合テスト追加
**目標**: CLI全体の動作確認
- [ ] **(Red)** エンドツーエンドテスト作成
- [ ] **(Green)** 実際のファイル変換テスト実装
- [ ] **(Refactor)** テストケース整理
- [ ] **(Commit)**

---

## Phase 5: CI/CD・ドキュメント（低優先度）

### ⬜ Task 5.1: GitHub Actions CI設定
**ファイル**: `.github/workflows/ci.yml`
- [ ] gofmt チェック
- [ ] golangci-lint 実行
- [ ] go test 実行
- [ ] **(Commit)**

### ⬜ Task 5.2: golangci-lint設定
**ファイル**: `.golangci.yml`
- [ ] リンタールール定義
- [ ] **(Commit)**

### ⬜ Task 5.3: GoReleaser設定
**目標**: タグプッシュでクロスプラットフォームバイナリリリース
- [ ] Windows、macOS、Linux対応
- [ ] **(Commit)**

### ⬜ Task 5.4: ドキュメント整備
**ファイル**: 
- [ ] `README.md`: プロジェクト概要、インストール、使用方法
- [ ] `LICENSE`: MITまたはApache 2.0
- [ ] `.gitignore`: Go プロジェクト用設定
- [ ] **(Commit)**

---

## 品質チェックリスト（各コミット前必須）

1. [ ] `goimports -w .` でコードフォーマット
2. [ ] `golangci-lint run` で静的解析
3. [ ] `go test ./...` で全テスト実行
4. [ ] Conventional Commits形式でコミット

## コミット例
```bash
feat: add version command support
test: add unit tests for heading conversion
refactor: extract AST walker pattern
docs: update README with usage examples
```