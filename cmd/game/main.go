package main

import (
	"log"

	"spacebattle/internal/fonts"
	"spacebattle/internal/game"
	"spacebattle/internal/i18n"
	"spacebattle/internal/progress"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// 初始化国际化系统
	if err := i18n.Init(); err != nil {
		log.Printf("警告: 国际化系统初始化失败: %v", err)
	}

	// 初始化字体系统
	if err := fonts.Init(); err != nil {
		log.Printf("警告: 字体系统初始化失败: %v", err)
	}

	// 初始化进度存储（SQLite, 内置驱动 modernc.org/sqlite）
	if err := progress.Init("game_progress.db"); err != nil {
		log.Printf("警告: 进度存储初始化失败: %v", err)
	} else {
		_ = progress.Load()
	}

	// 创建游戏实例
	g := game.NewGame()

	// 设置窗口属性
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Space Shooter Game")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// 设置TPS和FPS分离
	ebiten.SetTPS(60) // 逻辑更新频率：60 TPS (游戏逻辑、碰撞检测等)
	ebiten.SetVsyncEnabled(true)

	// 运行游戏
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
