package game

import (
	"strings"
	"testing"
)

func TestBotBlocksImmediateWin(t *testing.T) {
	result, err := PlayMove([]string{"X", "X", "", "", "O", "", "O", "", ""}, 8, "X")
	if err != nil {
		t.Fatal(err)
	}
	if result.BotPosition != 2 || result.Strategy != "MINIMAX" || result.NodesEvaluated == 0 {
		t.Fatalf("unexpected decision: %+v", result)
	}
}

func TestHumanWinReturnsImmediately(t *testing.T) {
	result, err := PlayMove([]string{"X", "X", "", "O", "O", "", "", "", ""}, 2, "X")
	if err != nil {
		t.Fatal(err)
	}
	if result.Outcome != "X" || result.BotPosition != -1 {
		t.Fatalf("unexpected decision: %+v", result)
	}
}

func TestOccupiedCellIsRejected(t *testing.T) {
	_, err := PlayMove([]string{"X", "", "", "", "", "", "", "", ""}, 0, "O")
	if err == nil || !strings.Contains(err.Error(), "occupied") {
		t.Fatalf("expected occupied error, got %v", err)
	}
}

func TestWinnerDetectsDraw(t *testing.T) {
	if got := Winner([]string{"X", "O", "X", "X", "O", "O", "O", "X", "X"}); got != Draw {
		t.Fatalf("got %q, want %q", got, Draw)
	}
}
