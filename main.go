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

	screen := "main_menu"
	//currentGame.GameBoard.AddObstacle(4, 1, Wall)
	//currentGame.GameBoard.AddObstacle(2, 2, Breakable)
	for !rl.WindowShouldClose() {
		audio.MainAudio()
		if screen == "main_menu" {
			screens.MenuHandleInput(&gameWindow, screens.OptionsCount_MainMenu)
		} else if screen == "game_init" {
			screens.MenuHandleInput(&gameWindow, screens.SelectedPlayerCount)
			if screens.EnteredKey != -1 {
				screens.UpdatePlayerName(&currentGame)
			}
		} else {
			gameWindow.HandleInput(&currentGame, rl.GetFrameTime())
		}
		currentGame.Update()

		gameWindow.GenerateGameTexture(&currentGame)

		gameOver, winner := currentGame.GameShouldEnd()
		rl.BeginDrawing()
		//rl.BeginShaderMode(gfx.Shader)
		if gameOver {
			screens.GameOverScreen(&gameWindow, &currentGame, winner)
		} else if screen == "main_menu" {
			switch screens.SelectedMenuOption {
			case 0:
				if screens.SelectedPlayerCount > 0 {
					currentGame.AddPlayer("PLAYER1", 0, 1)
				}
				if screens.SelectedPlayerCount > 1 {
					currentGame.AddPlayer("PLAYER2", globals.GLOBAL_TILE_SIZE*float32(currentGame.GameBoard.Size_x-1), 1)
				}
				if screens.SelectedPlayerCount > 2 {
					currentGame.AddPlayer("PLAYER3", 1, globals.GLOBAL_TILE_SIZE*float32(currentGame.GameBoard.Size_x-1))
				}

				screen = "game_init"
			case 1:
				screens.SelectedPlayerCount = screens.SelectedPlayerCount + 1
				if screens.SelectedPlayerCount > globals.MAX_PLAYERS {
					screens.SelectedPlayerCount = 2
				}

			case 2:
				rl.CloseWindow()
				return
			default:
				screens.MainMenuScreen(&gameWindow)
			}

		} else if screen == "game" {
			screens.GameScreen(&gameWindow, &currentGame)
		} else if screen == "game_init" {
			if screens.SelectedMenuOption == -1 {
				screens.GameInitMenuScreen(&gameWindow, &currentGame)
			} else {
				screen = "game"
			}
		}
		//rl.EndShaderMode()
		rl.EndDrawing()
	}
}
