package main

import (
	"testing"
)

func TestNewGame(t *testing.T) {
	game := NewGame()
	
	if game == nil {
		t.Fatal("NewGame should return a non-nil game")
	}
	
	// Count non-zero tiles
	tileCount := 0
	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
			if game.board[i][j] != 0 {
				tileCount++
			}
		}
	}
	
	if tileCount != 2 {
		t.Errorf("New game should have exactly 2 tiles, got %d", tileCount)
	}
	
	if game.score != 0 {
		t.Errorf("New game should have score 0, got %d", game.score)
	}
	
	if game.won || game.lost {
		t.Error("New game should not be won or lost")
	}
}

func TestMoveLeft(t *testing.T) {
	game := &Game{}
	
	// Test case 1: Simple move left
	game.board = [SIZE][SIZE]int{
		{0, 2, 0, 2},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}
	
	moved := game.moveLeft()
	if !moved {
		t.Error("Should have moved")
	}
	
	expected := [SIZE][SIZE]int{
		{4, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}
	
	if game.board != expected {
		t.Errorf("Board mismatch after moveLeft.\nExpected: %v\nGot: %v", expected, game.board)
	}
	
	if game.score != 4 {
		t.Errorf("Score should be 4, got %d", game.score)
	}
}

func TestMoveLeftNoMerge(t *testing.T) {
	game := &Game{}
	
	// Test case: Move without merge
	game.board = [SIZE][SIZE]int{
		{0, 2, 0, 4},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}
	
	moved := game.moveLeft()
	if !moved {
		t.Error("Should have moved")
	}
	
	expected := [SIZE][SIZE]int{
		{2, 4, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}
	
	if game.board != expected {
		t.Errorf("Board mismatch after moveLeft.\nExpected: %v\nGot: %v", expected, game.board)
	}
	
	if game.score != 0 {
		t.Errorf("Score should be 0, got %d", game.score)
	}
}

func TestMoveLeftMultipleMerges(t *testing.T) {
	game := &Game{}
	
	// Test case: Multiple merges in one row
	game.board = [SIZE][SIZE]int{
		{2, 2, 4, 4},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}
	
	moved := game.moveLeft()
	if !moved {
		t.Error("Should have moved")
	}
	
	expected := [SIZE][SIZE]int{
		{4, 8, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}
	
	if game.board != expected {
		t.Errorf("Board mismatch after moveLeft.\nExpected: %v\nGot: %v", expected, game.board)
	}
	
	if game.score != 12 {
		t.Errorf("Score should be 12 (4+8), got %d", game.score)
	}
}

func TestMoveRight(t *testing.T) {
	game := &Game{}
	
	game.board = [SIZE][SIZE]int{
		{2, 0, 0, 2},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}
	
	moved := game.moveRight()
	if !moved {
		t.Error("Should have moved")
	}
	
	expected := [SIZE][SIZE]int{
		{0, 0, 0, 4},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}
	
	if game.board != expected {
		t.Errorf("Board mismatch after moveRight.\nExpected: %v\nGot: %v", expected, game.board)
	}
}

func TestMoveUp(t *testing.T) {
	game := &Game{}
	
	game.board = [SIZE][SIZE]int{
		{2, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{2, 0, 0, 0},
	}
	
	moved := game.moveUp()
	if !moved {
		t.Error("Should have moved")
	}
	
	expected := [SIZE][SIZE]int{
		{4, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}
	
	if game.board != expected {
		t.Errorf("Board mismatch after moveUp.\nExpected: %v\nGot: %v", expected, game.board)
	}
}

func TestMoveDown(t *testing.T) {
	game := &Game{}
	
	game.board = [SIZE][SIZE]int{
		{2, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{2, 0, 0, 0},
	}
	
	moved := game.moveDown()
	if !moved {
		t.Error("Should have moved")
	}
	
	expected := [SIZE][SIZE]int{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{4, 0, 0, 0},
	}
	
	if game.board != expected {
		t.Errorf("Board mismatch after moveDown.\nExpected: %v\nGot: %v", expected, game.board)
	}
}

func TestCanMove(t *testing.T) {
	game := &Game{}
	
	// Test case: Board with empty cells
	game.board = [SIZE][SIZE]int{
		{2, 4, 8, 16},
		{32, 64, 128, 256},
		{512, 1024, 2, 0},
		{4, 8, 16, 32},
	}
	
	if !game.canMove() {
		t.Error("Should be able to move (has empty cell)")
	}
	
	// Test case: Board with possible merge
	game.board = [SIZE][SIZE]int{
		{2, 4, 8, 16},
		{32, 64, 128, 256},
		{512, 1024, 2, 4},
		{4, 8, 16, 16},
	}
	
	if !game.canMove() {
		t.Error("Should be able to move (has possible merge)")
	}
	
	// Test case: No moves possible
	game.board = [SIZE][SIZE]int{
		{2, 4, 8, 16},
		{32, 64, 128, 256},
		{512, 1024, 2, 4},
		{8, 16, 32, 64},
	}
	
	if game.canMove() {
		t.Error("Should not be able to move")
	}
}

func TestWinCondition(t *testing.T) {
	game := &Game{}
	
	// Set up a board where merging creates 2048
	game.board = [SIZE][SIZE]int{
		{1024, 1024, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}
	
	game.moveLeft()
	
	if !game.won {
		t.Error("Game should be won after creating 2048")
	}
	
	if game.score != 2048 {
		t.Errorf("Score should be 2048, got %d", game.score)
	}
}

func TestTranspose(t *testing.T) {
	game := &Game{}
	
	game.board = [SIZE][SIZE]int{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12},
		{13, 14, 15, 16},
	}
	
	game.transpose()
	
	expected := [SIZE][SIZE]int{
		{1, 5, 9, 13},
		{2, 6, 10, 14},
		{3, 7, 11, 15},
		{4, 8, 12, 16},
	}
	
	if game.board != expected {
		t.Errorf("Board mismatch after transpose.\nExpected: %v\nGot: %v", expected, game.board)
	}
}

func TestReverseRows(t *testing.T) {
	game := &Game{}
	
	game.board = [SIZE][SIZE]int{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12},
		{13, 14, 15, 16},
	}
	
	game.reverseRows()
	
	expected := [SIZE][SIZE]int{
		{4, 3, 2, 1},
		{8, 7, 6, 5},
		{12, 11, 10, 9},
		{16, 15, 14, 13},
	}
	
	if game.board != expected {
		t.Errorf("Board mismatch after reverseRows.\nExpected: %v\nGot: %v", expected, game.board)
	}
}

func TestNoMoveWhenAlreadyWon(t *testing.T) {
	game := &Game{won: true}
	
	game.board = [SIZE][SIZE]int{
		{2, 0, 0, 2},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}
	
	originalBoard := game.board
	game.move("a")
	
	if game.board != originalBoard {
		t.Error("Board should not change when game is already won")
	}
}

func TestNoMoveWhenAlreadyLost(t *testing.T) {
	game := &Game{lost: true}
	
	game.board = [SIZE][SIZE]int{
		{2, 0, 0, 2},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}
	
	originalBoard := game.board
	game.move("a")
	
	if game.board != originalBoard {
		t.Error("Board should not change when game is already lost")
	}
}