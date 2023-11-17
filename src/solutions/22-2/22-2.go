package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

// Notes
// - initial thought was to use graph of nodes with 4 directions, but...
//  - parsing the faces still requires cop-out definition of bounds, or unreasonable complexity
// 	- instructions want the direction to be maintained in terms of the input board format
// - will define specific edge regions in existing code that checked for edges in part 1
//  - hitting a certain edge in a certain direction teleports to a specific place and direction
//  - 24 face-edges on cube, connecting to give 12 cube-edges
//  - 5 cube-edges contained in board definition, 7 cube-edges (14 face-edges) remaining to connect
//
// - make array of seamPoints, each having a pair of entrance point+direction
//  - whenever it's going off an edge from part 1 code, refer to the seamPoints

type Coordinates struct{ x, y int }

type Seam struct {
	a_start    Coordinates
	a_end      Coordinates
	a_approach int
	b_start    Coordinates
	b_end      Coordinates
	b_approach int
}

type SeamPoint struct {
	a          Coordinates
	a_approach int
	b          Coordinates
	b_approach int
}

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
	seamPoints := makeSeamPoints()
	directions := parseDirections(input[len(input)-1])

	endCol, endRow, direction := tracePath(&board, &seamPoints, &directions)

	// fmt.Printf("end position: %d, %d, %d\n", endCol, endRow, direction)

	answer := 1000*(endRow+1) + 4*(endCol+1) + direction

	fmt.Printf("Answer: %d\n", answer)
}

func tracePath(board *[]string, seamPoints *[]SeamPoint, directions *[]string) (int, int, int) {
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
				position, facing = moveRight(board, seamPoints, position)
			case Down:
				traceBoard[position.y][position.x] = 'v'
				position, facing = moveDown(board, seamPoints, position)
			case Left:
				traceBoard[position.y][position.x] = '<'
				position, facing = moveLeft(board, seamPoints, position)
			case Up:
				traceBoard[position.y][position.x] = '^'
				position, facing = moveUp(board, seamPoints, position)
			}
		}
	}

	// printTrace(traceBoard)

	return position.x, position.y, facing
}

func moveRight(board *[]string, seamPoints *[]SeamPoint, position Coordinates) (Coordinates, int) {
	row := (*board)[position.y]

	if position.x+1 >= len(row) { // end of row, try to go to available space at start of row
		seamExit, exitDirection := getSeamExit(seamPoints, position, Right)

		if (*board)[seamExit.y][seamExit.x] == '.' {
			return seamExit, exitDirection
		}
		if (*board)[seamExit.y][seamExit.x] == '#' {
			return position, Right
		}
	} else if (row)[position.x+1] == '.' {
		return Coordinates{x: position.x + 1, y: position.y}, Right
	}

	return position, Right
}

func moveDown(board *[]string, seamPoints *[]SeamPoint, position Coordinates) (Coordinates, int) {
	if position.y+1 >= len(*board) || // position is on absolute bottom row
		len((*board)[position.y+1]) < position.x+1 || // row below is too short so position below doesn't exist
		(*board)[position.y+1][position.x] == ' ' { // position below isn't part of the board
		seamExit, exitDirection := getSeamExit(seamPoints, position, Down)

		if (*board)[seamExit.y][seamExit.x] == '.' {
			return seamExit, exitDirection
		}
		if (*board)[seamExit.y][seamExit.x] == '#' {
			return position, Down
		}

	} else if (*board)[position.y+1][position.x] == '.' {
		return Coordinates{x: position.x, y: position.y + 1}, Down
	}

	return position, Down
}

func moveLeft(board *[]string, seamPoints *[]SeamPoint, position Coordinates) (Coordinates, int) {
	row := (*board)[position.y]

	if position.x <= 0 || (row)[position.x-1] == ' ' { // start of row, try to go to end of row
		seamExit, exitDirection := getSeamExit(seamPoints, position, Left)

		if (*board)[seamExit.y][seamExit.x] == '.' {
			return seamExit, exitDirection
		}
		if (*board)[seamExit.y][seamExit.x] == '#' {
			return position, Left
		}
	} else if (row)[position.x-1] == '.' {
		return Coordinates{x: position.x - 1, y: position.y}, Left
	}

	return position, Left
}

func moveUp(board *[]string, seamPoints *[]SeamPoint, position Coordinates) (Coordinates, int) {
	if position.y <= 0 || // position is on absolute top row
		len((*board)[position.y-1]) < position.x+1 || // row above is too short so position above doesn't exist
		(*board)[position.y-1][position.x] == ' ' { // position above isn't part of the board
		seamExit, exitDirection := getSeamExit(seamPoints, position, Up)

		if (*board)[seamExit.y][seamExit.x] == '.' {
			return seamExit, exitDirection
		}
		if (*board)[seamExit.y][seamExit.x] == '#' {
			return position, Up
		}

	} else if (*board)[position.y-1][position.x] == '.' {
		return Coordinates{x: position.x, y: position.y - 1}, Up
	}

	return position, Up
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

func getSeamExit(seamPoints *[]SeamPoint, entry Coordinates, approach int) (Coordinates, int) {
	for _, seamPoint := range *seamPoints {
		if seamPoint.a.x == entry.x && seamPoint.a.y == entry.y && seamPoint.a_approach == approach {
			exitDirection := rotate(rotate(seamPoint.b_approach, true), true)
			return seamPoint.b, exitDirection
		} else if seamPoint.b.x == entry.x && seamPoint.b.y == entry.y && seamPoint.b_approach == approach {
			exitDirection := rotate(rotate(seamPoint.a_approach, true), true)
			return seamPoint.a, exitDirection
		}
	}
	log.Fatalf("Seam point not found for %v, %d\n", entry, approach)
	return Coordinates{}, 0
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

// bespoke to orientation/configuration/size of cube panels in this specific input file
func makeSeamPoints() []SeamPoint {
	sideLength := 50

	seams := []Seam{
		{a_start: Coordinates{50, 0}, a_end: Coordinates{99, 0}, a_approach: Up, b_start: Coordinates{0, 150}, b_end: Coordinates{0, 199}, b_approach: Left},          // top center to bottomost left side
		{a_start: Coordinates{100, 49}, a_end: Coordinates{149, 49}, a_approach: Down, b_start: Coordinates{99, 50}, b_end: Coordinates{99, 99}, b_approach: Right},   // first 90 degree
		{a_start: Coordinates{50, 50}, a_end: Coordinates{50, 99}, a_approach: Left, b_start: Coordinates{0, 100}, b_end: Coordinates{49, 100}, b_approach: Up},       // second 90 degree
		{a_start: Coordinates{50, 149}, a_end: Coordinates{99, 149}, a_approach: Down, b_start: Coordinates{49, 150}, b_end: Coordinates{49, 199}, b_approach: Right}, // third 90 degree
		{a_start: Coordinates{100, 0}, a_end: Coordinates{149, 0}, a_approach: Up, b_start: Coordinates{0, 199}, b_end: Coordinates{49, 199}, b_approach: Down},       // top right horizontal to bottom left horizontal
		{a_start: Coordinates{149, 49}, a_end: Coordinates{149, 0}, a_approach: Right, b_start: Coordinates{99, 100}, b_end: Coordinates{99, 149}, b_approach: Right}, // top right side to middle right side
		{a_start: Coordinates{50, 49}, a_end: Coordinates{50, 0}, a_approach: Left, b_start: Coordinates{0, 100}, b_end: Coordinates{0, 149}, b_approach: Left},       // top left side to middle left side
	}

	var seamPoints []SeamPoint

	for _, seam := range seams {
		a_horizontal := 0
		if seam.a_start.x-seam.a_end.x < 0 {
			a_horizontal = 1
		} else if seam.a_start.x-seam.a_end.x > 0 {
			a_horizontal = -1
		}
		a_vertical := 0
		if seam.a_start.y-seam.a_end.y < 0 {
			a_vertical = 1
		} else if seam.a_start.y-seam.a_end.y > 0 {
			a_vertical = -1
		}
		b_horizontal := 0
		if seam.b_start.x-seam.b_end.x < 0 {
			b_horizontal = 1
		} else if seam.b_start.x-seam.b_end.x > 0 {
			b_horizontal = -1
		}
		b_vertical := 0
		if seam.b_start.y-seam.b_end.y < 0 {
			b_vertical = 1
		} else if seam.b_start.y-seam.b_end.y > 0 {
			b_vertical = -1
		}

		for i := 0; i < sideLength; i++ {
			seamPoint := SeamPoint{
				a: Coordinates{
					x: seam.a_start.x + a_horizontal*i,
					y: seam.a_start.y + a_vertical*i,
				},
				a_approach: seam.a_approach,
				b: Coordinates{
					x: seam.b_start.x + b_horizontal*i,
					y: seam.b_start.y + b_vertical*i,
				},
				b_approach: seam.b_approach,
			}

			seamPoints = append(seamPoints, seamPoint)
		}
	}

	return seamPoints
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
