# Starlark API リファレンス

## 基本機能

### エンティティ操作

```python
# エンティティの作成
entity_id = create_entity()

# コンポーネントの追加
add_component(entity_id, "transform", {
    "x": 100,
    "y": 200,
    "rotation": 0
})

add_component(entity_id, "sprite", {
    "image": "player",
    "scale": 1.0
})

# コンポーネントの取得
transform = get_component(entity_id, "transform")
```

### グローバル関数

```python
# デバッグ出力
print("Hello from Starlark!")

# 画面サイズの取得
width, height = get_screen_size()

# 時間関連
delta_time = get_delta_time()  # 前フレームからの経過時間
```

## コンポーネント一覧

### TransformComponent
位置、回転、スケールを管理します。

```python
add_component(entity_id, "transform", {
    "x": float,        # X座標
    "y": float,        # Y座標
    "rotation": float, # 回転（ラジアン）
    "scale_x": float,  # X方向スケール（デフォルト: 1.0）
    "scale_y": float   # Y方向スケール（デフォルト: 1.0）
})
```

### SpriteComponent
画像の描画を管理します。

```python
add_component(entity_id, "sprite", {
    "image": string,     # 画像リソース名
    "z_index": int,      # 描画順序（デフォルト: 0）
    "visible": bool,     # 表示/非表示（デフォルト: true）
    "color": {           # 色調整（デフォルト: 白）
        "r": float,      # 赤 (0-1)
        "g": float,      # 緑 (0-1)
        "b": float,      # 青 (0-1)
        "a": float       # 不透明度 (0-1)
    }
})
```

### PhysicsComponent
物理演算を管理します。

```python
add_component(entity_id, "physics", {
    "velocity_x": float,    # X方向速度
    "velocity_y": float,    # Y方向速度
    "gravity": float,       # 重力（デフォルト: 0）
    "friction": float,      # 摩擦（デフォルト: 0）
    "solid": bool          # 衝突判定の有無（デフォルト: true）
})
```

## イベントシステム

```python
# イベントの発行
emit_event("player_damaged", {
    "damage": 10,
    "source": "enemy"
})

# イベントハンドラーの登録
def on_player_damaged(event):
    print(f"Player took {event.damage} damage from {event.source}")

register_event_handler("player_damaged", on_player_damaged)
```

## タイマー機能

```python
# 遅延実行
def delayed_function():
    print("3秒経過しました")

schedule_timer(3.0, delayed_function)

# 繰り返し実行
def repeated_function():
    print("1秒ごとに実行")

schedule_interval(1.0, repeated_function)
```

## リソース管理

```python
# 画像のロード
load_image("player", "assets/player.png")
load_image("enemy", "assets/enemy.png")

# 音声のロード
load_sound("jump", "assets/jump.wav", "se")
load_sound("bgm", "assets/bgm.ogg", "bgm")

# 音声の再生
play_sound("jump")
play_bgm("bgm", loop=True)
```

## 注意事項

1. すべての座標は画面左上を（0,0）とする座標系で指定します。

2. Z-indexは小さい値が手前、大きい値が奥になります。デフォルト値は以下の通りです：
   - メッセージウィンドウ: 100
   - UI要素: 50
   - 通常スプライト: 0
   - 背景: -100

3. コンポーネントの追加時、指定されていないパラメーターはデフォルト値が使用されます。

4. イベントハンドラーは非同期で実行される可能性があるため、状態の変更には注意が必要です。

5. タイマー関数は非同期で実行されます。メインスレッドの処理をブロックしないようにしてください。 