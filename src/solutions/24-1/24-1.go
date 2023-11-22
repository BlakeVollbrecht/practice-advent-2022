package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"time"
)

// Notes:
// - the blizzards could be modelled as 1-dimensional columns and rows
//   - i.e. could have a left-moving and right-moving row as bits and AND them together, but it's too long for bit flags in registers (20x150)
// - edge case where expedition might wait before moving onto the board
// - only need to know what's in 5 positions at any point - can get that from 3 rows, 3 columns, starting positions of blizzards, and current time
//   - can work it backward from the final position to find whether there's a blizzard in the starting position (instead of calculating position of everything in the 3 rows and columns)
//   - can cache which positions are clear or not for a given time (so this calculation doesn't happen in the hot loop)
//		 - it's basically finding paths through a constant 3D array by adding the dimension of time
//     - and subtrees of possibility are consistent from each point in the 3D, so shortest path found to a point can be cached

type Coordinates struct {
	x, y int
}

type Key struct {
	x, y, time int
}

const (
	Clear int = iota
	North
	South
	East
	West
)

const TIME_LIMIT = 327

func main() {
	input, err := readLines("inputs/24.txt")
	check(err)

	timer := time.Now()
	solve(input)
	fmt.Println("running time:", time.Since(timer))
}

func solve(input []string) {
	blizzardStart, startX, endX := parseInput(&input)

	blizzardFrames := getBlizzardFrames(&blizzardStart, TIME_LIMIT)

	// for _, frame := range blizzardFrames {
	// 	drawBlizzards(&frame)
	// 	fmt.Println()
	// }

	expeditionPosition := Coordinates{x: startX, y: -1}
	subtreeCache := make(map[Key]int)

	shortestPath, success := getShortestPath(&blizzardFrames, 1, &subtreeCache, expeditionPosition, startX, endX)

	if !success {
		fmt.Println("No path found")
		return
	}

	fmt.Printf("Shortest path: %d\n", shortestPath)
}

func getShortestPath(blizzardFrames *[][][]bool, time int, subtreeCache *map[Key]int, position Coordinates, startX int, endX int) (int, bool) {
	possibleMoves := getPossibleMoves(blizzardFrames, time, position, startX)

	if len(possibleMoves) == 0 { // dead end
		return 0, false
	}

	if time >= TIME_LIMIT {
		return 0, false
	}

	if position.y == len((*blizzardFrames)[0])-1 && position.x == endX { // is in position above the exit
		return 1, true
	}

	subtreeKey := Key{x: position.x, y: position.y, time: time}
	shortestPath, exists := (*subtreeCache)[subtreeKey]
	if exists {
		return shortestPath, shortestPath > 0
	}

	minPath := math.MaxInt
	someSuccess := false

	for _, move := range possibleMoves {
		shortestPath, success := getShortestPath(blizzardFrames, time+1, subtreeCache, move, startX, endX)
		if success && shortestPath < minPath {
			minPath = shortestPath + 1
			someSuccess = true
		}
	}

	cachePathLength := 0
	if someSuccess {
		cachePathLength = minPath
	}
	(*subtreeCache)[subtreeKey] = cachePathLength

	return minPath, someSuccess
}

// checks current and south positions, and if not at start position, also checks north, east, and west positions
func getPossibleMoves(blizzardFrames *[][][]bool, time int, position Coordinates, startX int) []Coordinates {
	nextFrame := (*blizzardFrames)[time]
	var moves []Coordinates

	if position.y == -1 && position.x == startX || // is in starting alcove
		!nextFrame[position.y][position.x] {
		moves = append(moves, position) // stay in current position
	}

	south := Coordinates{x: position.x, y: position.y + 1}
	if position.y < len(nextFrame)-1 && // off bottom edge
		!nextFrame[south.y][south.x] {
		moves = append(moves, south)
	}

	if position.y != -1 || position.x != startX { // if not at start position, also check north, east, west
		north := Coordinates{x: position.x, y: position.y - 1}
		if position.y > 0 && // off top edge
			!nextFrame[north.y][north.x] {
			moves = append(moves, north)
		}

		east := Coordinates{x: position.x + 1, y: position.y}
		if position.x < len(nextFrame[0])-1 && // off right edge
			!nextFrame[east.y][east.x] {
			moves = append(moves, east)
		}

		west := Coordinates{x: position.x - 1, y: position.y}
		if position.x > 0 && // off left edge
			!nextFrame[west.y][west.x] {
			moves = append(moves, west)
		}
	}

	return moves
}

// check whether blizzards from 4 directions exist at the point they would need to be in the original map to be at the given position at the given time
func isPositionClear(blizzardStart *[][]int, time int, position Coordinates) bool {
	height := len(*blizzardStart)
	width := len((*blizzardStart)[0])

	isClear := true

	northBlizzardStart := (position.y + time) % height
	if (*blizzardStart)[northBlizzardStart][position.x] == North { // there is a north-moving blizzard in the starting position where it would need to be to now be at this time & position
		isClear = false
	}

	southBlizzardStart := (position.y - time) % height
	if southBlizzardStart < 0 {
		southBlizzardStart += height
	}
	if (*blizzardStart)[southBlizzardStart][position.x] == South {
		isClear = false
	}

	eastBlizzardStart := (position.x - time) % width
	if eastBlizzardStart < 0 {
		eastBlizzardStart += width
	}
	if (*blizzardStart)[position.y][eastBlizzardStart] == East {
		isClear = false
	}

	westBlizzardStart := (position.x + time) % width
	if (*blizzardStart)[position.y][westBlizzardStart] == West {
		isClear = false
	}

	return isClear
}

func getBlizzardFrames(blizzardStart *[][]int, maxTime int) [][][]bool {
	blizzardFrames := make([][][]bool, maxTime+1)

	for t := range blizzardFrames {
		blizzardFrames[t] = getBlizzardPositions(blizzardStart, t)
	}

	return blizzardFrames
}

func getBlizzardPositions(blizzardStart *[][]int, time int) [][]bool {
	blizzardPositions := make([][]bool, len(*blizzardStart))

	for i := range blizzardPositions {
		blizzardPositions[i] = make([]bool, len((*blizzardStart)[0]))

		for j := range blizzardPositions[i] {
			blizzardPositions[i][j] = !isPositionClear(blizzardStart, time, Coordinates{x: j, y: i})
		}
	}

	return blizzardPositions
}

func drawBlizzards(blizzards *[][]bool) {
	for _, row := range *blizzards {
		for _, presence := range row {
			if presence {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func parseInput(input *[]string) ([][]int, int, int) {
	blizzards := make([][]int, len(*input)-2) // (-2) for skipping first and last row

	for i := 0; i < len(blizzards); i++ {
		row := (*input)[i+1][1 : len((*input)[i+1])-1] // remove the first and last column, (i+1) skips first row

		blizzards[i] = make([]int, len(row))

		for j := 0; j < len(row); j++ {
			switch row[j] {
			case '>':
				blizzards[i][j] = East
			case 'v':
				blizzards[i][j] = South
			case '<':
				blizzards[i][j] = West
			case '^':
				blizzards[i][j] = North
			case '.':
				blizzards[i][j] = Clear
			}
		}
	}

	startX := bytes.IndexByte([]byte((*input)[0]), '.') - 1
	endX := bytes.IndexByte([]byte((*input)[len(*input)-1]), '.') - 1

	return blizzards, startX, endX
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
