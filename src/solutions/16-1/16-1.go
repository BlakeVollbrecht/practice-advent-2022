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
// - there are 15 nodes with "rate" values and 15^2 = 225 routes between them (nodes with no rate value just add time to the routes)
//   - need to start on node AA and it has no rate (therefore, 16 relevant nodes and max 256 routes?)
// - on starting node, almost need to calculate total value of every possible future until the time limit and determine the highest
//   - exponential; order of 10^19 operations?
//   - post-completion: time limit cut the complex end off this tree, reducing to maybe 15^6 = 10^7 operations
// - could keep track of a highest-so-far full path, then at each step get a highest-possible-from-here and stop if it doesn't beat the highest-so-far
//	 - highest-possible-from-here could be estimated with getting remaining valves in descending order of importance min distance apart (1 min travel + 1 min open valve)
//   - i.e. greedy algorithm does the estimating but don't rely on it for the answer
//	 - post-completion: didn't end up trying this since it seems like computing these estimates every time wouldn't be worthwhile, especially since they're a high estimate
// - post-completion: re-implemented parts using map instead of array of struct; just wanted to try map in Go, some uses of it here aren't the most appropriate

const STARTING_VALVE_ID = "AA"
const TIME_LIMIT_MINUTES = 30

func main() {
	input, err := readLines("inputs/16.txt")
	check(err)

	timer := time.Now()
	solve(input)
	fmt.Println("running time:", time.Since(timer))
}

func solve(input []string) {
	valve_rates, valve_links := parseInput(input)
	valve_map := getValveMap(valve_rates, valve_links)

	// for _, v := range valve_map {
	// 	fmt.Println(v)
	// }

	total := getHighestScore(valve_map, valve_rates, STARTING_VALVE_ID, nil, 0, TIME_LIMIT_MINUTES)

	fmt.Printf("Answer: %d\n", total)
}

func getHighestScore(valve_map map[string]map[string]int, valve_rates map[string]int, current_id string, visited_valves []string, score int, t int) int {
	visited_valves = append(visited_valves, current_id)

	links, ok := valve_map[current_id]
	if !ok {
		log.Fatalf("Valve links not found for %s", current_id)
	}

	if !hasUnvisitedLinks(links, visited_valves) {
		return score
	}

	highest := 0

	present_values := getPresentValues(t, visited_valves, valve_rates, links)

	for next_id, present_value := range present_values {
		distance, ok := links[next_id]
		if !ok {
			log.Fatalf("Distance not found for %s", next_id)
		}
		updated_time := t - (distance + 1) // time taken to move to next valve and open it is distance + 1
		if updated_time < 0 {
			continue
		}

		prospective_highest := getHighestScore(valve_map, valve_rates, next_id, visited_valves, score+present_value, updated_time)

		if prospective_highest > highest {
			highest = prospective_highest
		}
	}

	if highest > 0 {
		return highest
	} else {
		return score
	}
}

func hasUnvisitedLinks(links map[string]int, visited_valves []string) bool {
	has_unvisited_link := false

	if len(links) == 0 {
		return false
	}

	for link_id := range links {
		if !containsString(visited_valves, link_id) {
			has_unvisited_link = true
		}
	}

	return has_unvisited_link
}

// The value from now until the time limit of travelling to each valve and opening it
func getPresentValues(time_remaining int, visited_valves []string, valve_rates map[string]int, links map[string]int) map[string]int {
	present_values := make(map[string]int)

	for id, distance := range links {
		if containsString(visited_valves, id) {
			continue
		}

		time_needed := distance + 1 // 1 minute for each unit of distance travelled and 1 minute to open valve

		flow_rate, ok := valve_rates[id]
		if !ok {
			log.Fatalf("Valve rate not found for %s", id)
		}
		value := flow_rate * (time_remaining - time_needed)

		present_values[id] = value
	}

	return present_values
}

// Turns the map of all valves and their links into a reduced map that excludes valves with no flow rate
func getValveMap(valve_rates map[string]int, valve_links map[string][]string) map[string]map[string]int {
	var flowing_valves []string
	for id, rate := range valve_rates {
		if rate > 0 {
			flowing_valves = append(flowing_valves, id)
		}
	}
	valve_map := make(map[string]map[string]int)

	valve_distances := getValveDistances(valve_links, flowing_valves, STARTING_VALVE_ID, nil, 0)
	valve_map[STARTING_VALVE_ID] = valve_distances

	for _, id := range flowing_valves {
		valve_distances := getValveDistances(valve_links, flowing_valves, id, nil, 0)
		valve_map[id] = valve_distances
	}

	return valve_map
}

// Use connections of all valves (including zero-rate valves) to determine distance between non-zero-rate valves
func getValveDistances(valve_links map[string][]string, flowing_valves []string, current_valve string, visited_valves []string, distance int) map[string]int {
	visited_valves = append(visited_valves, current_valve)
	valve_distances := make(map[string]int)

	if containsString(flowing_valves, current_valve) { // only making note of valves with a flow rate
		valve_distances[current_valve] = distance
	}

	current_valve_links, ok := valve_links[current_valve]
	if !ok {
		log.Fatalf("Valve links not found for %s", current_valve)
	}

	for _, next_valve := range current_valve_links {
		if containsString(visited_valves, next_valve) { // don't reverse or loop (but multiple branches may visit the same node)
			continue
		}

		next_valve_distances := getValveDistances(valve_links, flowing_valves, next_valve, visited_valves, distance+1)

		for next_valve_id, next_valve_distance := range next_valve_distances {
			stored_distance, exists := valve_distances[next_valve_id]
			if !exists || next_valve_distance < stored_distance {
				valve_distances[next_valve_id] = next_valve_distance
			}
		}
	}

	return valve_distances
}

func containsString(list []string, element string) bool {
	for _, item := range list {
		if item == element {
			return true
		}
	}
	return false
}

func parseInput(input []string) (map[string]int, map[string][]string) {
	valve_id_regex := regexp.MustCompile(`[A-Z]{2}`)
	number_regex := regexp.MustCompile(`\d+`)

	valve_rates := make(map[string]int)
	valve_links := make(map[string][]string)

	for _, line := range input {
		id_matches := valve_id_regex.FindAllString(line, -1)
		number_matches := number_regex.FindAllString(line, -1)

		rate, err := strconv.Atoi(number_matches[0])
		check(err)

		var links []string
		for _, id := range id_matches[1:] {
			links = append(links, id)
		}

		valve_rates[id_matches[0]] = rate
		valve_links[id_matches[0]] = links
	}

	return valve_rates, valve_links
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
