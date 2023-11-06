package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// Notes:
// - remove bottom of gameboard every time a row is filled (nothing could get past), and keep track of number of removed rows
//   - reduces memory requirement from terabytes+ to ~40MB

// - find repeating patterns?
// - amount of moves and shapes multiplied together is lowest common denominator for a cycle(?)
// - check how the last shape in that cycle would fit with the first (10091 moves x 5 shapes = 50455 -> doesn't make sense to do this since doesn't consider number of moves per rock, but as a random guess, might discover repeatable pattern)
//   - start of pattern would sit on top of how it ends: height of pattern is 78749
//   - 1,000,000,000,000 / 50,455 = 19,819,641 (and some decimals * 50,455 = 13,345 rocks in remainder tower)
//   - 19819641 repeats * 78749 height + <13345 rock tower height -> 20777>
//     = 1,560,776,929,886 isn't right (probably because the joins between patterns don't have the entire row filled and rocks could get past)
// - there's a clear pattern in the indices where rows are completely filled and bottom of gameboard is removed
//   - stable pattern begins happening at tower height of 1698 and is every 4409 - 1698 = 2711 height after that
//   - pattern is 1735 rocks that have height of 2711, start is 1131 rocks that have height 1698
//   - (1,000,000,000,000 - 1,131) / 1,735 = 576,368,875 (and some decimals * 1,735 = 744 rocks in remainder tower)
//   - 1,698 non-repeating tower height + 576,368,875 repeats * 2,711 height + <744 rock partial pattern height -> 2841 - 1698 = 1143>
//     = 1,562,536,022,966

var BOARD_WIDTH = 7

// var NUM_ROCKS = 1000000000000
var NUM_ROCKS = 10000

var SHAPES = [][][]byte{
	{{'#', '#', '#', '#'}},
	{
		{'.', '#', '.'},
		{'#', '#', '#'},
		{'.', '#', '.'}},
	{
		{'#', '#', '#'}, // shapes are vertically reversed, but this is the only one that's different
		{'.', '.', '#'},
		{'.', '.', '#'}},
	{
		{'#'},
		{'#'},
		{'#'},
		{'#'}},
	{
		{'#', '#'},
		{'#', '#'}},
}

type Point struct{ x, y int }

func main() {
	input, err := readLines("inputs/17.txt")
	check(err)

	timer := time.Now()
	solve(input)
	fmt.Println("running time:", time.Since(timer))
}

func solve(input []string) {
	movements := parseInput(input[0])
	movement_count := 0

	var gameboard [][]byte
	removed_row_count := 0
	first_clear_row := 0

	for i := 0; i < NUM_ROCKS; i++ {
		// if i%1000000000 == 0 {
		// 	fmt.Printf("progress: %d/1000\n", i/1000000000)
		// }

		shape_index := i % len(SHAPES)
		shape := SHAPES[shape_index]
		start_position := Point{x: 2, y: first_clear_row + 3}

		full_row_index, exists := dropRock(&gameboard, shape, start_position, movements, &movement_count)

		if exists {
			gameboard = gameboard[full_row_index+1:]
			removed_row_count += full_row_index + 1
			// fmt.Printf("Total removed: %d  %d  %d\n", removed_row_count, full_row_index+1, i) // repeating pattern found in this
		}

		first_clear_row = getFirstClearRow(&gameboard)
	}

	// drawGameboard(&gameboard)

	fmt.Printf("Tower height: %d\n", first_clear_row+removed_row_count)
}

func getFirstClearRow(gameboard *[][]byte) int {
	for i := len(*gameboard) - 1; i >= 0; i-- {
		row_occupied := false

		for _, position := range (*gameboard)[i] {
			if position == '#' {
				row_occupied = true
				break
			}
		}

		if row_occupied {
			return i + 1
		}
	}

	return 0
}

func dropRock(gameboard *[][]byte, shape [][]byte, start Point, movements []bool, movement_count *int) (int, bool) {
	shape_height := len(shape)
	for len(*gameboard) < start.y+shape_height {
		*gameboard = append(*gameboard, make([]byte, BOARD_WIDTH))
	}

	position := start

	in_motion := true
	for in_motion {
		movement_index := *movement_count % len(movements)
		moving_right := movements[movement_index]

		if moving_right && !detectCollision(gameboard, shape, Point{x: position.x + 1, y: position.y}) {
			position.x++
		} else if !moving_right && !detectCollision(gameboard, shape, Point{x: position.x - 1, y: position.y}) {
			position.x--
		}

		lower_position := Point{x: position.x, y: position.y - 1}
		if detectCollision(gameboard, shape, lower_position) {
			in_motion = false
			addRock(gameboard, shape, position)
		} else {
			position = lower_position
		}

		*movement_count++
	}

	// check if any of the rows occupied by the newly resting rock are completely filled left to right
	for i := position.y; i < position.y+shape_height; i++ {
		row_filled := true
		for _, symbol := range (*gameboard)[i] {
			if symbol != '#' {
				row_filled = false
				break
			}
		}

		if row_filled {
			return i, true // return the index of the filled row in the gameboard
		}
	}

	return 0, false
}

func addRock(gameboard *[][]byte, shape [][]byte, position Point) {
	for i, row := range shape {
		for j, shape_symbol := range row {

			if shape_symbol == '#' {
				(*gameboard)[position.y+i][position.x+j] = '#'
			}
		}
	}
}

func detectCollision(gameboard *[][]byte, shape [][]byte, position Point) bool {
	shape_width := len(shape[0])

	if position.x < 0 || position.x+shape_width > BOARD_WIDTH { // collides with left or right edge of gameboard
		return true
	}
	if position.y < 0 { // collides with bottom of gameboard
		return true
	}

	for i, row := range shape {
		for j, shape_symbol := range row {
			gameboard_symbol := (*gameboard)[position.y+i][position.x+j]

			if shape_symbol == '#' && gameboard_symbol == '#' {
				return true
			}
		}
	}

	return false
}

func drawGameboard(gameboard *[][]byte) {
	for i := len(*gameboard) - 1; i >= 0; i-- {
		row := (*gameboard)[i]

		for _, symbol := range row {
			if symbol == '#' {
				fmt.Printf("%c", symbol)
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
}

func parseInput(input string) []bool {
	movements := make([]bool, len(input))

	for i, char := range input {
		movements[i] = char == '>'
	}

	return movements
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
