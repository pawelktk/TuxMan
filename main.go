package main

import "fmt"

import rl "github.com/gen2brain/raylib-go/raylib"

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

func (board *Board) PutPlayers(players []Player) {
	for i := range players {
		board.board_matrix[players[i].Position_x][players[i].Position_y] = 2
	}
}

type Player struct {
	Name           string
	Position_x     int32
	Position_y     int32
	Points         int32
	AvailableBombs int32
	Status         bool
}

func NewPlayer(name string, position_x, position_y int32) Player {
	player := Player{}
	player.Name = name
	player.Position_x = position_x
	player.Position_y = position_y
	player.Points = 0
	player.AvailableBombs = 1
	player.Status = true
	return player
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
func (game *Game) AddPlayer(name string, position_x, position_y int32) {
	player := NewPlayer(name, position_x, position_y)
	game.Players = append(game.Players, player)
}

func (game *Game) GetNextPosition(player *Player, direction string) [2]int32 {
	var moveVector [2]int32
	switch direction {
	case "up":
		moveVector[0] = 1
	case "down":
		moveVector[0] = -1
	case "left":
		moveVector[1] = -1
	case "right":
		moveVector[1] = 1
	}
	moveVector[1] += player.Position_x
	moveVector[0] += player.Position_y
	return moveVector
}
func (game *Game) PositionIsValid(position [2]int32) bool {
	if position[0] <= game.GameBoard.Size_x && position[1] <= game.GameBoard.Size_y && position[0] >= 0 && position[1] >= 0 {
		return true
	} else {
		return false
	}
}

func (game *Game) PositionCollidesWithBomb(position [2]int32) bool {
	for _, bomb := range game.Bombs {
		if bomb.Position_x == position[0] && bomb.Position_y == position[1] {
			return true
		}
	}
	return false
}
func (game *Game) PositionCollidesWithPlayer(position [2]int32) bool {
	for _, player := range game.Players {
		if player.Position_x == position[0] && player.Position_y == position[1] {
			return true
		}
	}
	return false
}

func (game *Game) PositionCollidesWithObstacle(position [2]int32) bool {
	//WARNING: PositionIsValid MUST be run first
	fmt.Printf("position[0]: %v, position[1]: %v\n", position[0], position[1])
	if game.GameBoard.board_matrix[position[0]][position[1]] > 100 {
		return true
	} else {
		return false
	}
}

func (game *Game) MovePlayer(player *Player, direction string) {
	nextPosition := game.GetNextPosition(player, direction)
	if game.PositionIsValid(nextPosition) && !game.PositionCollidesWithBomb(nextPosition) && !game.PositionCollidesWithPlayer(nextPosition) && !game.PositionCollidesWithObstacle(nextPosition) {
		player.Position_y = nextPosition[0]
		player.Position_x = nextPosition[1]
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
		rl.DrawRectangle(int32(v.Position_x)*gfx.Tile_size, int32(v.Position_y)*gfx.Tile_size, gfx.Tile_size, gfx.Tile_size, rl.Red)
	}
}

func (gfx *Gfx) HandleInput(game *Game) {
	var playerAlreadyChecked [4]bool

	player1Key, player1KeyIsPressed := gfx.GetPlayer1Key()

	if player1KeyIsPressed && !playerAlreadyChecked[0] {
		playerAlreadyChecked[0] = true
		game.MovePlayer(&game.Players[0], player1Key)
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
		gfx.HandleInput(&game)
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		gfx.DrawBoard(&game)
		gfx.DrawPlayers(&game)
		rl.EndDrawing()
	}
}
