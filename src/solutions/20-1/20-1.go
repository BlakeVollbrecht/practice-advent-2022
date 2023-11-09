package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

func main() {
	input, err := readLines("inputs/20.txt")
	check(err)

	timer := time.Now()
	solve(input)
	fmt.Println("running time:", time.Since(timer))
}

func solve(input []string) {
	original := parseInput(input)

	mixed := mix(&original)

	coordinates := getCoordinates(mixed)

	var answer int16 = 0
	for i := 0; i < len(coordinates); i++ {
		answer += coordinates[i]
	}

	fmt.Printf("Coordinates: %v\n", coordinates)
	fmt.Printf("Answer: %d\n", answer)
}

func getCoordinates(mixed *[]int16) [3]int16 {
	zeroPosition := 0
	for i := 0; i < len(*mixed); i++ {
		if (*mixed)[i] == 0 {
			zeroPosition = i
			break
		}
	}

	position1 := (zeroPosition + 1000) % len(*mixed)
	position2 := (zeroPosition + 2000) % len(*mixed)
	position3 := (zeroPosition + 3000) % len(*mixed)

	return [3]int16{(*mixed)[position1], (*mixed)[position2], (*mixed)[position3]}
}

func mix(original *[]int16) *[]int16 {
	workspace := list.New()
	refs := make([]*list.Element, len(*original))
	for i, num := range *original {
		workspace.PushBack(num)
		refs[i] = workspace.Back()
	}

	for i, num := range *original {
		node := refs[i]
		destination := getDestination(workspace, node, num)

		if num > 0 {
			workspace.MoveAfter(node, destination)

		} else if num < 0 { // specifically ignoring the 0 case so move isn't called for 0
			workspace.MoveBefore(node, destination)
		}
	}

	var mixed []int16
	for element := workspace.Front(); element != nil; element = element.Next() {
		mixed = append(mixed, element.Value.(int16))
	}

	return &mixed
}

func formatWorkspace(workspace *list.List, start *list.Element) string {
	str := ""

	element := start
	startCount := 0
	for true {
		if element == start && startCount > 0 {
			break
		} else if element == start {
			startCount++
		}

		str = fmt.Sprintf("%v, %d", str, element.Value)

		element = element.Next()
		if element == nil {
			element = workspace.Front()
		}
	}

	return str
}

func getDestination(workspace *list.List, element *list.Element, distance int16) *list.Element {
	destination := element

	// modulus used to prevent unnecessary cycling
	//   "-1" is to correct for the fact that element is supposed to "be moving" (i.e. not in-place when coming back around)
	for i := 0; i < int(math.Abs(float64(distance)))%(workspace.Len()-1); i++ {
		if distance > 0 {
			destination = destination.Next()
			if destination == nil {
				destination = workspace.Front()
			}
		} else {
			destination = destination.Prev()
			if destination == nil {
				destination = workspace.Back()
			}
		}
	}

	return destination
}

func find(arr []int16, num int16) []int {
	var locations []int

	for i, element := range arr {
		if element == num {
			locations = append(locations, i)
		}
	}

	return locations
}

func parseInput(input []string) []int16 {
	numbers := make([]int16, len(input))

	for i, line := range input {
		number, err := strconv.Atoi(line)
		check(err)

		if number > math.MaxInt16 || number < math.MinInt16 {
			log.Fatalf("Parsed number %d out of int16 range\n", number)
		}

		numbers[i] = int16(number)
	}

	return numbers
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
