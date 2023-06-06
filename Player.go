package main
import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

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
	player.Speed = 70
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
