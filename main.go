package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"

	"golang.org/x/term"
)

const (
	SIZE = 4
)

// ANSI color codes for terminal colors
const (
	ColorReset = "\033[0m"
	ColorBold  = "\033[1m"

	// Improved tile colors for better readability
	Color2    = "\033[48;5;253;38;5;0;1m"   // Light gray background, bold black text
	Color4    = "\033[48;5;250;38;5;0;1m"   // Medium gray background, bold black text
	Color8    = "\033[48;5;214;38;5;0;1m"   // Orange background, bold black text
	Color16   = "\033[48;5;209;38;5;255;1m" // Dark orange background, bold white text
	Color32   = "\033[48;5;196;38;5;255;1m" // Red background, bold white text
	Color64   = "\033[48;5;202;38;5;255;1m" // Red-orange background, bold white text
	Color128  = "\033[48;5;226;38;5;0;1m"   // Bright yellow background, bold black text
	Color256  = "\033[48;5;220;38;5;0;1m"   // Gold background, bold black text
	Color512  = "\033[48;5;208;38;5;0;1m"   // Dark gold background, bold black text
	Color1024 = "\033[48;5;166;38;5;255;1m" // Dark orange background, bold white text
	Color2048 = "\033[48;5;196;38;5;255;1m" // Bright red background, bold white text
	ColorHigh = "\033[48;5;93;38;5;255;1m"  // Purple background, bold white text
)

type Game struct {
	board [SIZE][SIZE]int
	score int
	won   bool
	lost  bool
}

func getTileColor(value int) string {
	switch value {
	case 2:
		return Color2
	case 4:
		return Color4
	case 8:
		return Color8
	case 16:
		return Color16
	case 32:
		return Color32
	case 64:
		return Color64
	case 128:
		return Color128
	case 256:
		return Color256
	case 512:
		return Color512
	case 1024:
		return Color1024
	case 2048:
		return Color2048
	default:
		return ColorHigh
	}
}

func NewGame() *Game {
	g := &Game{}
	g.addRandomTile()
	g.addRandomTile()
	return g
}

func (g *Game) addRandomTile() {
	emptyCells := make([][2]int, 0)
	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
			if g.board[i][j] == 0 {
				emptyCells = append(emptyCells, [2]int{i, j})
			}
		}
	}

	if len(emptyCells) == 0 {
		return
	}

	randomCell := emptyCells[rand.Intn(len(emptyCells))]
	value := 2
	if rand.Float32() < 0.1 {
		value = 4
	}
	g.board[randomCell[0]][randomCell[1]] = value
}

func (g *Game) display() {
	clearScreen()
	fmt.Printf("%sScore: %d%s\r\n\r\n", ColorBold, g.score, ColorReset)

	// Top border
	fmt.Print("┌──────┬──────┬──────┬──────┐\r\n")

	for i := 0; i < SIZE; i++ {
		// Row content
		fmt.Print("│")
		for j := 0; j < SIZE; j++ {
			if g.board[i][j] == 0 {
				fmt.Print("      │")
			} else {
				color := getTileColor(g.board[i][j])
				fmt.Printf("%s %4d %s│", color, g.board[i][j], ColorReset)
			}
		}
		fmt.Print("\r\n")

		// Middle or bottom border
		if i < SIZE-1 {
			fmt.Print("├──────┼──────┼──────┼──────┤\r\n")
		} else {
			fmt.Print("└──────┴──────┴──────┴──────┘\r\n")
		}
	}

	fmt.Print("\r\nControls: Arrow keys or WASD to move, Q to quit\r\n")

	if g.won {
		fmt.Print("🎉 You won! Congratulations!\r\n")
	} else if g.lost {
		fmt.Print("💀 Game Over! No more moves available.\r\n")
	}
}

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (g *Game) moveLeft() bool {
	moved := false
	for i := 0; i < SIZE; i++ {
		row := make([]int, 0)
		for j := 0; j < SIZE; j++ {
			if g.board[i][j] != 0 {
				row = append(row, g.board[i][j])
			}
		}

		for k := 0; k < len(row)-1; k++ {
			if row[k] == row[k+1] {
				row[k] *= 2
				g.score += row[k]
				if row[k] == 2048 {
					g.won = true
				}
				row = append(row[:k+1], row[k+2:]...)
			}
		}

		newRow := make([]int, SIZE)
		copy(newRow, row)

		for j := 0; j < SIZE; j++ {
			if g.board[i][j] != newRow[j] {
				moved = true
			}
			g.board[i][j] = newRow[j]
		}
	}
	return moved
}

func (g *Game) moveRight() bool {
	g.reverseRows()
	moved := g.moveLeft()
	g.reverseRows()
	return moved
}

func (g *Game) moveUp() bool {
	g.transpose()
	moved := g.moveLeft()
	g.transpose()
	return moved
}

func (g *Game) moveDown() bool {
	g.transpose()
	moved := g.moveRight()
	g.transpose()
	return moved
}

func (g *Game) transpose() {
	for i := 0; i < SIZE; i++ {
		for j := i + 1; j < SIZE; j++ {
			g.board[i][j], g.board[j][i] = g.board[j][i], g.board[i][j]
		}
	}
}

func (g *Game) reverseRows() {
	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE/2; j++ {
			g.board[i][j], g.board[i][SIZE-1-j] = g.board[i][SIZE-1-j], g.board[i][j]
		}
	}
}

func (g *Game) canMove() bool {
	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
			if g.board[i][j] == 0 {
				return true
			}

			if i < SIZE-1 && g.board[i][j] == g.board[i+1][j] {
				return true
			}
			if j < SIZE-1 && g.board[i][j] == g.board[i][j+1] {
				return true
			}
		}
	}
	return false
}

func (g *Game) move(direction string) {
	if g.won || g.lost {
		return
	}

	var moved bool
	switch direction {
	case "w":
		moved = g.moveUp()
	case "s":
		moved = g.moveDown()
	case "a":
		moved = g.moveLeft()
	case "d":
		moved = g.moveRight()
	}

	if moved {
		g.addRandomTile()
		if !g.canMove() {
			g.lost = true
		}
	}
}

// enableRawMode sets terminal to raw mode for single key input
func enableRawMode() (*term.State, error) {
	oldState, err := term.MakeRaw(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	return oldState, nil
}

// disableRawMode restores terminal to normal mode
func disableRawMode(oldState *term.State) error {
	return term.Restore(int(syscall.Stdin), oldState)
}

// readKey reads a single key press
func readKey() (string, error) {
	var buf [3]byte
	n, err := os.Stdin.Read(buf[:])
	if err != nil {
		return "", err
	}

	if n == 1 {
		switch buf[0] {
		case 3: // Ctrl+C
			return "quit", nil
		case 'q', 'Q':
			return "quit", nil
		case 'w', 'W':
			return "up", nil
		case 's', 'S':
			return "down", nil
		case 'a', 'A':
			return "left", nil
		case 'd', 'D':
			return "right", nil
		}
	}

	// Handle arrow keys (escape sequences)
	if n >= 3 && buf[0] == 27 && buf[1] == 91 {
		switch buf[2] {
		case 65: // Up arrow
			return "up", nil
		case 66: // Down arrow
			return "down", nil
		case 67: // Right arrow
			return "right", nil
		case 68: // Left arrow
			return "left", nil
		}
	}

	return "", nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	game := NewGame()

	fmt.Println("Welcome to 2048!")
	fmt.Println("Combine tiles with the same number to reach 2048!")
	fmt.Println("Press any key to start...")

	// Set up terminal for raw input after displaying welcome messages
	oldState, err := enableRawMode()
	if err != nil {
		fmt.Printf("Error enabling raw mode: %v\n", err)
		return
	}
	defer disableRawMode(oldState)

	readKey()

	for {
		game.display()

		if game.won || game.lost {
			fmt.Print("\r\nPlay again? (y/n): ")
			for {
				key, err := readKey()
				if err != nil {
					return
				}
				if key == "quit" {
					return
				}
				if len(key) == 1 && (key[0] == 'n' || key[0] == 'N') {
					fmt.Print("\r\n")
					return
				}
				if len(key) == 1 && (key[0] == 'y' || key[0] == 'Y') {
					fmt.Print("\r\n")
					game = NewGame()
					break
				}
				// Ignore other keys and continue waiting for y/n
			}
			continue
		}

		key, err := readKey()
		if err != nil {
			break
		}

		if key == "quit" {
			break
		}

		switch key {
		case "up":
			game.move("w")
		case "down":
			game.move("s")
		case "left":
			game.move("a")
		case "right":
			game.move("d")
		}
	}

	fmt.Print("\r\nThanks for playing!\r\n")
}
