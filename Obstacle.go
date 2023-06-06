package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

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