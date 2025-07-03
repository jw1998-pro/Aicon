# 所持品管理 API

高級品やコレクションアイテムを管理する REST API サーバーです。

## 📋 API 仕様

### ベース URL

```
http://localhost:8080
```

### エンドポイント一覧

| メソッド | パス             | 説明             | ステータスコード |
| -------- | ---------------- | ---------------- | ---------------- |
| GET      | `/health`        | ヘルスチェック   | 200              |
| GET      | `/items`         | 全アイテム取得   | 200              |
| POST     | `/items`         | アイテム登録     | 201, 400         |
| GET      | `/items/{id}`    | 特定アイテム取得 | 200, 404         |
| PATCH    | `/items/{id}`    | アイテム更新     | 200, 400, 404    |
| DELETE   | `/items/{id}`    | アイテム削除     | 204, 404         |
| GET      | `/items/summary` | カテゴリー別集計 | 200              |

### データ形式

#### アイテム (Item)

```json
{
  "id": 1,
  "name": "ロレックス デイトナ",
  "category": "時計",
  "brand": "ROLEX",
  "purchase_price": 1500000,
  "purchase_date": "2023-01-15",
  "created_at": "2023-01-15T10:00:00Z",
  "updated_at": "2023-01-15T10:00:00Z"
}
```

#### 有効なカテゴリー

- `時計`
- `バッグ`
- `ジュエリー`
- `靴`
- `その他`

### バリデーションルール

| フィールド     | 必須 | 制限                 |
| -------------- | ---- | -------------------- |
| name           | ✓    | 100 文字以内         |
| category       | ✓    | 有効なカテゴリーのみ |
| brand          | ✓    | 100 文字以内         |
| purchase_price | ✓    | 0 以上の整数         |
| purchase_date  | ✓    | YYYY-MM-DD 形式      |

### API 使用例

#### 1. 全アイテム取得

```bash
curl -X GET http://localhost:8080/items
```

**レスポンス:**

```json
[
  {
    "id": 1,
    "name": "ロレックス デイトナ",
    "category": "時計",
    "brand": "ROLEX",
    "purchase_price": 1500000,
    "purchase_date": "2023-01-15",
    "created_at": "2023-01-15T10:00:00Z",
    "updated_at": "2023-01-15T10:00:00Z"
  }
]
```

#### 2. アイテム登録

```bash
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{
    "name": "エルメス バーキン",
    "category": "バッグ",
    "brand": "HERMÈS",
    "purchase_price": 2000000,
    "purchase_date": "2023-02-20"
  }'
```

#### 3. 特定アイテム取得

```bash
curl -X GET http://localhost:8080/items/1
```

#### 4. アイテム更新

```bash
# 名前のみ更新
curl -X PATCH http://localhost:8080/items/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "更新された名前"
  }'

# 複数フィールド更新
curl -X PATCH http://localhost:8080/items/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "更新された名前",
    "brand": "更新されたブランド",
    "purchase_price": 2000000
  }'
```

**レスポンス:**

```json
{
  "id": 1,
  "name": "更新された名前",
  "category": "時計",
  "brand": "更新されたブランド",
  "purchase_price": 2000000,
  "purchase_date": "2023-01-15",
  "created_at": "2023-01-15T10:00:00Z",
  "updated_at": "2023-12-01T15:30:00Z"
}
```

#### 5. アイテム削除

```bash
curl -X DELETE http://localhost:8080/items/1
```

#### 6. カテゴリー別集計

```bash
curl -X GET http://localhost:8080/items/summary
```

**レスポンス:**

```json
{
  "categories": {
    "時計": 2,
    "バッグ": 1,
    "ジュエリー": 3,
    "靴": 0,
    "その他": 1
  },
  "total": 7
}
```

### エラーレスポンス形式

```json
{
  "error": "validation failed",
  "details": ["name is required", "purchase_price must be 0 or greater"]
}
```

## 🛠️ 技術スタック

- **言語**: Go 1.23
- **フレームワーク**: Echo v4
- **データベース**: MySQL 8.0
- **コンテナ**: Docker & Docker Compose

## 📁 プロジェクト構成

```
.
├── cmd/
│   └── main.go                 # エントリーポイント
├── internal/
│   ├── domain/
│   │   ├── entity/            # ドメインエンティティ
│   │   └── errors/            # ドメインエラー
│   ├── infrastructure/
│   │   ├── config/            # 設定管理
│   │   ├── database/          # データベース接続
│   │   └── server/            # HTTPサーバー
│   ├── interfaces/
│   │   ├── controller/        # HTTPハンドラー
│   │   └── database/          # リポジトリ
│   └── usecase/              # ビジネスロジック
├── sql/
│   └── init.sql              # データベース初期化
├── docker-compose.yml
├── Dockerfile
├── .env.example
└── README.md
```

## 🔧 開発環境

### 前提条件

- Docker
- Docker Compose

### ローカル開発

```bash
# 依存関係をインストール
go mod download

# ローカルでMySQLを起動（docker-compose経由）
docker-compose up -d mysql

# 環境変数を設定（ローカル開発用）
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=password
export DB_NAME=items_db

# アプリケーションを起動
go run cmd/main.go
```

### テストデータ

初期データとして以下のアイテムが登録されています：

1. ロレックス デイトナ (時計)
2. エルメス バーキン (バッグ)
3. ティファニー ネックレス (ジュエリー)
4. ルブタン パンプス (靴)
5. アップルウォッチ (その他)
