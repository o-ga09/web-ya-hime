# E2Eテスト

Scenarigoを使用したE2Eテストシナリオです。

## セットアップ

### Scenarigoのインストール

#### Homebrewの場合
```bash
brew install scenarigo/tap/scenarigo
```

#### Goの場合
```bash
go install github.com/zoncoen/scenarigo/cmd/scenarigo@latest
```

#### バイナリダウンロード
https://github.com/zoncoen/scenarigo/releases から最新版をダウンロード

## テストの実行

### アプリケーションの起動
```bash
# Docker Composeでアプリケーションとデータベースを起動
make up

# マイグレーション実行
make migrate-up

# シードデータ投入（オプション）
make seed
```

### E2Eテストの実行

#### すべてのシナリオを実行
```bash
make e2e-test
```

または直接実行:
```bash
cd e2e
scenarigo run scenarios
```

#### 特定のシナリオのみ実行
```bash
cd e2e
scenarigo run scenarios/health_check.yaml
scenarigo run scenarios/category.yaml
scenarigo run scenarios/subcategory.yaml
scenarigo run scenarios/summary_with_category.yaml
scenarigo run scenarios/validation.yaml
```

#### ベースURLを変更して実行
```bash
BASE_URL=http://localhost:9090 scenarigo run scenarios
```

## テストシナリオ

### health_check.yaml
- アプリケーションヘルスチェック (`/health`)
- データベースヘルスチェック (`/db-health`)

### category.yaml
カテゴリ管理APIの基本的なCRUD操作:
- カテゴリ作成 (`POST /categories`)
- カテゴリ一覧取得 (`GET /categories`)
- カテゴリ詳細取得 (`GET /categories/{id}`)
- カテゴリ削除 (`DELETE /categories/{id}`)

### subcategory.yaml
サブカテゴリ管理APIの基本的なCRUD操作:
- サブカテゴリ作成 (`POST /subcategories`)
- サブカテゴリ一覧取得 (`GET /subcategories`)
- サブカテゴリ一覧取得（カテゴリ指定） (`GET /subcategories?category_id=xxx`)
- サブカテゴリ詳細取得 (`GET /subcategories/{id}`)
- サブカテゴリ削除 (`DELETE /subcategories/{id}`)

### summary_with_category.yaml
カテゴリ・サブカテゴリを含むサマリー管理のテスト:
- カテゴリとサブカテゴリを指定したサマリー作成
- サマリー一覧取得（フィルタなし）
- カテゴリIDでフィルタリング
- サブカテゴリIDでフィルタリング
- サマリー詳細取得（カテゴリ・サブカテゴリ情報を含む）

### validation.yaml
カテゴリとサブカテゴリの組み合わせバリデーションテスト:
- サブカテゴリのみ指定（親カテゴリ自動補完）
- 正しいカテゴリ+サブカテゴリの組み合わせ
- 誤ったカテゴリ+サブカテゴリの組み合わせ（400エラー）
- サマリー検索時のバリデーション

## テストデータのクリーンアップ

各シナリオは実行後に作成したデータを自動的にクリーンアップします。

テストデータベースを完全にリセットしたい場合:
```bash
make seed  # TRUNCATEして初期データを再投入
```

## トラブルシューティング

### アプリケーションが起動しない
```bash
# ログを確認
docker-compose logs -f api

# コンテナの状態を確認
docker-compose ps
```

### テストが失敗する
1. アプリケーションが起動しているか確認
2. マイグレーションが実行されているか確認
3. ベースURLが正しいか確認
4. シナリオファイルの構文エラーがないか確認

### 詳細なログを出力
```bash
scenarigo run --verbose scenarios
```

## CI/CDでの実行

GitHub Actionsなどで実行する場合の例:
```yaml
- name: Run E2E Tests
  run: |
    make up
    make migrate-up
    sleep 5  # アプリケーションの起動を待つ
    make e2e-test
  env:
    BASE_URL: http://localhost:8080
```
