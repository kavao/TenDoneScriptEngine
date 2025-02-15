# 画面設定の定数
HD = {
    "width": 1280,
    "height": 720,
    "scale": 1.0
}

FULLHD = {
    "width": 1920,
    "height": 1080,
    "scale": 1.0
}

def init_screen(mode="HD"):
    """画面設定を初期化します"""
    config = HD if mode == "HD" else FULLHD
    
    # 画面設定エンティティを作成
    config_entity = create_entity()
    add_component(config_entity, "screen_config", config)