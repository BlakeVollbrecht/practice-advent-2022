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
    var prev_line: []u8 = try allocator.alloc(u8, 0);
    var index: u32 = 1;
    var sum_indices: u32 = 0;

    while (try in_stream.readUntilDelimiterOrEof(&buf, '\n')) |line| {
        if (!std.mem.eql(u8, line, "")) {
            if (!std.mem.eql(u8, prev_line, "")) {
                const in_order = try checkOrder(prev_line, line);

                if (in_order != OrderStatus.not_in_order) {
                    sum_indices += index;
                }
            }
        } else {
            index += 1;
        }

        allocator.free(prev_line);
        prev_line = try allocator.alloc(u8, line.len);
        @memcpy(prev_line, line);
    }

    try stdout.print("Sum of in-order pair indices: {d}\n", .{sum_indices});
    try bw.flush();
}

fn printParsed(label: []const u8, parsed: [][]const u8) void {
    std.debug.print("{s}: ", .{label});
    for (parsed) |string| {
        std.debug.print("{s}, ", .{string});
    }
    std.debug.print("\n", .{});
}

fn checkOrder(left: []const u8, right: []const u8) !OrderStatus {
    const left_items = try parseArray(left);
    const right_items = try parseArray(right);

    // printParsed("L", left_items);
    // printParsed("R", right_items);

    var in_order = OrderStatus.unknown;

    for (0..left_items.len) |i| {
        if (right_items.len < i + 1) {
            in_order = OrderStatus.not_in_order;
            break;
        }

        if (left_items[i][0] != '[' and right_items[i][0] != '[') { // left and right are both numbers
            const left_number = try std.fmt.parseInt(u4, left_items[i], 10);
            const right_number = try std.fmt.parseInt(u4, right_items[i], 10);

            if (left_number < right_number) {
                in_order = OrderStatus.in_order;
                break;
            } else if (left_number > right_number) {
                in_order = OrderStatus.not_in_order;
                break;
            }
        } else if (left_items[i][0] != '[') { // left side is number, right side is array
            const order = try checkOrder(try std.fmt.allocPrint(allocator, "[{s}]", .{left_items[i]}), right_items[i]);
            in_order = order;
            if (order != OrderStatus.unknown) {
                break;
            }
        } else if (right_items[i][0] != '[') { // left side is array, right side is number
            const order = try checkOrder(left_items[i], try std.fmt.allocPrint(allocator, "[{s}]", .{right_items[i]}));
            in_order = order;
            if (order != OrderStatus.unknown) {
                break;
            }
        } else { // left and right are both arrays
            const order = try checkOrder(left_items[i], right_items[i]);
            in_order = order;
            if (order != OrderStatus.unknown) {
                break;
            }
        }
    }

    if (in_order == OrderStatus.unknown and left_items.len < right_items.len) {
        in_order = OrderStatus.in_order;
        return in_order;
    }

    return in_order;
}

fn parseArray(string: []const u8) ![][]const u8 {
    const s = string[1 .. string.len - 1]; // remove leading and trailing "[" "]"

    var substrings = ArrayList([]const u8).init(allocator);
    // defer substrings.deinit();

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

    return substrings.items;
}
