package main

type Vector2int32 struct {
	X int32
	Y int32
}

func NewVector2int32(x, y int32) Vector2int32 {
	vector := Vector2int32{}
	vector.X = x
	vector.Y = y
	return vector
}