from collections import deque


def min_black_cells(grid):
    if not grid or not grid[0]:
        return -1
    m, n = len(grid), len(grid[0])
    min_black = [[float('inf')] * n for _ in range(m)]
    min_black[0][0] = grid[0][0]
    white_queue = deque([(0, 0, grid[0][0])])
    black_queue = deque()
    while white_queue or black_queue:
        if white_queue:
            x, y, count = white_queue.popleft()
        else:
            x, y, count = black_queue.popleft()

        if (x, y) == (m - 1, n - 1):
            return count
        for dx, dy in [(0, 1), (1, 0), (0, -1), (-1, 0)]:
            nx, ny = x + dx, y + dy
            if 0 <= nx < m and 0 <= ny < n and count + grid[nx][ny] < min_black[nx][ny]:
                min_black[nx][ny] = count + grid[nx][ny]
                if grid[nx][ny] == 0:
                    white_queue.append((nx, ny, count))
                else:
                    black_queue.append((nx, ny, count + 1))
    return -1


m, n = map(int, input().split())
grid = []
for _ in range(m):
    row = list(map(int, input().split()))
    grid.append(row)
print(min_black_cells(grid))