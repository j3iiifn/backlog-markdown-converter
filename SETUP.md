# Go開発環境構築手順書 (Claude Code対応)

このドキュメントでは、macOS上でDockerとVS Code Dev Containersを使用してGo言語の開発環境をサンドボックス内に構築する手順を説明します。この環境ではClaude Codeが事前にインストールされており、AI支援開発が可能です。

## 前提条件

以下のソフトウェアがmacOSにインストールされている必要があります：

- **Docker Desktop for Mac** (https://www.docker.com/products/docker-desktop/)
- **Visual Studio Code** (https://code.visualstudio.com/)
- **Dev Containers拡張機能** (VS Code内でインストール)
- **Claude Code** (既にコンテナ内にインストール済み)

## 1. 必要なソフトウェアのインストール

### Docker Desktop for Mac
1. https://www.docker.com/products/docker-desktop/ からDocker Desktop for Macをダウンロード
2. dmgファイルを開いてApplicationsフォルダにドラッグ&ドロップ
3. Dockerを起動し、初期設定を完了

### Visual Studio Code
1. https://code.visualstudio.com/ からVS Codeをダウンロード
2. zipファイルを展開してApplicationsフォルダに移動

### Dev Containers拡張機能
1. VS Codeを起動
2. 拡張機能パネル（⌘+Shift+X）を開く
3. \"Dev Containers\"で検索
4. \"Dev Containers\" (Microsoft製) をインストール

## 2. プロジェクトのセットアップ

### プロジェクトディレクトリの確認
現在のプロジェクトディレクトリ構造：
```
backlog-markdown-converter/
├── LICENSE
├── README.md
├── CLAUDE.md
├── SETUP.md (このファイル)
└── doc/
    └── md2backlog/
        ├── plan.md
        ├── requirements.ja.md
        └── requirements.md
```

## 3. Dev Container設定ファイルについて

このプロジェクトには以下のDev Container設定ファイルが既に含まれています：

### `.devcontainer/devcontainer.json`
Claude Code対応のGo開発環境設定ファイルです。以下が含まれます：
- Go言語サポート拡張機能
- 自動フォーマット設定（goimports）
- 静的解析設定（golangci-lint）
- Claude Code設定とボリュームマウント

### `.devcontainer/Dockerfile`
以下が事前にインストールされたカスタムDockerイメージ：
- Go 1.24.5（最新安定版）
- Claude Code CLI
- 必要な開発ツール（goimports, golangci-lint等）
- ネットワークセキュリティ設定
- 日本標準時（JST）設定

### `.devcontainer/init-firewall.sh`
セキュリティ強化されたネットワーク設定で、Go開発に必要なドメインのみアクセス許可

## 4. 開発環境の起動手順

### Step 1: VS Codeでプロジェクトを開く
```bash
code /path/to/backlog-markdown-converter
```

### Step 2: Dev Containerで再開
1. VS Codeでプロジェクトフォルダを開く
2. コマンドパレット（⌘+Shift+P）を開く
3. \"Dev Containers: Reopen in Container\" を選択
4. 初回は数分かかる場合があります（Dockerイメージのダウンロード）

### Step 3: 環境の確認
Dev Container内でターミナルを開き、以下のコマンドで環境を確認：

```bash
# Go言語のバージョン確認
go version

# Claude Codeの確認
claude --version

# 作業ディレクトリの確認
pwd

# 必要なツールの確認
goimports -h
golangci-lint --version

# Claude Codeでの開発開始
claude
```

## 5. 開発開始手順

### Go モジュールの初期化
```bash
go mod init github.com/yourusername/backlog-markdown-converter
```

### 依存関係の追加
```bash
go get github.com/spf13/cobra
go get github.com/yuin/goldmark
```

### ディレクトリ構造の作成
```bash
mkdir -p cmd/md2backlog
mkdir -p internal/converter
mkdir -p test
```

## 6. 開発環境の特徴

Dev Container内では以下の機能が利用できます：

### Go言語開発支援
- **Go言語サポート**: シンタックスハイライト、インテリセンス
- **デバッグ機能**: ブレークポイント設定、ステップ実行
- **自動フォーマット**: ファイル保存時に`goimports`で自動整形
- **Linting**: リアルタイムで`golangci-lint`によるコードの問題検出
- **テスト実行**: エディタ内からテストの実行・デバッグ

### Claude Code AI開発支援
- **インタラクティブAI**: `claude`コマンドでClaude Codeとの対話セッション
- **ファイル操作**: AIによるコード生成、編集、リファクタリング
- **TDD支援**: テスト駆動開発のガイダンス
- **Git連携**: Conventional Commitsでの適切なコミットメッセージ生成

## 7. 作業終了時

Dev Containerでの作業を終了する場合：

1. VS Codeを閉じる
2. 必要に応じてDocker Desktopでコンテナを停止

## 8. トラブルシューティング

### Dockerの権限エラー
```bash
sudo chmod 666 /var/run/docker.sock
```

### Dev Container起動失敗
- Docker Desktopが起動していることを確認
- VS CodeのDev Containers拡張機能が最新版であることを確認
- `.devcontainer/devcontainer.json`の書式エラーをチェック

### Go言語ツールのインストールエラー
Dev Container内で以下を実行：
```bash
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## 次のステップ

環境構築が完了したら、以下の手順で開発を開始してください：

### 1. Claude Codeセッション開始
```bash
# Dev Container内でClaude Codeを起動
claude
```

### 2. AI支援での開発進行
Claude Codeと協力して以下の順序で開発を進めます：
1. プロジェクトセットアップ（go mod, 依存関係）
2. 基本CLI構造の作成
3. TDDによる機能実装
4. テストとLintingの実行
5. Conventional Commitsでのコミット

### 3. 開発ガイドライン
- `CLAUDE.md`に記載されたプロジェクト指針に従う
- `doc/md2backlog/plan.md`の実装計画を参照
- 厳密なRed-Green-RefactorサイクルでのTDD

この環境ではGo言語の知識がなくても、Claude Codeが適切にガイドしながら開発を進めることができます。Claude Codeはプロジェクトの構造を理解し、CLAUDE.mdの指針に従った開発を支援します。