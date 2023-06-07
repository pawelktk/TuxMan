package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func GameOverScreen(gfx *Gfx, game *Game, winner *Player) {
	rl.ClearBackground(rl.RayWhite)
	gfx.DrawTextOnScreenPart("Game Over!", 50, 0.5, 0.3)
	gfx.DrawTextOnScreenPart("Winner: "+winner.Name, 40, 0.5, 0.4)
	gfx.DrawTextOnScreenPart("Kills: "+fmt.Sprint(winner.Points), 40, 0.5, 0.5)

}

func GameScreen(gfx *Gfx, game *Game) {
	rl.ClearBackground(rl.RayWhite)
	gfx.DrawTextCenterX("TuxMan", 40, 10)

	texturePosition := rl.NewVector2(float32((gfx.Size_x-gfx.Game_Texture_Size_x)/2), 60)
	rl.DrawTextureEx(gfx.Game_Texture.Texture, texturePosition, 0, 1, rl.White)
}

func main() {
	//rl.SetConfigFlags(rl.FlagWindowResizable)
	game := NewGame()
	gfx := NewGfx(int32(1920*0.7), int32(1080*0.7))
	gfx.InitGameTextureBox(&game)
	game.AddPlayer("Pablo", 0, 1)
	game.AddPlayer("SecondPlayer", GLOBAL_TILE_SIZE*float32(game.GameBoard.Size_x-1), 1)

	//game.GameBoard.AddObstacle(4, 1, Wall)
	//game.GameBoard.AddObstacle(2, 2, Breakable)

	for !rl.WindowShouldClose() {
		gfx.HandleInput(&game, rl.GetFrameTime())
		game.Update()

		gfx.GenerateGameTexture(&game)

		gameOver, winner := game.GameShouldEnd()
		rl.BeginDrawing()
		//rl.BeginShaderMode(gfx.Shader)
		if gameOver {
			GameOverScreen(&gfx, &game, winner)
		} else {
			GameScreen(&gfx, &game)
		}
		//rl.EndShaderMode()
		rl.EndDrawing()
	}
}
