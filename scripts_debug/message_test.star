def init():
    # メッセージウィンドウの作成
    window_id = create_entity()
    
    add_component(window_id, "transform", {
        "x": 50,
        "y": 400,
        "scale_x": 1.0,
        "scale_y": 1.0
    })
    
    add_component(window_id, "message_window", {
        "width": 700,
        "height": 150,
        "padding": 20,
        "line_height": 24,
        "background_color": "rgba(0, 0, 0, 0.8)",
        "text_color": "rgb(255, 255, 255)",
        "font_size": 18
    })
    
    # メッセージシステムの初期化
    message_system_id = create_entity()
    add_component(message_system_id, "message_system", {
        "messages": [
            "これはメッセージウィンドウのテストです。",
            "文字を1文字ずつ表示することができます。",
            "また、選択肢を表示することもできます。\n[選択肢を表示しますか？]",
            {
                "type": "choice",
                "text": "選択してください：",
                "options": [
                    "はい、表示します",
                    "いいえ、終了します"
                ]
            }
        ],
        "char_delay": 0.05,  # 1文字の表示間隔（秒）
        "window_id": window_id
    })
    
    # 入力コンポーネントの追加
    add_component(message_system_id, "input", {
        "keys": {
            "next": "Space",
            "up": "ArrowUp",
            "down": "ArrowDown",
            "select": "Enter"
        }
    })

    print("Message window sample initialized!")

init() 