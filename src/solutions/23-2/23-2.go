package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"time"
)

type Coordinates struct {
	x, y int
}

const ( // these are in a specific order for the specific pattern in the exercise
	North int = iota
	South
	West
	East
)
const NUM_DIRECTIONS = 4

const MOVE_DISTANCE = 1

func main() {
	input, err := readLines("inputs/23.txt")
	check(err)

	timer := time.Now()
	solve(input)
	fmt.Println("running time:", time.Since(timer))
}

func solve(input []string) {
	elves := parseInput(&input)

	elvesMoved := true
	count := 0

	for elvesMoved {
		elvesMoved = moveElves(&elves, count)
		count++
	}

	answer := count
	fmt.Printf("Answer: %d\n", answer)
}

func moveElves(elves *map[int]Coordinates, roundCounter int) bool {
	moves, haltedElves, gridwidth, gridOffset := getMoves(elves, roundCounter)

	if len(moves) == 0 {
		return false
	}

	for scanlinePosition, elfId := range moves {
		if containsInt(haltedElves, elfId) {
			continue
		}

		newY := (scanlinePosition / gridwidth) - gridOffset.y
		newX := (scanlinePosition % gridwidth) - gridOffset.x

		(*elves)[elfId] = Coordinates{newX, newY}
	}

	return true
}

func getMoves(elves *map[int]Coordinates, roundCounter int) (map[int]int, []int, int, Coordinates) {
	grid, gridOffset := makeElfGrid(elves, roundCounter)
	gridWidth := len(grid[0])

	moves := make(map[int]int)
	var haltedElves []int

	for elfId, position := range *elves {
		move := getElfMove(&grid, gridOffset, position, roundCounter)

		if move.x == position.x && move.y == position.y { // (no move)
			continue
		}

		scanlinePosition := (move.y+gridOffset.y)*gridWidth + move.x + gridOffset.x

		existingElfId, exists := moves[scanlinePosition]

		if exists {
			if !containsInt(haltedElves, existingElfId) { // check if another collision has already added this elfId to haltedElves
				haltedElves = append(haltedElves, existingElfId)
			}
		} else {
			moves[scanlinePosition] = elfId
		}
	}

	return moves, haltedElves, gridWidth, gridOffset
}

// logic assumes grid has a margin of 1 position around the outside of the area containing all the elves,
// and therefore the given "position" is not on the edge of any array in "grid"
func getElfMove(g *[][]bool, gridOffset Coordinates, position Coordinates, roundCounter int) Coordinates {
	grid := *g
	firstDirection := roundCounter % NUM_DIRECTIONS

	move := position

	gridPosition := Coordinates{x: position.x + gridOffset.x, y: position.y + gridOffset.y}

	for i := 0; i < 4; i++ {
		active := (firstDirection + i) % NUM_DIRECTIONS

		if active == North {
			northClear := grid[gridPosition.y-1][gridPosition.x-1] == false &&
				grid[gridPosition.y-1][gridPosition.x] == false &&
				grid[gridPosition.y-1][gridPosition.x+1] == false
			if hasAdjacentElf(g, gridOffset, position) && northClear {
				move.y--
				break
			}
		} else if active == South {
			southClear := grid[gridPosition.y+1][gridPosition.x-1] == false &&
				grid[gridPosition.y+1][gridPosition.x] == false &&
				grid[gridPosition.y+1][gridPosition.x+1] == false
			if hasAdjacentElf(g, gridOffset, position) && southClear {
				move.y++
				break
			}
		} else if active == West {
			westClear := grid[gridPosition.y+1][gridPosition.x-1] == false &&
				grid[gridPosition.y][gridPosition.x-1] == false &&
				grid[gridPosition.y-1][gridPosition.x-1] == false
			if hasAdjacentElf(g, gridOffset, position) && westClear {
				move.x--
				break
			}
		} else if active == East {
			eastClear := grid[gridPosition.y+1][gridPosition.x+1] == false &&
				grid[gridPosition.y][gridPosition.x+1] == false &&
				grid[gridPosition.y-1][gridPosition.x+1] == false
			if hasAdjacentElf(g, gridOffset, position) && eastClear {
				move.x++
				break
			}
		}
	}

	return move
}

func hasAdjacentElf(g *[][]bool, gridOffset Coordinates, position Coordinates) bool {
	grid := *g
	gridPosition := Coordinates{x: position.x + gridOffset.x, y: position.y + gridOffset.y}

	return grid[gridPosition.y-1][gridPosition.x] == true ||
		grid[gridPosition.y-1][gridPosition.x+1] == true ||
		grid[gridPosition.y][gridPosition.x+1] == true ||
		grid[gridPosition.y+1][gridPosition.x+1] == true ||
		grid[gridPosition.y+1][gridPosition.x] == true ||
		grid[gridPosition.y+1][gridPosition.x-1] == true ||
		grid[gridPosition.y][gridPosition.x-1] == true ||
		grid[gridPosition.y-1][gridPosition.x-1] == true
}

func makeElfGrid(elves *map[int]Coordinates, roundCounter int) ([][]bool, Coordinates) {
	nwBound, seBound := getElfBounds(elves)
	boundLatitude := seBound.y - nwBound.y + 1 + 2*MOVE_DISTANCE // margin added on all sides; enough for one round only
	boundLongitude := seBound.x - nwBound.x + 1 + 2*MOVE_DISTANCE

	grid := make([][]bool, boundLatitude)
	for i := range grid {
		grid[i] = make([]bool, boundLongitude)
	}

	gridOffset := Coordinates{x: MOVE_DISTANCE - nwBound.x, y: MOVE_DISTANCE - nwBound.y} // align to 0,0 plus a margin of MOVE_DISTANCE

	for _, position := range *elves {
		grid[position.y+gridOffset.y][position.x+gridOffset.x] = true
	}

	return grid, gridOffset
}

func drawGrid(grid *[][]bool) {
	for _, row := range *grid {
		for _, position := range row {
			if position == true {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}

		fmt.Print("\n")
	}
}

func getElfBounds(elves *map[int]Coordinates) (Coordinates, Coordinates) {
	northwest := Coordinates{math.MaxInt, math.MaxInt}
	southeast := Coordinates{0, 0}

	for _, elf := range *elves {
		if elf.x < northwest.x {
			northwest.x = elf.x
		}
		if elf.x > southeast.x {
			southeast.x = elf.x
		}
		if elf.y < northwest.y {
			northwest.y = elf.y
		}
		if elf.y > southeast.y {
			southeast.y = elf.y
		}
	}

	return northwest, southeast
}

func containsInt(list []int, element int) bool {
	for _, item := range list {
		if item == element {
			return true
		}
	}
	return false
}

func parseInput(input *[]string) map[int]Coordinates {
	elves := make(map[int]Coordinates)
	elfCount := 0

	for i, line := range *input {
		for j, char := range line {
			if char == '#' {
				elves[elfCount] = Coordinates{x: j, y: i}
				elfCount++
			}
		}
	}

	return elves
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
