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

	var gap_x, gap_y int
	gap_found := false

	for i := 0; i < 4000000; i++ {
		collision_ranges := findRowCollisionRanges(i, sensors)

		if len(collision_ranges) > 1 {
			fmt.Printf("Gap found: %v\n", collision_ranges)

			gap_found = true
			gap_y = i
			gap_x = getRangeGapPosition(collision_ranges)
		}
	}

	if gap_found {
		tuning_freq := gap_x*4000000 + gap_y

		fmt.Printf("Tuning Freq: %d\n", tuning_freq)
	} else {
		fmt.Println("Gap not found")
	}
}

func getRangeGapPosition(ranges []Range) int {
	if len(ranges) != 2 {
		log.Fatalln("Only accepts 2 ranges")
	}

	high_range := ranges[0]
	low_range := ranges[1]

	if high_range.min < low_range.max {
		high_range = ranges[1]
		low_range = ranges[0]
	}

	if high_range.min-low_range.max != 2 {
		log.Fatalln("Ranges can only have gap of 1 integer between them")
	}

	return low_range.max + 1
}

func findRowCollisionRanges(row_index int, sensors []Sensor) []Range {
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

	return collision_ranges
}

func mergeRanges(ranges []Range) []Range {
	var merged_ranges []Range

	for _, r := range ranges {
		merged := false

		for i, m := range merged_ranges {
			if r.min <= m.min && r.max >= m.max { // r contains m
				// fmt.Printf("r contains m: %v %v\n", r, m)
				merged_ranges[i].min = r.min
				merged_ranges[i].max = r.max
				merged = true
			} else if r.min >= m.min && r.max <= m.max { // r contained within m
				// fmt.Printf("r contained in m: %v %v\n", r, m)
				merged = true
			} else if r.min <= m.max && r.max >= m.max { // r overlaps m.max
				// fmt.Printf("r overlaps m.max: %v %v\n", r, m)
				merged_ranges[i].max = r.max
				merged = true
			} else if r.min <= m.min && r.max >= m.min { // r overlaps m.min
				// fmt.Printf("r overlaps m.min: %v %v\n", r, m)
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
