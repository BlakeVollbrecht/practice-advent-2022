const std = @import("std");
const ArrayList = std.ArrayList;
const allocator = std.heap.page_allocator;

pub const Monkey = struct {
    inspections: u64,
    items: ArrayList(u64),
    operation_multiply: bool,
    operation_value: ?u64,
    test_modulus: u32,
    test_true_monkey: u32,
    test_false_monkey: u32,
    fn init(items: []const u64, operation_multiply: bool, operation_value: ?u64, test_modulus: u32, test_true_monkey: u32, test_false_monkey: u32) Monkey {
        var formatted_items = ArrayList(u64).init(allocator);
        formatted_items.appendSlice(items) catch |err| {
            std.debug.print("couldn't append item: {any}", .{err});
        };

        return Monkey{
            .inspections = 0,
            .items = formatted_items,
            .operation_multiply = operation_multiply,
            .operation_value = operation_value,
            .test_modulus = test_modulus,
            .test_true_monkey = test_true_monkey,
            .test_false_monkey = test_false_monkey,
        };
    }
};

fn getMonkeys() !ArrayList(Monkey) {
    var monkeys = ArrayList(Monkey).init(allocator);

    var monkey_0_items = [_]u64{ 53, 89, 62, 57, 74, 51, 83, 97 };
    var monkey_1_items = [_]u64{ 85, 94, 97, 92, 56 };
    var monkey_2_items = [_]u64{ 86, 82, 82 };
    var monkey_3_items = [_]u64{ 94, 68 };
    var monkey_4_items = [_]u64{ 83, 62, 74, 58, 96, 68, 85 };
    var monkey_5_items = [_]u64{ 50, 68, 95, 82 };
    var monkey_6_items = [_]u64{75};
    var monkey_7_items = [_]u64{ 92, 52, 85, 89, 68, 82 };

    try monkeys.append(Monkey.init(&monkey_0_items, true, 3, 13, 1, 5));
    try monkeys.append(Monkey.init(&monkey_1_items, false, 2, 19, 5, 2));
    try monkeys.append(Monkey.init(&monkey_2_items, false, 1, 11, 3, 4));
    try monkeys.append(Monkey.init(&monkey_3_items, false, 5, 17, 7, 6));
    try monkeys.append(Monkey.init(&monkey_4_items, false, 4, 3, 3, 6));
    try monkeys.append(Monkey.init(&monkey_5_items, false, 8, 7, 2, 4));
    try monkeys.append(Monkey.init(&monkey_6_items, true, 7, 5, 7, 0));
    try monkeys.append(Monkey.init(&monkey_7_items, true, null, 2, 0, 1));

    return monkeys;
}

pub fn solve() !void {
    const start_time = std.time.milliTimestamp();

    const stdout_file = std.io.getStdOut().writer();
    var bw = std.io.bufferedWriter(stdout_file);
    const stdout = bw.writer();

    var monkeys = try getMonkeys();
    defer monkeys.deinit();
    const NUM_ROUNDS = 10000;

    for (0..NUM_ROUNDS) |_| {
        try runRound(monkeys);
    }

    // for (monkeys.items, 0..) |monkey, i| {
    //     // try stdout.print("{d}: {any} {?d} {d} {d} {d}\n", .{ i, monkey.operation_multiply, monkey.operation_value, monkey.test_modulus, monkey.test_true_monkey, monkey.test_false_monkey });
    //     try stdout.print("{d}: {d} {any}\n", .{ i, monkey.inspections, monkey.items.items });
    //     try bw.flush();
    // }

    var monkey_inspections = ArrayList(u64).init(allocator);
    for (monkeys.items) |monkey| {
        try monkey_inspections.append(monkey.inspections);
    }

    std.mem.sort(u64, monkey_inspections.items, {}, std.sort.desc(u64));
    var monkey_business: u64 = monkey_inspections.items[0] * monkey_inspections.items[1];

    const timer = std.time.milliTimestamp() - start_time;
    try stdout.print("{d}ms\n", .{timer});
    try bw.flush();

    try stdout.print("{d} x {d} = {d}\n", .{ monkey_inspections.items[0], monkey_inspections.items[1], monkey_business });
    try bw.flush();
}

fn runRound(monkeys: ArrayList(Monkey)) !void {
    for (0..monkeys.items.len) |i| {
        try runTurn(monkeys, i);
    }
}

fn runTurn(monkeys: ArrayList(Monkey), monkey_id: usize) !void {
    const common_denominator = getCommonDenominator(monkeys);

    var monkey = &monkeys.items[monkey_id];

    for (monkey.items.items) |item_worry_level| {
        monkey.inspections += 1;

        var updated_worry_level = getUpdatedWorryLevel(item_worry_level, monkey.operation_multiply, monkey.operation_value);
        updated_worry_level %= common_denominator;

        if (updated_worry_level % monkey.test_modulus == 0) {
            try monkeys.items[monkey.test_true_monkey].items.append(updated_worry_level);
        } else {
            try monkeys.items[monkey.test_false_monkey].items.append(updated_worry_level);
        }
    }

    // all items are thrown by end of turn, but removed at end to avoid indexing issues above
    for (0..monkey.items.items.len) |_| {
        _ = monkey.items.pop();
    }
}

fn getUpdatedWorryLevel(worryLevel: u64, isMultiply: bool, operation_value: ?u64) u64 {
    const value = operation_value orelse worryLevel;

    // std.debug.print("getUpdatedWorryLevel: {d} {d}\n", .{ worryLevel, value });

    const updated_worry_level: u64 = switch (isMultiply) {
        true => worryLevel * value,
        false => worryLevel + value,
    };

    // std.debug.print("updated_worry_level: {d}\n", .{updated_worry_level});

    return updated_worry_level;
}

fn getCommonDenominator(monkeys: ArrayList(Monkey)) u64 {
    var common_denominator: u64 = 1;

    for (monkeys.items) |monkey| {
        common_denominator *= monkey.test_modulus;
    }

    return common_denominator;
}
