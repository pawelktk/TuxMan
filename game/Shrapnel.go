package game

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
