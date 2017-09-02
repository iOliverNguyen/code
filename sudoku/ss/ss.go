package ss

import (
	"bytes"
	"errors"
	"fmt"
)

type Board [9][9]int

func (b Board) String() string {
	var buf bytes.Buffer
	for y := 0; y < 9; y++ {
		if y == 3 || y == 6 {
			buf.WriteRune('\n')
		}
		for x := 0; x < 9; x++ {
			if x%3 == 0 {
				buf.WriteRune(' ')
			}
			buf.WriteRune('0' + rune(b[y][x]))
		}
		buf.WriteRune('\n')
	}
	return buf.String()
}

/*
Parse parses input string and return a Board

  000 507 000
  002 406 300
  090 010 020

  270 000 068
  003 000 100
  140 000 093

  060 040 050
  009 205 600
  000 903 000

*/
func Parse(input string) (Board, error) {
	var b Board
	c := 0
	for _, ch := range input {
		if ch >= '0' && ch <= '9' {
			if c >= 9*9 {
				return b, errors.New(fmt.Sprintf("Unexpected numbers of value: %v", c))
			}
			b[c/9][c%9] = int(ch - '0')
			c++
		}
	}
	if c != 9*9 {
		return b, errors.New(fmt.Sprintf("Unexpected numbers of value: %v", c))
	}
	return b, nil
}

// IsValid checks that the provided Board is a valid sudoku
func isValid(b Board) error {
	var flags [10][3]bool

	for k := 0; k < 9; k++ {
		for i := 1; i < 10; i++ {
			flags[i][0] = false
			flags[i][1] = false
			flags[i][2] = false
		}
		for h := 0; h < 9; h++ {
			// Check in row
			v := b[k][h]
			if v != 0 && flags[v][0] {
				return errors.New(fmt.Sprintf("Duplicated cell: (%v, %v)", k, h))
			}
			flags[v][0] = true

			// Check in column
			v = b[h][k]
			if v != 0 && flags[v][1] {
				return errors.New(fmt.Sprintf("Duplicated cell: (%v, %v)", h, k))
			}
			flags[v][1] = true

			// Check in 3x3 square
			y := k/3*3 + h/3
			x := k%3*3 + h%3
			v = b[y][x]
			if v != 0 && flags[v][2] {
				return errors.New(fmt.Sprintf("Duplicated cell: (%v, %v)", y, x))
			}
			flags[v][2] = true
		}
	}
	return nil
}

func IsValid(b Board) error {
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			v := b[y][x]
			if v < 0 && v > 9 {
				return errors.New(fmt.Sprintf("Invalid cell: (%v, %v)", y, x))
			}
		}
	}
	return isValid(b)
}

func IsCompleted(b Board) error {
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			v := b[y][x]
			if v < 0 && v > 9 {
				return errors.New(fmt.Sprintf("Invalid cell: (%v, %v)", y, x))
			}
			if v == 0 {
				return errors.New(fmt.Sprintf("Incompleted cell: (%v, %v)", y, x))
			}
		}
	}
	return isValid(b)
}

func FindOne(b Board) (Board, bool) {
	var flags, minFlags [10]bool
	var minX, minY, minCount int

	for {
		minX = -1
		minY = -1
		minCount = 10

		for y := 0; y < 9; y++ {
			for x := 0; x < 9; x++ {
				if b[y][x] != 0 {
					continue
				}
				for i := 1; i < 10; i++ {
					flags[i] = false
				}

				y0 := y / 3 * 3
				x0 := x / 3 * 3
				for i := 0; i < 9; i++ {
					// Check in row
					v := b[y][i]
					flags[v] = true

					// Check in column
					v = b[i][x]
					flags[v] = true

					// Check in 3x3 square
					y1 := y0 + i/3
					x1 := x0 + i%3
					v = b[y1][x1]
					flags[v] = true
				}

				c := 0
				v := 0
				for i := 1; i < 10; i++ {
					if !flags[i] {
						c++
						v = i
					}
				}
				if c == 0 {
					// No result
					return Board{}, false
				}
				if c == 1 {
					// Fill the cell
					b[y][x] = v
					minCount = 1

				} else if c < minCount {
					// No cell was filled, we store least posibility cell for guessing
					minX = x
					minY = y
					minCount = c
					minFlags = flags
				}
			}
		}

		if minCount > 1 {
			// No cell was filled in last round, let's guess
			break
		}
	}

	if IsCompleted(b) == nil {
		return b, true
	}

	// Recursive
	for i := 1; i < 10; i++ {
		if !minFlags[i] {
			b[minY][minX] = i
			if result, ok := FindOne(b); ok {
				return result, ok
			}
			b[minY][minX] = 0
		}
	}
	return Board{}, false
}
