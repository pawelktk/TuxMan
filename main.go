package main

import (
	"fmt"
	"strings"
	//	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func RoundDownXtoY(x float32, y int32) int32 {
	x2 := int32(x)
	return x2 - x2%y
}
func RoundUpXtoY(x float32, y int32) int32 {
	x2 := int32(x)

	if x2%y == 0 {
		return x2
	} else {
		return (y - x2%y) + x2
	}
}
func RoundXtoY(x float32, y int32) int32 {
	x2 := int32(x)
	if x2%y == 0 {
		return x2
	} else if x2%y >= int32(0.5*float32(y)) {
		return RoundUpXtoY(x, y)
	} else {
		return RoundDownXtoY(x, y)
	}
}

type TileType int32

const (
	Blank TileType = iota
	Wall
	Breakable
)

type Board struct {
	Size_x       int32
	Size_y       int32
	board_matrix [][]TileType
	Obstacles    []Obstacle
}

func NewBoard(size_x, size_y int32) Board {
	board := Board{}
	board.Size_x = size_x
	board.Size_y = size_y
	board.board_matrix = make([][]TileType, size_y)
	for i := int32(0); i < size_y; i++ {
		board.board_matrix[i] = make([]TileType, size_x)
	}
	board.Clear()
	return board
}

func (board *Board) Clear() {
	for i := range board.board_matrix {
		for j := range board.board_matrix[i] {
			board.board_matrix[i][j] = 0
		}
	}
}
func (board *Board) Print() {
	for i := range board.board_matrix {
		for j := range board.board_matrix[i] {
			fmt.Printf("%v", board.board_matrix[i][j])
		}
		fmt.Printf("\n")
	}
}
func (board *Board) AddObstacle(position_x, position_y int32, tileType TileType) {
	obstacle := NewObstacle(position_x, position_y, tileType)
	board.Obstacles = append(board.Obstacles, obstacle)
}

type Obstacle struct {
	Position_x   int32
	Position_y   int32
	ObstacleType TileType
	HitBox       rl.Rectangle
}

func NewObstacle(position_x, position_y int32, tileType TileType) Obstacle {
	obstacle := Obstacle{}
	obstacle.Position_x = position_x
	obstacle.Position_y = position_y
	obstacle.ObstacleType = tileType
	obstacle.HitBox = rl.NewRectangle(float32(position_x*GLOBAL_TILE_SIZE), float32(position_y*GLOBAL_TILE_SIZE), GLOBAL_TILE_SIZE, GLOBAL_TILE_SIZE)
	return obstacle
}

type Player struct {
	Name           string
	Position       rl.Vector2
	Points         int32
	AvailableBombs int32
	Status         bool
	HitBox         rl.Rectangle
	Speed          float32
	PlayerSize     float32
}

func NewPlayer(name string, position_x, position_y float32) Player {
	player := Player{}
	player.Name = name
	player.Position.X = position_x
	player.Position.Y = position_y
	player.Points = 0
	player.AvailableBombs = 13
	player.Status = true
	player.PlayerSize = GLOBAL_TILE_SIZE - GLOBAL_TILE_SIZE*0.2
	player.HitBox = rl.NewRectangle(float32(position_x), float32(position_y), player.PlayerSize, player.PlayerSize)
	player.Speed = 50
	return player
}

func (player *Player) UpdatePosition(position rl.Vector2) {
	player.Position = position
	player.HitBox.X = position.X
	player.HitBox.Y = position.Y
}

func (player *Player) UpdatePositionFloat32(position_x, position_y float32) {
	player.Position.X = position_x
	player.Position.Y = position_y
	player.HitBox.X = position_x
	player.HitBox.Y = position_y
}

//func (player *Player) PlaceBomb()

type Bomb struct {
	Position_x     int32
	Position_y     int32
	RemainingTicks int32
	Radius         int32
	Owner          *Player
}

func NewBomb(owner *Player, position_x int32, position_y int32, radius int32) Bomb {
	bomb := Bomb{}
	bomb.Owner = owner
	bomb.Position_x = position_x
	bomb.Position_y = position_y
	bomb.RemainingTicks = 60
	bomb.Radius = 2

	return bomb
}

type Shrapnel struct {
	Position_x     int32
	Position_y     int32
	RemainingTicks int32
	Owner          *Player
}

func NewShrapnel(owner *Player, position_x int32, position_y int32) Shrapnel {
	shrapnel := Shrapnel{}
	shrapnel.Owner = owner
	shrapnel.Position_x = position_x
	shrapnel.Position_y = position_y
	shrapnel.RemainingTicks = 10
	return shrapnel
}

type Game struct {
	Ticks     int32
	GameBoard Board
	Players   []Player
	Bombs     []Bomb
	Shrapnels []Shrapnel
}

func NewGame() Game {
	game := Game{}
	game.Ticks = 0
	game.GameBoard = NewBoard(12, 6)
	return game
}
func (game *Game) AddPlayer(name string, position_x, position_y float32) {
	player := NewPlayer(name, position_x, position_y)
	game.Players = append(game.Players, player)
}

func (game *Game) GetNextPosition(player *Player, direction string, deltatime float32) rl.Vector2 {
	var nextPosition rl.Vector2
	nextPosition = player.Position

	direction_list := strings.Split(direction, "-")

	// diagonal movement is √2 faster
	movementVector := player.Speed * deltatime
	if len(direction_list) > 1 {
		movementVector /= 1.414213562

	}
	for _, v := range direction_list {
		switch v {
		case "up":
			nextPosition.Y += movementVector
		case "down":
			nextPosition.Y -= movementVector
		case "left":
			nextPosition.X -= movementVector
		case "right":
			nextPosition.X += movementVector

		}
	}

	return nextPosition
}

func (game *Game) GetNextPositionHitbox(player *Player, direction string, deltatime float32) rl.Rectangle {
	nextHitbox := player.HitBox
	nextPosition := game.GetNextPosition(player, direction, deltatime)
	nextHitbox.Y = nextPosition.Y
	nextHitbox.X = nextPosition.X
	return nextHitbox
}

func (game *Game) PositionIsValid(position rl.Vector2) bool {
	if position.X <= float32(game.GameBoard.Size_x-1)*GLOBAL_TILE_SIZE+0.2*GLOBAL_TILE_SIZE && position.Y <= float32(game.GameBoard.Size_y-1)*GLOBAL_TILE_SIZE+0.2*GLOBAL_TILE_SIZE && position.X >= 0 && position.Y >= 0 {
		return true
	} else {
		return false
	}
}

/*
	func (game *Game) PositionCollidesWithBomb(position [2]int32) bool {
		for _, bomb := range game.Bombs {
			if bomb.Position_x == position[0] && bomb.Position_y == position[1] {
				return true
			}
		}
		return false
	}
*/
func (game *Game) HitboxCollidesWithOtherPlayer(sourcePlayer *Player, hitbox rl.Rectangle) bool {
	for i := range game.Players {
		if game.Players[i] != *sourcePlayer && rl.CheckCollisionRecs(hitbox, game.Players[i].HitBox) {
			return true
		}
	}
	return false
}

/*
	func (game *Game) PositionCollidesWithObstacle(position [2]int32) bool {
		//WARNING: PositionIsValid MUST be run first
		fmt.Printf("position[0]: %v, position[1]: %v\n", position[0], position[1])
		if game.GameBoard.board_matrix[position[0]][position[1]] > 100 {
			return true
		} else {
			return false
		}
	}
*/

func (game *Game) HitboxCollidesWithObstacle(hitbox rl.Rectangle) bool {
	/*
		//WARNING: PositionIsValid MUST be run first
		positionTile := game.GameBoard.board_matrix[int(position.Y-GLOBAL_TILE_SIZE*0.2)/GLOBAL_TILE_SIZE][int(position.X)/GLOBAL_TILE_SIZE]
		if positionTile == Wall || positionTile == Breakable {
			return true
		} else {
			return false
		}*/
	for _, v := range game.GameBoard.Obstacles {
		if rl.CheckCollisionRecs(hitbox, v.HitBox) {
			return true
		}
	}
	return false
}

func (game *Game) MovePlayer(player *Player, direction string, deltatime float32) {
	nextPosition := game.GetNextPosition(player, direction, deltatime)
	nextHitbox := game.GetNextPositionHitbox(player, direction, deltatime)
	if game.PositionIsValid(nextPosition) && !game.HitboxCollidesWithObstacle(nextHitbox) && !game.HitboxCollidesWithOtherPlayer(player, nextHitbox) {
		fmt.Println("Updating position ", nextPosition)
		player.UpdatePosition(nextPosition)
		fmt.Println("New Player position: ", player.Position)
		fmt.Println("New Player hitbox: ", player.HitBox)

	}
}

// TODO Check for multiple bombs on the same tile
func (game *Game) PlaceBomb(sourcePlayer *Player, location rl.Vector2, radius int32) {
	if !game.IsBombPlacedHere(RoundXtoY(location.X, GLOBAL_TILE_SIZE), RoundXtoY(location.Y, GLOBAL_TILE_SIZE)) {

		bomb := NewBomb(sourcePlayer, RoundXtoY(location.X, GLOBAL_TILE_SIZE), RoundXtoY(location.Y, GLOBAL_TILE_SIZE), radius)
		sourcePlayer.AvailableBombs--
		game.Bombs = append(game.Bombs, bomb)
	}
}

func (game *Game) IsBombPlacedHere(position_x, position_y int32) bool {
	for _, v := range game.Bombs {
		if v.Position_x == position_x && v.Position_y == position_y {
			return true
		}
	}
	return false
}

func (game *Game) ExplodeBomb(bomb_index int) {
	game.GenerateShrapnel(&game.Bombs[bomb_index])
	game.Bombs[bomb_index].Owner.AvailableBombs++
	game.Bombs = append(game.Bombs[:bomb_index], game.Bombs[bomb_index+1:]...)
}

func (game *Game) GenerateShrapnel(sourceBomb *Bomb) {
	//TODO make it break stuff
	game.PlaceShrapnel(sourceBomb.Owner, sourceBomb.Position_x, sourceBomb.Position_y)
	var up_blocked, down_blocked, left_blocked, right_blocked bool
	for i := GLOBAL_TILE_SIZE; i <= int(sourceBomb.Radius)*GLOBAL_TILE_SIZE; i += GLOBAL_TILE_SIZE {
		nextpos_x_right := int(sourceBomb.Position_x) + i
		nextpos_x_left := int(sourceBomb.Position_x) - i

		nextpos_y_up := int(sourceBomb.Position_y) + i
		nextpos_y_down := int(sourceBomb.Position_y) - i

		up_blocked = nextpos_y_up > int(game.GameBoard.Size_y)*GLOBAL_TILE_SIZE
		down_blocked = nextpos_y_down < 0
		right_blocked = nextpos_x_right > int(game.GameBoard.Size_x)*GLOBAL_TILE_SIZE
		left_blocked = nextpos_x_left < 0

		//TODO block at obstacles

		if !up_blocked {
			game.PlaceShrapnel(sourceBomb.Owner, sourceBomb.Position_x, int32(nextpos_y_up))
		}
		if !down_blocked {
			game.PlaceShrapnel(sourceBomb.Owner, sourceBomb.Position_x, int32(nextpos_y_down))
		}
		if !left_blocked {
			game.PlaceShrapnel(sourceBomb.Owner, int32(nextpos_x_left), sourceBomb.Position_y)
		}
		if !right_blocked {
			game.PlaceShrapnel(sourceBomb.Owner, int32(nextpos_x_right), sourceBomb.Position_y)
		}

	}
}
func (game *Game) UpdateBombs() {
	for i := 0; i < len(game.Bombs); i++ {
		game.Bombs[i].RemainingTicks--
		if game.Bombs[i].RemainingTicks <= 0 {
			game.ExplodeBomb(i)
			i--
		}
	}
}
func (game *Game) RemoveShrapnel(shrapnel_index int) {
	game.Shrapnels[shrapnel_index].Owner.AvailableBombs++
	game.Shrapnels = append(game.Shrapnels[:shrapnel_index], game.Shrapnels[shrapnel_index+1:]...)
}
func (game *Game) UpdateShrapnels() {
	for i := 0; i < len(game.Shrapnels); i++ {
		game.Shrapnels[i].RemainingTicks--
		if game.Shrapnels[i].RemainingTicks <= 0 {
			game.KillPlayersUsingShrapnel(&game.Shrapnels[i])
			game.RemoveShrapnel(i)
			i--
		}
	}
}

func (game *Game) KillPlayersUsingShrapnel(sourceShrapnel *Shrapnel) {
	tempHitBox := rl.NewRectangle(float32(sourceShrapnel.Position_x), float32(sourceShrapnel.Position_y), GLOBAL_TILE_SIZE, GLOBAL_TILE_SIZE)
	for i := range game.Players {
		if rl.CheckCollisionRecs(tempHitBox, game.Players[i].HitBox) {
			game.PlayerDeath(&game.Players[i])
		}
	}
}

func (game *Game) PlayerDeath(player *Player) {
	player.Status = false
}

func (game *Game) PlaceShrapnel(owner *Player, position_x, position_y int32) {
	//TODO check if shrapnel is already placed here
	shrapnel := NewShrapnel(owner, position_x, position_y)
	game.Shrapnels = append(game.Shrapnels, shrapnel)
}
func (game *Game) IsShrapnelPlacedHere(position_x, position_y int32) bool {
	for _, v := range game.Shrapnels {
		if v.Position_x == position_x && v.Position_y == position_y {
			return true
		}
	}
	return false
}
func (game *Game) Update() {
	game.Ticks++
	game.UpdateBombs()
	game.UpdateShrapnels()

}

func (game *Game) GameShouldEnd() (bool, *Player) {
	players_left := 0
	winner := &game.Players[0]
	for i, v := range game.Players {
		if v.Status {
			players_left++
			winner = &game.Players[i]
		}
	}
	if players_left == 1 {
		return true, winner
	} else {
		return false, nil
	}
}

type Gfx struct {
	Size_x              int32
	Size_y              int32
	Tile_size           int32
	Game_Texture        rl.RenderTexture2D
	Game_Texture_Size_x int32
	Game_Texture_Size_y int32
}

const GLOBAL_TILE_SIZE = 20

func NewGfx(size_x, size_y int32) Gfx {
	gfx := Gfx{}
	gfx.Size_x = size_x
	gfx.Size_y = size_y
	gfx.Tile_size = GLOBAL_TILE_SIZE //= 30
	rl.InitWindow(size_x, size_y, "TuxMan")
	//defer rl.CloseWindow()
	rl.SetTargetFPS(30)
	return gfx
}
func (gfx *Gfx) DrawBoard(game *Game) {
	/*
		for i := range game.GameBoard.board_matrix {
			for j := range game.GameBoard.board_matrix[i] {
				if game.GameBoard.board_matrix[i][j] == Blank {
					rl.DrawRectangle(int32(j)*gfx.Tile_size, int32(i)*gfx.Tile_size, gfx.Tile_size, gfx.Tile_size, rl.Gray)
				} else if game.GameBoard.board_matrix[i][j] == Wall {
					rl.DrawRectangle(int32(j)*gfx.Tile_size, int32(i)*gfx.Tile_size, gfx.Tile_size, gfx.Tile_size, rl.Black)
				}
			}
		}*/
	/*
	  for i:=0;i<game.GameBoard.Size_y;i++{
	    for j:=0;j<game.G
	  }*/
	rl.DrawRectangle(0, 0, game.GameBoard.Size_x*gfx.Tile_size, game.GameBoard.Size_y*gfx.Tile_size, rl.Gray)

}

func (gfx *Gfx) DrawObstacles(game *Game) {
	for _, v := range game.GameBoard.Obstacles {
		if v.ObstacleType == Wall {
			rl.DrawRectangleRec(v.HitBox, rl.Black)
		} else if v.ObstacleType == Breakable {
			rl.DrawRectangleRec(v.HitBox, rl.Brown)
		}
	}
}
func (gfx *Gfx) DrawPlayers(game *Game) {
	for _, v := range game.Players {
		if v.Status {
			rl.DrawRectangleRec(v.HitBox, rl.Beige)
		}
		//rl.DrawRectangle(int32(v.Position.X)*gfx.Tile_size, int32(v.Position.Y)*gfx.Tile_size, gfx.Tile_size, gfx.Tile_size, rl.Red)
	}
}

func (gfx *Gfx) DrawBombs(game *Game) {
	for _, v := range game.Bombs {
		rl.DrawRectangle(v.Position_x, v.Position_y, GLOBAL_TILE_SIZE, GLOBAL_TILE_SIZE, rl.Red)
	}
}
func (gfx *Gfx) DrawShrapnel(game *Game) {
	for _, v := range game.Shrapnels {
		rl.DrawRectangle(v.Position_x, v.Position_y, GLOBAL_TILE_SIZE, GLOBAL_TILE_SIZE, rl.Orange)
	}
}
func (gfx *Gfx) GenerateGameTexture(game *Game) {
	//TODO check for texture init
	rl.BeginTextureMode(gfx.Game_Texture)
	rl.ClearBackground(rl.White)
	gfx.DrawBoard(game)
	gfx.DrawObstacles(game)
	gfx.DrawBombs(game)
	gfx.DrawPlayers(game)
	gfx.DrawShrapnel(game)
	rl.EndTextureMode()

}

func (gfx *Gfx) HandleInput(game *Game, deltatime float32) {
	var playerAlreadyChecked [4]bool
	playerAlreadyChecked = [4]bool{false, false, false, false}

	player1Key, player1KeyIsPressed := gfx.GetPlayer1Key()

	if gfx.IsPlayer1BombKeyPressed() && game.Players[0].AvailableBombs > 0 {
		game.PlaceBomb(&game.Players[0], game.Players[0].Position, 3)
	}

	if player1KeyIsPressed && !playerAlreadyChecked[0] {
		playerAlreadyChecked[0] = true
		game.MovePlayer(&game.Players[0], player1Key, deltatime)
		fmt.Println("Key pressed: ", player1Key)
	}
	if len(game.Players) > 1 {
		player2Key, player2KeyIsPressed := gfx.GetPlayer2Key()

		if gfx.IsPlayer2BombKeyPressed() && game.Players[1].AvailableBombs > 0 {
			game.PlaceBomb(&game.Players[1], game.Players[1].Position, 3)
		}

		if player2KeyIsPressed && !playerAlreadyChecked[1] {
			playerAlreadyChecked[1] = true
			game.MovePlayer(&game.Players[1], player2Key, deltatime)
			fmt.Println("Key pressed: ", player2Key)
		}

	}

}

// TODO Make placing bombs and movement more responsive

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

func (gfx *Gfx) InitGameTextureBox(game *Game) {
	gfx.Game_Texture_Size_x = gfx.Tile_size * game.GameBoard.Size_x
	gfx.Game_Texture_Size_y = gfx.Tile_size * game.GameBoard.Size_y
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
	gfx := NewGfx(600, 600)
	gfx.InitGameTextureBox(&game)
	game.AddPlayer("Pablo", 0, 1)
	game.AddPlayer("SecondPlayer", GLOBAL_TILE_SIZE*float32(game.GameBoard.Size_x-1), 1)

	game.GameBoard.AddObstacle(4, 1, Wall)
	game.GameBoard.AddObstacle(2, 2, Breakable)

	for !rl.WindowShouldClose() {
		gfx.HandleInput(&game, rl.GetFrameTime())
		game.Update()

		gfx.GenerateGameTexture(&game)

		gameOver, winner := game.GameShouldEnd()
		rl.BeginDrawing()
		if gameOver {
			GameOverScreen(&gfx, &game, winner)
		} else {
			GameScreen(&gfx, &game)
		}
		rl.EndDrawing()
	}
}
