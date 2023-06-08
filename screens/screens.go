package screens

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/pawelktk/TuxMan/game"
	"github.com/pawelktk/TuxMan/gfx"
)

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

type MainMenuScreen struct {
	HighlightedOption int
}

func NewMainMenuScreen() MainMenuScreen {
	screen := MainMenuScreen{}
	screen.HighlightedOption = 0
	return screen
}

func (screen *MainMenuScreen) Display(gfx *gfx.Gfx) {
	rl.ClearBackground(rl.RayWhite)
	gfx.DrawTextCenterX("TuxMan", 40, 10)

	//pressedKey := rl.GetKeyPressed()

}

