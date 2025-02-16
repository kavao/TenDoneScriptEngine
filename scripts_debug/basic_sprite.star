# グローバル変数を辞書として管理
vars = {
    "player_id": None,
    "stats_text_id": None
}

print("Defining update function") # デバッグ出力

def update():
    # プレイヤーの移動処理
    transform = get_component(vars["player_id"], "transform")
    if not transform:
        print("Warning: Player transform not found")
        return

    speed = 8.0
    x = transform["x"]
    y = transform["y"]
    
    # カーソルキーでの移動
    moved = False
    if is_key_pressed("ArrowLeft"):
        x -= speed
        moved = True
        if x < 0:
            x = 0
            
    if is_key_pressed("ArrowRight"):
        x += speed
        moved = True
        if x > 1280-32:  # HD幅に合わせて調整
            x = 1280-32
            
    if is_key_pressed("ArrowUp"):
        y -= speed
        moved = True
        if y < 0:
            y = 0
            
    if is_key_pressed("ArrowDown"):
        y += speed
        moved = True
        if y > 720-32:  # HD高さに合わせて調整
            y = 720-32

    if moved:
        set_component(vars["player_id"], "transform", {
            "x": x,
            "y": y
        })
    
    # 弾の発射（位置を調整）
    if is_key_pressed("Space"):
        cooldown = get_state(vars["player_id"], "shot_cooldown") or 0
        if cooldown <= 0:
            bullets = get_bullets()
            if len(bullets) < 10:
                bullet_x = x + (32/2) - (8/2)
                create_bullet(bullet_x, y)
                set_state(vars["player_id"], "shot_cooldown", 6)
        
    # クールダウンを減少
    cooldown = get_state(vars["player_id"], "shot_cooldown") or 0
    if cooldown > 0:
        set_state(vars["player_id"], "shot_cooldown", cooldown - 1)

    # 統計情報の更新
    bullets = get_bullets()
    total_entities = get_total_entities()
    # print("Debug - Total entities:", total_entities, "Bullets:", len(bullets))

    # テキストコンポーネントの動的更新
    if vars["stats_text_id"] == None:
        vars["stats_text_id"] = create_entity()
        add_component(vars["stats_text_id"], "text", {
            "text": "",
            "x": 1000,  # HD幅に合わせて調整
            "y": 650    # HD高さに合わせて調整
        })

    # 毎フレーム統計情報を更新
    stats_text = "Objects: " + str(total_entities) + " | Bullets: " + str(len(bullets))
    set_component(vars["stats_text_id"], "text", {
        "text": stats_text,
        "x": 1000,  # HD幅に合わせて調整
        "y": 650    # HD高さに合わせて調整
    })

print("Update function defined") # デバッグ出力

def create_bullet(x, y):
    bullet_id = create_entity()
    # print("Creating bullet with ID:", bullet_id)
    add_tag(bullet_id, "bullet")
    # print("Added 'bullet' tag to entity:", bullet_id)
    
    # タグが正しく追加されたか確認
    bullets = find_entities_by_tag("bullet")
    # print("Current bullet count:", len(bullets))
    
    # 弾のTransform
    add_component(bullet_id, "transform", {
        "x": x,
        "y": y
    })
    
    # 弾のSprite
    add_component(bullet_id, "sprite", {
        "width": 8,
        "height": 8,
        "color": "cyan"
    })
    
    # 弾の物理挙動を修正（速度をさらに上げる）
    add_component(bullet_id, "physics", {
        "velocity_x": 0,
        "velocity_y": -600.0,  # 速度を-に増加（より速く上に移動）
        "gravity": 0.0
    })
    
    return bullet_id

def get_bullets():
    bullets = find_entities_by_tag("bullet")
    # print("Found bullets:", len(bullets))  # デバッグ出力を削除
    return bullets

print("Starting initialization") # デバッグ出力

def init():
    # 初期化時にFull HDに変更
    # set_screen_resolution("FULL_HD")
    # プレイヤーエンティティの作成
    vars["player_id"] = create_entity()
    print("Created player entity:", vars["player_id"])

    # Transformコンポーネントの追加
    add_component(vars["player_id"], "transform", {
        "x": 600,  # HD幅の中央付近
        "y": 600   # HD高さの下部
    })
    
    # Spriteコンポーネントの追加
    add_component(vars["player_id"], "sprite", {
        "width": 32,
        "height": 32,
        "color": "white"
    })
        
    # 操作説明テキストの追加
    text_id = create_entity()
    add_component(text_id, "text", {
        "text": "矢印キーで移動できます\nスペースキーで弾を発射（最大10発）",
        "x": 20,
        "y": 50
    })

    # 統計情報表示用のテキスト
    vars["stats_text_id"] = create_entity()
    print("Created stats text entity:", vars["stats_text_id"])
    
    add_component(vars["stats_text_id"], "text", {
        "text": "Objects: "+str(get_total_entities())+" | Bullets: "+str(len(get_bullets())),
        "x": 400,
        "y": 550
    })
    print("Added text component to stats entity")  # デバッグ出力を追加

    # 初期状態をまとめて設定
    set_states(vars["player_id"], {
        "shot_cooldown": 0,
        "health": 100,
        "score": 0
    })

print("Initialization complete")
# 設定の読み込み
load("scripts_debug/config.star", "init_screen")  # 必要な関数だけをインポート
init()  # init()を直接呼び出す 

# オブジェクトの状態を設定
set_state(vars["player_id"], "health", 100)
set_state(vars["player_id"], "score", 0)

# 状態を取得
health = get_state(vars["player_id"], "health")
score = get_state(vars["player_id"], "score")

