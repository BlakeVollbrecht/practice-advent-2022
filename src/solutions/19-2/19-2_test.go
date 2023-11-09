package main

import (
	"testing"
)

var mockBlueprint = Blueprint{id: 1, oreRobotCost: Resources{3, 0, 0, 0}, clayRobotCost: Resources{4, 0, 0, 0}, obsidianRobotCost: Resources{3, 10, 0, 0}, geodeRobotCost: Resources{5, 0, 8, 0}, resourceThresholds: Resources{9, 8, 17, 16}}
var mockResources = Resources{1, 3, 1, 0}

func BenchmarkGetRationalBuilds(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getRationalBuilds(&mockBlueprint, Resources{1, 3, 1, 0}, nil)
	}
}

func BenchmarkFindMaxGeodes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		findMaxGeodes(&mockBlueprint, mockResources, mockResources, 22, nil)
	}
}
