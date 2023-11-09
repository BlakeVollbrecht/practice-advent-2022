package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

// Notes:
//  - the mightAffordGeodeRobot stuff could be precalculated for each value of t and stored so it doesn't get calculated a zillion times?
//    - pre-compute resource thresholds as well
//    - mightAffordGeodeRobot can't really be precomputed unless it's an even bigger overestimate, and it already isn't saving enough iterations to be worth having in the hot loop
//       - deleting it is arount 5% faster after the following improvement
//  - 50x faster performance by adding rejectedBuilds to prevent building a robot after a waiting if it could have been built right away
//  - just calculate the lifetime geodes and add to resources when buying a geode robot, rather than running a simulation of it

type Blueprint struct {
	id                                                             int
	oreRobotCost, clayRobotCost, obsidianRobotCost, geodeRobotCost Resources
	resourceThresholds                                             Resources
}

type Resources struct{ ore, clay, obsidian, geode int }

func (a Resources) Add(b Resources) Resources {
	a.ore += b.ore
	a.clay += b.clay
	a.obsidian += b.obsidian
	a.geode += b.geode
	return a
}
func (a Resources) Subtract(b Resources) Resources {
	a.ore -= b.ore
	a.clay -= b.clay
	a.obsidian -= b.obsidian
	a.geode -= b.geode
	return a
}
func (a Resources) GreaterThanEqual(b Resources) bool {
	return a.ore >= b.ore &&
		a.clay >= b.clay &&
		a.obsidian >= b.obsidian &&
		a.geode >= b.geode
}

const (
	Ore int = iota
	Clay
	Obsidian
	Geode
)

const TIME_MINUTES = 32

const EXCESS_RESOURCE_THRESHOLD float32 = 1.67

func main() {
	input, err := readLines("inputs/19.txt")
	check(err)

	timer := time.Now()
	solve(input)
	fmt.Println("running time:", time.Since(timer))
}

func solve(input []string) {
	blueprints := parseInput(input[:3])

	var maxGeodeCounts []int

	for _, blueprint := range blueprints {
		maxGeodes := findMaxGeodes(&blueprint, Resources{0, 0, 0, 0}, Resources{1, 0, 0, 0}, TIME_MINUTES, nil)

		fmt.Printf("Max geodes blueprint %d: %d\n", blueprint.id, maxGeodes)

		maxGeodeCounts = append(maxGeodeCounts, maxGeodes)
	}

	multiplied := 1
	for _, count := range maxGeodeCounts {
		multiplied *= count
	}

	fmt.Printf("Answer: %d\n", multiplied)
}

func findMaxGeodes(blueprint *Blueprint, resources Resources, robots Resources, timeRemaining int, rejectedBuilds []int) int {
	updatedTime := timeRemaining - 1
	if updatedTime < 0 {
		return resources.geode
	}

	updatedResources := resources.Add(robots) // each robot collects one of its resource

	maxGeodes := 0

	rationalRobotBuilds := getRationalBuilds(blueprint, resources, rejectedBuilds)
	for _, robotType := range rationalRobotBuilds {
		var nextResources Resources
		nextRobots := robots
		switch robotType {
		case Ore:
			nextResources = updatedResources.Subtract(blueprint.oreRobotCost)
			nextRobots.ore++
		case Clay:
			nextResources = updatedResources.Subtract(blueprint.clayRobotCost)
			nextRobots.clay++
		case Obsidian:
			nextResources = updatedResources.Subtract(blueprint.obsidianRobotCost)
			nextRobots.obsidian++
		case Geode:
			nextResources = updatedResources.Subtract(blueprint.geodeRobotCost)
			nextRobots.geode++
		}

		localMaxGeodes := findMaxGeodes(blueprint, nextResources, nextRobots, updatedTime, nil)
		if localMaxGeodes > maxGeodes {
			maxGeodes = localMaxGeodes
		}
	}

	// include case of building nothing, unless every robot type can be afforded, because not building one when it can be afforded won't lead to best outcome
	if len(rationalRobotBuilds) < 4 {
		rejectedBuilds := rationalRobotBuilds // if this is a path where we could have just built a robot but didn't, then don't build one of that type next but later on
		localMaxGeodes := findMaxGeodes(blueprint, updatedResources, robots, updatedTime, rejectedBuilds)
		if localMaxGeodes > maxGeodes {
			maxGeodes = localMaxGeodes
		}
	}

	return maxGeodes
}

func getRationalBuilds(blueprint *Blueprint, resources Resources, rejectedBuilds []int) []int {
	if resources.GreaterThanEqual(blueprint.geodeRobotCost) {
		return []int{Geode} // always buy geode robot if it can be afforded
	}

	var rationalBuilds []int

	if !containsInt(rejectedBuilds, Obsidian) &&
		resources.GreaterThanEqual(blueprint.obsidianRobotCost) &&
		resources.obsidian < blueprint.resourceThresholds.obsidian { // must require more of the resource to build the robot
		rationalBuilds = append(rationalBuilds, Obsidian)
	}

	if !containsInt(rejectedBuilds, Clay) &&
		resources.GreaterThanEqual(blueprint.clayRobotCost) &&
		resources.clay < blueprint.resourceThresholds.clay { // must require more of the resource to build the robot
		rationalBuilds = append(rationalBuilds, Clay)
	}

	if !containsInt(rejectedBuilds, Ore) &&
		resources.GreaterThanEqual(blueprint.oreRobotCost) &&
		resources.ore < blueprint.resourceThresholds.ore { // must require more of the resource to build the robot
		rationalBuilds = append(rationalBuilds, Ore)
	}

	return rationalBuilds
}

func containsInt(list []int, element int) bool {
	for _, item := range list {
		if item == element {
			return true
		}
	}
	return false
}

func parseInput(input []string) []Blueprint {
	number_regex := regexp.MustCompile(`\d+`)

	var blueprints []Blueprint

	for _, line := range input {
		number_matches := number_regex.FindAllString(line, -1)

		id, err := strconv.Atoi(number_matches[0])
		check(err)

		oreRobotOre, err := strconv.Atoi(number_matches[1])
		check(err)
		oreRobotCost := Resources{ore: oreRobotOre, clay: 0, obsidian: 0}

		clayRobotOre, err := strconv.Atoi(number_matches[2])
		check(err)
		clayRobotCost := Resources{ore: clayRobotOre, clay: 0, obsidian: 0}

		obsidianRobotOre, err := strconv.Atoi(number_matches[3])
		check(err)
		obsidianRobotClay, err := strconv.Atoi(number_matches[4])
		check(err)
		obsidianRobotCost := Resources{ore: obsidianRobotOre, clay: obsidianRobotClay, obsidian: 0}

		geodeRobotOre, err := strconv.Atoi(number_matches[5])
		check(err)
		geodeRobotObsidian, err := strconv.Atoi(number_matches[6])
		check(err)
		geodeRobotCost := Resources{ore: geodeRobotOre, clay: 0, obsidian: geodeRobotObsidian}

		// pre-compute resource thresholds to keep them out of the hot loop
		maxRobotOreCost := oreRobotOre
		if clayRobotOre > maxRobotOreCost {
			maxRobotOreCost = clayRobotOre
		}
		if obsidianRobotOre > maxRobotOreCost {
			maxRobotOreCost = obsidianRobotOre
		}
		if geodeRobotOre > maxRobotOreCost {
			maxRobotOreCost = geodeRobotOre
		}
		resourceThresholds := Resources{
			ore:      int(float32(maxRobotOreCost) * EXCESS_RESOURCE_THRESHOLD),
			clay:     int(float32(obsidianRobotClay) * EXCESS_RESOURCE_THRESHOLD),
			obsidian: int(float32(geodeRobotObsidian) * EXCESS_RESOURCE_THRESHOLD),
			geode:    -1, // no threshold
		}

		blueprints = append(blueprints, Blueprint{
			id:                 id,
			oreRobotCost:       oreRobotCost,
			clayRobotCost:      clayRobotCost,
			obsidianRobotCost:  obsidianRobotCost,
			geodeRobotCost:     geodeRobotCost,
			resourceThresholds: resourceThresholds,
		})
	}

	return blueprints
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
