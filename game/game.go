package game

import "fmt"

type Game struct {
	Board    [3][3]string
	Winner   string
	Turn     string
	Finished bool
	Count    int
}

func (g *Game) PrintBoard() {
	fmt.Println("\nCurrent Board:")
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			cell := g.Board[i][j]
			if cell == "" {
				cell = " "
			}
			fmt.Printf(" %s ", cell)
			if j < 2 {
				fmt.Print("|")
			}
		}
		if i < 2 {
			fmt.Println("\n-----------")
		}
	}
	fmt.Println("\n")
}

func (g *Game) CheckWinner() bool {
	board := g.Board
	
	for i := 0; i < 3; i++ {
		if board[i][0] != "" && board[i][0] == board[i][1] && board[i][1] == board[i][2] {
			g.Winner = board[i][0]
			g.Finished = true
			return true
		}
		if board[0][i] != "" && board[0][i] == board[1][i] && board[1][i] == board[2][i] {
			g.Winner = board[0][i]
			g.Finished = true
			return true
		}
		if board[0][0] != "" && board[0][0] == board[1][1] && board[1][1] == board[2][2] {
			g.Winner = board[0][0]
			g.Finished = true
			return true
		}
		if board[0][2] != "" && board[0][2] == board[1][1] && board[1][1] == board[2][0] {
			g.Winner = board[0][2]
			g.Finished = true
			return true
		}

	}
	return false
}

func (g *Game) CheckDraw() bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if g.Board[i][j] == "" {
				return false
			}
		}
	}
	if g.Winner == "" {
		g.Finished = true
		return true
	}
	return false
}
