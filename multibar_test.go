package multibar

import (
	"testing"

	"github.com/sethgrid/curse"
)

func TestNewLineDetection_EmptyPrintln(t *testing.T) {
	bc := &BarContainer{}
	bc.Println()
	expected := 1
	resetLines(expected)

	if bc.historicNewlinesCounter != expected {
		t.Errorf("empty Println - got %d, want %d newline", bc.historicNewlinesCounter, expected)
	}
}

func TestNewLineDetection_Println(t *testing.T) {
	bc := &BarContainer{}
	bc.Println("Look, ma!\n I printed somethin'\n")
	expected := 3
	resetLines(expected)

	if bc.historicNewlinesCounter != expected {
		t.Errorf("Println - got %d, want %d newline", bc.historicNewlinesCounter, expected)
	}
}

func TestNewLineDetection_Printf(t *testing.T) {
	bc := &BarContainer{}
	bc.Printf("Look, ma!\n I printed somethin'\n %s %s", "and\n", "\nI can keep printing")
	expected := 4
	resetLines(expected)

	if bc.historicNewlinesCounter != expected {
		t.Errorf("Printf - got %d, want %d newline", bc.historicNewlinesCounter, expected)
	}
}

func TestNewLineDetection_Print(t *testing.T) {
	bc := &BarContainer{}
	bc.Print("Look, ma!\n I printed somethin'\n", "and\n", "\nI can keep printing")
	expected := 4
	resetLines(expected)

	if bc.historicNewlinesCounter != expected {
		t.Errorf("Print - got %d, want %d newline", bc.historicNewlinesCounter, expected)
	}
}

func TestNewLineDetection_ManyPrints(t *testing.T) {
	bc := &BarContainer{}
	bc.Print("Look, ma!\n I printed somethin'\n", "and\n", "\nI can keep printing")
	bc.Println("and more!")
	expected := 5
	resetLines(expected)

	if bc.historicNewlinesCounter != expected {
		t.Errorf("Print - got %d, want %d newline", bc.historicNewlinesCounter, expected)
	}
}

// resetLines erases the print output and resets the cursor to the pre-print line
func resetLines(n int) {
	c := &curse.Cursor{}
	for i := 0; i < n; i++ {
		c.EraseCurrentLine()
		c.MoveUp(1)
	}
	c.EraseCurrentLine()
}
