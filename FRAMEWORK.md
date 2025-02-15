# フレームワーク構成ドキュメント

## 使用ミドルウェア

### コアライブラリ
- **Ebiten v2** (github.com/hajimehoshi/ebiten/v2)
  - ゲームエンジン本体
  - 画面描画、入力処理、オーディオ処理の基盤
  - バージョン: v2.8.6

### オーディオ関連
- **oggvorbis** (github.com/jfreymuth/oggvorbis)
  - OGGファイルのデコード
- **go-mp3** (github.com/hajimehoshi/go-mp3)
  - MP3ファイルのデコード

### スクリプトエンジン
- **Starlark** (go.starlark.net)
  - Pythonライクなスクリプト言語
  - ゲームロジックの記述に使用

## アーキテクチャ概要

### レイヤー構造
1. **View層** (`view/`)
   - 画面描画処理
   - UIコンポーネントのレンダリング

2. **Model層** (`model/`)
   - ゲームデータの管理
   - セーブ/ロード機能

3. **Controller層** (`game/`)
   - ゲームのメインロジック
   - 各コンポーネントの統合

4. **UI層** (`ui/`)
   - ボタン、メニュー等のUI管理
   - レイアウト計算

### 主要コンポーネント

#### スクリプトエンジン (`script/`)
```go
scriptEngine, err := script.NewScriptEngine(
    addMessage,      // メッセージ表示用コールバック
    registerSE,      // SE登録用コールバック
    playSE,         // SE再生用コールバック
    updateStats,    // ステータス更新用コールバック
    animationManager // アニメーション管理
)
```

#### オーディオマネージャー (`audio/`)
```go
audio := audio.NewJukebox()
audio.SetBGMVolume(0.5)
audio.SetSEVolume(0.5)
```

#### アニメーションマネージャー (`animation/`)
```go
animationManager := animation.NewAnimationManager()
```

## 設定ファイル構造

### 解像度設定
```go
type Resolution struct {
    Name   string
    Width  int
    Height int
}
```

### 設定保存形式
```json
{
    "current_resolution": 2,
    "is_fullscreen": false,
    "bgm_volume": 0.5,
    "se_volume": 0.5
}
```

## アセット構造

```
assets/
├── audio/
│   ├── bgm/
│   └── se/
├── images/
└── scripts/
    ├── logic/
    └── vn/
```

## 注意点

1. **スクリプトエンジンの初期化**
   - 必ずゲーム状態の初期化後に行う
   - エラーハンドリングを適切に実装

2. **オーディオ処理**
   - MP3とOGGの両方のデコーダーが必要
   - メモリリークを防ぐため、適切なリソース解放が必要

3. **レイアウト管理**
   - 解像度変更時に全UIコンポーネントの再計算が必要
   - フォントサイズは画面サイズに応じて動的に計算

4. **状態管理**
   - ゲーム状態の変更は必ずStateManagerを通して行う
   - 直接的な状態変更を避ける

## 移行時の注意点

1. **依存関係**
   - `go.mod`ファイルの完全なコピー
   - バージョン互換性の確認

2. **アセット管理**
   - 埋め込みリソースの適切な移行
   - フォントファイルの取り扱い

3. **設定ファイル**
   - 既存の設定ファイルとの互換性維持
   - 新規設定項目追加時の後方互換性 