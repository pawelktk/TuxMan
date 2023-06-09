package game

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/pawelktk/TuxMan/globals"
)

type TileType int32

const (
	Blank TileType = iota
	Wall
	Breakable
)

type Board struct {
	Size_x       int32
	Size_y       int32
	board_matrix [][]TileType
	Obstacles    map[Vector2int32]Obstacle
}

func NewBoard(size_x, size_y int32) Board {
	board := Board{}
	board.Size_x = size_x
	board.Size_y = size_y
	board.Obstacles = make(map[Vector2int32]Obstacle)
	board.board_matrix = make([][]TileType, size_y)
	for i := int32(0); i < size_y; i++ {
		board.board_matrix[i] = make([]TileType, size_x)
	}
	board.Clear()
	return board
}
func NewRandomBoard(size_x3, size_y3 int32) Board {
	board := NewBoard(size_x3*3, size_y3*3)
	board.GenerateRandomMap(size_x3, size_y3)
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

func GetInt32FromScanner(scanner *bufio.Scanner) int32 {
	scanner.Scan()
	scannedInt, err := strconv.Atoi(scanner.Text())
	if err != nil {
		log.Fatal(err)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return int32(scannedInt)

}

func (board *Board) LoadFromFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	board.Size_x = GetInt32FromScanner(scanner)
	board.Size_y = GetInt32FromScanner(scanner)

	// x :=0
	y := 0

	for scanner.Scan() {
		line := scanner.Text()
		for x := 0; x < int(board.Size_x); x++ {
			tileType := TileType(line[x] - '0')
			if tileType != Blank {
				board.AddObstacle(int32(x), int32(y), TileType(tileType))
			}
		}
		y++
	}
	if err != nil {
		log.Fatal(err)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

/*
	func (board *Board) GenerateRandom15x15Map() {
		board.Size_x = 15
		board.Size_y = 15
		generatorSource := rand.NewSource(time.Now().UnixNano())
		generator := rand.New(generatorSource)
		var wg sync.WaitGroup
		mapMutex := sync.RWMutex{}
		for i0 := 0; i0 < 15-2; i0 += 3 {
			for j0 := 0; j0 < 15-2; j0 += 3 {
				wg.Add(1)
				go func(i, j int) {
					defer wg.Done()
					for k := 0; k < 2; k++ {
						x := i + generator.Intn(3)
						y := j + generator.Intn(3)
						mapMutex.Lock()
						obstacleAlreadyPlaced := board.ObstacleExist(NewVector2int32(int32(x), int32(y)))
						mapMutex.Unlock()

						if !obstacleAlreadyPlaced && !((x < 3 && y < 3) || (x < 3 && y > 15-4) || (x > 15-4 && y < 3) || (x > 15-4 && y > 15-4)) {
							mapMutex.Lock()
							board.AddObstacle(int32(x), int32(y), Wall)
							mapMutex.Unlock()
						}
					}
					for k := 0; k < 5; k++ {
						x := i + generator.Intn(3)
						y := j + generator.Intn(3)
						mapMutex.Lock()
						obstacleAlreadyPlaced := board.ObstacleExist(NewVector2int32(int32(x), int32(y)))
						mapMutex.Unlock()

						if !obstacleAlreadyPlaced && !((x < 3 && y < 3) || (x < 3 && y > 15-4) || (x > 15-4 && y < 3) || (x > 15-4 && y > 15-4)) {
							mapMutex.Lock()
							board.AddObstacle(int32(x), int32(y), Breakable)
							mapMutex.Unlock()
						}
					}
				}(i0, j0)
			}
		}
		wg.Wait()

}
*/
func (board *Board) GenerateRandomMap(size_x3, size_y3 int32) {
	size_x := size_x3 * 3
	size_y := size_y3 * 3
	board.Size_x = size_x
	board.Size_y = size_y
	generatorSource := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(generatorSource)
	var wg sync.WaitGroup
	mapMutex := sync.RWMutex{}
	for i0 := 0; i0 < int(size_x)-2; i0 += 3 {
		for j0 := 0; j0 < int(size_y)-2; j0 += 3 {
			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				for k := 0; k < 2; k++ {
					x := i + generator.Intn(3)
					y := j + generator.Intn(3)
					mapMutex.Lock()
					obstacleAlreadyPlaced := board.ObstacleExist(NewVector2int32(int32(x), int32(y)))
					mapMutex.Unlock()

					if !obstacleAlreadyPlaced && !((x < 3 && y < 3) || (x < 3 && y > int(size_y)-4) || (x > int(size_x)-4 && y < 3) || (x > int(size_x)-4 && y > int(size_y)-4)) {
						mapMutex.Lock()
						board.AddObstacle(int32(x), int32(y), Wall)
						mapMutex.Unlock()
					}
				}
				for k := 0; k < 5; k++ {
					x := i + generator.Intn(3)
					y := j + generator.Intn(3)
					mapMutex.Lock()
					obstacleAlreadyPlaced := board.ObstacleExist(NewVector2int32(int32(x), int32(y)))
					mapMutex.Unlock()

					if !obstacleAlreadyPlaced && !((x < 3 && y < 3) || (x < 3 && y > int(size_y)-4) || (x > int(size_x)-4 && y < 3) || (x > int(size_x)-4 && y > int(size_y)-4)) {
						mapMutex.Lock()
						board.AddObstacle(int32(x), int32(y), Breakable)
						mapMutex.Unlock()
					}
				}
			}(i0, j0)
		}
	}
	wg.Wait()

}
func (board *Board) AddObstacle(position_x, position_y int32, tileType TileType) {
	obstacle := NewObstacle(position_x, position_y, tileType)
	board.Obstacles[NewVector2int32(position_x*globals.GLOBAL_TILE_SIZE, position_y*globals.GLOBAL_TILE_SIZE)] = obstacle
	//board.Obstacles = append(board.Obstacles, obstacle)
}
func (board *Board) RemoveObstacle(position Vector2int32) {
	//board.Obstacles = append(board.Obstacles[:obstacle_index], board.Obstacles[obstacle_index+1:]...)
	fmt.Println("Removing obstacle at", position)
	delete(board.Obstacles, position)
}

func (board *Board) RemoveObstacleIfBreakable(position Vector2int32) {
	fmt.Println("Considering removal of obstacle at", position)
	obstacleType, exists := board.GetObstacleType(position)
	if exists && obstacleType == Breakable {
		board.RemoveObstacle(position)
	}
}

func (board *Board) GetObstacleType(position Vector2int32) (TileType, bool) {
	val, ok := board.Obstacles[position]
	if ok {
		return val.ObstacleType, true
	} else {
		return Blank, false

	}
}
func (board *Board) ObstacleExist(position Vector2int32) bool {
	_, ok := board.Obstacles[position]
	if !ok {
		return false
	} else {
		return true
	}
}
