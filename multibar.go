package multibar

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sethgrid/curse"
)

type ProgressFunc func(progress int)

type BarContainer struct {
	Bars []*ProgressBar

	screenLines             int
	screenWidth             int
	startingLine            int
	totalNewlines           int
	historicNewlinesCounter int

	history map[int]string
	sync.Mutex
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

	progressChan chan int
}

func New() (*BarContainer, error) {
	// can swallow err because sensible defaults are returned from curse
	width, lines, _ := curse.GetScreenDimensions()
	_, line, _ := curse.GetCursorPosition()

	history := make(map[int]string)

	b := &BarContainer{screenWidth: width, screenLines: lines, startingLine: line, history: history}
	// todo: need to figure out a way to deal with additional progressbars while the listener
	// is listening. for the time being, the calling app will have to call listen after
	// all bars are declared
	//go b.Listen()
	return b, nil
}

func (b *BarContainer) Listen() {
	for len(b.Bars) == 0 {
		// wait until we have some bars to work with
		time.Sleep(time.Millisecond * 100)
	}
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
	b.Println()
}

func (b *BarContainer) MakeBar(total int, prepend string) ProgressFunc {
	ch := make(chan int)
	bar := &ProgressBar{
		Width:           b.screenWidth - len(prepend) - 20,
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
	bar.Line = b.startingLine + b.totalNewlines
	b.history[bar.Line] = ""
	bar.Update(0)
	b.Println()

	return func(progress int) { bar.progressChan <- progress }
}

func (p *ProgressBar) AddPrepend(str string) {
	p.Prepend = str
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

	// record where we are, jump to the progress bar, update it, jump back
	c, _ := curse.New()
	c.Move(1, p.Line)
	c.EraseCurrentLine()
	fmt.Printf("\r%s %s%c%s%c%s", p.Prepend, percent, p.LeftEnd, strings.Join(bar, ""), p.RightEnd, timeElapsed)
	c.Move(c.StartingPosition.X, c.StartingPosition.Y)
}

func prettyTime(t time.Duration) string {
	re, err := regexp.Compile(`(\d+).(\d+)(\w+)`)
	if err != nil {
		return err.Error()
	}
	parts := re.FindSubmatch([]byte(t.String()))
	if len(parts) != 4 {
		return "---"
	}
	return string(parts[1]) + string(parts[3])
}

func (b *BarContainer) addedNewlines(count int) {
	b.totalNewlines += count
	b.historicNewlinesCounter += count

	// if we hit the bottom of the screen, we "scroll" our bar displays by pushing
	// them up count lines (closer to line 0)
	if b.startingLine+b.totalNewlines > b.screenLines {
		b.totalNewlines -= count
		for _, bar := range b.Bars {
			bar.Line -= count
		}
		b.redrawAll(count)
	}
}

func (b *BarContainer) redrawAll(moveUp int) {
	c, _ := curse.New()

	newHistory := make(map[int]string)
	for line, printed := range b.history {
		newHistory[line+moveUp] = printed
		c.Move(1, line)
		c.EraseCurrentLine()
		c.Move(1, line+moveUp)
		c.EraseCurrentLine()
		fmt.Print(printed)
	}
	b.history = newHistory
	c.Move(c.StartingPosition.X, c.StartingPosition.Y)
}

// print wrappers to capture newlines to adjust line positions on bars

func (b *BarContainer) Print(a ...interface{}) (n int, err error) {
	b.Lock()
	defer b.Unlock()
	newlines := countAllNewlines(a...)
	b.addedNewlines(newlines)
	thisLine := b.startingLine + b.totalNewlines
	b.history[thisLine] = fmt.Sprint(a...)
	return fmt.Print(a...)
}

func (b *BarContainer) Printf(format string, a ...interface{}) (n int, err error) {
	b.Lock()
	defer b.Unlock()
	newlines := strings.Count(format, "\n")
	newlines += countAllNewlines(a...)
	b.addedNewlines(newlines)
	thisLine := b.startingLine + b.totalNewlines
	b.history[thisLine] = fmt.Sprintf(format, a...)
	return fmt.Printf(format, a...)
}

func (b *BarContainer) Println(a ...interface{}) (n int, err error) {
	b.Lock()
	defer b.Unlock()
	newlines := countAllNewlines(a...) + 1
	b.addedNewlines(newlines)
	thisLine := b.startingLine + b.totalNewlines
	b.history[thisLine] = fmt.Sprint(a...)
	return fmt.Println(a...)
}

func countAllNewlines(interfaces ...interface{}) int {
	count := 0
	for _, iface := range interfaces {
		switch s := iface.(type) {
		case string:
			count += strings.Count(s, "\n")
		}
	}
	return count
}
