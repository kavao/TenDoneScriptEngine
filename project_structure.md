project/
├── src/
│   ├── main.go                 # エントリーポイント
│   ├── engine/
│   │   ├── script/            # Starlarkエンジン
│   │   │   ├── engine.go      # スクリプトエンジン本体
│   │   │   ├── api.go         # Starlark API定義
│   │   │   └── coroutine.go   # コルーチン実装
│   │   ├── render/            # 描画エンジン
│   │   │   ├── manager.go     # 描画マネージャー
│   │   │   ├── cache.go       # 描画キャッシュ
│   │   │   └── zbuffer.go     # Z-buffer実装
│   │   ├── audio/             # オーディオエンジン
│   │   │   ├── manager.go     # オーディオマネージャー
│   │   │   ├── bgm.go         # BGM制御
│   │   │   └── se.go          # SE制御
│   │   └── save/              # セーブ/ロード
│   │       ├── manager.go     # セーブマネージャー
│   │       └── serializer.go  # シリアライズ処理
│   ├── game/
│   │   ├── state/             # ゲーム状態管理
│   │   ├── scene/             # シーン管理
│   │   └── ui/                # UI管理
│   └── assets/
│       ├── scripts/           # Starlarkスクリプト
│       ├── images/            # 画像リソース
│       └── audio/             # 音声リソース
├── tests/                     # テストコード
└── docs/                      # ドキュメント 