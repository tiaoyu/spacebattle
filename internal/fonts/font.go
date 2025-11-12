package fonts

import (
	"image/color"

	"github.com/hajimehoshi/bitmapfont/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

// FontManager 字体管理器
type FontManager struct {
	regularFont font.Face
	largeFont   font.Face
}

var globalFontManager *FontManager

// Init 初始化字体管理器
func Init() error {
	globalFontManager = &FontManager{}

	// 尝试加载自定义字体
	if err := globalFontManager.loadCustomFonts(); err != nil {
		// 如果加载失败，使用默认字体
		globalFontManager.loadDefaultFonts()
	}

	// 确保字体已正确设置
	if globalFontManager.regularFont == nil {
		globalFontManager.regularFont = bitmapfont.Face
	}
	if globalFontManager.largeFont == nil {
		globalFontManager.largeFont = bitmapfont.Face
	}

	return nil
}

// loadCustomFonts 加载自定义字体
func (fm *FontManager) loadCustomFonts() error {
	// 这里可以加载自定义字体文件
	// 例如：加载支持中文的字体文件
	return nil
}

// loadDefaultFonts 加载默认字体
func (fm *FontManager) loadDefaultFonts() {
	// 使用bitmapfont，支持更多字符
	fm.regularFont = bitmapfont.Face
	fm.largeFont = bitmapfont.Face
}

// GetRegularFont 获取常规字体
func GetRegularFont() font.Face {
	if globalFontManager != nil && globalFontManager.regularFont != nil {
		return globalFontManager.regularFont
	}
	return bitmapfont.Face
}

// GetLargeFont 获取大字体
func GetLargeFont() font.Face {
	if globalFontManager != nil && globalFontManager.largeFont != nil {
		return globalFontManager.largeFont
	}
	return bitmapfont.Face
}

// DrawText 绘制文本（支持中文）
func DrawText(screen *ebiten.Image, textStr string, x, y int, clr color.Color) {
	// 使用text包绘制文本，支持更多字符
	text.Draw(screen, textStr, GetRegularFont(), x, y, clr)
}

// DrawTextLarge 绘制大文本
func DrawTextLarge(screen *ebiten.Image, textStr string, x, y int, clr color.Color) {
	text.Draw(screen, textStr, GetLargeFont(), x, y, clr)
}

// MeasureText 测量文本宽度
func MeasureText(textStr string) int {
	bounds, _ := font.BoundString(GetRegularFont(), textStr)
	return (bounds.Max.X - bounds.Min.X).Ceil()
}

// MeasureTextLarge 测量大文本宽度
func MeasureTextLarge(textStr string) int {
	bounds, _ := font.BoundString(GetLargeFont(), textStr)
	return (bounds.Max.X - bounds.Min.X).Ceil()
}

// DrawTextCentered 绘制居中文本
func DrawTextCentered(screen *ebiten.Image, textStr string, x, y, width int, clr color.Color) {
	textWidth := MeasureText(textStr)
	startX := x + (width-textWidth)/2
	DrawText(screen, textStr, startX, y, clr)
}

// DrawTextCenteredLarge 绘制居中大文本
func DrawTextCenteredLarge(screen *ebiten.Image, textStr string, x, y, width int, clr color.Color) {
	textWidth := MeasureTextLarge(textStr)
	startX := x + (width-textWidth)/2
	DrawTextLarge(screen, textStr, startX, y, clr)
}
