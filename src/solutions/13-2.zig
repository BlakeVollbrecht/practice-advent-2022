const std = @import("std");
const ArrayList = std.ArrayList;
const allocator = std.heap.page_allocator;

const OrderStatus = enum { in_order, not_in_order, unknown };

pub fn solve() !void {
    const stdout_file = std.io.getStdOut().writer();
    var bw = std.io.bufferedWriter(stdout_file);
    const stdout = bw.writer();

    var file = try std.fs.cwd().openFile("inputs/13.txt", .{});
    defer file.close();

    var buf_reader = std.io.bufferedReader(file.reader());
    var in_stream = buf_reader.reader();

    var buf: [1024]u8 = undefined;
    var lines = ArrayList([]const u8).init(allocator);
    defer lines.deinit();

    while (try in_stream.readUntilDelimiterOrEof(&buf, '\n')) |line| {
        if (!std.mem.eql(u8, line, "")) {
            var x = try allocator.alloc(u8, line.len);
            // defer allocator.free(x);
            @memcpy(x, line);
            try lines.append(x);
        }
    }

    const DIVIDERS = [_][]const u8{ "[[2]]", "[[6]]" };
    try lines.appendSlice(&DIVIDERS);

    std.sort.heap([]const u8, lines.items, {}, lessThan);

    // for (lines.items) |line| {
    //     std.debug.print("{s}\n", .{line});
    // }

    const index_1 = findIndex(lines.items, DIVIDERS[0]);
    const index_2 = findIndex(lines.items, DIVIDERS[1]);

    if (index_1 != null and index_2 != null) {
        const key = (index_1.? + 1) * (index_2.? + 1);
        try stdout.print("Divider indices: {d}, {d}\n", .{ index_1.? + 1, index_2.? + 1 });
        try stdout.print("Decoder key: {d}\n", .{key});
        try bw.flush();
    } else {
        try stdout.print("Index of divider(s) was not found\n", .{});
        try bw.flush();
    }
}

fn findIndex(list: [][]const u8, token: []const u8) ?u32 {
    for (list, 0..) |item, i| {
        if (std.mem.eql(u8, item, token)) {
            return @intCast(i);
        }
    }

    return null;
}

fn lessThan(_: void, left: []const u8, right: []const u8) bool {
    return compare(left, right) == OrderStatus.in_order;
}

fn compare(left: []const u8, right: []const u8) OrderStatus {
    const left_items = parseArray(left) catch unreachable;
    const right_items = parseArray(right) catch unreachable;

    var in_order = OrderStatus.unknown;

    for (0..left_items.len) |i| {
        if (right_items.len < i + 1) {
            in_order = OrderStatus.not_in_order;
            break;
        }

        if (left_items[i][0] != '[' and right_items[i][0] != '[') { // left and right are both numbers
            const left_number = std.fmt.parseInt(u4, left_items[i], 10) catch unreachable;
            const right_number = std.fmt.parseInt(u4, right_items[i], 10) catch unreachable;

            if (left_number < right_number) {
                in_order = OrderStatus.in_order;
                break;
            } else if (left_number > right_number) {
                in_order = OrderStatus.not_in_order;
                break;
            }
        } else if (left_items[i][0] != '[') { // left side is number, right side is array
            const order = compare(std.fmt.allocPrint(allocator, "[{s}]", .{left_items[i]}) catch unreachable, right_items[i]);
            in_order = order;
            if (order != OrderStatus.unknown) {
                break;
            }
        } else if (right_items[i][0] != '[') { // left side is array, right side is number
            const order = compare(left_items[i], std.fmt.allocPrint(allocator, "[{s}]", .{right_items[i]}) catch unreachable);
            in_order = order;
            if (order != OrderStatus.unknown) {
                break;
            }
        } else { // left and right are both arrays
            const order = compare(left_items[i], right_items[i]);
            in_order = order;
            if (order != OrderStatus.unknown) {
                break;
            }
        }
    }

    if (in_order == OrderStatus.unknown and left_items.len < right_items.len) {
        in_order = OrderStatus.in_order;
    }

    return in_order;
}

fn parseArray(string: []const u8) ![][]const u8 {
    const s = string[1 .. string.len - 1]; // remove leading and trailing "[" "]"

    var substrings = ArrayList([]const u8).init(allocator);
    defer substrings.deinit();

    var bracketCount: u32 = 0;
    var start: u32 = 0;

    for (0..s.len) |i| {
        if (s[i] == '[') {
            bracketCount += 1;
        } else if (s[i] == ']') {
            bracketCount -= 1;
        } else if (s[i] == ',' and bracketCount == 0) {
            try substrings.append(s[start..i]);
            start = @as(u32, @intCast(i)) + 1;
        }

        if (i == s.len - 1) {
            if (bracketCount == 0) {
                try substrings.append(s[start..s.len]);
            } else {
                std.debug.print("bracketCount incorrect: {s}\n", .{s});
            }
        }
    }

    var x = substrings.toOwnedSlice();
    return x;
}
