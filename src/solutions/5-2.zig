const std = @import("std");
const ArrayList = std.ArrayList;

pub fn solve() !void {
    const stdout_file = std.io.getStdOut().writer();
    var bw = std.io.bufferedWriter(stdout_file);
    const stdout = bw.writer();

    var file = try std.fs.cwd().openFile("inputs/5.txt", .{});
    defer file.close();

    var buf_reader = std.io.bufferedReader(file.reader());
    var in_stream = buf_reader.reader();

    var buf: [1024]u8 = undefined;
    const allocator = std.heap.page_allocator;
    var crates: ArrayList(ArrayList(u8)) = ArrayList(ArrayList(u8)).init(allocator);
    var onMoves = false;

    while (try in_stream.readUntilDelimiterOrEof(&buf, '\n')) |line| {
        if (std.mem.eql(u8, line, "")) {
            onMoves = true;
            continue;
        }

        if (!onMoves) {
            try readCrates(line, &crates, allocator);
        } else {
            var tokens = std.mem.split(u8, line, " ");
            var i: u32 = 0;
            var count: u32 = undefined;
            var source: u32 = undefined;
            var dest: u32 = undefined;

            while (tokens.next()) |token| {
                if (i == 1) {
                    count = try std.fmt.parseInt(u32, token, 10);
                } else if (i == 3) {
                    source = try std.fmt.parseInt(u32, token, 10) - 1;
                } else if (i == 5) {
                    dest = try std.fmt.parseInt(u32, token, 10) - 1;
                }

                i += 1;
            }

            // std.debug.print("{s}\n", .{line});

            try moveCrates(crates, count, source, dest, allocator);
        }
    }

    const top_crates = getTopCrates(crates, allocator) catch "";

    try stdout.print("{s}\n", .{top_crates});
    try bw.flush();

    for (crates.items) |stack| {
        defer stack.deinit();
    }
    defer crates.deinit();
}

fn readCrates(line: []const u8, crates: *ArrayList(ArrayList(u8)), allocator: std.mem.Allocator) !void {
    if (std.mem.eql(u8, line, "")) {
        return;
    }

    // read characters in positions 1, 5, 9, etc (4 * i + 1) until end
    for (0..(line.len + 1) / 4) |i| {
        if (crates.items.len < i + 1) {
            var new_stack = ArrayList(u8).init(allocator);
            try crates.*.append(new_stack);
        }

        var current_stack: *ArrayList(u8) = &crates.items[i];
        var potential_item: u8 = line[4 * i + 1];

        if (potential_item >= 'A' and potential_item <= 'Z') {
            try current_stack.insert(0, potential_item);
        }
    }
}

fn moveCrates(crates: ArrayList(ArrayList(u8)), count: u32, source: u32, dest: u32, allocator: std.mem.Allocator) !void {
    var source_stack: *ArrayList(u8) = &crates.items[source];
    var dest_stack: *ArrayList(u8) = &crates.items[dest];

    var temp_stack = ArrayList(u8).init(allocator);

    var j: u32 = 0;
    while (j < count) {
        j += 1;
        try temp_stack.append(source_stack.pop());
    }

    j = 0;
    while (j < count) {
        j += 1;
        try dest_stack.append(temp_stack.pop());
    }
}

fn getTopCrates(crates: ArrayList(ArrayList(u8)), allocator: std.mem.Allocator) ![]const u8 {
    var top_crates = ArrayList(u8).init(allocator);

    for (crates.items) |stack| {
        try top_crates.append(stack.getLast());
    }

    return top_crates.items;
}
