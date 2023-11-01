# Advent of Code 2022

- trying out Zig and then Golang
  - (When the going gets tough, the tough get Go-ing)

## Notes

- Zig solutions don't deal with windows line endings:
  - Convert any CRLF txt files to unix line endings: ```tr -d '\015' <DOS-file >UNIX-file```

- To run a Zig solution, open main.zig, set the "solution" import path, and ```zig build run```

- To run a Go solution: ```go run src/solutions/14-1/14-1.go```

## Performance

- Intel i5 12600K, 32GB RAM

|    Exercise    |  Running Time  |
| -------------- | -------------- |
|   Day 1 Part 1 |         1 ms   |
|         Part 2 |     64 ±1 ms   |
|                |                |
|   Day 2 Part 1 |       < 1 ms   |
|         Part 2 |       < 1 ms   |
|                |                |
|   Day 3 Part 1 |         1 ms   |
|         Part 2 |         1 ms   |
|                |                |
|   Day 4 Part 1 |         1 ms   |
|         Part 2 |         1 ms   |
|                |                |
|   Day 5 Part 1 |         1 ms   |
|         Part 2 |         2 ms   |
|                |                |
|   Day 6 Part 1 |       < 1 ms   |
|         Part 2 |       < 1 ms   |
|                |                |
|   Day 7 Part 1 |         4 ms   |
|         Part 2 |         4 ms   |
|                |                |
|   Day 8 Part 1 |     41 ±2 ms   |
|         Part 2 |    102 ±3 ms   |
|                |                |
|   Day 9 Part 1 |    127 ±2 ms   |
|         Part 2 |     73 ±4 ms   |
|                |                |
|  Day 10 Part 1 |       < 1 ms   |
|         Part 2 |         4 ms   |
|                |                |
|  Day 11 Part 1 |       < 1 ms   |
|         Part 2 |        15 ms   |
|                |                |
|  Day 12 Part 1 |     10 ±1 ms   | * not drawing maps (adds 5ms)
|         Part 2 |      7 ±1 ms   | * not drawing maps (adds 4ms)
|                |                |
|  Day 13 Part 1 |         4 ms   |
|         Part 2 |     44 ±3 ms   |
|                |                |
|  Day 14 Part 1 |       < 1 ms   | * not drawing grid
|         Part 2 |         3 ms   | * not drawing grid
|                |                |
|  Day 15 Part 1 |       < 1 ms   |
|         Part 2 |   933 ±11 ms   |
|                |                |
|  Day 16 Part 1 |    549 ±7 ms   |
|         Part 2 |   437 ±25 ms   | * manually adjusted complexity
|                |                |
|  Day 17 Part 1 |                |
|         Part 2 |                |
|                |                |
|  Day 18 Part 1 |                |
|         Part 2 |                |
|                |                |
|  Day 19 Part 1 |                |
|         Part 2 |                |
