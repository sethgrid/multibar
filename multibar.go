package multibar

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sethgrid/curse"
)

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
	prepend         string
}

func New(total int) (*ProgressBar, error) {
	// can swallow err because sensible defaults are returned
	width, _, _ := curse.GetScreenDimensions()

	bar := &ProgressBar{
		Width:           width * 3 / 5,
		Total:           total,
		LeftEnd:         '[',
		RightEnd:        ']',
		Fill:            '=',
		Head:            '>',
		Empty:           '-',
		ShowPercent:     true,
		ShowTimeElapsed: true,
		StartTime:       time.Now(),
	}
	return bar, nil
}

func (p *ProgressBar) Prepend(str string) {
	p.prepend = str
}

func (p *ProgressBar) Display(progress int) {
	/*
	   notes:
	   consider a prepend string
	   handle show time and percent
	*/
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
	c := &curse.Cursor{}
	c.EraseCurrentLine()
	fmt.Printf("\r%s%s%c%s%c%s", p.prepend, percent, p.LeftEnd, strings.Join(bar, ""), p.RightEnd, timeElapsed)
}

func prettyTime(t time.Duration) string {
	re, err := regexp.Compile(`(\d+).(\d+)(\w+)`)
	if err != nil {
		return err.Error()
	}
	parts := re.FindSubmatch([]byte(t.String()))

	return string(parts[1]) + string(parts[3])
}
