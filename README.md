# Advent of Code 2022

- trying out Zig and then Golang
  - (When the going gets tough, the tough get Go-ing)

## Notes

- Zig solutions don't deal with windows line endings:
  - Convert any CRLF txt files to unix line endings: ```tr -d '\015' <DOS-file >UNIX-file```

- To run Zig a solution, open main.zig, set the "solution" import path, and ```zig build run```

- To run Go a solution: ```go run src/solutions/14-1/14-1.go```
