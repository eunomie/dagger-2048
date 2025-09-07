package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	SIZE = 4
)

type Game struct {
	board [SIZE][SIZE]int
	score int
	won   bool
	lost  bool
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
	fmt.Printf("Score: %d\n\n", g.score)
	
	fmt.Println("+------+------+------+------+")
	for i := 0; i < SIZE; i++ {
		fmt.Print("|")
		for j := 0; j < SIZE; j++ {
			if g.board[i][j] == 0 {
				fmt.Print("      |")
			} else {
				fmt.Printf(" %4d |", g.board[i][j])
			}
		}
		fmt.Println()
		fmt.Println("+------+------+------+------+")
	}
	
	fmt.Println("\nControls: w(up), s(down), a(left), d(right), q(quit)")
	
	if g.won {
		fmt.Println("🎉 You won! Congratulations!")
	} else if g.lost {
		fmt.Println("💀 Game Over! No more moves available.")
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

func main() {
	rand.Seed(time.Now().UnixNano())
	game := NewGame()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to 2048!")
	fmt.Println("Combine tiles with the same number to reach 2048!")
	fmt.Print("Press Enter to start...")
	scanner.Scan()

	for {
		game.display()
		
		if game.won || game.lost {
			fmt.Print("Play again? (y/n): ")
			scanner.Scan()
			input := strings.ToLower(strings.TrimSpace(scanner.Text()))
			if input == "y" || input == "yes" {
				game = NewGame()
				continue
			} else {
				break
			}
		}

		fmt.Print("Enter move: ")
		scanner.Scan()
		input := strings.ToLower(strings.TrimSpace(scanner.Text()))
		
		if input == "q" || input == "quit" {
			break
		}
		
		if input == "w" || input == "s" || input == "a" || input == "d" {
			game.move(input)
		} else {
			fmt.Println("Invalid input! Use w/a/s/d for movement or q to quit.")
			time.Sleep(1 * time.Second)
		}
	}
	
	fmt.Println("Thanks for playing!")
}
