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
// - just implement with 2 current_valves and a shared visited_valves?
//   - need 2 of everything except for new shared closed_valves
//   - it explodes the complexity of the brute force approach, effectively squaring the exponent
//   - example input is correct in ~300ms of running time, but actual input will never complete
// - basically implement idea that wasn't tried from part 1, doing it as a manually adjusted depth + score cutoff
//   - i.e. if score is below a threshold at a given depth threshold into the tree, give up (i.e. assume branches beyond that point can no longer reach the highest possible score)
//   - will approximate depth cutoff with time cutoff, because there are already too many params
//   - process of manually adjusting CUTTOFF_TIME lower to get better answers but high enough to complete quickly, and setting CUTOFF_SCORE higher as new high scores are found

const STARTING_VALVE_ID = "AA"
const TIME_LIMIT_MINUTES = 26

// Settings for a fast run to get the correct answer ( running time ~450ms )
const CUTOFF_TIME = 13
const CUTOFF_SCORE = 2500

// Settings used to confidently submit the correct answer (running time 3-4 minutes)
// const CUTOFF_TIME = 5
// const CUTOFF_SCORE = 2658

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

	total := getHighestScore(valve_map, valve_rates, nil, STARTING_VALVE_ID, STARTING_VALVE_ID, nil, nil, 0, TIME_LIMIT_MINUTES, TIME_LIMIT_MINUTES)

	fmt.Printf("Answer: %d\n", total)
}

func getHighestScore(valve_map map[string]map[string]int, valve_rates map[string]int, closed_valves []string, current_id_1 string, current_id_2 string, visited_valves_1 []string, visited_valves_2 []string, score int, t_1 int, t_2 int) int {
	if (t_1 < CUTOFF_TIME || t_2 < CUTOFF_TIME) && score <= CUTOFF_SCORE {
		// give up when a certain score hasn't been achieved by a certain time (to significantly reduce tree size and running time)
		return score
	}

	visited_valves_1 = append(visited_valves_1, current_id_1)
	visited_valves_2 = append(visited_valves_2, current_id_2)

	closed_valves = append(closed_valves, current_id_1)
	closed_valves = append(closed_valves, current_id_2)

	links_1, ok := valve_map[current_id_1]
	if !ok {
		log.Fatalf("Valve links not found for %s", current_id_1)
	}
	links_2, ok := valve_map[current_id_2]
	if !ok {
		log.Fatalf("Valve links not found for %s", current_id_2)
	}

	present_values_1 := getPresentValues(t_1, visited_valves_1, valve_rates, closed_valves, links_1)
	present_values_2 := getPresentValues(t_2, visited_valves_2, valve_rates, closed_valves, links_2)

	if len(present_values_1) == 0 && len(present_values_2) == 0 {
		return score
	}

	highest := 0

	for next_id_1, present_value_1 := range present_values_1 {
		for next_id_2, present_value_2 := range present_values_2 {
			// if current_id_1 == "AA" {
			// 	fmt.Printf("Progress: %s %s\n", next_id_1, next_id_2)
			// }

			if next_id_2 == next_id_1 {
				continue
			}

			distance_1, ok := links_1[next_id_1]
			if !ok {
				log.Fatalf("Distance not found for %s", next_id_1)
			}
			distance_2, ok := links_2[next_id_2]
			if !ok {
				log.Fatalf("Distance not found for %s", next_id_2)
			}

			updated_time_1 := t_1 - (distance_1 + 1) // time taken to move to next valve and open it is distance + 1
			updated_time_2 := t_2 - (distance_2 + 1) // time taken to move to next valve and open it is distance + 1

			if updated_time_1 >= 0 || updated_time_2 >= 0 {
				prospective_highest := getHighestScore(valve_map, valve_rates, closed_valves, next_id_1, next_id_2, visited_valves_1, visited_valves_2, score+present_value_1+present_value_2, updated_time_1, updated_time_2)

				if prospective_highest > highest {
					highest = prospective_highest
				}
			}
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
func getPresentValues(time_remaining int, visited_valves []string, valve_rates map[string]int, closed_valves []string, links map[string]int) map[string]int {
	present_values := make(map[string]int)

	for id, distance := range links {
		if containsString(visited_valves, id) {
			continue
		}
		if containsString(closed_valves, id) {
			present_values[id] = 0 // this will allow one of the partners to continue but get no score so the other can keep going
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
