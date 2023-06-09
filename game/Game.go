package game

import (
	"fmt"
	"math/rand"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/pawelktk/TuxMan/globals"
)

type Game struct {
	Ticks     int32
	GameBoard Board
	Players   []Player
	Bombs     []Bomb
	Shrapnels []Shrapnel
	PowerUps  []PowerUp
}

func NewGame() Game {
	game := Game{}
	game.Ticks = 0
	game.GameBoard = NewRandomBoard(7, 5)
	//game.GameBoard.LoadFromFile("tuxman16.map")
	//game.GameBoard.GenerateRandomMap(5, 5)

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
	movementVector := player.Speed*deltatime + (player.Speed * globals.POWERUP_SPEED_BOOST * float32(player.SpeedBoost) * deltatime)
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
	if position.X <= float32(game.GameBoard.Size_x-1)*globals.GLOBAL_TILE_SIZE+0.2*globals.GLOBAL_TILE_SIZE && position.Y <= float32(game.GameBoard.Size_y-1)*globals.GLOBAL_TILE_SIZE+0.2*globals.GLOBAL_TILE_SIZE && position.X >= 0 && position.Y >= 0 {
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
	for _, v := range game.GameBoard.Obstacles {
		if rl.CheckCollisionRecs(hitbox, v.HitBox) {
			return true
		}
	}
	return false
}
func (game *Game) HandleCollisionWithWithPowerUp(sourcePlayer *Player, hitbox rl.Rectangle) {
	for i, v := range game.PowerUps {
		if rl.CheckCollisionRecs(hitbox, game.PowerUps[i].HitBox) {
			fmt.Println("[[Player collided with PowerUp!]]")
			if v.BoostType == Speed {

				sourcePlayer.SpeedBoost++
			}
			game.RemovePowerUp(i)
			return
		}
	}
}

func (game *Game) MovePlayer(player *Player, direction string, deltatime float32) {
	nextPosition := game.GetNextPosition(player, direction, deltatime)
	nextHitbox := game.GetNextPositionHitbox(player, direction, deltatime)
	if game.PositionIsValid(nextPosition) && !game.HitboxCollidesWithObstacle(nextHitbox) && !game.HitboxCollidesWithOtherPlayer(player, nextHitbox) {
		fmt.Println("Updating position ", nextPosition)
		player.UpdatePosition(nextPosition)
		fmt.Println("New Player position: ", player.Position)
		fmt.Println("New Player hitbox: ", player.HitBox)
		game.HandleCollisionWithWithPowerUp(player, nextHitbox)

	}
}

// DONE Check for multiple bombs on the same tile
func (game *Game) PlaceBomb(sourcePlayer *Player, location rl.Vector2, radius int32) {
	if !game.IsBombPlacedHere(RoundXtoY(location.X, globals.GLOBAL_TILE_SIZE), RoundXtoY(location.Y, globals.GLOBAL_TILE_SIZE)) {

		bomb := NewBomb(sourcePlayer, RoundXtoY(location.X, globals.GLOBAL_TILE_SIZE), RoundXtoY(location.Y, globals.GLOBAL_TILE_SIZE), radius)
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

func (game *Game) PlacePowerUp(position_x, position_y int32) {
	powerup := NewPowerUp(position_x, position_y, Speed)
	game.PowerUps = append(game.PowerUps, powerup)
}
func (game *Game) PlacePowerUpRandom(position_x, position_y int32) {
	if rand.Float32() < globals.POWERUP_CHANCE {
		game.PlacePowerUp(position_x, position_y)
	}
}

func (game *Game) RemovePowerUp(pu_index int) {
	game.PowerUps = append(game.PowerUps[:pu_index], game.PowerUps[pu_index+1:]...)
}

func (game *Game) GenerateShrapnel(sourceBomb *Bomb) {
	//DONE make it break stuff
	game.PlaceShrapnel(sourceBomb.Owner, sourceBomb.Position_x, sourceBomb.Position_y)
	var up_blocked, down_blocked, left_blocked, right_blocked bool
	for i := globals.GLOBAL_TILE_SIZE; i <= int(sourceBomb.Radius)*globals.GLOBAL_TILE_SIZE; i += globals.GLOBAL_TILE_SIZE {

		fmt.Println("#Shrapnel spread round ", i/globals.GLOBAL_TILE_SIZE)

		nextpos_x_right := int(sourceBomb.Position_x) + i
		nextpos_x_left := int(sourceBomb.Position_x) - i

		nextpos_y_up := int(sourceBomb.Position_y) + i
		nextpos_y_down := int(sourceBomb.Position_y) - i

		up_blocked = up_blocked || (nextpos_y_up > int(game.GameBoard.Size_y)*globals.GLOBAL_TILE_SIZE)
		down_blocked = down_blocked || (nextpos_y_down < 0)
		right_blocked = right_blocked || (nextpos_x_right > int(game.GameBoard.Size_x)*globals.GLOBAL_TILE_SIZE)
		left_blocked = left_blocked || (nextpos_x_left < 0)

		fmt.Printf("Initial check for map boundaries: \n\tup_blocked: %v\n\tdown_blocked: %v\n\tright_blocked: %v\n\tleft_blocked: %v\n", up_blocked, down_blocked, right_blocked, left_blocked)
		//DONE block at obstacles

		if !up_blocked {
			vecPosition := NewVector2int32(sourceBomb.Position_x, int32(nextpos_y_up))
			obstacleType, _ := game.GameBoard.GetObstacleType(vecPosition)
			//fmt.Printf("->vecPosition: %v\n->obstacleType: %v\n", vecPosition, obstacleType)
			if obstacleType == Wall {
				up_blocked = true
			} else {
				game.PlaceShrapnelDestructive(sourceBomb.Owner, sourceBomb.Position_x, int32(nextpos_y_up))
				game.PlacePowerUpRandom(vecPosition.X, vecPosition.Y)
			}
		}
		if !down_blocked {
			vecPosition := NewVector2int32(sourceBomb.Position_x, int32(nextpos_y_down))
			obstacleType, _ := game.GameBoard.GetObstacleType(vecPosition)
			if obstacleType == Wall {
				down_blocked = true
			} else {
				game.PlaceShrapnelDestructive(sourceBomb.Owner, sourceBomb.Position_x, int32(nextpos_y_down))
				game.PlacePowerUpRandom(vecPosition.X, vecPosition.Y)

			}
		}
		if !left_blocked {
			vecPosition := NewVector2int32(int32(nextpos_x_left), sourceBomb.Position_y)
			obstacleType, _ := game.GameBoard.GetObstacleType(vecPosition)
			if obstacleType == Wall {
				left_blocked = true
			} else {
				game.PlaceShrapnelDestructive(sourceBomb.Owner, int32(nextpos_x_left), sourceBomb.Position_y)
				game.PlacePowerUpRandom(vecPosition.X, vecPosition.Y)
			}
		}
		if !right_blocked {

			vecPosition := NewVector2int32(int32(nextpos_x_right), sourceBomb.Position_y)
			obstacleType, _ := game.GameBoard.GetObstacleType(vecPosition)
			if obstacleType == Wall {
				right_blocked = true
			} else {
				game.PlaceShrapnelDestructive(sourceBomb.Owner, int32(nextpos_x_right), sourceBomb.Position_y)
				game.PlacePowerUpRandom(vecPosition.X, vecPosition.Y)
			}
		}
		fmt.Printf("Check for obstacles: \n\tup_blocked: %v\n\tdown_blocked: %v\n\tright_blocked: %v\n\tleft_blocked: %v\n", up_blocked, down_blocked, right_blocked, left_blocked)

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
	tempHitBox := rl.NewRectangle(float32(sourceShrapnel.Position_x), float32(sourceShrapnel.Position_y), globals.GLOBAL_TILE_SIZE, globals.GLOBAL_TILE_SIZE)
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

func (game *Game) PlaceShrapnelDestructive(owner *Player, position_x, position_y int32) {
	game.PlaceShrapnel(owner, position_x, position_y)
	vecPosition := NewVector2int32(position_x, position_y)
	game.GameBoard.RemoveObstacleIfBreakable(vecPosition)
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
	var winner *Player //&game.Players[0]
	winner = nil
	for i, v := range game.Players {
		if v.Status {
			players_left++
			winner = &game.Players[i]
		}
	}
	if players_left <= 1 && len(game.Players) > 0 {
		if players_left == 0 {
			winner = &game.Players[0]
		}
		return true, winner

	} else {
		return false, nil
	}
}
