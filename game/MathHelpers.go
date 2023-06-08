package game

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
