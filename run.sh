#!/bin/bash

echo "正在启动游戏..."
echo "=================================="
echo "太空射击游戏 - 方向键移动，空格射击"
echo "=================================="
echo "控制说明："
echo "- 方向键：选择游戏/移动"
echo "- 回车/空格：确认/射击"
echo "- ESC键：返回主菜单/退出"
echo "- R键：重新开始游戏"
echo "=================================="

# 检查是否已编译
if [ ! -f "./game" ]; then
    echo "正在编译游戏..."
    go build -o game cmd/game/main.go
    if [ $? -ne 0 ]; then
        echo "编译失败！"
        exit 1
    fi
fi

# 运行游戏
./game
