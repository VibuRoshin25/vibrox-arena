package game

import (
	"errors"
	"fmt"
	"time"
)

const (
	Empty      = ""
	InProgress = "IN_PROGRESS"
	Draw       = "DRAW"
)

var winningLines = [][3]int{
	{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
	{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
	{0, 4, 8}, {2, 4, 6},
}

type Decision struct {
	Board              []string
	BotPosition        int32
	Outcome            string
	Strategy           string
	Score              int32
	NodesEvaluated     int32
	SearchDepth        int32
	DecisionTimeMicros int64
}

type metrics struct {
	nodes    int32
	maxDepth int32
}

func PlayMove(input []string, position int32, humanMark string) (Decision, error) {
	started := time.Now()
	board := append([]string(nil), input...)
	if humanMark == "" {
		humanMark = "X"
	}
	if err := validateBoard(board); err != nil {
		return Decision{}, err
	}
	if humanMark != "X" && humanMark != "O" {
		return Decision{}, errors.New("human_mark must be X or O")
	}

	xCount, oCount := markCounts(board)
	if (humanMark == "X" && xCount != oCount) || (humanMark == "O" && xCount != oCount+1) {
		return Decision{}, fmt.Errorf("it is not %s's turn", humanMark)
	}
	if position < 0 || position > 8 {
		return Decision{}, errors.New("position must be between 0 and 8")
	}
	if board[position] != Empty {
		return Decision{}, errors.New("position is already occupied")
	}

	botMark := "O"
	if humanMark == "O" {
		botMark = "X"
	}
	board[position] = humanMark
	outcome := Winner(board)
	decision := Decision{Board: board, BotPosition: -1, Outcome: outcome, Strategy: "MINIMAX"}
	if outcome == InProgress {
		move, score, stats := chooseBotMove(board, botMark, humanMark)
		board[move] = botMark
		decision.BotPosition = int32(move)
		decision.Score = score
		decision.NodesEvaluated = stats.nodes
		decision.SearchDepth = stats.maxDepth
		decision.Outcome = Winner(board)
	}
	decision.DecisionTimeMicros = time.Since(started).Microseconds()
	return decision, nil
}

func Winner(board []string) string {
	for _, line := range winningLines {
		if board[line[0]] != Empty && board[line[0]] == board[line[1]] && board[line[0]] == board[line[2]] {
			return board[line[0]]
		}
	}
	for _, cell := range board {
		if cell == Empty {
			return InProgress
		}
	}
	return Draw
}

func validateBoard(board []string) error {
	if len(board) != 9 {
		return errors.New("board must contain exactly 9 cells")
	}
	for _, cell := range board {
		if cell != Empty && cell != "X" && cell != "O" {
			return errors.New("board cells must be X, O, or empty")
		}
	}
	xCount, oCount := markCounts(board)
	if xCount < oCount || xCount > oCount+1 {
		return errors.New("board has impossible move counts")
	}
	if Winner(board) != InProgress {
		return errors.New("game is already complete")
	}
	return nil
}

func markCounts(board []string) (int, int) {
	var xCount, oCount int
	for _, cell := range board {
		switch cell {
		case "X":
			xCount++
		case "O":
			oCount++
		}
	}
	return xCount, oCount
}

func chooseBotMove(board []string, botMark, humanMark string) (int, int32, *metrics) {
	stats := &metrics{}
	bestMove, bestScore := -1, int32(-100)
	for position, cell := range board {
		if cell != Empty {
			continue
		}
		board[position] = botMark
		score := minimax(board, botMark, humanMark, false, 1, stats)
		board[position] = Empty
		if score > bestScore {
			bestMove, bestScore = position, score
		}
	}
	return bestMove, bestScore, stats
}

func minimax(board []string, botMark, humanMark string, maximizing bool, depth int32, stats *metrics) int32 {
	stats.nodes++
	if depth > stats.maxDepth {
		stats.maxDepth = depth
	}
	switch Winner(board) {
	case botMark:
		return 10 - depth
	case humanMark:
		return depth - 10
	case Draw:
		return 0
	}

	best := int32(100)
	mark := humanMark
	if maximizing {
		best, mark = -100, botMark
	}
	for position, cell := range board {
		if cell != Empty {
			continue
		}
		board[position] = mark
		score := minimax(board, botMark, humanMark, !maximizing, depth+1, stats)
		board[position] = Empty
		if maximizing && score > best || !maximizing && score < best {
			best = score
		}
	}
	return best
}
