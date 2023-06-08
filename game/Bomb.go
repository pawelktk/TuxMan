package game

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
