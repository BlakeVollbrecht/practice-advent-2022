const std = @import("std");
const ArrayList = std.ArrayList;
const allocator = std.heap.page_allocator;

const Direction = enum { north, south, east, west };

pub const PathNode = struct {
    coordinates: @Vector(2, u32),
    direction: ?Direction,
    fn init(coordinates: @Vector(2, u32), direction: ?Direction) PathNode {
        return PathNode{ .coordinates = coordinates, .direction = direction };
    }
};

pub fn solve() !void {
    var file = try std.fs.cwd().openFile("inputs/12.txt", .{});
    defer file.close();

    var buf_reader = std.io.bufferedReader(file.reader());
    var in_stream = buf_reader.reader();
    var buf: [1024]u8 = undefined;

    var map = ArrayList([]const u8).init(allocator);
    var start: @Vector(2, u32) = undefined;
    var destination: @Vector(2, u32) = undefined;

    var line_count: u32 = 0;

    while (try in_stream.readUntilDelimiterOrEof(&buf, '\n')) |line| {
        var line_copy = ArrayList(u8).init(allocator);
        try line_copy.appendSlice(line);
        try map.append(line_copy.items);

        for (line, 0..) |position, i| {
            if (position == 'S') {
                start = @Vector(2, u32){ @intCast(i), line_count };
                line_copy.items[i] = 'a'; // replace with its height value
            } else if (position == 'E') {
                destination = @Vector(2, u32){ @intCast(i), line_count };
                line_copy.items[i] = 'z'; // replace with its height value
            }
        }

        line_count += 1;
    }

    // const start_node = PathNode.init(start, null);
    var distance_grid = try makeDistanceGrid(@intCast(map.items[0].len), @intCast(map.items.len), start);
    const shortest_path = try travel(map, distance_grid, start, destination);

    const mock_path = ArrayList(PathNode).init(allocator);
    try printPath(map, mock_path);

    try printDistanceGrid(distance_grid);

    std.debug.print("start: {any}  end: {any}\n", .{ start, destination });
    std.debug.print("shortest path: {d}\n", .{shortest_path});
}

fn travel(map: ArrayList([]const u8), distance_grid: ArrayList(ArrayList(u32)), start: @Vector(2, u32), destination: @Vector(2, u32)) !u32 {
    var current_coords = start;
    var possible_moves = try getPossibleMoves(map, distance_grid, current_coords);

    while (possible_moves.items.len > 0) {
        var next_possible_moves = ArrayList(@Vector(2, u32)).init(allocator);

        for (possible_moves.items) |move| {
            if (move[0] == destination[0] and move[1] == destination[1]) {
                return distance_grid.items[destination[1]].items[destination[0]];
            }

            const more_moves = try getPossibleMoves(map, distance_grid, move);
            try next_possible_moves.appendSlice(more_moves.items);
        }

        possible_moves = next_possible_moves;
    }

    return 0;

    // // will need an algorithm to follow the -1, -1, -1, etc. path in the distance grid to trace a shortest path (may be more than one of same length)
    // var shortest_path = ArrayList(PathNode).init(allocator);
    // // try shortest_path.append(current_node);

    // // std.debug.print("shortest: {any}\n", .{shortest_path.items});
    // return shortest_path;
}

// fn travel(map: ArrayList([]const u8), distance_grid: ArrayList(ArrayList(u32)), current_node: PathNode, destination: @Vector(2, u32)) !ArrayList(PathNode) {
//     // initialize a path to return if at the destination
//     if (current_node.coordinates[0] == destination[0] and current_node.coordinates[1] == destination[1]) {
//         var path = ArrayList(PathNode).init(allocator);
//         try path.append(current_node);
//         return path;
//     }

//     const possible_moves = try getPossibleMoves(map, distance_grid, current_node);
//     var shortest_path = ArrayList(PathNode).init(allocator);

//     for (possible_moves.items) |new_node| {
//         const path = try travel(map, distance_grid, new_node, destination);
//         if (path.items.len == 0) {
//             continue;
//         }

//         if (shortest_path.items.len == 0 or path.items.len < shortest_path.items.len) {
//             shortest_path = path;
//         }
//     }

//     // shortest_path being empty is the condition of returning from a dead end; only append if it's returning from the destination
//     if (shortest_path.items.len > 0) {
//         try shortest_path.append(current_node);
//     }

//     // std.debug.print("shortest: {any}\n", .{shortest_path.items});
//     return shortest_path;
// }

fn getPossibleMoves(map: ArrayList([]const u8), distance_grid: ArrayList(ArrayList(u32)), current_coords: @Vector(2, u32)) !ArrayList(@Vector(2, u32)) {
    var possible_moves = ArrayList(@Vector(2, u32)).init(allocator);

    const current_x = current_coords[0];
    const current_y = current_coords[1];
    const current_height = map.items[current_y][current_x];

    const offNorthEdge = current_y > 0;
    const offSouthEdge = current_y < map.items.len - 1;
    const offEastEdge = current_x < map.items[0].len - 1;
    const offWestEdge = current_x > 0;

    if (offNorthEdge) {
        const next_height = map.items[current_y - 1][current_x];
        const next_height_fine = next_height <= current_height + 1;
        const next_unvisited = distance_grid.items[current_y - 1].items[current_x] == std.math.maxInt(u32);

        if (next_height_fine and next_unvisited) {
            distance_grid.items[current_y - 1].items[current_x] = distance_grid.items[current_y].items[current_x] + 1;
            try possible_moves.append(@Vector(2, u32){ current_x, current_y - 1 });
        }
    }
    if (offSouthEdge) {
        const next_height = map.items[current_y + 1][current_x];
        const next_height_fine = next_height <= current_height + 1;
        const next_unvisited = distance_grid.items[current_y + 1].items[current_x] == std.math.maxInt(u32);

        if (next_height_fine and next_unvisited) {
            distance_grid.items[current_y + 1].items[current_x] = distance_grid.items[current_y].items[current_x] + 1;
            try possible_moves.append(@Vector(2, u32){ current_x, current_y + 1 });
        }
    }
    if (offEastEdge) {
        const next_height = map.items[current_y][current_x + 1];
        const next_height_fine = next_height <= current_height + 1;
        const next_unvisited = distance_grid.items[current_y].items[current_x + 1] == std.math.maxInt(u32);

        // std.debug.print("{any}, {any}, {any}\n", .{ current_node.coordinates, next_height_fine, next_unvisited });

        if (next_height_fine and next_unvisited) {
            // std.debug.print("{d}, {d}\n", .{ distance_grid.items[current_y].items[current_x + 1], distance_grid.items[current_y].items[current_x] });
            distance_grid.items[current_y].items[current_x + 1] = distance_grid.items[current_y].items[current_x] + 1;
            try possible_moves.append(@Vector(2, u32){ current_x + 1, current_y });
        }
    }
    if (offWestEdge) {
        const next_height = map.items[current_y][current_x - 1];
        const next_height_fine = next_height <= current_height + 1;
        const next_unvisited = distance_grid.items[current_y].items[current_x - 1] == std.math.maxInt(u32);

        if (next_height_fine and next_unvisited) {
            distance_grid.items[current_y].items[current_x - 1] = distance_grid.items[current_y].items[current_x] + 1;
            try possible_moves.append(@Vector(2, u32){ current_x - 1, current_y });
        }
    }

    return possible_moves;
}

fn makeDistanceGrid(width: u32, height: u32, starting_point: @Vector(2, u32)) !ArrayList(ArrayList(u32)) {
    var grid = ArrayList(ArrayList(u32)).init(allocator);

    for (0..height) |_| {
        var row = ArrayList(u32).init(allocator);

        for (0..width) |_| {
            try row.append(std.math.maxInt(u32));
        }

        try grid.append(row);
    }

    grid.items[starting_point[1]].items[starting_point[0]] = 0;

    return grid;
}

fn printPath(map: ArrayList([]const u8), path: ArrayList(PathNode)) !void {
    const stdout_file = std.io.getStdOut().writer();
    var bw = std.io.bufferedWriter(stdout_file);
    const stdout = bw.writer();

    for (map.items, 0..) |row, i| {
        for (row, 0..) |position, j| {
            const current_coordinates = @Vector(2, u32){ @intCast(i), @intCast(j) };
            const path_direction = pathAtCoordinates(path, current_coordinates);

            var character = position;

            if (path_direction != null) {
                character = switch (path_direction orelse unreachable) {
                    Direction.north => '^',
                    Direction.south => 'v',
                    Direction.east => '>',
                    Direction.west => '<',
                };
            }

            try stdout.print("{c}", .{character});
        }
        try stdout.print("\n", .{});
        try bw.flush();
    }
}

fn printDistanceGrid(distance_grid: ArrayList(ArrayList(u32))) !void {
    const stdout_file = std.io.getStdOut().writer();
    var bw = std.io.bufferedWriter(stdout_file);
    const stdout = bw.writer();

    for (distance_grid.items) |row| {
        for (row.items, 0..) |distance, i| {
            if (i > row.items.len) {
                break;
            }

            var output: u8 = '-';

            if (distance != std.math.maxInt(u32)) {
                const distance_2nd_last_digit = distance % 100 / 10;
                output = @intCast(distance_2nd_last_digit + 48);
            }
            try stdout.print("{c}", .{output});
        }
        try stdout.print("\n", .{});
        try bw.flush();
    }
}

fn pathAtCoordinates(path: ArrayList(PathNode), coordinates: @Vector(2, u32)) ?Direction {
    for (path.items) |node| {
        if (node.coordinates[0] == coordinates[0] and node.coordinates[1] == coordinates[1]) {
            return node.direction;
        }
    }
    return null;
}
