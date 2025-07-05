# 組織階層管理システム

階層型組織構造を時系列で管理するWebアプリケーションです。組織の再編や移動を日付指定で記録し、任意の時点での組織構造を確認できます。

## 基本的な考え方

部門マスタの発行日・失効日を管理しつつ、自己参照によって階層構造を表現できるようにする。
失効日は期間別組織属性テーブルのレコードごとに、次のレコードの発行日-1で計算される。

## 機能概要

- **部門管理**: 部門マスタのCRUD操作
- **期間別組織属性**: 組織階層の時系列管理（失効年月の自動導出）
- **組織階層ビュー**: 特定日付時点の組織構造をツリー形式で表示
- **組織再編**: 部門の所属変更を簡単に実行

## 技術スタック

- **バックエンド**: Go (Echo Framework)
- **データベース**: PostgreSQL
- **フロントエンド**: Vanilla JavaScript, HTML, CSS
- **インフラ**: Docker, Docker Compose

## データモデル

### 1. 部門テーブル (departments)
| カラム名 | 型 | 説明 |
|---------|---|------|
| department_id | VARCHAR(50) | 部門ID (PK) |
| department_name | VARCHAR(255) | 部門名 |
| created_at | TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | 更新日時 |

### 2. 期間別組織属性テーブル (organization_attributes)
| カラム名 | 型 | 説明 |
|---------|---|------|
| department_id | VARCHAR(50) | 部門ID (PK) |
| effective_date | DATE | 発効年月日 (PK) |
| parent_department_id | VARCHAR(50) | 上位部門ID (NULL可) |
| created_at | TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | 更新日時 |

※ 失効年月は次の発効日の前日として自動的に導出されます。

## セットアップ

### 前提条件
- Docker
- Docker Compose

### 起動方法

```bash
# リポジトリをクローン
git clone [repository-url]
cd pyramid-structure

# Docker Composeで起動（デフォルトポート: 8080）
docker-compose up -d

# カスタムポートで起動する場合
APP_PORT=3000 docker-compose up -d

# または.envファイルを使用
cp .env.example .env
# .envファイルでAPP_PORTを編集
docker-compose up -d
```

アプリケーションは http://localhost:8080 （またはカスタムポート）でアクセスできます。

## サンプルデータ

初期状態で以下のサンプル組織データが登録されています：

### 部門一覧
- **HQ**: 本社
- **STRATEGY**: 経営企画部
- **GENERAL**: 総務部
- **HR**: 人事部
- **FINANCE**: 財務部
- **SALES_HQ**: 営業本部
  - **SALES_1**: 営業1部
  - **SALES_2**: 営業2部
  - **SALES_SUPPORT**: 営業支援部
- **TECH_HQ**: 技術本部
  - **DEV**: 開発部
  - **RESEARCH**: 研究部
  - **QA**: 品質管理部
- **MFG_HQ**: 製造本部
  - **MFG_1**: 製造1部
  - **MFG_2**: 製造2部
- **IT**: IT推進部
- **DIGITAL**: デジタル戦略部

### 組織改編履歴
1. **2023年1月1日**: 初期組織構造
2. **2023年7月1日**: IT推進部を総務部配下に新設
3. **2024年1月1日**: IT推進部を技術本部に移管
4. **2024年4月1日**: 
   - デジタル戦略部を経営企画部配下に新設
   - 営業支援部を本社直轄に変更
5. **2024年10月1日**: 製造本部を技術本部に統合

## 使い方

### 1. 部門管理タブ

部門マスタの管理を行います。

- **新規部門追加**: 「新規部門追加」ボタンから部門ID・部門名を入力
- **編集**: 各行の「編集」ボタンから部門名を変更
- **削除**: 各行の「削除」ボタンから部門を削除（関連する組織属性も削除されます）

### 2. 期間別組織属性タブ

組織階層の履歴を管理します。

- **新規属性追加**: 組織階層の変更を登録
  - 部門IDを選択
  - 発効年月日を指定
  - 上位部門を選択（空の場合は最上位）
- **フィルタ**: 部門別に履歴を絞り込み
- **失効年月**: 次の発効日の前日が自動表示されます

### 3. 組織階層タブ

特定日付時点の組織構造を階層表示します。

- **基準日**: 表示したい時点の日付を選択
- **展開/折りたたみ**: 階層の表示を制御
- **JSONエクスポート**: 現在の組織構造をJSON形式でダウンロード

エクスポートされるJSONの形式：
```json
{
  "exportDate": "2024-12-01T10:00:00.000Z",
  "hierarchyDate": "2024-10-01",
  "departments": [...],
  "hierarchy": [
    {
      "id": "HQ",
      "name": "本社",
      "effectiveDate": "2023-01-01",
      "expirationDate": null,
      "children": [...]
    }
  ]
}
```

### 4. 組織再編タブ

組織の所属変更を簡単に実行できます。

**左側：組織移動フォーム**
1. 移動する部門を選択
2. 新しい上位部門を選択（空の場合は最上位へ移動）
3. 発効日を指定
4. プレビューで内容を確認
5. 「組織移動を実行」で登録

**右側：現在の組織構造**
- 基準日時点の組織階層を表示
- 移動前の状態を確認可能

**循環参照チェック**
- 部門を自身の配下に移動させることはできません
- 子部門を親部門の上位に移動させることはできません

## API エンドポイント

### 部門管理
- `GET /api/departments` - 全部門取得
- `GET /api/departments/:id` - 部門取得
- `POST /api/departments` - 部門作成
- `PUT /api/departments/:id` - 部門更新
- `DELETE /api/departments/:id` - 部門削除

### 組織属性管理
- `GET /api/organization-attributes` - 全組織属性取得
- `GET /api/organization-attributes/:department_id/:effective_date` - 組織属性取得
- `POST /api/organization-attributes` - 組織属性作成
- `PUT /api/organization-attributes/:department_id/:effective_date` - 組織属性更新
- `DELETE /api/organization-attributes/:department_id/:effective_date` - 組織属性削除
- `GET /api/departments/:id/history` - 部門履歴取得
- `GET /api/hierarchy?date=YYYY-MM-DD` - 特定日付の組織階層取得

## 開発

### ローカル開発環境

```bash
# データベースのみ起動
docker-compose up -d postgres

# Goアプリケーションの起動
go mod download
go run main.go
```

### テスト

```bash
go test ./...
```
