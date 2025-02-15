# Starlark API Reference

## エンティティ管理

### create_entity()
新しいエンティティを作成します。
- 引数: なし
- 戻り値: エンティティID (整数)
- 例:
```python
entity_id = create_entity()
```

### add_tag(entity_id, tag)
エンティティにタグを追加します。
- 引数:
  - entity_id: エンティティID (整数)
  - tag: タグ名 (文字列)
- 戻り値: なし
- 例:
```python
add_tag(entity_id, "bullet")
```

### find_entities_by_tag(tag)
指定したタグを持つすべてのエンティティを検索します。
- 引数:
  - tag: タグ名 (文字列)
- 戻り値: エンティティIDのリスト
- 例:
```python
bullets = find_entities_by_tag("bullet")
```

## コンポーネント管理

### add_component(entity_id, component_type, properties)
エンティティにコンポーネントを追加します。
- 引数:
  - entity_id: エンティティID (整数)
  - component_type: コンポーネントの種類 (文字列)
  - properties: コンポーネントのプロパティ (辞書)
- 戻り値: なし
- サポートされているコンポーネントタイプ:
  - "transform": 位置情報
  - "sprite": 描画情報
  - "text": テキスト表示
  - "physics": 物理演算
- 例:
```python
# Transform コンポーネント
add_component(entity_id, "transform", {
    "x": 100,
    "y": 200
})

# Sprite コンポーネント
add_component(entity_id, "sprite", {
    "width": 32,
    "height": 32,
    "color": "white"  # "white" または "cyan" がサポート
})

# Text コンポーネント
add_component(entity_id, "text", {
    "text": "Hello, World!",
    "x": 100,
    "y": 100
})

# Physics コンポーネント
add_component(entity_id, "physics", {
    "velocity_x": 0,
    "velocity_y": -300.0,
    "gravity": 0.0
})
```

### get_component(entity_id, component_type)
エンティティのコンポーネントを取得します。
- 引数:
  - entity_id: エンティティID (整数)
  - component_type: コンポーネントの種類 (文字列)
- 戻り値: コンポーネントのプロパティ (辞書) または None
- 例:
```python
transform = get_component(entity_id, "transform")
if transform:
    x = transform["x"]
    y = transform["y"]
```

### set_component(entity_id, component_type, properties)
エンティティのコンポーネントのプロパティを更新します。
- 引数:
  - entity_id: エンティティID (整数)
  - component_type: コンポーネントの種類 (文字列)
  - properties: 更新するプロパティ (辞書)
- 戻り値: なし
- 例:
```python
set_component(entity_id, "transform", {
    "x": new_x,
    "y": new_y
})
```

## 入力管理

### is_key_pressed(key)
指定したキーが押されているかを確認します。
- 引数:
  - key: キー名 (文字列)
- 戻り値: 真偽値
- サポートされているキー:
  - "Space": スペースキー
  - "ArrowLeft": 左矢印キー
  - "ArrowRight": 右矢印キー
  - "ArrowUp": 上矢印キー
  - "ArrowDown": 下矢印キー
- 例:
```python
if is_key_pressed("Space"):
    # スペースキーが押された時の処理
```

## システム情報

### get_total_entities()
アクティブなエンティティの総数を取得します。
- 引数: なし
- 戻り値: エンティティの数 (整数)
- 例:
```python
total = get_total_entities()
```

## デバッグ

### print(*args)
デバッグ情報を出力します。
- 引数: 任意の数の引数
- 戻り値: なし
- 例:
```python
print("Debug:", value)
```

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

## 状態管理

### set_state(entity_id, key, value)
エンティティの状態値を設定します。
- 引数:
  - entity_id: エンティティID（整数）
  - key: 状態のキー（文字列）
  - value: 設定する値（数値、文字列、真偽値）
- 戻り値: なし
- 例:
```python
set_state(player_id, "health", 100)
set_state(player_id, "name", "Player 1")
```

### get_state(entity_id, key)
エンティティの状態値を取得します。
- 引数:
  - entity_id: エンティティID（整数）
  - key: 状態のキー（文字列）
- 戻り値: 保存された値またはNone
- 例:
```python
health = get_state(player_id, "health")
if health > 0:
    # プレイヤーが生存している場合の処理
``` 