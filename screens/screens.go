package screens

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/pawelktk/TuxMan/game"
	"github.com/pawelktk/TuxMan/gfx"
)

var HighlightedMenuOption = 0
var SelectedMenuOption = -1

var SelectedPlayerCount = 2

var OptionsCount_MainMenu = 3

func GameOverScreen(gfx *gfx.Gfx, game *game.Game, winner *game.Player) {
	rl.ClearBackground(rl.RayWhite)
	gfx.DrawTextOnScreenPart("Game Over!", 50, 0.5, 0.3)
	gfx.DrawTextOnScreenPart("Winner: "+winner.Name, 40, 0.5, 0.4)
	gfx.DrawTextOnScreenPart("Kills: "+fmt.Sprint(winner.Points), 40, 0.5, 0.5)

}

func GameScreen(gfx *gfx.Gfx, game *game.Game) {
	rl.ClearBackground(rl.RayWhite)
	gfx.DrawTextCenterX("TuxMan", 40, 10)

	texturePosition := rl.NewVector2(float32((gfx.Size_x-gfx.Game_Texture_Size_x)/2), 60)
	rl.DrawTextureEx(gfx.Game_Texture.Texture, texturePosition, 0, 1, rl.White)
}

func MenuHandleInput(gameWindow *gfx.Gfx, optionsCount int) {
	pressedKey := rl.GetKeyPressed()
	SelectedMenuOption = -1

	switch pressedKey {
	case rl.KeyDown:
		fallthrough
	case rl.KeyS:
		HighlightedMenuOption = (HighlightedMenuOption + 1) % optionsCount
	case rl.KeyUp:
		fallthrough
	case rl.KeyW:
		HighlightedMenuOption = int(math.Abs(float64(HighlightedMenuOption-1))) % optionsCount
	case rl.KeySpace:
		fallthrough
	case rl.KeyEnter:
		SelectedMenuOption = HighlightedMenuOption
	}
}

func MainMenuScreen(gameWindow *gfx.Gfx) {

	rl.ClearBackground(rl.RayWhite)
	gameWindow.DrawTextCenterX("TuxMan", 80, 30)

	drawButtonCenterX(gameWindow, 200, 300, 70, "Play", HighlightedMenuOption == 0)
	drawButtonCenterX(gameWindow, 300, 300, 70, "Players: "+fmt.Sprint(SelectedPlayerCount), HighlightedMenuOption == 1)
	drawButtonCenterX(gameWindow, 400, 300, 70, "Quit", HighlightedMenuOption == 2)
}

func GameInitMenuScreen(gameWindow *gfx.Gfx) {

	rl.ClearBackground(rl.RayWhite)
	gameWindow.DrawTextCenterX("TuxMan", 80, 30)

	drawPlayerConfig(gameWindow, 150, 400, 200)

}

func drawPlayerConfig(gameWindow *gfx.Gfx, y int32, width, height int32) {
	box := rl.NewRectangle(float32((gfx.ScreenWidth-width)/2), float32(y), float32(width), float32(height))

	rl.DrawRectangleRoundedLines(box, 0.1, 0, 5, rl.LightGray)
}

func drawButtonCenterX(gameWindow *gfx.Gfx, y, width, height int32, text string, selected bool) {
	drawButton(gameWindow, (gfx.ScreenWidth-width)/2, y, width, height, text, selected)
}

func drawButton(gameWindow *gfx.Gfx, x, y, width, height int32, text string, selected bool) {
	button := rl.NewRectangle(float32(x), float32(y), float32(width), float32(height))
	if selected {
		rl.DrawRectangleRoundedLines(button, 10, 0, 5, rl.Orange)
	} else {
		rl.DrawRectangleRoundedLines(button, 10, 0, 5, rl.LightGray)
	}
	textSize := rl.MeasureTextEx(rl.GetFontDefault(), text, float32(height/2), 0)

	gameWindow.DrawTextCenterX(text, height/2, y+int32(textSize.Y/2))
	//textSize := rl.MeasureTextEx(rl.GetFontDefault(), text, float32(height/2), 0)
	//gameWindow.DrawText(text, x+int32(textSize.X/2), y+int32(textSize.Y/2), height/2, rl.Black)
}
