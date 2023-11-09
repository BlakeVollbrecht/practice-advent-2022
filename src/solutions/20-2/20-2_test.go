package main

import (
	"container/list"
	"testing"
)

var input = []string{"-1382", "7577", "-629", "-169", "5825", "-1331", "-5150", "7385", "2795", "-1382", "7577", "-629", "-169", "5825", "-1331", "-5150", "7385", "2795", "4589", "-1382", "7577", "-629", "-169", "5825", "-1331", "-5150", "7385", "2795", "3546", "-1382", "7577", "-629", "-169", "5825", "-1331", "-5150", "7385", "2795", "2167", "-1382", "7577", "-629", "-169", "5825", "-1331", "-5150", "7385", "2795", "9854", "0"}

func BenchmarkListCreation(b *testing.B) {
	original := parseInput(input)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		workspace := list.New()
		refs := make([]*list.Element, len(original))
		for i, num := range original {
			workspace.PushBack(num)
			refs[i] = workspace.Back()
		}
	}
}

func BenchmarkMix(b *testing.B) {
	original := parseInput(input)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mix(&original, 1)
	}
}
