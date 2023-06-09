package gfx

import (
	"fmt"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/pawelktk/TuxMan/game"
	"github.com/pawelktk/TuxMan/globals"
)

var ScreenWidth int32
var ScreenHeight int32

type Gfx struct {
	Size_x                      int32
	Size_y                      int32
	Tile_size                   int32
	Game_Texture                rl.RenderTexture2D
	Game_Texture_Size_x         int32
	Game_Texture_Size_y         int32
	SpriteSheet                 rl.Texture2D
	AnimatePlayer               [4]bool
	PlayerAnimationFrameCounter int32
	PlayerAnimationCurrentFrame int32
	Shader                      rl.Shader
}

func NewGfx(size_x, size_y int32) Gfx {
	gfx := Gfx{}
	gfx.Size_x = size_x
	gfx.Size_y = size_y
	ScreenWidth = size_x
	ScreenHeight = size_y
	gfx.Tile_size = globals.GLOBAL_TILE_SIZE //= 30

	gfx.PlayerAnimationFrameCounter = 0
	rl.InitWindow(size_x, size_y, "TuxMan")
	gfx.SpriteSheet = rl.LoadTexture("resources/spritesheet.png")
	//defer rl.CloseWindow()
	rl.SetTargetFPS(30)
	//gfx.Shader = rl.LoadShader("", "scanline.fs")
	return gfx
}
func (gfx *Gfx) DrawBoard(currentGame *game.Game) {
	/*
		for i := range currentGame.GameBoard.board_matrix {
			for j := range currentGame.GameBoard.board_matrix[i] {
				if currentGame.GameBoard.board_matrix[i][j] == Blank {
					rl.DrawRectangle(int32(j)*gfx.Tile_size, int32(i)*gfx.Tile_size, gfx.Tile_size, gfx.Tile_size, rl.Gray)
				} else if currentGame.GameBoard.board_matrix[i][j] == Wall {
					rl.DrawRectangle(int32(j)*gfx.Tile_size, int32(i)*gfx.Tile_size, gfx.Tile_size, gfx.Tile_size, rl.Black)
				}
			}
		}*/
	rand.Seed(int64(currentGame.GameBoard.Size_x + currentGame.GameBoard.Size_y))
	for i := 0; i < int(currentGame.GameBoard.Size_y); i++ {
		for j := 0; j < int(currentGame.GameBoard.Size_x); j++ {
			if rand.Float32() > 0.8 {
				gfx.DrawStaticTexture("tile_blank", int32(j*globals.GLOBAL_TILE_SIZE), int32(i*globals.GLOBAL_TILE_SIZE))
			} else {
				gfx.DrawStaticTexture("tile_grass", int32(j*globals.GLOBAL_TILE_SIZE), int32(i*globals.GLOBAL_TILE_SIZE))

			}

		}
	}
	//rl.DrawRectangle(0, 0, currentGame.GameBoard.Size_x*gfx.Tile_size, currentGame.GameBoard.Size_y*gfx.Tile_size, rl.Gray)

}

func (gfx *Gfx) DrawObstacles(currentGame *game.Game) {
	for _, v := range currentGame.GameBoard.Obstacles {
		if v.ObstacleType == game.Wall {
			//rl.DrawRectangleRec(v.HitBox, rl.Black)
			gfx.DrawStaticTexture("boulder", int32(v.HitBox.X), int32(v.HitBox.Y))
		} else if v.ObstacleType == game.Breakable {
			//rl.DrawRectangleRec(v.HitBox, rl.Brown)
			gfx.DrawStaticTexture("bush", int32(v.HitBox.X), int32(v.HitBox.Y))
		}
	}
}

func (gfx *Gfx) UpdatePlayerFrameCounter() {
	gfx.PlayerAnimationFrameCounter++
	if gfx.PlayerAnimationFrameCounter >= 5 {
		gfx.PlayerAnimationFrameCounter = 0
		gfx.PlayerAnimationCurrentFrame++
		if gfx.PlayerAnimationCurrentFrame > 1 {
			gfx.PlayerAnimationCurrentFrame = 0
		}
	}
}

func (gfx *Gfx) DrawPlayers(currentGame *game.Game) {
	gfx.UpdatePlayerFrameCounter()
	for i, v := range currentGame.Players {
		if v.Status {

			if gfx.AnimatePlayer[i] == true {
				if gfx.PlayerAnimationCurrentFrame == 0 {
					gfx.DrawDynamicTexture("player"+fmt.Sprint(i+1)+"_1", v.Position.X, v.Position.Y)

				} else {
					gfx.DrawDynamicTexture("player"+fmt.Sprint(i+1)+"_2", v.Position.X, v.Position.Y)

				}
			} else {
				gfx.DrawDynamicTexture("player"+fmt.Sprint(i+1)+"_0", v.Position.X, v.Position.Y)
			}
		}
		//rl.DrawRectangle(int32(v.Position.X)*gfx.Tile_size, int32(v.Position.Y)*gfx.Tile_size, gfx.Tile_size, gfx.Tile_size, rl.Red)
	}
}

func (gfx *Gfx) DrawBombs(currentGame *game.Game) {
	for _, v := range currentGame.Bombs {
		//rl.DrawRectangle(v.Position_x, v.Position_y, globals.GLOBAL_TILE_SIZE, globals.GLOBAL_TILE_SIZE, rl.Red)
		//rl.DrawTextureRec(gfx.SpriteSheet, gfx.GetTextureRec("bomb"), rl.NewVector2(float32(v.Position_x), float32(v.Position_y)), rl.White)

		/*
			destRect := rl.NewRectangle(float32(v.Position_x), float32(v.Position_y), globals.GLOBAL_TILE_SIZE, globals.GLOBAL_TILE_SIZE)
			rl.DrawTexturePro(gfx.SpriteSheet, gfx.GetTextureRec("bomb"), destRect, rl.NewVector2(globals.GLOBAL_TILE_SIZE, globals.GLOBAL_TILE_SIZE), 180, rl.White)
		*/
		if v.RemainingTicks <= 20 {
			gfx.DrawStaticTexture("bomb_end", v.Position_x, v.Position_y)
		} else if v.RemainingTicks <= 40 {
			gfx.DrawStaticTexture("bomb_mid", v.Position_x, v.Position_y)
		} else {

			gfx.DrawStaticTexture("bomb", v.Position_x, v.Position_y)
		}
	}
}
func (gfx *Gfx) DrawShrapnel(currentGame *game.Game) {
	for _, v := range currentGame.Shrapnels {
		//rl.DrawRectangle(v.Position_x, v.Position_y, globals.GLOBAL_TILE_SIZE, globals.GLOBAL_TILE_SIZE, rl.Orange)
		if v.RemainingTicks <= 4 {
			gfx.DrawStaticTexture("shrapnel", v.Position_x, v.Position_y)
		} else if v.RemainingTicks <= 7 {
			gfx.DrawStaticTexture("shrapnel_medium", v.Position_x, v.Position_y)
		} else {
			gfx.DrawStaticTexture("shrapnel_small", v.Position_x, v.Position_y)
		}
	}
}

func (gfx *Gfx) DrawPowerUps(currentGame *game.Game) {
	for _, v := range currentGame.PowerUps {
		if v.BoostType == game.Speed {
			gfx.DrawStaticTexture("powerup_speed", v.Position_x, v.Position_y)
		} else {
			gfx.DrawStaticTexture("powerup_blank", v.Position_x, v.Position_y)
		}
	}
}
func (gfx *Gfx) DrawStaticTexture(texture_name string, position_x, position_y int32) {
	destRect := rl.NewRectangle(float32(position_x), float32(position_y), globals.GLOBAL_TILE_SIZE, globals.GLOBAL_TILE_SIZE)
	rl.DrawTexturePro(gfx.SpriteSheet, gfx.GetTextureRec(texture_name), destRect, rl.NewVector2(globals.GLOBAL_TILE_SIZE, globals.GLOBAL_TILE_SIZE), 180, rl.White)
}

func (gfx *Gfx) DrawDynamicTexture(texture_name string, position_x, position_y float32) {
	destRect := rl.NewRectangle(position_x, position_y, globals.GLOBAL_TILE_SIZE, globals.GLOBAL_TILE_SIZE)
	rl.DrawTexturePro(gfx.SpriteSheet, gfx.GetTextureRec(texture_name), destRect, rl.NewVector2(globals.GLOBAL_TILE_SIZE*1.0, globals.GLOBAL_TILE_SIZE*1.0), 180, rl.White)
}

func (gfx *Gfx) DrawTexture(texture_name string, position_x, position_y int32, dest_size float32) {
	destRect := rl.NewRectangle(float32(position_x), float32(position_y), dest_size, dest_size)
	rl.DrawTexturePro(gfx.SpriteSheet, gfx.GetTextureRec(texture_name), destRect, rl.NewVector2(0, 0), 0, rl.White)
}
func (gfx *Gfx) GenerateGameTexture(currentGame *game.Game) {
	//TODO check for texture init
	rl.BeginTextureMode(gfx.Game_Texture)
	rl.ClearBackground(rl.White)
	gfx.DrawBoard(currentGame)
	gfx.DrawObstacles(currentGame)
	gfx.DrawBombs(currentGame)
	gfx.DrawPowerUps(currentGame)
	gfx.DrawPlayers(currentGame)
	gfx.DrawShrapnel(currentGame)
	rl.EndTextureMode()

}

func (gfx *Gfx) HandleInput(currentGame *game.Game, deltatime float32) {
	var playerAlreadyChecked [4]bool
	playerAlreadyChecked = [4]bool{false, false, false, false}

	player1Key, player1KeyIsPressed := gfx.GetPlayer1Key()

	if gfx.IsPlayer1BombKeyPressed() && currentGame.Players[0].AvailableBombs > 0 {
		currentGame.PlaceBomb(&currentGame.Players[0], currentGame.Players[0].Position, 3)
	}

	if player1KeyIsPressed && !playerAlreadyChecked[0] {
		playerAlreadyChecked[0] = true
		currentGame.MovePlayer(&currentGame.Players[0], player1Key, deltatime)
		fmt.Println("Key pressed: ", player1Key)
		gfx.AnimatePlayer[0] = true
	} else {
		gfx.AnimatePlayer[0] = false
	}
	if len(currentGame.Players) > 1 {
		player2Key, player2KeyIsPressed := gfx.GetPlayer2Key()

		if gfx.IsPlayer2BombKeyPressed() && currentGame.Players[1].AvailableBombs > 0 {
			currentGame.PlaceBomb(&currentGame.Players[1], currentGame.Players[1].Position, 3)
		}

		if player2KeyIsPressed && !playerAlreadyChecked[1] {
			playerAlreadyChecked[1] = true
			currentGame.MovePlayer(&currentGame.Players[1], player2Key, deltatime)
			fmt.Println("Key pressed: ", player2Key)
			gfx.AnimatePlayer[1] = true
		} else {
			gfx.AnimatePlayer[1] = false
		}

	}
	if len(currentGame.Players) > 2 {
		player3Key, player3KeyIsPressed := gfx.GetPlayer3Key()

		if gfx.IsPlayer3BombKeyPressed() && currentGame.Players[2].AvailableBombs > 0 {
			currentGame.PlaceBomb(&currentGame.Players[2], currentGame.Players[2].Position, 3)
		}

		if player3KeyIsPressed && !playerAlreadyChecked[2] {
			playerAlreadyChecked[2] = true
			currentGame.MovePlayer(&currentGame.Players[2], player3Key, deltatime)
			fmt.Println("Key pressed: ", player3Key)
			gfx.AnimatePlayer[2] = true
		} else {
			gfx.AnimatePlayer[2] = false
		}

	}
}

// DONE Make placing bombs and movement more responsive

func (gfx *Gfx) IsPlayer1BombKeyPressed() bool {
	if rl.IsKeyPressed(rl.KeyComma) {
		return true
	} else {
		return false
	}
}

func (gfx *Gfx) GetPlayer1Key() (string, bool) {
	switch {
	case rl.IsKeyDown(rl.KeyLeft):
		if rl.IsKeyDown(rl.KeyUp) {
			return "left-up", true
		} else if rl.IsKeyDown(rl.KeyDown) {
			return "left-down", true
		}
		return "left", true
	case rl.IsKeyDown(rl.KeyRight):
		if rl.IsKeyDown(rl.KeyUp) {
			return "right-up", true
		} else if rl.IsKeyDown(rl.KeyDown) {
			return "right-down", true
		}
		return "right", true
	case rl.IsKeyDown(rl.KeyUp):
		return "up", true
	case rl.IsKeyDown(rl.KeyDown):
		return "down", true
	default:
		return "", false
	}
}

func (gfx *Gfx) IsPlayer2BombKeyPressed() bool {
	if rl.IsKeyPressed(rl.KeySpace) {
		return true
	} else {
		return false
	}
}

func (gfx *Gfx) GetPlayer2Key() (string, bool) {
	switch {
	case rl.IsKeyDown(rl.KeyA):
		if rl.IsKeyDown(rl.KeyW) {
			return "left-up", true
		} else if rl.IsKeyDown(rl.KeyS) {
			return "left-down", true
		}
		return "left", true
	case rl.IsKeyDown(rl.KeyD):
		if rl.IsKeyDown(rl.KeyW) {
			return "right-up", true
		} else if rl.IsKeyDown(rl.KeyS) {
			return "right-down", true
		}
		return "right", true
	case rl.IsKeyDown(rl.KeyW):
		return "up", true
	case rl.IsKeyDown(rl.KeyS):
		return "down", true
	default:
		return "", false
	}
}

func (gfx *Gfx) IsPlayer3BombKeyPressed() bool {
	if rl.IsKeyPressed(rl.KeyV) {
		return true
	} else {
		return false
	}
}

func (gfx *Gfx) GetPlayer3Key() (string, bool) {
	switch {
	case rl.IsKeyDown(rl.KeyG):
		if rl.IsKeyDown(rl.KeyY) {
			return "left-up", true
		} else if rl.IsKeyDown(rl.KeyH) {
			return "left-down", true
		}
		return "left", true
	case rl.IsKeyDown(rl.KeyJ):
		if rl.IsKeyDown(rl.KeyY) {
			return "right-up", true
		} else if rl.IsKeyDown(rl.KeyH) {
			return "right-down", true
		}
		return "right", true
	case rl.IsKeyDown(rl.KeyY):
		return "up", true
	case rl.IsKeyDown(rl.KeyH):
		return "down", true
	default:
		return "", false
	}
}

func (gfx *Gfx) GetTextureRec(texture_name string) rl.Rectangle {
	var x, y float32
	switch texture_name {
	case "bomb":
		x = 4
		y = 2
	case "bomb_mid":
		x = 4
		y = 3
	case "bomb_end":
		x = 4
		y = 4

	case "tile_blank":
		x = 0
		y = 0
	case "tile_grass":
		x = 1
		y = 0
	case "bush":
		x = 3
		y = 0
	case "boulder":
		x = 6
		y = 3

	case "shrapnel_small":
		x = 7
		y = 3

	case "shrapnel_medium":
		x = 7
		y = 4
	case "shrapnel":
		x = 7
		y = 5

	case "player1_0":
		x = 1
		y = 1
	case "player1_1":
		x = 2
		y = 1
	case "player1_2":
		x = 3
		y = 1

	case "player2_0":
		x = 1
		y = 2
	case "player2_1":
		x = 2
		y = 2
	case "player2_2":
		x = 3
		y = 2

	case "player3_0":
		x = 1
		y = 3
	case "player3_1":
		x = 2
		y = 3
	case "player3_2":
		x = 3
		y = 3

	case "powerup_blank":
		x = 8
		y = 0
	case "powerup_speed":
		x = 8
		y = 2

	default:
		x = 0
		y = 0
	}
	return rl.NewRectangle(x*16, y*16, 16, 16)
}

func (gfx *Gfx) InitGameTextureBox(currentGame *game.Game) {
	gfx.Game_Texture_Size_x = gfx.Tile_size * currentGame.GameBoard.Size_x
	gfx.Game_Texture_Size_y = gfx.Tile_size * currentGame.GameBoard.Size_y
	gfx.Game_Texture = rl.LoadRenderTexture(gfx.Game_Texture_Size_x, gfx.Game_Texture_Size_y) // rl.LoadRenderTexture(int32(float32(size_x)*0.8), int32(float32(size_y)*0.8))
	rl.SetTextureFilter(gfx.Game_Texture.Texture, rl.TextureFilterMode(rl.RL_TEXTURE_FILTER_BILINEAR))
}

func (gfx *Gfx) TextCenterX(text string, fontSize int32) int32 {
	return (gfx.Size_x - rl.MeasureText(text, fontSize)) / 2
}
func (gfx *Gfx) TextCenterY(text string, fontSize int32) int32 {
	sizeVector := rl.MeasureTextEx(rl.GetFontDefault(), text, float32(fontSize), 0)
	return (gfx.Size_y - int32(sizeVector.Y)) / 2
}

func (gfx *Gfx) DrawTextCenterX(text string, fontSize int32, y int32) {
	rl.DrawText(text, gfx.TextCenterX(text, fontSize), y, fontSize, rl.Black)
}

func (gfx *Gfx) DrawTextOnScreenPart(text string, fontSize int32, xPart, yPart float32) {
	rl.DrawText(text, int32(float32(gfx.TextCenterX(text, fontSize)*2)*xPart), int32(float32(gfx.TextCenterY(text, fontSize)*2)*yPart), fontSize, rl.Black)
}
