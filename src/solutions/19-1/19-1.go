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
// - can speed things up by just finishing-out a branch when there's not enough time to buy another geode robot (just return geode robots * remaining time)
// - can this just be worked backward to put all resources in terms of lower resources? (time?)

type Blueprint struct {
	id                                                             int
	oreRobotCost, clayRobotCost, obsidianRobotCost, geodeRobotCost Resources
	maxRobotOreCost                                                int
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

const TIME_MINUTES = 24

const EXCESS_RESOURCE_THRESHOLD float32 = 1.35

func main() {
	input, err := readLines("inputs/19.txt")
	check(err)

	timer := time.Now()
	solve(input)
	fmt.Println("running time:", time.Since(timer))
}

func solve(input []string) {
	blueprints := parseInput(input)

	totalQualityLevel := 0

	for _, blueprint := range blueprints {
		maxGeodes := findMaxGeodes(&blueprint, Resources{0, 0, 0, 0}, Resources{1, 0, 0, 0}, TIME_MINUTES)

		// fmt.Printf("Max geodes blueprint %d: %d\n", blueprint.id, maxGeodes)

		qualityLevel := blueprint.id * maxGeodes
		totalQualityLevel += qualityLevel
	}

	fmt.Printf("Total quality level: %d\n", totalQualityLevel)
}

func findMaxGeodes(blueprint *Blueprint, resources Resources, robots Resources, timeRemaining int) int {
	updatedTime := timeRemaining - 1
	if updatedTime < 0 {
		return resources.geode
	}

	updatedResources := resources.Add(robots) // each robot collects one of its resource

	if !mightAffordGeodeRobot(blueprint, &updatedResources, &robots, updatedTime) {
		return updatedResources.geode + robots.geode*updatedTime // we know how much geode will be collected with current geode robots
	}

	maxGeodes := 0

	rationalRobotBuilds := getRationalBuilds(blueprint, resources)
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

		localMaxGeodes := findMaxGeodes(blueprint, nextResources, nextRobots, updatedTime)
		if localMaxGeodes > maxGeodes {
			maxGeodes = localMaxGeodes
		}
	}

	// include case of building nothing, unless every robot type can be afforded, because not building one when it can be afforded won't lead to best outcome
	if len(rationalRobotBuilds) < 4 {
		localMaxGeodes := findMaxGeodes(blueprint, updatedResources, robots, updatedTime)
		if localMaxGeodes > maxGeodes {
			maxGeodes = localMaxGeodes
		}
	}

	return maxGeodes
}

func mightAffordGeodeRobot(blueprint *Blueprint, resources *Resources, robots *Resources, time int) bool {
	totalOre := resources.ore + robots.ore*time
	possibleNewOreRobots := totalOre / blueprint.oreRobotCost.ore // high estimate not even considering things like ore cost or timing of build
	totalOre += possibleNewOreRobots * time

	possibleNewClayRobots := totalOre / blueprint.clayRobotCost.ore // high estimate not even considering things like ore cost or timing of build
	totalClay := resources.clay + robots.clay*time + possibleNewClayRobots*time

	possibleNewObsidianRobots := totalClay / blueprint.obsidianRobotCost.clay // high estimate not even considering things like ore cost or timing of build
	totalObsidian := resources.obsidian + robots.obsidian*time + possibleNewObsidianRobots*time

	if totalOre > blueprint.geodeRobotCost.ore && totalObsidian > blueprint.geodeRobotCost.obsidian {
		return true
	}
	return false
}

func getRationalBuilds(blueprint *Blueprint, resources Resources) []int {
	if resources.GreaterThanEqual(blueprint.geodeRobotCost) {
		return []int{Geode} // always buy geode robot if it can be afforded
	}

	var rationalBuilds []int

	if resources.GreaterThanEqual(blueprint.obsidianRobotCost) &&
		resources.obsidian < int(EXCESS_RESOURCE_THRESHOLD*float32(blueprint.geodeRobotCost.obsidian)) { // must require more of the resource to build the robot
		rationalBuilds = append(rationalBuilds, Obsidian)
	}

	if resources.GreaterThanEqual(blueprint.clayRobotCost) &&
		resources.clay < int(EXCESS_RESOURCE_THRESHOLD*float32(blueprint.obsidianRobotCost.clay)) { // must require more of the resource to build the robot
		rationalBuilds = append(rationalBuilds, Clay)
	}

	if resources.GreaterThanEqual(blueprint.oreRobotCost) &&
		resources.ore < int(EXCESS_RESOURCE_THRESHOLD*float32(blueprint.maxRobotOreCost)) { // must require more of the resource to build the robot
		rationalBuilds = append(rationalBuilds, Ore)
	}

	return rationalBuilds
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

		blueprints = append(blueprints, Blueprint{
			id:                id,
			oreRobotCost:      oreRobotCost,
			clayRobotCost:     clayRobotCost,
			obsidianRobotCost: obsidianRobotCost,
			geodeRobotCost:    geodeRobotCost,
			maxRobotOreCost:   maxRobotOreCost,
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
