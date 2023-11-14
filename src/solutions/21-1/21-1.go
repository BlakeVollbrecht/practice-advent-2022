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

func main() {
	input, err := readLines("inputs/21.txt")
	check(err)

	timer := time.Now()
	solve(input)
	fmt.Println("running time:", time.Since(timer))
}

func solve(input []string) {
	monkeyYells := parseInput((input))

	answer, ok := calculate(&monkeyYells, "root")
	if !ok {
		log.Fatalln("Error while calculating")
	}

	fmt.Printf("Answer: %d\n", answer)
}

func calculate(monkeyYells *map[string]string, current string) (int, bool) {
	yell, ok := (*monkeyYells)[current]
	if !ok {
		log.Fatalf("Yell not found for %s\n", current)
	}

	elements := strings.Split(yell, " ")
	if len(elements) == 1 { // this yell is a number
		number, err := strconv.Atoi(elements[0])
		check(err)
		return number, true
	} else if len(elements) == 3 { // this yell is an operation
		left, leftOk := calculate(monkeyYells, elements[0])
		right, rightOk := calculate(monkeyYells, elements[2])
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
