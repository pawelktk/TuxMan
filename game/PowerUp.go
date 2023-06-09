package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/pawelktk/TuxMan/globals"
)

type PowerUpType int32

const (
	Speed PowerUpType = iota
)

type PowerUp struct {
	Position_x int32
	Position_y int32
	BoostType  PowerUpType
	HitBox     rl.Rectangle
}

func NewPowerUp(position_x int32, position_y int32, boostType PowerUpType) PowerUp {
	powerup := PowerUp{}
	powerup.Position_x = position_x
	powerup.Position_y = position_y
	powerup.BoostType = boostType
	powerup.HitBox = rl.NewRectangle(float32(position_x), float32(position_y), globals.GLOBAL_TILE_SIZE, globals.GLOBAL_TILE_SIZE)

	return powerup
}
