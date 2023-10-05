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
    var visible_trees: u64 = 0;

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
            var visible_west = true;
            var visible_east = true;
            var visible_north = true;
            var visible_south = true;

            if (j > 0) { // ignore first column
                visible_west = isLargest(tree, row.items[0..j]);
            }

            if (j < row.items.len - 1) { // ignore last column
                visible_east = isLargest(tree, row.items[j + 1 ..]);
            }

            if (i > 0) { // ignore first row
                var trees_north = ArrayList(u4).init(allocator);
                for (0..i) |k| {
                    try trees_north.append(forest.items[k].items[j]);
                }
                visible_north = isLargest(tree, trees_north.items);
            }

            if (i < forest.items.len - 1) { // ignore last row
                var trees_south = ArrayList(u4).init(allocator);
                for (i + 1..forest.items.len) |k| {
                    try trees_south.append(forest.items[k].items[j]);
                }
                visible_south = isLargest(tree, trees_south.items);
            }

            if (visible_west or visible_east or visible_north or visible_south) {
                visible_trees += 1;
            }
        }
    }

    try stdout.print("{d}\n", .{visible_trees});
    try bw.flush();

    for (forest.items) |row| {
        defer row.deinit();
    }
    defer forest.deinit();
}

fn isLargest(num: u4, contenders: []const u4) bool {
    var is_largest = true;

    for (contenders) |contender| {
        if (contender >= num) {
            is_largest = false;
            break;
        }
    }

    return is_largest;
}
