package tests

import (
	"testing"
	"spacebattle/internal/game"
)

func TestNewGame(t *testing.T) {
	g := game.NewGame()
	if g == nil {
		t.Error("NewGame() 返回了 nil")
	}
}

func TestGameLayout(t *testing.T) {
	g := game.NewGame()
	width, height := g.Layout(800, 600)
	
	if width != 800 {
		t.Errorf("期望宽度 800，实际得到 %d", width)
	}
	
	if height != 600 {
		t.Errorf("期望高度 600，实际得到 %d", height)
	}
}
