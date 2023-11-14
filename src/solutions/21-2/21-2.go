package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Notes:
// - change "humn" by 1 to see how much it affects "qpct", then calculate required change to "qpct" to satisfy equivalence with "dthc" in terms of "humn"
//   - integer division is potential source of non-linearity that could make this fail (just check for all even-divisions, or switch to floats?)
//     - there are no divisions with a remainder in this input
//     - but then adding a delta to "humn" causes uneven divisions to occur
//   - needed to add a large OFFSET to deal with loss of precision of float64 on the large numbers involved
//   - ugly compared to rearranging the equation in terms of "humn", but it works and is interesting in terms of working with linear regression

func main() {
	input, err := readLines("inputs/21.txt")
	check(err)

	timer := time.Now()
	solve(input)
	fmt.Println("running time:", time.Since(timer))
}

func solve(input []string) {
	monkeyYells := parseInput((input))

	left, right := getRootSides(&monkeyYells, 0)
	const OFFSET float64 = 1000000000000
	deltaLeft, deltaRight := getRootSides(&monkeyYells, OFFSET)

	totalAdjustment := right - left
	adjustmentDelta := deltaLeft - left
	if left == deltaLeft {
		totalAdjustment = left - right
		adjustmentDelta = deltaRight - right
	}

	humnOriginalValue, err := strconv.ParseFloat(monkeyYells["humn"], 64)
	check(err)

	answer := humnOriginalValue + totalAdjustment/(adjustmentDelta/OFFSET)

	fmt.Printf("Answer: %f\n", answer)
}

func getRootSides(monkeyYells *map[string]string, humnOffset float64) (float64, float64) {
	root, ok := (*monkeyYells)["root"]
	if !ok {
		log.Fatalln("Yell not found for root")
	}

	elements := strings.Split(root, " ")
	if len(elements) != 3 {

		log.Fatalln("Unexpected number of input elements")
	}

	left, leftOk := calculate(monkeyYells, elements[0], humnOffset)
	right, rightOk := calculate(monkeyYells, elements[2], humnOffset)
	if !leftOk || !rightOk {
		log.Fatalln("Yell not found for root")
	}

	return left, right
}

func calculate(monkeyYells *map[string]string, current string, humnOffset float64) (float64, bool) {
	yell, ok := (*monkeyYells)[current]
	if !ok {
		log.Fatalf("Yell not found for %s\n", current)
	}

	elements := strings.Split(yell, " ")
	if len(elements) == 1 { // this yell is a number
		number, err := strconv.ParseFloat(elements[0], 64)
		check(err)

		if current == "humn" {
			number += humnOffset
		}

		return number, true
	} else if len(elements) == 3 { // this yell is an operation
		left, leftOk := calculate(monkeyYells, elements[0], humnOffset)
		right, rightOk := calculate(monkeyYells, elements[2], humnOffset)
		if !leftOk || !rightOk {
			return 0, false
		}

		operation := elements[1]

		switch operation {
		case "+":
			return left + right, true
		case "-":
			return left - right, true
		case "*":
			return left * right, true
		case "/":
			return left / right, true
		}
	} else {
		log.Fatalln("Unexpected number of input elements")
	}

	return 0, false
}

func parseInput(input []string) map[string]string {
	monkeyYells := make(map[string]string)

	for _, line := range input {
		elements := strings.Split(line, ": ")
		monkeyYells[elements[0]] = elements[1]
	}

	return monkeyYells
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
