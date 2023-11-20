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

|    Exercise    |  Running Time  |                           Notes                          |
| -------------- | -------------- | -------------------------------------------------------- |
|   Day 1 Part 1 |         1 ms   |                                                          |
|   Day 1 Part 2 |     64 ±1 ms   |                                                          |
|                |                |                                                          |
|   Day 2 Part 1 |       < 1 ms   |                                                          |
|   Day 2 Part 2 |       < 1 ms   |                                                          |
|                |                |                                                          |
|   Day 3 Part 1 |         1 ms   |                                                          |
|   Day 3 Part 2 |         1 ms   |                                                          |
|                |                |                                                          |
|   Day 4 Part 1 |         1 ms   |                                                          |
|   Day 4 Part 2 |         1 ms   |                                                          |
|                |                |                                                          |
|   Day 5 Part 1 |         1 ms   |                                                          |
|   Day 5 Part 2 |         2 ms   |                                                          |
|                |                |                                                          |
|   Day 6 Part 1 |       < 1 ms   |                                                          |
|   Day 6 Part 2 |       < 1 ms   |                                                          |
|                |                |                                                          |
|   Day 7 Part 1 |         4 ms   |                                                          |
|   Day 7 Part 2 |         4 ms   |                                                          |
|                |                |                                                          |
|   Day 8 Part 1 |     41 ±2 ms   |                                                          |
|   Day 8 Part 2 |    102 ±3 ms   |                                                          |
|                |                |                                                          |
|   Day 9 Part 1 |    127 ±2 ms   |                                                          |
|   Day 9 Part 2 |     73 ±4 ms   |                                                          |
|                |                |                                                          |
|  Day 10 Part 1 |       < 1 ms   |                                                          |
|  Day 10 Part 2 |         4 ms   |                                                          |
|                |                |                                                          |
|  Day 11 Part 1 |       < 1 ms   |                                                          |
|  Day 11 Part 2 |        15 ms   |                                                          |
|                |                |                                                          |
|  Day 12 Part 1 |     10 ±1 ms   | * not drawing maps (adds 5ms)                            |
|  Day 12 Part 2 |      7 ±1 ms   | * not drawing maps (adds 4ms)                            |
|                |                |                                                          |
|  Day 13 Part 1 |         4 ms   |                                                          |
|  Day 13 Part 2 |     44 ±3 ms   |                                                          |
|                |                |                                                          |
|  Day 14 Part 1 |       < 1 ms   | * not drawing grid                                       |
|  Day 14 Part 2 |         3 ms   | * not drawing grid                                       |
|                |                |                                                          |
|  Day 15 Part 1 |       < 1 ms   |                                                          |
|  Day 15 Part 2 |   933 ±11 ms   | * no optimization done                                   |
|                |                |                                                          |
|  Day 16 Part 1 |    549 ±7 ms   |                                                          |
|  Day 16 Part 2 |   437 ±25 ms   | * manually-adjusted threshold                            |
|                |                |                                                          |
|  Day 17 Part 1 |     29 ±1 ms   | * < 1ms if just using improved code from part 2          |
|  Day 17 Part 2 |      2 ±1 ms   | * to see pattern. full run has correct answer in 55h9m31s|
|                |                |                                                          |
|  Day 18 Part 1 |     35 ±1 ms   | * 25ms if using face-counting code from part 2           |
|  Day 18 Part 2 |     61 ±2 ms   |                                                          |
|                |                |                                                          |
|  Day 19 Part 1 |   975 ±38 ms   | * 119 ±3 ms if using optimized code from part 2          |
|  Day 19 Part 2 |  5100 ±90 ms   | * these have manually-adjusted EXCESS_RESOURCE_THRESHOLD |
|                |                |                                                          |
|  Day 20 Part 1 |     37 ±2 ms   |                                                          |
|  Day 20 Part 2 |    430 ±7 ms   | * no optimization done                                   |
|                |                |                                                          |
|  Day 21 Part 1 |       < 1 ms   |                                                          |
|  Day 21 Part 2 |       < 1 ms   |                                                          |
|                |                |                                                          |
|  Day 22 Part 1 |       < 1 ms   | * 13 ±2 ms if outputting trace of path                   |
|  Day 22 Part 2 |       < 1 ms   | * 14 ±2 ms if outputting trace of path                   |
|                |                |                                                          |
|  Day 23 Part 1 |      3 ±1 ms   |                                                          |
|  Day 23 Part 2 |    215 ±6 ms   |                                                          |
|                |                |                                                          |
|  Day 24 Part 1 |                |                                                          |
|  Day 24 Part 2 |                |                                                          |
|                |                |                                                          |
|  Day 25 Part 1 |                |                                                          |
|  Day 25 Part 2 |                |                                                          |
