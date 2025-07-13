# 権限管理システム

組織ベースの権限管理を実装したGo Webアプリケーションです。ユーザーの所属組織に基づいてデータアクセスを制御します。

## 主な機能

- **組織ベースのアクセス制御**: ユーザーは所属する組織のデータのみ閲覧・編集可能
- **ロールベースの画面アクセス**: ロールに応じて利用可能な画面を制御
- **業務アプリケーション**: 予算管理、販売管理、在庫管理
- **管理機能**: ユーザー管理、権限設定、メニュー設定

## アーキテクチャ

### 認証・認可の仕組み

1. **ユーザー認証**: セッションベースの認証（gorilla/sessions）
2. **画面アクセス制御**: `role_screen_permission`テーブルによる制御
3. **データアクセス制御**: ユーザーの所属組織に基づくフィルタリング

### 主要コンポーネント

- **フロントエンド**: 静的HTML + JavaScript
- **バックエンド**: Go (net/http)
- **データベース**: PostgreSQL
- **コンテナ**: Docker & Docker Compose

## 前提条件

- Docker と Docker Compose
- Go 1.20以上（ローカル開発の場合）

## クイックスタート

1. リポジトリをクローン:
```bash
git clone <repository-url>
cd authority-matrix
```

2. Docker Composeで起動:
```bash
docker compose up -d --build
```

3. ブラウザでアクセス:
```
http://localhost:8080
```

4. テストユーザーでログイン:
- user001 (管理職 - 本社、北日本支店)
- user002 (営業担当 - 東京支店)
- user003 (在庫管理者 - 本社)
- admin (システム管理者 - 本社)

## データベーススキーマ

詳細なER図とテーブル説明は [docs/ER-diagram.md](docs/ER-diagram.md) を参照してください。

### テーブル概要

#### 認証・権限管理
- **user**: システムユーザー
- **role**: ロール（管理職、営業担当、在庫管理者、システム管理者）
- **organization**: 組織マスタ（本社、支店など）
- **user_role**: ユーザーとロールの関連
- **user_organization**: ユーザーと組織の関連
- **role_screen_permission**: ロールごとの画面アクセス権限

#### 画面・API管理
- **screen**: 画面マスタ
- **api**: APIマスタ
- **menu_setting**: メニュー設定
- **api_screen_mapping**: APIと画面のマッピング

#### 業務データ
- **budget**: 予算データ（組織別）
- **sales**: 販売データ（組織別）
- **inventory**: 在庫データ（組織別）

## API エンドポイント

### 認証
- `GET /login` - ログイン画面
- `POST /login` - ログイン処理
- `POST /logout` - ログアウト処理
- `GET /api/session` - セッション確認
- `GET /api/current-user` - 現在のユーザー情報取得

### 業務API
- `GET /back-budget1` - 予算データ取得
- `GET /back-sales` - 販売データ取得
- `GET /back-inventory` - 在庫データ取得

### 管理API
- `GET /api/admin/users` - ユーザー一覧
- `POST /api/admin/users` - ユーザー作成
- `PUT /api/admin/users/{id}` - ユーザー更新
- `DELETE /api/admin/users/{id}` - ユーザー削除

## 開発

### ローカル環境でのデバッグ

1. PostgreSQLのみ起動:
```bash
docker compose up -d postgres
```

2. 環境変数設定:
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=authority_user
export DB_PASSWORD=authority_pass
export DB_NAME=authority_db
```

3. アプリケーション起動:
```bash
go run main.go
```

### データベース接続

```bash
docker compose exec postgres psql -U authority_user -d authority_db
```

### ログ確認

```bash
docker compose logs -f
```

## トラブルシューティング

### テーブルが存在しないエラー

```bash
# コンテナとボリュームを完全に削除して再構築
docker compose down -v
docker compose up -d --build
```

### ログインできない場合

1. user_roleテーブルにデータが存在するか確認
2. セッションCookieをクリア
3. ブラウザのデベロッパーツールでエラーを確認

## ライセンス

MIT License