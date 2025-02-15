# 2Dゲームエンジン EbiTendon

シンプルで拡張可能な2Dゲーム開発エンジンです。Ebitenをベースに、ECSアーキテクチャとStarlarkスクリプティングを採用しています。

## 特徴

- ECSベースのゲームロジック
- Starlarkによるスクリプティング
- シーン管理システム
- UIコンポーネント
- オーディオ管理
- パーティクルシステム
- 衝突判定
- ステートマシン
- アセット管理

## 必要条件

- Go 1.20以上

## クイックスタート

1. プロジェクトのダウンロード:
```bash
git clone [リポジトリURL]
cd ebitendon
```

2. ゲームの実行:
```bash
go run src/main.go
```

## ゲーム開発の始め方

1. スクリプトの編集 (`scripts/main.star`):
```python
def init():
    # プレイヤーエンティティの作成
    player_id = create_entity()
    add_component(player_id, "transform", {
        "x": 100,
        "y": 200
    })
    add_component(player_id, "sprite", {
        "image": "player"
    })

init()
```

2. アセットの追加:
- 画像を `assets/images/` に配置
- 音声を `assets/audio/` に配置
- フォントを `assets/fonts/` に配置

3. アセットマニフェストの更新 (`assets/manifest.json`):
```json
{
  "images": {
    "player": {
      "path": "images/player.png"
    }
  }
}
```

## デフォルトで含まれる機能

### 基本システム
- ウィンドウ管理（800x600、タイトル "Game"）
- 基本的なECS実装
- アセット管理システム
- 簡単なUIシステム

### サンプルコンテンツ
- 基本的なメッセージウィンドウ
- シンプルなタイトル画面
- デバッグコンソール（`~`キーで表示）

## 注意事項

- フォントファイルは`src/engine/font/assets/`に配置済み
- アセットマニフェスト（`assets/manifest.json`）でリソースを管理
- 開発中のため、APIは変更される可能性があります

## ライセンス

MITライセンス

## サポート

バグ報告や機能要望はIssuesでお願いします。
プルリクエストも歓迎します。 