---

name: Web屋姫 - プロジェクト全体のコーディング規約と設計方針
description: このファイルは、Web屋姫プロジェクトにおけるコーディング規約、設計方針、開発ワークフローを説明します。
applyTo: **

---

# Web 屋姫 - Copilot Instructions

## プロジェクト概要

Go 標準パッケージのみで構築する RESTful API プロジェクト(外部依存は最小限: MySQL driver, UUID, sql-migrate, テストライブラリのみ)。
Vtuber 宮乃やみさんの動画企画を実装したもの。

## アーキテクチャパターン

### レイヤー構造

- **cmd/**: エントリーポイント (`api/main.go`, `migration/main.go`)
- **internal/domain/**: ドメインモデルとリポジトリインターフェース
- **internal/handler/**: HTTP ハンドラー (request/response 変換を担当)
- **internal/infra/**: インフラ実装 (MySQL 接続など)
- **internal/server/**: サーバー設定、ルーティング、ミドルウェア
- **pkg/**: 汎用ユーティリティ (context, logger, httputil 等)

### 重要な設計原則

1. **Context ベースの依存注入**: DB 接続や設定は`context.Context`経由で渡す

   - DB 接続: `Ctx.GetDB(ctx)` ([pkg/context/context.go](pkg/context/context.go))
   - 設定: `Ctx.GetCtxCfg(ctx)` ([pkg/config/config.go](pkg/config/config.go))
   - リクエスト ID: `Ctx.GetRequestID(ctx)`

2. **インターフェース駆動開発**: 全レイヤーでインターフェース定義

   - Handler: `IUserHandler`, `ISummaryHandler` ([internal/handler/](internal/handler/))
   - Repository: `IUserRepository`, `ISummaryRepository` ([internal/domain/](internal/domain/))

3. **標準 HTTP ルーター**: Go 1.22+の`http.ServeMux`パターンマッチング使用

   ```go
   engine.HandleFunc("POST /users", userSaveHandler)
   engine.HandleFunc("GET /users/{id}", userDetailHandler)
   ```

4. **ミドルウェアチェーン**: `UseMiddleware(ctx, handler)` ([internal/server/middleware.go](internal/server/middleware.go#L168))
   - RequestID 付与、タイムアウト、ログ、DB 接続を自動付与

## コーディング規約

### リクエスト/レスポンス処理

- リクエストバインド: `request.Bind(r, &req)` - リフレクションベースのパラメータバインディング
- バリデーション: `request.Validate(&req)` - タグベース検証 (`validate:"required"`)
- レスポンス: `httputil.Response(&w, status, data)` - 常に JSON 形式

例: [internal/handler/user/user.go](internal/handler/user/user.go#L27-L56)

### データベース操作

- 接続取得: 必ず `Ctx.GetDB(ctx)` で nil チェック
- クエリ実行: `db.ExecContext(ctx, query, args...)` / `db.QueryContext(ctx, query, args...)`
- エラーハンドリング: `sql.ErrNoRows` を明示的にチェック

例: [internal/infra/database/mysql/user.go](internal/infra/database/mysql/user.go#L23-L34)

### エラー処理パターン

```go
if err != nil {
    return fmt.Errorf("descriptive message: %w", err)  // エラーラップ
}
```

## 開発ワークフロー

### ローカル開発コマンド

```bash
make up              # Docker Compose起動 (DB + App)
make migrate-up      # マイグレーション実行
make seed            # テストデータ投入
make test            # テスト実行 (-race -parallel 1)
make lint            # go vet + golangci-lint
make migrate-new name=<name>  # 新規マイグレーションファイル生成
```

### データベース

- DSN 形式: `user:P@ssw0rd@tcp(host:3306)/develop_web_ya_hime?parseTime=true`
- 環境変数: `DATABASE_URL` (デフォルト: localhost:3306)
- マイグレーション: `db/migrations/*.sql` (sql-migrate 使用)
- シード: `db/seed/*.sql` (手動実行順序: 00_trancate.sql → 01_seed.sql)

### テスト

- テーブルテスト推奨: `tests := []struct { name string; input X; want Y }{...}`
- モック: `sqlmock` 使用 ([internal/infra/database/mysql/\*\_test.go](internal/infra/database/mysql/))
- ゴールデンファイル: `pkg/testutil/golden.go` でレスポンステスト

## 新機能追加手順

1. **ドメイン層**: `internal/domain/<entity>/` にモデルとリポジトリインターフェース定義
2. **インフラ層**: `internal/infra/database/mysql/<entity>.go` でリポジトリ実装
3. **ハンドラー層**:
   - `internal/handler/request/<entity>.go` - リクエスト構造体
   - `internal/handler/response/<entity>.go` - レスポンス構造体
   - `internal/handler/<entity>/<entity>.go` - ビジネスロジック
4. **ルーティング**: `internal/server/server.go` の`Run()`メソッドに追加
5. **OpenAPI**: `openapi.yaml` にエンドポイント定義追加

## プロジェクト固有の注意点

- **グレースフルシャットダウン**: SIGINT シグナル処理済み ([internal/server/server.go](internal/server/server.go#L88-L104))
- **リクエストタイムアウト**: デフォルト 5 秒 ([internal/server/middleware.go](internal/server/middleware.go#L34-L56))
- **UUIDv4**: 独自実装 `pkg/uuid/uuid.go` (crypto/rand ベース)
- **カスタム設定ローダー**: リフレクションで環境変数をロード ([pkg/config/config.go](pkg/config/config.go#L25-L43))
- **ベースモデル**: 全エンティティは `domain.WYHBaseModel` を埋め込み (ID, CreatedAt, UpdatedAt)

## よくある実装パターン

### 新しいエンドポイント追加

```go
// 1. ハンドラー作成 (internal/handler/<entity>/<entity>.go)
func (h *handler) Action(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    var req request.ActionRequest
    if err := request.Bind(r, &req); err != nil { /*...*/ }
    if err := request.Validate(&req); err != nil { /*...*/ }
    // ビジネスロジック
    httputil.Response(&w, http.StatusOK, response)
}

// 2. ルーティング登録 (internal/server/server.go)
handler := UseMiddleware(ctx, h.Action)
engine.HandleFunc("POST /entities", handler)
```

### カスタムバリデーション追加

`internal/handler/request/request.go` の `validateField()` に新ルール追加

### 新しいミドルウェア追加

`internal/server/middleware.go` に関数定義 → `UseMiddleware()` 内で適用

## Git & GitHub ワークフロー

### ブランチ戦略

- **main**: 本番環境用のメインブランチ
- **story-#<番号>**: 機能実装用のフィーチャーブランチ（例: `story-#2`）

### 開発フロー（機能実装 → PR 作成 → レビュー対応）

#### 1. 機能実装とコミット

```bash
# 実装完了後、変更をステージング
git add -A

# わかりやすいコミットメッセージでコミット
git commit -m "feat: 機能の説明"
# 例: git commit -m "feat: カテゴリ機能・ページネーション・論理削除の実装"

# リモートブランチにプッシュ
git push origin <ブランチ名>
```

#### 2. MCP ツールで PR 作成

GitHub の MCP ツールを使用してプルリクエストを作成:

```
# Copilotに依頼する例:
"o-ga09/web-ya-himeにPRテンプレートを使用してPRを作成して"
```

MCP ツール内部では以下が実行される:

- `mcp_github_create_pull_request` を使用
- `base`: "main"（マージ先ブランチ）
- `head`: 現在のブランチ（例: "story-#2"）
- `title`: わかりやすい PR タイトル
- `body`: PR の概要、変更内容、技術詳細、使用例などを含む詳細な説明

#### 3. レビューコメント対応

```bash
# レビューコメントを確認
"MCPを使用してPR #<番号>のレビューコメントを取得して"

# コメントに対応後、再度コミット
git add -A
git commit -m "fix: レビューコメント対応 - 修正内容の説明"
# 例: git commit -m "fix: レビューコメント対応 - カテゴリをNULL可能に変更、ポインタ型使用"

# プッシュ（PRが自動更新される）
git push origin <ブランチ名>
```

#### 4. レビューコメント取得の方法

MCP ツールでレビューコメントを確認:

```
# 方法1: レビューコメント取得
mcp_github_pull_request_read(method="get_review_comments", owner="o-ga09", repo="web-ya-hime", pullNumber=<PR番号>)

# 方法2: 通常のコメント取得
mcp_github_pull_request_read(method="get_comments", owner="o-ga09", repo="web-ya-hime", pullNumber=<PR番号>)

# 方法3: PR詳細取得
mcp_github_pull_request_read(method="get", owner="o-ga09", repo="web-ya-hime", pullNumber=<PR番号>)
```

### コミットメッセージ規約

プレフィックスを使用した明確なメッセージ:

- `feat:` - 新機能追加
- `fix:` - バグ修正
- `refactor:` - リファクタリング
- `test:` - テスト追加・修正
- `docs:` - ドキュメント更新
- `style:` - コードフォーマット修正（機能変更なし）
- `chore:` - ビルド処理、補助ツールの変更

例:

```bash
git commit -m "feat: カテゴリ機能とページネーション実装"
git commit -m "fix: レビューコメント対応 - NULL対応とポインタ型使用"
git commit -m "style: gofmtでコードフォーマットを修正"
```

### よくある Git 操作

```bash
# 現在のブランチ確認
git branch

# 変更されたファイル確認
git status

# 差分確認
git diff

# 特定ファイルの差分確認
git diff <ファイルパス>

# コミット履歴確認
git log --oneline

# 最新のコミットを修正（まだpushしていない場合）
git commit --amend
```

### フォーマットと Lint

コミット前に必ず実行:

```bash
# コードフォーマット
gofmt -w .

# Lint実行
make lint

# テスト実行
make test
```

これらのチェックは、PR の CI でも自動実行されます。
