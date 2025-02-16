# Game Engine

シンプルなECSベースのゲームエンジン

## 特徴

- ECSアーキテクチャ
- Starlarkスクリプティング
- 解像度設定のサポート（HD, Full HD, SD, モバイル）
- コンポーネントベースの設計
- 物理演算システム
- テキストレンダリング（日本語対応）

## 予定
- シーン管理システム
- UIコンポーネント
- オーディオ管理
- パーティクルシステム
- 衝突判定
- ステートマシン
- アセット管理
  
## 使用方法

1. 実行:
```bash
go run src/main.go
```

2. デバッグモードで実行:
```bash
go run src/main.go -debug
```

## スクリプティング

ゲームロジックはStarlarkスクリプトで記述します。
詳細なAPIリファレンスは[STARLARK_API_REFERENCE.md](docs/STARLARK_API_REFERENCE.md)を参照してください。

### 基本的なスクリプト例:

```python
def init():
    # HD解像度で初期化
    set_screen_resolution("HD")
    
    # プレイヤーの作成
    player_id = create_entity()
    add_component(player_id, "transform", {"x": 100, "y": 100})
    add_component(player_id, "sprite", {"width": 32, "height": 32})

def update():
    # ゲームループの更新処理
    pass
```

## ライセンス

MITライセンス

## サポート

バグ報告や機能要望はIssuesでお願いします。
プルリクエストも歓迎します。

## 開発モード

### デバッグモードの使用

デバッグモードでは、開発用のスクリプトを選択して実行できます：

```bash
go run src/main.go -debug
```

デバッグモードでは以下の機能が利用可能です：
- スクリプトセレクター（複数のテストスクリプトから選択可能）
- デバッグコンソール（`~`キーで表示）
- パフォーマンスモニター

### サンプルスクリプト

`scripts_debug/`ディレクトリには以下のサンプルが用意されています：

1. `basic_sprite.star`
   - 基本的なスプライト表示
   - キーボード入力による移動
   - 簡単なアニメーション

2. `message_test.star`
   - メッセージウィンドウの表示
   - テキストの段階的表示
   - 選択肢の実装

各サンプルの実行方法：
1. デバッグモードで起動
2. スクリプトセレクターから実行したいサンプルを選択
3. 各サンプルのREADMEにしたがって操作 