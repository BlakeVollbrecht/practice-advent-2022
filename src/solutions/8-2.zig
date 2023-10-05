const std = @import("std");
const ArrayList = std.ArrayList;

pub fn solve() !void {
    const stdout_file = std.io.getStdOut().writer();
    var bw = std.io.bufferedWriter(stdout_file);
    const stdout = bw.writer();

    var file = try std.fs.cwd().openFile("inputs/8.txt", .{});
    defer file.close();

    var buf_reader = std.io.bufferedReader(file.reader());
    var in_stream = buf_reader.reader();

    var buf: [1024]u8 = undefined;
    const allocator = std.heap.page_allocator;
    var forest: ArrayList(ArrayList(u4)) = ArrayList(ArrayList(u4)).init(allocator);
    var max_visibility_score: u32 = 0;

    // read input into 2D array
    while (try in_stream.readUntilDelimiterOrEof(&buf, '\n')) |line| {
        var new_row = ArrayList(u4).init(allocator);

        for (line) |tree| {
            var tree_height: u4 = try std.fmt.parseInt(u4, &[1]u8{tree}, 10);
            try new_row.append(tree_height);
        }

        try forest.append(new_row);
    }

    for (forest.items, 0..) |row, i| {
        for (row.items, 0..) |tree, j| {
            var visibility_west: u32 = 0;
            var visibility_east: u32 = 0;
            var visibility_north: u32 = 0;
            var visibility_south: u32 = 0;

            if (j > 0) { // ignore first column
                var trees_west = ArrayList(u4).init(allocator);
                for (row.items[0..j]) |t| {
                    try trees_west.insert(0, t);
                }
                visibility_west = getViewDistance(tree, trees_west.items);
            }

            if (j < row.items.len - 1) { // ignore last column
                visibility_east = getViewDistance(tree, row.items[j + 1 ..]);
            }

            if (i > 0) { // ignore first row
                var trees_north = ArrayList(u4).init(allocator);
                for (0..i) |k| {
                    try trees_north.insert(0, forest.items[k].items[j]);
                }
                visibility_north = getViewDistance(tree, trees_north.items);
            }

            if (i < forest.items.len - 1) { // ignore last row
                var trees_south = ArrayList(u4).init(allocator);
                for (i + 1..forest.items.len) |k| {
                    try trees_south.append(forest.items[k].items[j]);
                }
                visibility_south = getViewDistance(tree, trees_south.items);
            }

            const visibility_score = visibility_west * visibility_east * visibility_north * visibility_south;
            if (visibility_score > max_visibility_score) {
                max_visibility_score = visibility_score;
            }
        }
    }

    try stdout.print("{d}\n", .{max_visibility_score});
    try bw.flush();

    for (forest.items) |row| {
        defer row.deinit();
    }
    defer forest.deinit();
}

fn getViewDistance(view_height: u4, heights: []const u4) u32 {
    var view_distance: u32 = 0;

    for (heights) |height| {
        view_distance += 1;
        if (height >= view_height) {
            break;
        }
    }

    return view_distance;
}
