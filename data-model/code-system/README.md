# 品目コード体系管理システム

PostgreSQL on DockerとGo言語のEchoフレームワークを使用した品目管理システムです。

## システム概要

品目を品種別（A品種・B品種）に管理し、それぞれの品種特有の属性を持つコード体系をモデリングしています。

## 基本的な考え方

品目IDはただのユニークな識別子として定義する。ただし、品目コードを二次識別子とすることで、コードを見ただけで何の品目かわかる人のユーザビリティを損なわないようにする。このとき、ユーザーには品目コードだけが見えるようにする。  
品目ごとに異なる属性は、別途テーブルを切り出して管理する。  
これで、システムは品目コードの各桁を特別扱いするようなロジックを持たずに済む。  
もちろん、既存ユーザーへの配慮のための構成であるから、コードから仕様を類推することはできなくなる。  

### データベース構造

- **品目基本属性テーブル**: 品目の基本情報を管理
  - 品目ID (PK)
  - 品目名
  - 品種区分
  - 品目コード (SK)

- **A品種品目属性テーブル**: A品種特有の属性
  - 品目ID (PK/FK)
  - 容量
  - 材質

- **B品種品目属性テーブル**: B品種特有の属性
  - 品目ID (PK/FK)
  - 内径
  - 外径

## 起動方法

```bash
# Dockerコンテナの起動
docker-compose up -d

# ログの確認
docker-compose logs -f

# 停止
docker-compose down

# データを含めて削除
docker-compose down -v
```

## API エンドポイント

### 品目一覧取得
```
GET /api/items?page=1&page_size=10&category_type=A
```

### 品目詳細取得
```
GET /api/items/:id
```

### 品目作成
```
POST /api/items
Content-Type: application/json

{
  "item_id": "A004",
  "item_name": "プラスチックボトル750ml",
  "category_type": "A",
  "item_code": "PBOT-750",
  "capacity": 750.00,
  "material": "PET"
}
```

### 品目更新
```
PUT /api/items/:id
Content-Type: application/json

{
  "item_name": "新しい品目名",
  "capacity": 800.00
}
```

### 品目削除
```
DELETE /api/items/:id
```

## アクセス方法

- Webアプリケーション: http://localhost:8080
- API: http://localhost:8080/api

## 技術スタック

- Backend: Go 1.21 + Echo Framework
- Database: PostgreSQL 15
- Frontend: HTML/CSS/JavaScript (Vanilla)
- Container: Docker & Docker Compose

## 開発環境

```bash
# 依存関係のインストール
go mod download

# アプリケーションの実行（ローカル）
go run main.go
```

## テストデータ

初期状態で以下のテストデータが投入されています：

### A品種
- A001: プラスチックボトル500ml (PET, 500ml)
- A002: ガラスボトル1000ml (ガラス, 1000ml)
- A003: アルミ缶350ml (アルミニウム, 350ml)

### B品種
- B001: ステンレスパイプ20mm (内径15mm, 外径20mm)
- B002: 塩ビパイプ50mm (内径44mm, 外径50mm)
- B003: 銅管15mm (内径13mm, 外径15mm)