package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

var BOARD_WIDTH = 7

var NUM_ROCKS = 2022

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
	first_clear_row := 0

	for i := 0; i < NUM_ROCKS; i++ {
		shape_index := i % len(SHAPES)
		shape := SHAPES[shape_index]
		start_position := Point{x: 2, y: first_clear_row + 3}

		dropRock(&gameboard, shape, start_position, movements, &movement_count)

		first_clear_row = getFirstClearRow(&gameboard)
	}

	// drawGameboard(&gameboard)

	fmt.Printf("First clear row index = tower height: %d\n", first_clear_row)
}

func getFirstClearRow(gameboard *[][]byte) int {
	for i, row := range *gameboard {
		row_clear := true

		for _, position := range row {
			if position == '#' {
				row_clear = false
			}
		}

		if row_clear {
			return i
		}
	}

	return len(*gameboard)
}

func dropRock(gameboard *[][]byte, shape [][]byte, start Point, movements []bool, movement_count *int) {
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
