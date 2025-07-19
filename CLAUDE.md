# CLAUDE.md

このファイルは、このリポジトリでコードを扱う際にClaude Code (claude.ai/code)へのガイダンスを提供します。

## 言語設定
このプロジェクトでの全ての会話とドキュメントは**日本語**で行ってください。

## プロジェクト概要

これは`md2backlog`というGo製CLIツールで、MarkdownテキストをASTパースを使用してBacklog記法に変換します。このプロジェクトはテスト駆動開発（TDD）の原則に従い、オープンソース公開を目的として設計されています。

## 開発コマンド

### セットアップ
```bash
go mod init                    # Goモジュールの初期化
go get github.com/spf13/cobra  # CLIフレームワークの追加
go get github.com/yuin/goldmark # Markdownパーサーの追加
```

### ビルド & 実行
```bash
go build -o md2backlog ./cmd/md2backlog  # バイナリのビルド
./md2backlog --version                   # バージョンコマンドのテスト
./md2backlog -i input.md -o output.txt   # ファイルの変換
cat input.md | ./md2backlog              # 標準入力からの変換
```

### コード品質チェック
```bash
goimports -w .          # コードフォーマット（コミット前必須）
golangci-lint run       # 静的解析（コミット前必須）
go test ./...           # 全テストの実行
go test -v ./internal/converter  # 特定パッケージのテスト実行
```

## Git コミット規則

### Conventional Commits形式を必須とする
```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### 使用するtype
- `feat`: 新機能
- `fix`: バグ修正
- `docs`: ドキュメントのみの変更
- `style`: コードの意味に影響しない変更（空白、フォーマット、セミコロンなど）
- `refactor`: バグ修正や機能追加ではないコード変更
- `test`: テストの追加や既存テストの修正
- `chore`: ビルドプロセスや補助ツールの変更

### コミット例
```bash
feat: add version command support
test: add unit tests for version command
refactor: extract common CLI setup logic
docs: update README with installation instructions
```

### コミット前チェックリスト
各コミット前に以下を必ず実行：
1. `goimports -w .` でコードフォーマット
2. `golangci-lint run` で静的解析チェック
3. `go test ./...` でテスト実行
4. 全てパスしてからConventional Commits形式でコミット

### コミット粒度
- できるだけ細かく、意味のある単位でコミット
- 1つのコミットは1つの論理的変更のみ
- テストとプロダクションコードは可能な限り分離してコミット
- 各コミットは独立してビルド・テストが通る状態を保つ

## アーキテクチャ

### ディレクトリ構造
- `cmd/` - CLIエントリーポイントとcobraコマンド
- `internal/converter` - goldmark ASTを使用するコア変換ロジック
- `doc/md2backlog/` - プロジェクトドキュメントと要件

### 主要コンポーネント
- **CLIフレームワーク**: コマンド構造とフラグ処理にCobraを使用
- **Markdownパーサー**: 変換にAST walking (`ast.Walk`)を使用するgoldmark
- **コンバーター**: `internal/converter`パッケージのコアロジック、`Convert(markdown string) (string, error)`関数

### TDD開発プロセス
厳密なRed-Green-Refactorサイクルに従い、**各タスク完了時に必ずコミット**：
1. **Red**: 最小機能の失敗テストを書く
2. **Green**: テストを通す最小限のコードを書く
3. **Refactor**: テストを通した状態でクリーンアップ
4. **Commit**: 品質チェック後にConventional Commits形式でコミット

### 変換ルール実装順序
1. バージョンコマンド (`--version`)
2. 標準入出力のパススルー
3. 基本変換関数の骨格
4. 見出し (`#` → `*`)
5. テキスト装飾（太字、斜体、打ち消し線）
6. リスト（番号なし、ネスト、番号付き）
7. コードブロックとインラインコード
8. リンクその他の要素
9. 非対応フォーマットの処理

### テスト戦略
- `internal/converter`の各変換ルールの単体テスト
- CLI機能の統合テスト
- エッジケースとエラー条件のテストカバレッジ

### 品質要件
- 全コードはコミット前に`goimports`でフォーマット必須
- 全コードは`golangci-lint`チェックを通す必須
- 包括的なテストカバレッジが必要
- 意味のある変更は個別にコミット

### CI/CD設定
- `.github/workflows/ci.yml`のGitHub Actionsワークフロー
- 自動チェック: gofmt, golangci-lint, go test
- gitタグでのクロスプラットフォームバイナリリリース用GoReleaser

## タスク進捗管理

### 進捗管理ルール
`doc/md2backlog/implementation-tasks.md` でタスク進捗を管理する：

**タスクレベル**（メインタスク）:
- 未着手: `⬜ Task X.Y: タスク名`
- 完了: `✅ Task X.Y: タスク名`

**サブタスクレベル**（TDDステップ）:
- 未着手: `- [ ] (Red/Green/Refactor) 作業内容`
- 完了: `- [x] (Red/Green/Refactor) 作業内容`

### 進捗更新タイミング
- TDDの各ステップ（Red/Green/Refactor）完了時にサブタスクを更新
- メインタスク完了時にタスクレベルを更新
- コミット前に必ず進捗を最新状態に更新