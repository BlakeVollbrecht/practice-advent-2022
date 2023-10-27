package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Sensor struct{ x, y, beacon_x, beacon_y, radius int }
type Range struct{ min, max int }

func main() {
	input, err := readLines("inputs/15.txt")
	check(err)

	timer := time.Now()
	solve(input)
	fmt.Println("Time:", time.Since(timer))
}

func solve(input []string) {
	sensors := parseInput(input)

	// for _, sensor := range sensors {
	// 	fmt.Printf("Sensor: %+v\n", sensor)
	// }

	ROW := 2000000
	collisionCount := findCollisionsWithRow(ROW, sensors)
	beaconCount := findBeaconsInRow(ROW, sensors)

	fmt.Printf("Row %d Collisions: %d Beacons: %d\n", ROW, collisionCount, beaconCount)
	fmt.Printf("Answer: %d\n", collisionCount-beaconCount)
}

func findCollisionsWithRow(row_index int, sensors []Sensor) int {
	var collision_ranges []Range

	for _, sensor := range sensors {
		sensor_row_distance := absInt(sensor.y - row_index)
		if sensor_row_distance <= sensor.radius { // sensor area of effect collides with row
			min := sensor.x - (sensor.radius - sensor_row_distance)
			max := sensor.x + (sensor.radius - sensor_row_distance)
			collision_ranges = append(collision_ranges, Range{min: min, max: max})
		}
	}

	// fmt.Printf("ranges: %v\n", collision_ranges)
	collision_ranges = mergeRanges(collision_ranges)
	// fmt.Printf("ranges: %v\n", collision_ranges)

	collisions := sumRangeWidths(collision_ranges)

	return collisions
}

func sumRangeWidths(ranges []Range) int {
	sum := 0

	for _, r := range ranges {
		sum += r.max - r.min + 1
	}

	return sum
}

func mergeRanges(ranges []Range) []Range {
	var merged_ranges []Range

	for _, r := range ranges {
		merged := false

		for i, m := range merged_ranges {
			if r.min <= m.min && r.max >= m.max { // r contains m
				merged_ranges[i].min = r.min
				merged_ranges[i].max = r.max
				merged = true
			} else if r.min >= m.min && r.max <= m.max { // r contained within m
				merged = true
			} else if r.min <= m.max && r.max >= m.max { // r overlaps m.max
				merged_ranges[i].max = r.max
				merged = true
			} else if r.min <= m.min && r.max >= m.min { // r overlaps m.min
				merged_ranges[i].min = r.min
				merged = true
			}

			if merged {
				break
			}
		}

		if !merged {
			merged_ranges = append(merged_ranges, r)
		}
	}

	if len(ranges) > len(merged_ranges) { // recursively merge ranges until they don't merge anymore
		merged_ranges = mergeRanges(merged_ranges)
	}

	return merged_ranges
}

func absInt(a int) int {
	if a >= 0 {
		return a
	}
	return -a
}

func findBeaconsInRow(row_index int, sensors []Sensor) int {
	var beacon_x_coords []int

	for _, sensor := range sensors {
		if sensor.beacon_y == row_index {
			if !containsInt(beacon_x_coords, sensor.beacon_x) { // multiple sensors can have the same beacon; need to check uniqueness
				beacon_x_coords = append(beacon_x_coords, sensor.beacon_x)
			}
		}
	}

	return len(beacon_x_coords)
}

func containsInt(list []int, element int) bool {
	contains := false
	for _, item := range list {
		if item == element {
			contains = true
		}
	}
	return contains
}

func getManhattanDistance(a_x int, a_y int, b_x int, b_y int) int {
	x_distance := a_x - b_x
	if x_distance < 0 {
		x_distance = b_x - a_x
	}

	y_distance := a_y - b_y
	if y_distance < 0 {
		y_distance = b_y - a_y
	}

	return x_distance + y_distance
}

func parseInput(input []string) []Sensor {
	re := regexp.MustCompile(`[-]?\d+`)
	sensors := make([]Sensor, len(input))

	for i, line := range input {
		matches := re.FindAllString(line, -1)

		x, x_err := strconv.Atoi(matches[0])
		check(x_err)
		y, y_err := strconv.Atoi(matches[1])
		check(y_err)
		beacon_x, beacon_x_err := strconv.Atoi(matches[2])
		check(beacon_x_err)
		beacon_y, beacon_y_err := strconv.Atoi(matches[3])
		check(beacon_y_err)

		radius := getManhattanDistance(x, y, beacon_x, beacon_y)

		sensors[i] = Sensor{x: x, y: y, beacon_x: beacon_x, beacon_y: beacon_y, radius: radius}
	}

	return sensors
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
