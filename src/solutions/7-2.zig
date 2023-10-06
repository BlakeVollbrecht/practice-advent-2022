const std = @import("std");
const ArrayList = std.ArrayList;
const StringHashMap = std.StringHashMap;
const allocator = std.heap.page_allocator;

pub const Node = struct {
    name: []const u8,
    is_dir: bool,
    size: u32,
    parent: ?[]const u8,
    children: ArrayList([]const u8),
    fn init(name: []const u8, is_dir: bool, size: u32, parent: ?[]const u8) Node {
        return Node{
            .name = name,
            .is_dir = is_dir,
            .size = size,
            .parent = parent,
            .children = ArrayList([]const u8).init(allocator),
        };
    }
};

pub fn solve() !void {
    const stdout_file = std.io.getStdOut().writer();
    var bw = std.io.bufferedWriter(stdout_file);
    const stdout = bw.writer();

    var file = try std.fs.cwd().openFile("inputs/7.txt", .{});
    defer file.close();

    var buf_reader = std.io.bufferedReader(file.reader());
    var in_stream = buf_reader.reader();

    var buf: [1024]u8 = undefined;

    var nodes: StringHashMap(Node) = StringHashMap(Node).init(allocator);
    defer nodes.deinit();
    try nodes.put("/", Node.init("root", true, 0, null));
    var current: []const u8 = "/";

    while (try in_stream.readUntilDelimiterOrEof(&buf, '\n')) |line| {
        var values = std.mem.split(u8, line, " ");
        var parts: [3][]const u8 = undefined;
        var i: u32 = 0;

        while (values.next()) |value| {
            parts[i] = value;
            i += 1;
        }
        // std.debug.print("LINE: {s}\n", .{line});

        if (std.mem.eql(u8, parts[0], "$")) {
            if (std.mem.eql(u8, parts[1], "cd")) {
                if (std.mem.eql(u8, parts[2], "/")) {
                    // std.debug.print("changing dir to root\n", .{});
                    current = "/";
                } else if (std.mem.eql(u8, parts[2], "..")) {
                    // std.debug.print("changing dir up\n", .{});
                    const current_node: ?Node = nodes.get(current);

                    if (current_node) |c| {
                        current = c.parent orelse continue;
                    }
                } else { // change to named dir
                    const dir_name: []const u8 = parts[2];
                    // std.debug.print("changing dir down: {s}\n", .{dir_name});

                    const file_path = try makeFilePath(current, dir_name);

                    if (nodes.contains(file_path)) {
                        current = file_path;
                    } else {
                        std.debug.print("encountered cd to non-existant dir\n", .{});
                    }
                }
            } else {
                continue; // ignore "ls"
            }
        } else if (std.mem.eql(u8, parts[0], "dir")) {
            // check for child dir and create if doesn't exist
            const dir_name: []const u8 = parts[1];

            const file_path = try makeFilePath(current, dir_name);

            if (!nodes.contains(file_path)) {
                var new_node = Node.init(dir_name, true, 0, current);
                try nodes.put(file_path, new_node);
                // std.debug.print("added dir: {s}\n", .{file_path});

                // add child id to children array since nodes are flat mapped
                var current_node: ?Node = nodes.get(current);
                if (current_node != null) {
                    try current_node.?.children.append(file_path);
                    try nodes.put(current, current_node.?);
                }
            }
        } else {
            // check for file and create if it doesn't exist
            const file_name = parts[1];
            const file_size = try std.fmt.parseInt(u32, parts[0], 10);

            const file_path = try makeFilePath(current, file_name);

            if (!nodes.contains(file_path)) {
                var new_node = Node.init(file_name, false, file_size, current);
                try nodes.put(file_path, new_node);
                // std.debug.print("added file: {s}\n", .{file_path});

                // add child id to children array since nodes are flat mapped
                var current_node: ?Node = nodes.get(current);
                if (current_node != null) {
                    try current_node.?.children.append(file_path);
                    try nodes.put(current, current_node.?);
                }
            }
        }
    }

    const root_size = try setDirectorySize(&nodes, "/");
    const target_removal_size = root_size - 40000000;

    var sum = getSmallestRemovalDir(target_removal_size, &nodes, "/");

    try stdout.print("{d}, {d}\n", .{ root_size, sum });
    try bw.flush();
}

fn makeFilePath(prefix: []const u8, suffix: []const u8) ![]const u8 {
    var result = ArrayList(u8).init(allocator);
    try result.appendSlice(prefix);
    try result.append('/');
    try result.appendSlice(suffix);

    return result.items;
}

fn setDirectorySize(nodes: *StringHashMap(Node), node_id: []const u8) !u32 {
    var node: Node = nodes.get(node_id) orelse unreachable;

    for (node.children.items) |child_id| {
        node.size += try setDirectorySize(nodes, child_id);
        try nodes.put(node_id, node);
    }

    return node.size;
}

fn getSmallestRemovalDir(target_size: u32, nodes: *StringHashMap(Node), node_id: []const u8) u32 {
    var smallest: u32 = std.math.maxInt(u32);
    var node: Node = nodes.get(node_id) orelse unreachable;

    for (node.children.items) |child_id| {
        var candidate = getSmallestRemovalDir(target_size, nodes, child_id);

        if (candidate < smallest) {
            smallest = candidate;
        }
    }

    if (node.is_dir and node.size >= target_size and node.size < smallest) {
        return node.size;
    } else {
        return smallest;
    }
}
