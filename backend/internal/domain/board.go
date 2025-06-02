package domain

import "fmt"

type Board struct {
	Size int
	Data [][]int32
}

func NewBoard(size int) (*Board, error) {
	if size < 3 {
		return nil, fmt.Errorf("board size must be at least 3x3")
	}

	g := &Board{
		Size: size,
		Data: make([][]int32, size),
	}

	for i := 0; i < size; i++ {
		g.Data[i] = make([]int32, size)
	}

	return g, nil
}

func (g *Board) Put(row, col int, player *Player) error {
	if g.IsOutOfBounds(row, col) {
		return fmt.Errorf("invalid position (%d, %d)", row, col)
	}

	if g.IsOccupied(row, col) {
		return fmt.Errorf("position (%d, %d) is already occupied", row, col)
	}

	g.Data[row][col] = player.ID
	return nil
}

func (g *Board) GetSize() int {
	return g.Size
}

func (g *Board) CheckWin(row, col int, player *Player, maxLine int) bool {
	maxLine = min(maxLine, g.Size)

	// horizontal, vertical and diagonal
	directions := [...][2]int{{0, 1}, {1, 0}, {1, 1}, {1, -1}}

	for _, dir := range directions {
		dr, dc := dir[0], dir[1]
		count := 1

		// Check in both directions
		for _, step := range []int{-1, 1} {
			r, c := row+step*dr, col+step*dc
			for !g.IsOutOfBounds(r, c) && g.IsOccupied(r, c, player) {
				count++
				r += step * dr
				c += step * dc
			}
		}

		if count >= maxLine {
			return true
		}
	}

	return false
}

func (g *Board) IsOutOfBounds(row, col int) bool {
	return row < 0 || row >= g.Size || col < 0 || col >= g.Size
}

func (g *Board) IsOccupied(row, col int, player ...*Player) bool {
	if g.IsOutOfBounds(row, col) {
		return false
	}

	if len(player) > 0 && player[0] != nil {
		return g.Data[row][col] == player[0].ID
	}

	return g.Data[row][col] != 0
}

// min is a helper function for CheckWin
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
