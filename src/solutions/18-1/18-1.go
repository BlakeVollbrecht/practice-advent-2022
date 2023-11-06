package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Voxel struct{ x, y, z int }

func main() {
	input, err := readLines("inputs/18.txt")
	check(err)

	timer := time.Now()
	solve(input)
	fmt.Println("running time:", time.Since(timer))
}

func solve(input []string) {
	voxels := parseInput(input)
	surface_area := 0

	for _, voxel := range voxels {
		adjacency := findAdjacency(&voxels, voxel)
		surface_area += 6 - adjacency
	}

	fmt.Printf("Answer: %d\n", surface_area)
}

func findAdjacency(voxels *[]Voxel, voxel Voxel) int {
	adjacent_count := 0

	adjecent_voxels := []Voxel{
		{x: voxel.x + 1, y: voxel.y, z: voxel.z},
		{x: voxel.x - 1, y: voxel.y, z: voxel.z},
		{x: voxel.x, y: voxel.y + 1, z: voxel.z},
		{x: voxel.x, y: voxel.y - 1, z: voxel.z},
		{x: voxel.x, y: voxel.y, z: voxel.z + 1},
		{x: voxel.x, y: voxel.y, z: voxel.z - 1},
	}

	for _, voxel := range *voxels {
		for _, adjecent_voxel := range adjecent_voxels {
			if voxel == adjecent_voxel {
				adjacent_count++
			}
		}
	}

	return adjacent_count
}

func parseInput(input []string) []Voxel {
	var voxels []Voxel

	for _, line := range input {
		coordinates := strings.Split(line, ",")

		x, err_x := strconv.Atoi(coordinates[0])
		check(err_x)
		y, err_y := strconv.Atoi(coordinates[1])
		check(err_y)
		z, err_z := strconv.Atoi(coordinates[2])
		check(err_z)

		voxels = append(voxels, Voxel{x: x, y: y, z: z})
	}

	return voxels
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
