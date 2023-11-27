package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

// Notes:
//  - don't need to convert it to decimal; input and output are same number format and it may be easier to add this than to convert
//    - might need to implement conversion anyway depending on what's in part 2, but adding seems interesting
//    - when a carry happens (column value > 2), need to adjust the current column to current column's sum minus 5
//      - there's a negative carry, in which case adjust current column to current column's sum plus 5
//  - can store snafu digits in integer array by replacing - and = with actual -1 and -2
//    - store backward so least significant digit is at index 0

func main() {
	input, err := readLines("inputs/25.txt")
	check(err)

	timer := time.Now()
	solve(input)
	fmt.Println("running time:", time.Since(timer))
}

func solve(input []string) {
	fuelCapacities := parseInput(&input)

	// fmt.Printf("%v\n", fuelCapacities)

	var total []int
	for _, number := range fuelCapacities {
		total = add(total, number)
	}

	fmt.Printf("Total fuel: %s\n", formatSnafu(total))
}

func add(a, b []int) []int {
	highestLength := len(a)
	if len(b) > highestLength {
		highestLength = len(b)
	}

	result := make([]int, highestLength)

	for i := 0; i < highestLength; i++ {
		sum := 0
		carry := 0

		if len(a) > i && len(b) > i {
			sum, carry = digitAdd(a[i], b[i], result[i])
		} else if len(a) <= i {
			sum, carry = digitAdd(b[i], 0, result[i])
		} else if len(b) <= i {
			sum, carry = digitAdd(a[i], 0, result[i])
		}

		result[i] = sum

		if i == highestLength-1 && carry != 0 { // if on the last digit of result, append the carry
			result = append(result, carry)
		} else if i < highestLength-1 {
			result[i+1] = carry
		}
	}

	return result
}

func digitAdd(a, b, c int) (int, int) {
	sum := a + b + c
	carry := 0

	if sum > 2 {
		carry = 1
		sum -= 5
	} else if sum < -2 {
		carry = -1
		sum += 5
	}

	return sum, carry
}

func formatSnafu(number []int) string {
	var result string

	for _, digit := range number {
		switch digit {
		case -2:
			result = "=" + result
		case -1:
			result = "-" + result
		case 0:
			result = "0" + result
		case 1:
			result = "1" + result
		case 2:
			result = "2" + result
		default:
			log.Fatalln("Unexpected digit while formatting")
		}
	}

	return result
}

func parseInput(input *[]string) [][]int {
	numbers := make([][]int, len(*input))

	for i, row := range *input {
		numbers[i] = make([]int, len(row))

		reverseIndex := len(row) - 1
		for _, position := range row {
			switch position {
			case '=':
				numbers[i][reverseIndex] = -2
			case '-':
				numbers[i][reverseIndex] = -1
			case '0':
				numbers[i][reverseIndex] = 0
			case '1':
				numbers[i][reverseIndex] = 1
			case '2':
				numbers[i][reverseIndex] = 2
			default:
				log.Fatalln("Invalid character in input")
			}

			reverseIndex--
		}
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
