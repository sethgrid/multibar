package multibar

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sethgrid/curse"
)

type progressFunc func(progress int)

type BarContainer struct {
	Bars []*ProgressBar
}

type ProgressBar struct {
	Width           int
	Total           int
	LeftEnd         byte
	RightEnd        byte
	Fill            byte
	Head            byte
	Empty           byte
	ShowPercent     bool
	ShowTimeElapsed bool
	StartTime       time.Time
	Line            int
	Prepend         string
	progressChan    chan int
}

func New() (*BarContainer, error) {
	return &BarContainer{}, nil
}

func (b *BarContainer) Listen() {
	cases := make([]reflect.SelectCase, len(b.Bars))
	for i, bar := range b.Bars {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(bar.progressChan)}
	}

	remaining := len(cases)
	for remaining > 0 {
		chosen, value, ok := reflect.Select(cases)
		if !ok {
			// The chosen channel has been closed, so zero out the channel to disable the case
			cases[chosen].Chan = reflect.ValueOf(nil)
			remaining -= 1
			continue
		}

		b.Bars[chosen].Update(int(value.Int()))
	}
	fmt.Println()
}

func (b *BarContainer) MakeBar(total int, prepend string) progressFunc {
	// can swallow err because sensible defaults are returned
	fmt.Println("\n")
	width, _, _ := curse.GetScreenDimensions()
	ch := make(chan int)
	bar := &ProgressBar{
		Width:           (width - len(prepend)) * 3 / 5,
		Total:           total,
		Prepend:         prepend,
		LeftEnd:         '[',
		RightEnd:        ']',
		Fill:            '=',
		Head:            '>',
		Empty:           '-',
		ShowPercent:     true,
		ShowTimeElapsed: true,
		StartTime:       time.Now(),
		progressChan:    ch,
	}

	b.Bars = append(b.Bars, bar)
	bar.Display()
	return func(progress int) { bar.progressChan <- progress }
}

func (p *ProgressBar) AddPrepend(str string) {
	width, _, _ := curse.GetScreenDimensions()
	p.Prepend = str
	p.Width = (width - len(str)) * 3 / 5
}

func (p *ProgressBar) Display() {
	_, line, _ := curse.GetCursorPosition()
	p.Line = line
	p.Update(0)
}

func (p *ProgressBar) Update(progress int) {
	bar := make([]string, p.Width)

	// avoid division by zero errors on non-properly constructed progressbars
	if p.Width == 0 {
		p.Width = 1
	}
	if p.Total == 0 {
		p.Total = 1
	}
	justGotToFirstEmptySpace := true
	for i, _ := range bar {
		if float32(progress)/float32(p.Total) > float32(i)/float32(p.Width) {
			bar[i] = string(p.Fill)
		} else {
			bar[i] = string(p.Empty)
			if justGotToFirstEmptySpace {
				bar[i] = string(p.Head)
				justGotToFirstEmptySpace = false
			}
		}
	}

	percent := ""
	if p.ShowPercent {
		asInt := int(100 * (float32(progress) / float32(p.Total)))
		padding := ""
		if asInt < 10 {
			padding = "  "
		} else if asInt < 99 {
			padding = " "
		}
		percent = padding + strconv.Itoa(asInt) + "% "
	}

	timeElapsed := ""
	if p.ShowTimeElapsed {
		timeElapsed = " " + prettyTime(time.Since(p.StartTime))
	}
	currentRow, currentLine, _ := curse.GetCursorPosition()
	c := &curse.Cursor{}
	c.Move(1, p.Line)
	c.EraseCurrentLine()
	fmt.Printf("\r%s%s%c%s%c%s", p.Prepend, percent, p.LeftEnd, strings.Join(bar, ""), p.RightEnd, timeElapsed)
	c.Move(currentRow, currentLine)
}

func prettyTime(t time.Duration) string {
	re, err := regexp.Compile(`(\d+).(\d+)(\w+)`)
	if err != nil {
		return err.Error()
	}
	parts := re.FindSubmatch([]byte(t.String()))

	return string(parts[1]) + string(parts[3])
}
