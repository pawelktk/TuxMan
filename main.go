package main

import (
	"fmt"
	//	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Board struct {
	Size_x       int32
	Size_y       int32
	board_matrix [][]int32
}

func NewBoard(size_x, size_y int32) Board {
	board := Board{}
	board.Size_x = size_x
	board.Size_y = size_y
	board.board_matrix = make([][]int32, size_y)
	for i := int32(0); i < size_y; i++ {
		board.board_matrix[i] = make([]int32, size_x)
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

type Player struct {
	Name           string
	Position       rl.Vector2
	Points         int32
	AvailableBombs int32
	Status         bool
	HitBox         rl.Rectangle
	Speed          float32
}

func NewPlayer(name string, position_x, position_y float32) Player {
	player := Player{}
	player.Name = name
	player.Position.X = position_x
	player.Position.Y = position_y
	player.Points = 0
	player.AvailableBombs = 1
	player.Status = true
	player.HitBox = rl.NewRectangle(float32(position_x), float32(position_y), 10, 10)
	player.Speed = 20
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
	bomb.RemainingTicks = 10
	bomb.Radius = 2

	return bomb
}

type Game struct {
	Ticks     int32
	GameBoard Board
	Players   []Player
	Bombs     []Bomb
}

func NewGame() Game {
	game := Game{}
	game.Ticks = 0
	game.GameBoard = NewBoard(20, 10)
	return game
}
func (game *Game) AddPlayer(name string, position_x, position_y float32) {
	player := NewPlayer(name, position_x, position_y)
	game.Players = append(game.Players, player)
}

func (game *Game) GetNextPosition(player *Player, direction string, deltatime float32) rl.Vector2 {
	var nextPosition rl.Vector2
	nextPosition = player.Position
	switch direction {
	case "up":
		nextPosition.Y -= player.Speed * deltatime
	case "down":
		nextPosition.Y += player.Speed * deltatime
	case "left":
		nextPosition.X -= player.Speed * deltatime
	case "right":
		nextPosition.X += player.Speed * deltatime
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
	if position.X <= float32(game.GameBoard.Size_x)*20 && position.Y <= float32(game.GameBoard.Size_y)*20 && position.X >= 0 && position.Y >= 0 {
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
func (game *Game) MovePlayer(player *Player, direction string, deltatime float32) {
	nextPosition := game.GetNextPosition(player, direction, deltatime)
	nextHitbox := game.GetNextPositionHitbox(player, direction, deltatime)
	if game.PositionIsValid(nextPosition) && !game.HitboxCollidesWithOtherPlayer(player, nextHitbox) {
		fmt.Println("Updating position ", nextPosition)
		player.UpdatePosition(nextPosition)
		fmt.Println("New Player position: ", player.Position)
		fmt.Println("New Player hitbox: ", player.HitBox)

	}
}

func (game *Game) Update() {
	game.Ticks++

}

type Gfx struct {
	Size_x    int32
	Size_y    int32
	Tile_size int32
}

func NewGfx(size_x, size_y int32) Gfx {
	gfx := Gfx{}
	gfx.Size_x = size_x
	gfx.Size_y = size_y
	gfx.Tile_size = 20
	rl.InitWindow(size_x, size_y, "TuxMan")
	//defer rl.CloseWindow()
	rl.SetTargetFPS(30)
	return gfx
}
func (gfx *Gfx) DrawBoard(game *Game) {
	for i := range game.GameBoard.board_matrix {
		for j := range game.GameBoard.board_matrix[i] {
			if game.GameBoard.board_matrix[i][j] == 0 {
				rl.DrawRectangle(int32(j)*gfx.Tile_size, int32(i)*gfx.Tile_size, gfx.Tile_size, gfx.Tile_size, rl.Gray)
			}
		}
	}
}

func (gfx *Gfx) DrawPlayers(game *Game) {
	for _, v := range game.Players {
		rl.DrawRectangleRec(v.HitBox, rl.Beige)
		//rl.DrawRectangle(int32(v.Position.X)*gfx.Tile_size, int32(v.Position.Y)*gfx.Tile_size, gfx.Tile_size, gfx.Tile_size, rl.Red)
	}
}

func (gfx *Gfx) HandleInput(game *Game, deltatime float32) {
	var playerAlreadyChecked [4]bool
	playerAlreadyChecked = [4]bool{false, false, false, false}

	player1Key, player1KeyIsPressed := gfx.GetPlayer1Key()

	if player1KeyIsPressed && !playerAlreadyChecked[0] {
		playerAlreadyChecked[0] = true
		game.MovePlayer(&game.Players[0], player1Key, deltatime)
		fmt.Println("Key pressed: ", player1Key)
	}
}

func (gfx *Gfx) GetPlayer1Key() (string, bool) {
	switch {
	case rl.IsKeyDown(rl.KeyLeft):
		return "left", true
	case rl.IsKeyDown(rl.KeyRight):
		return "right", true
	case rl.IsKeyDown(rl.KeyUp):
		return "up", true
	case rl.IsKeyDown(rl.KeyDown):
		return "down", true
	default:
		return "", false
	}
}

func main() {

	game := NewGame()
	gfx := NewGfx(800, 400)
	game.AddPlayer("aaa", 0, 1)
	for !rl.WindowShouldClose() {
		gfx.HandleInput(&game, rl.GetFrameTime())
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		gfx.DrawBoard(&game)
		gfx.DrawPlayers(&game)
		rl.EndDrawing()
	}
}
