package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type Point struct{ x, y int }

func main() {
	input, err := readLines("inputs/14.txt")
	check(err)

	timer := time.Now()
	solve(input)
	fmt.Println("Time:", time.Since(timer))
}

func solve(input []string) {
	// for _, line := range input {
	// 	fmt.Printf("%s\n", line)
	// }

	paths := getPaths(input)

	// for _, path := range paths {
	// 	for _, point := range path {
	// 		fmt.Printf("%v -> ", point)
	// 	}
	// 	fmt.Print("\n")
	// }

	min_corner, max_corner := getBounds(paths)

	// fmt.Printf("min: (%d, %d) max: (%d, %d)\n", min_corner.x, min_corner.y, max_corner.x, max_corner.y)

	grid, offset_x := makeGrid(min_corner, max_corner, paths)

	// drawGrid(grid)

	const ABSOLUTE_DROP_POSITION_X = 500
	const ABSOLUTE_DROP_POSITION_Y = 0
	drop_position := Point{x: ABSOLUTE_DROP_POSITION_X - offset_x, y: ABSOLUTE_DROP_POSITION_Y}

	fmt.Printf("drop position: %+v\n", drop_position)

	sand_escaped := false
	sand_count := 0
	for !sand_escaped {
		sand_escaped = dropSand(grid, drop_position)

		if !sand_escaped {
			sand_count++
		}
	}

	drawGrid(grid)

	fmt.Printf("Sand count: %d\n", sand_count)
}

func dropSand(grid [][]byte, drop_position Point) bool {
	sand := drop_position
	in_motion := true

	for in_motion {
		if sand.y+1 > len(grid)-1 { // sand escapes grid at the bottom
			return true
		}

		if grid[sand.y+1][sand.x] == 0 { // position below sand is clear
			sand.y++
			continue
		}

		if sand.x-1 < 0 { // sand escapes grid on left side (and therefore falls past bottom)
			return true
		}

		if grid[sand.y+1][sand.x-1] == 0 { // position below and to left is clear
			sand.y++
			sand.x--
			continue
		}

		if sand.x+1 > len(grid[0])-1 { // sand escapes on right side (therefore falls past bottom)
			return true
		}

		if grid[sand.y+1][sand.x+1] == 0 { // position below and to right is clear
			sand.y++
			sand.x++
			continue
		}

		in_motion = false // no valid positions were clear, sand at rest
	}

	grid[sand.y][sand.x] = 'o'

	return false
}

func getPaths(input []string) [][]Point {
	paths := make([][]Point, len(input))

	for i, line := range input {
		points := strings.Split(line, " -> ")

		for _, coordinates := range points {
			point := getPoint(strings.Split(coordinates, ","))
			paths[i] = append(paths[i], point)
		}
	}

	return paths
}

func getPoint(coordinates []string) Point {
	x, x_err := strconv.Atoi(coordinates[0])
	check(x_err)
	y, y_err := strconv.Atoi(coordinates[1])
	check(y_err)

	return Point{x: x, y: y}
}

func getBounds(paths [][]Point) (Point, Point) {
	min_x := math.MaxInt
	max_x := 0
	min_y := 0 // due to the nature of gravity/stacking being simulated, the top needs to be full height
	max_y := 0

	for _, path := range paths {
		for _, point := range path {
			if point.x < min_x {
				min_x = point.x
			}
			if point.x > max_x {
				max_x = point.x
			}
			if point.y > max_y {
				max_y = point.y
			}
		}
	}

	return Point{x: min_x, y: min_y}, Point{x: max_x, y: max_y}
}

func makeGrid(min_corner Point, max_corner Point, paths [][]Point) ([][]byte, int) {
	offset_x := min_corner.x

	grid := make([][]byte, max_corner.y+1)

	for i := 0; i < max_corner.y+1; i++ {
		grid[i] = make([]byte, max_corner.x-offset_x+1)
	}

	for _, path := range paths {
		for j := 1; j < len(path); j++ {
			point := path[j]
			prev_point := path[j-1]

			line_x_min, line_x_max, line_x_len := sortLine(prev_point.x-offset_x, point.x-offset_x)
			line_y_min, line_y_max, line_y_len := sortLine(prev_point.y, point.y)

			if line_y_min == line_y_max {
				for k := 0; k < line_x_len; k++ {
					grid[line_y_min][line_x_min+k] = '#'
				}
			} else if line_x_min == line_x_max {
				for k := 0; k < line_y_len; k++ {
					grid[line_y_min+k][line_x_min] = '#'
				}
			} else {
				log.Fatalf("Diagonal path: %v", path)
			}
		}
	}

	return grid, offset_x
}

func sortLine(a int, b int) (int, int, int) {
	if a < b {
		return a, b, b - a + 1
	}
	return b, a, a - b + 1
}

func drawGrid(grid [][]byte) {
	for _, row := range grid {
		for _, point := range row {
			if point != 0 {
				fmt.Printf("%c", point)
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
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
		if !contains(lines, scanner.Text()) { // dedupe input since most of the lines are non-unique
			lines = append(lines, scanner.Text())
		}
	}

	return lines, scanner.Err()
}

func contains(list []string, token string) bool {
	for _, item := range list {
		if item == token {
			return true
		}
	}
	return false
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
