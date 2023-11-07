package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Notes:
// - big problem is that a single voxel can have both exterior-exposed and interior-exposed faces
//   - may need to break this up into faces instead of voxels, and could try to "walk" exterior surface
//   - can determine exterior surface by observing there's a clear shot to infinity from some face (i.e. just find one, then traverse from there)
//     - can't conclude inverse: i.e. that a voxel is interior based on lack of straight shot to inifinity (could be in a tunnel)
// - could do something less compute efficient but more straightforward(?):
//   - get bounding box 1 unit larger in each dimension and fill it from an edge voxel, then exterior surface is summed based on voxel contact with known exterior voxels

type Coordinates struct{ x, y, z int }

func (a Coordinates) Add(b Coordinates) Coordinates {
	a.x += b.x
	a.y += b.y
	a.z += b.z
	return a
}

type Face struct {
	voxel     Coordinates
	direction int
}

const (
	X_asc int = iota
	X_desc
	Y_asc
	Y_desc
	Z_asc
	Z_desc
)

func main() {
	input, err := readLines("inputs/18.txt")
	check(err)

	timer := time.Now()
	solve(input)
	fmt.Println("running time:", time.Since(timer))
}

func solve(input []string) {
	voxels := parseInput(input)
	faces := *getFaces(&voxels)

	var exterior_face Face
	for _, face := range faces {
		if !isVoxelInPath(&voxels, face.voxel, face.direction) {
			exterior_face = face // get first known exterior face
			break
		}
	}

	var visited_faces []Face
	surface_area := traverseExterior(&voxels, &faces, &visited_faces, exterior_face)

	fmt.Printf("ext face: %v\n", exterior_face)
	fmt.Printf("faces: %d\n", len(faces))
	fmt.Printf("Answer: %d\n", surface_area)
}

func traverseExterior(voxels *[]Coordinates, faces *[]Face, visited_faces *[]Face, current Face) int {
	*visited_faces = append(*visited_faces, current)

	adjacent_faces := getAdjacentFaces(voxels, faces, current)

	face_count := 1

	for _, adjacent_face := range adjacent_faces {
		if !faceExists(visited_faces, adjacent_face) {
			face_count += traverseExterior(voxels, faces, visited_faces, adjacent_face)
		}
	}

	return face_count
}

// it's essentially looking at the 3 possible faces (forward angle, parallel, back angle) that can be attached to each of the 4 sides of the given face
func getAdjacentFaces(voxels *[]Coordinates, faces *[]Face, face Face) []Face {
	var adjacent_faces []Face

	switch face.direction {
	case X_asc, X_desc:
		adjacent_faces = append(adjacent_faces, getAdjacentFace(voxels, faces, face, Coordinates{0, 1, 0}))
		adjacent_faces = append(adjacent_faces, getAdjacentFace(voxels, faces, face, Coordinates{0, 0, 1}))
		adjacent_faces = append(adjacent_faces, getAdjacentFace(voxels, faces, face, Coordinates{0, -1, 0}))
		adjacent_faces = append(adjacent_faces, getAdjacentFace(voxels, faces, face, Coordinates{0, 0, -1}))
	case Y_asc, Y_desc:
		adjacent_faces = append(adjacent_faces, getAdjacentFace(voxels, faces, face, Coordinates{1, 0, 0}))
		adjacent_faces = append(adjacent_faces, getAdjacentFace(voxels, faces, face, Coordinates{0, 0, 1}))
		adjacent_faces = append(adjacent_faces, getAdjacentFace(voxels, faces, face, Coordinates{-1, 0, 0}))
		adjacent_faces = append(adjacent_faces, getAdjacentFace(voxels, faces, face, Coordinates{0, 0, -1}))
	case Z_asc, Z_desc:
		adjacent_faces = append(adjacent_faces, getAdjacentFace(voxels, faces, face, Coordinates{1, 0, 0}))
		adjacent_faces = append(adjacent_faces, getAdjacentFace(voxels, faces, face, Coordinates{0, 1, 0}))
		adjacent_faces = append(adjacent_faces, getAdjacentFace(voxels, faces, face, Coordinates{-1, 0, 0}))
		adjacent_faces = append(adjacent_faces, getAdjacentFace(voxels, faces, face, Coordinates{0, -1, 0}))
	}

	return adjacent_faces
}

// adjacent face off one side of a given face. tries forward-angled, then parallel, then backward-angled depending on presence of adjacent voxels
func getAdjacentFace(voxels *[]Coordinates, faces *[]Face, face Face, side_offset Coordinates) Face {
	forward_offset := Coordinates{0, 0, 0}
	switch face.direction {
	case X_asc:
		forward_offset.x = 1
	case X_desc:
		forward_offset.x = -1
	case Y_asc:
		forward_offset.y = 1
	case Y_desc:
		forward_offset.y = -1
	case Z_asc:
		forward_offset.z = 1
	case Z_desc:
		forward_offset.z = -1
	}

	voxel_ahead := face.voxel.Add(forward_offset.Add(side_offset))
	if voxelExists(voxels, voxel_ahead) {
		return Face{voxel: voxel_ahead, direction: getDirectionFromOffset(face.direction, true, side_offset)}
	}

	voxel_beside := face.voxel.Add(side_offset)
	if voxelExists(voxels, voxel_beside) {
		return Face{voxel: voxel_beside, direction: face.direction}
	}

	return Face{voxel: face.voxel, direction: getDirectionFromOffset(face.direction, false, side_offset)}
}

func getDirectionFromOffset(face_direction int, isAhead bool, side_offset Coordinates) int {
	// adjacent face angled ahead
	if isAhead && side_offset.x > 0 {
		return X_desc
	} else if isAhead && side_offset.x < 0 {
		return X_asc
	} else if isAhead && side_offset.y > 0 {
		return Y_desc
	} else if isAhead && side_offset.y < 0 {
		return Y_asc
	} else if isAhead && side_offset.z > 0 {
		return Z_desc
	} else if isAhead && side_offset.z < 0 {
		return Z_asc
		// adjacent face angled behind
	} else if !isAhead && side_offset.x > 0 {
		return X_asc
	} else if !isAhead && side_offset.x < 0 {
		return X_desc
	} else if !isAhead && side_offset.y > 0 {
		return Y_asc
	} else if !isAhead && side_offset.y < 0 {
		return Y_desc
	} else if !isAhead && side_offset.z > 0 {
		return Z_asc
	} else if !isAhead && side_offset.z < 0 {
		return Z_desc
	}

	// adjacent face parallel
	return face_direction
}

func isVoxelInPath(voxels *[]Coordinates, start Coordinates, direction int) bool {
	for _, voxel := range *voxels {
		switch direction {
		case X_asc:
			if voxel.x > start.x && voxel.y == start.y && voxel.z == start.z {
				return true
			}
		case X_desc:
			if voxel.x < start.x && voxel.y == start.y && voxel.z == start.z {
				return true
			}
		case Y_asc:
			if voxel.x == start.x && voxel.y > start.y && voxel.z == start.z {
				return true
			}
		case Y_desc:
			if voxel.x == start.x && voxel.y < start.y && voxel.z == start.z {
				return true
			}
		case Z_asc:
			if voxel.x == start.x && voxel.y == start.y && voxel.z > start.z {
				return true
			}
		case Z_desc:
			if voxel.x == start.x && voxel.y == start.y && voxel.z < start.z {
				return true
			}
		}
	}

	return false
}

func faceExists(faces *[]Face, face Face) bool {
	for _, candidate := range *faces {
		if candidate == face {
			return true
		}
	}
	return false
}

func voxelExists(voxels *[]Coordinates, voxel Coordinates) bool {
	for _, candidate := range *voxels {
		if candidate == voxel {
			return true
		}
	}
	return false
}

func getFaces(voxels *[]Coordinates) *[]Face {
	var faces []Face

	for _, voxel := range *voxels {
		directions := []int{X_asc, X_desc, Y_asc, Y_desc, Z_asc, Z_desc}

		for _, direction := range directions {
			var adjacent_voxel Coordinates
			switch direction {
			case X_asc:
				adjacent_voxel = Coordinates{x: voxel.x + 1, y: voxel.y, z: voxel.z}
			case X_desc:
				adjacent_voxel = Coordinates{x: voxel.x - 1, y: voxel.y, z: voxel.z}
			case Y_asc:
				adjacent_voxel = Coordinates{x: voxel.x, y: voxel.y + 1, z: voxel.z}
			case Y_desc:
				adjacent_voxel = Coordinates{x: voxel.x, y: voxel.y - 1, z: voxel.z}
			case Z_asc:
				adjacent_voxel = Coordinates{x: voxel.x, y: voxel.y, z: voxel.z + 1}
			case Z_desc:
				adjacent_voxel = Coordinates{x: voxel.x, y: voxel.y, z: voxel.z - 1}
			}

			if !voxelExists(voxels, adjacent_voxel) {
				faces = append(faces, Face{voxel: voxel, direction: direction})
			}
		}
	}

	return &faces
}

func parseInput(input []string) []Coordinates {
	var voxels []Coordinates

	for _, line := range input {
		coordinates := strings.Split(line, ",")

		x, err_x := strconv.Atoi(coordinates[0])
		check(err_x)
		y, err_y := strconv.Atoi(coordinates[1])
		check(err_y)
		z, err_z := strconv.Atoi(coordinates[2])
		check(err_z)

		voxels = append(voxels, Coordinates{x: x, y: y, z: z})
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
