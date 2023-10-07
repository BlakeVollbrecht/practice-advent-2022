const std = @import("std");
const ArrayList = std.ArrayList;
const Vector = std.meta.Vector;
const allocator = std.heap.page_allocator;

pub const PositionCounter = struct {
    count: u32,
    positions_counted: ArrayList([]const u8),
    fn init() PositionCounter {
        return PositionCounter{
            .count = 0,
            .positions_counted = ArrayList([]const u8).init(allocator),
        };
    }
    fn addPosition(self: *PositionCounter, position: @Vector(2, i32)) !void {
        const position_key = try std.fmt.allocPrint(allocator, "{d},{d}", .{ position[0], position[1] });

        const is_new = for (self.positions_counted.items) |item| {
            if (std.mem.eql(u8, item, position_key)) {
                break false;
            }
        } else true;

        if (is_new) {
            try self.positions_counted.append(position_key);
            self.count += 1;
        }
    }
};

pub fn solve() !void {
    const stdout_file = std.io.getStdOut().writer();
    var bw = std.io.bufferedWriter(stdout_file);
    const stdout = bw.writer();

    var file = try std.fs.cwd().openFile("inputs/9.txt", .{});
    defer file.close();

    var buf_reader = std.io.bufferedReader(file.reader());
    var in_stream = buf_reader.reader();

    var buf: [1024]u8 = undefined;

    var tail_positions = PositionCounter.init();
    try tail_positions.addPosition(@Vector(2, i32){ 0, 0 });

    var head: @Vector(2, i32) = .{ 0, 0 };
    var knots = ArrayList(@Vector(2, i32)).init(allocator);
    try knots.appendNTimes(@Vector(2, i32){ 0, 0 }, 9);

    while (try in_stream.readUntilDelimiterOrEof(&buf, '\n')) |line| {
        var values = std.mem.split(u8, line, " ");
        var parts: [2][]const u8 = undefined;
        var i: u32 = 0;

        while (values.next()) |value| {
            parts[i] = value;
            i += 1;
        }
        // std.debug.print("LINE: {s}\n", .{line});

        const direction = parts[0];
        const distance = try std.fmt.parseInt(u32, parts[1], 10);

        for (0..distance) |_| {
            const head_movement: @Vector(2, i32) = switch (direction[0]) {
                'L' => .{ -1, 0 },
                'R' => .{ 1, 0 },
                'U' => .{ 0, 1 },
                'D' => .{ 0, -1 },
                else => unreachable,
            };
            head += head_movement;

            var prev_knot = head;

            for (0..knots.items.len) |j| {
                const offset: @Vector(2, i32) = prev_knot - knots.items[j];
                knots.items[j] += try getMovementForOffset(offset);
                prev_knot = knots.items[j];
            }

            try tail_positions.addPosition(knots.getLast());

            // std.debug.print("{d} {d} --- {any}\n", .{ head[0], head[1], knots.items });
        }
    }

    try stdout.print("{any} | {any} | {d}\n", .{ head, knots.getLast(), tail_positions.count });
    try bw.flush();
}

fn getMovementForOffset(offset: @Vector(2, i32)) !@Vector(2, i32) {
    const is_touching = try std.math.absInt(offset[0]) <= 1 and try std.math.absInt(offset[1]) <= 1;
    if (is_touching) { // tail is within the 9 tiles around the head
        return @Vector(2, i32){ 0, 0 };
    }

    if (offset[0] == 0) { // head and tail are in a vertical line
        if (offset[1] > 1) {
            return @Vector(2, i32){ 0, 1 };
        } else if (offset[1] < -1) {
            return @Vector(2, i32){ 0, -1 };
        }
    } else if (offset[1] == 0) { // head and tail are in a horizontal line
        if (offset[0] > 1) {
            return @Vector(2, i32){ 1, 0 };
        } else if (offset[0] < -1) {
            return @Vector(2, i32){ -1, 0 };
        }
    } else if (offset[0] > 0 and offset[1] > 0) { // head in top right quadrant from tail
        return @Vector(2, i32){ 1, 1 };
    } else if (offset[0] > 0 and offset[1] < 0) { // head in bottom right quadrant from tail
        return @Vector(2, i32){ 1, -1 };
    } else if (offset[0] < 0 and offset[1] < 0) { // head in bottom left quadrant from tail
        return @Vector(2, i32){ -1, -1 };
    } else if (offset[0] < 0 and offset[1] > 0) { // head in top left quadrant from tail
        return @Vector(2, i32){ -1, 1 };
    }

    unreachable;
}
