package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Coordinates struct{ x, y int }

const (
	Right int = iota
	Down
	Left
	Up
)

func main() {
	input, err := readLines("inputs/22.txt")
	check(err)

	timer := time.Now()
	solve(input)
	fmt.Println("running time:", time.Since(timer))
}

func solve(input []string) {
	board := input[:len(input)-2]
	directions := parseDirections(input[len(input)-1])

	endCol, endRow, direction := tracePath(&board, &directions)

	// fmt.Printf("end position: %d, %d, %d\n", endCol, endRow, direction)

	answer := 1000*(endRow+1) + 4*(endCol+1) + direction

	fmt.Printf("Answer: %d\n", answer)
}

func tracePath(board *[]string, directions *[]string) (int, int, int) {
	traceBoard := make([][]byte, len(*board))
	for i, row := range *board {
		traceBoard[i] = []byte(row)
	}

	position := getStartPosition(board)
	facing := Right

	for _, element := range *directions {
		if element == "L" || element == "R" {
			clockwise := element == "R"
			facing = rotate(facing, clockwise)
			continue
		}

		num, err := strconv.Atoi(element)
		check(err)

		for i := 0; i < num; i++ {
			switch facing {
			case Right:
				traceBoard[position.y][position.x] = '>'
				position = moveRight(&(*board)[position.y], position)
			case Down:
				traceBoard[position.y][position.x] = 'v'
				position = moveDown(board, position)
			case Left:
				traceBoard[position.y][position.x] = '<'
				position = moveLeft(&(*board)[position.y], position)
			case Up:
				traceBoard[position.y][position.x] = '^'
				position = moveUp(board, position)
			}
		}
	}

	// printTrace(traceBoard)

	return position.x, position.y, facing
}

func moveRight(row *string, position Coordinates) Coordinates {
	if position.x+1 >= len(*row) { // end of row, try to go to available space at start of row
		for i, char := range *row {
			if char == '.' {
				return Coordinates{x: i, y: position.y}
			}
			if char == '#' {
				return position
			}
		}
	} else if (*row)[position.x+1] == '.' {
		return Coordinates{x: position.x + 1, y: position.y}
	}

	return position
}

func moveDown(board *[]string, position Coordinates) Coordinates {
	if position.y+1 >= len(*board) || // position is on absolute bottom row
		len((*board)[position.y+1]) < position.x+1 || // row below is too short so position below doesn't exist
		(*board)[position.y+1][position.x] == ' ' { // position below isn't part of the board
		for i, row := range *board { // find first position from top
			if len(row) < position.x+1 { // row is too short to include the column
				continue
			}
			if row[position.x] == '.' {
				return Coordinates{x: position.x, y: i}
			}
			if row[position.x] == '#' {
				return position
			}
		}
	} else if (*board)[position.y+1][position.x] == '.' {
		return Coordinates{x: position.x, y: position.y + 1}
	}

	return position
}

func moveLeft(row *string, position Coordinates) Coordinates {
	if position.x <= 0 || (*row)[position.x-1] == ' ' { // start of row, try to go to end of row
		lastPosition := len(*row) - 1
		if (*row)[lastPosition] == '.' {
			return Coordinates{x: lastPosition, y: position.y}
		}
		if (*row)[lastPosition] == '#' {
			return position
		}
	} else if (*row)[position.x-1] == '.' {
		return Coordinates{x: position.x - 1, y: position.y}
	}

	return position
}

func moveUp(board *[]string, position Coordinates) Coordinates {
	if position.y <= 0 || // position is on absolute top row
		len((*board)[position.y-1]) < position.x+1 || // row above is too short so position above doesn't exist
		(*board)[position.y-1][position.x] == ' ' { // position above isn't part of the board
		for i := len(*board) - 1; i > 0; i-- { // find first position from bottom
			row := (*board)[i]
			if len(row) < position.x+1 { // row is too short to include the column
				continue
			}
			if row[position.x] == '.' {
				return Coordinates{x: position.x, y: i}
			}
			if row[position.x] == '#' {
				return position
			}
		}
	} else if (*board)[position.y-1][position.x] == '.' {
		return Coordinates{x: position.x, y: position.y - 1}
	}

	return position
}

func rotate(facing int, clockwise bool) int {
	result := facing

	if clockwise {
		result++
	} else {
		result--
	}

	if result < 0 {
		result = 3
	} else if result > 3 {
		result = 0
	}

	return result
}

// "You begin the path in the leftmost open tile of the top row of tiles."
func getStartPosition(board *[]string) Coordinates {
	firstRow := (*board)[0]

	startX := 0

	for i, char := range firstRow {
		if char == '.' {
			startX = i
			break
		}
	}

	return Coordinates{x: startX, y: 0}
}

func printTrace(traceBoard [][]byte) {
	for _, row := range traceBoard {
		fmt.Println(string(row))
	}
}

func parseDirections(input string) []string {
	re := regexp.MustCompile(`[A-Z]|\d+`)
	return re.FindAllString(input, -1)
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
