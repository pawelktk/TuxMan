package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/pawelktk/TuxMan/audio"
	"github.com/pawelktk/TuxMan/game"
	"github.com/pawelktk/TuxMan/gfx"
	"github.com/pawelktk/TuxMan/globals"
	"github.com/pawelktk/TuxMan/screens"
)

func main() {
	//rl.SetConfigFlags(rl.FlagWindowResizable)
	currentGame := game.NewGame()
	gameWindow := gfx.NewGfx(int32(1920*0.7), int32(1080*0.7))
	audio.InitAudio()
	gameWindow.InitGameTextureBox(&currentGame)
	currentGame.AddPlayer("Pablo", 0, 1)
	currentGame.AddPlayer("SecondPlayer", globals.GLOBAL_TILE_SIZE*float32(currentGame.GameBoard.Size_x-1), 1)

	mainMenuScreen := screens.NewMainMenuScreen()
	screen := "NOT_main_menu" //TODO finish main menu
	//currentGame.GameBoard.AddObstacle(4, 1, Wall)
	//currentGame.GameBoard.AddObstacle(2, 2, Breakable)
	go audio.MainAudioLoop()
	for !rl.WindowShouldClose() {
		gameWindow.HandleInput(&currentGame, rl.GetFrameTime())
		currentGame.Update()

		gameWindow.GenerateGameTexture(&currentGame)

		gameOver, winner := currentGame.GameShouldEnd()
		rl.BeginDrawing()
		//rl.BeginShaderMode(gfx.Shader)
		if gameOver {
			screens.GameOverScreen(&gameWindow, &currentGame, winner)
		} else if screen == "main_menu" {
			mainMenuScreen.Display(&gameWindow)
		} else {
			screens.GameScreen(&gameWindow, &currentGame)
		}
		//rl.EndShaderMode()
		rl.EndDrawing()
	}
}
