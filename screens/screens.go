package screens

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/pawelktk/TuxMan/game"
	"github.com/pawelktk/TuxMan/gfx"
	"github.com/pawelktk/TuxMan/globals"
)

var HighlightedMenuOption = 0
var SelectedMenuOption = -1
var EnteredKey = -1
var AnimationFrame = 0

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
	EnteredKey = -1
	pressedKey := rl.GetKeyPressed()
	SelectedMenuOption = -1

	switch pressedKey {
	case rl.KeyDown:
		//	fallthrough
		//case rl.KeyS:
		HighlightedMenuOption = (HighlightedMenuOption + 1) % optionsCount
	case rl.KeyUp:
		//	fallthrough
		//case rl.KeyW:
		HighlightedMenuOption = int(math.Abs(float64(HighlightedMenuOption-1))) % optionsCount
	case rl.KeySpace:
		fallthrough
	case rl.KeyEnter:
		SelectedMenuOption = HighlightedMenuOption
	default:
		if (pressedKey >= 32 && pressedKey <= 126) || pressedKey == rl.KeyBackspace {
			EnteredKey = int(pressedKey)

		}
	}
}

func UpdatePlayerName(currentGame *game.Game) {
	if EnteredKey == rl.KeyBackspace && len(currentGame.Players[HighlightedMenuOption].Name) > 0 {
		currentGame.Players[HighlightedMenuOption].Name = currentGame.Players[HighlightedMenuOption].Name[:len(currentGame.Players[HighlightedMenuOption].Name)-1]
	} else if EnteredKey != rl.KeyBackspace && len(currentGame.Players[HighlightedMenuOption].Name) < globals.MAX_PLAYERNAME_SIZE {
		currentGame.Players[HighlightedMenuOption].Name = currentGame.Players[HighlightedMenuOption].Name + string(EnteredKey)
	}
}

func MainMenuScreen(gameWindow *gfx.Gfx) {

	rl.ClearBackground(rl.RayWhite)
	gameWindow.DrawTextCenterX("TuxMan", 80, 30)

	drawButtonCenterX(gameWindow, 200, 300, 70, "Play", HighlightedMenuOption == 0)
	drawButtonCenterX(gameWindow, 300, 300, 70, "Players: "+fmt.Sprint(SelectedPlayerCount), HighlightedMenuOption == 1)
	drawButtonCenterX(gameWindow, 400, 300, 70, "Quit", HighlightedMenuOption == 2)
}

func GameInitMenuScreen(gameWindow *gfx.Gfx, currentGame *game.Game) {

	rl.ClearBackground(rl.RayWhite)
	gameWindow.DrawTextCenterX("TuxMan", 80, 30)

	drawPlayerConfigBox(gameWindow, &currentGame.Players[0], 150, 600, 150, HighlightedMenuOption == 0)
	drawPlayerConfigBox(gameWindow, &currentGame.Players[1], 320, 600, 150, HighlightedMenuOption == 1)

	if SelectedPlayerCount > 2 {
		drawPlayerConfigBox(gameWindow, &currentGame.Players[2], 490, 600, 150, HighlightedMenuOption == 2)
	}
}

func drawPlayerConfigBox(gameWindow *gfx.Gfx, player *game.Player, y int32, width, height int32, selected bool) {
	box := rl.NewRectangle(float32((gfx.ScreenWidth-width)/2), float32(y), float32(width), float32(height))
	if selected {
		rl.DrawRectangleRoundedLines(box, 0.1, 0, 5, rl.Orange)
	} else {
		rl.DrawRectangleRoundedLines(box, 0.1, 0, 5, rl.LightGray)
	}
	prompt := ""
	AnimationFrame++
	if AnimationFrame > 30 && selected {
		prompt = "_"
	}
	AnimationFrame %= 60
	rl.DrawText("PLAYER "+fmt.Sprint(player.ID+1), gameWindow.TextCenterX("PLAYER 1", 30)+50, y+5, 30, rl.Black)
	rl.DrawText(player.Name+prompt, int32((float32(gfx.ScreenWidth)-float32(width)/3)/2)+15, y+55, 30, rl.Black)
	nameBox := rl.NewRectangle((float32(gfx.ScreenWidth)-float32(width)/3)/2, float32(y+40), float32(width)*0.5, float32(height)*0.4)

	gameWindow.DrawTexture("player"+fmt.Sprint(player.ID+1)+"_0", int32((gfx.ScreenWidth-width)/2)+int32(float32(width)*0.05), y+5, box.Height*0.8)
	rl.DrawRectangleLinesEx(nameBox, 5, rl.Gray)
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
