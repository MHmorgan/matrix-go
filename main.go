package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode"
)

const (
	DrawInterval    = 80 * time.Millisecond // Seconds
	ColumnMinHeight = 20
)

var (
	RandChar = make(chan rune)

	NoLatin    = flag.Bool("no-latin", false, "Disable latin extended characters")
	NoGreek    = flag.Bool("no-greek", false, "Disable greek characters")
	NoCyrillic = flag.Bool("no-cyril", false, "Disable cyrillic characters")
	Help       = flag.Bool("help", false, "Show help")
)

// -----------------------------------------------------------------------------
// Column

type Column struct {
	WindowHeight int
	Offset       int
	Chars        []rune
	Height       int
	Finished     bool
}

func NewColumn(height int) *Column {
	wh := ColumnMinHeight + rand.Intn(height-ColumnMinHeight)

	chars := make([]rune, height)
	for i := range chars {
		chars[i] = ' '
	}

	return &Column{
		WindowHeight: wh,
		Offset:       -wh,
		Chars:        chars,
		Height:       height,
	}
}

func (c *Column) EndOffset() int {
	return c.Offset + c.WindowHeight
}

func (c *Column) Char(i int) string {
	ch := c.Chars[i]

	if i == c.EndOffset() {
		return fmt.Sprintf("%s%s%c%s%s", FgWhite, Bold, ch, Reset, FgDefault)
	}

	if c.EndOffset()-2 <= i && i < c.EndOffset() {
		return fmt.Sprintf("%s%c%s", FgWhite, ch, FgDefault)
	}

	if c.Offset <= i && i <= c.Offset+2 {
		return fmt.Sprintf("%s%s%c%s%s", FgGreen, Dim, ch, Reset, FgDefault)
	}

	return fmt.Sprintf("%s%c%s", FgGreen, ch, FgDefault)
}

func (c *Column) Update() {
	if c.Finished {
		return
	}

	// Remove old top character
	if c.Offset >= 0 {
		c.Chars[c.Offset] = ' '
	}

	// Update offset
	c.Offset++
	if c.Offset >= c.Height {
		c.Finished = true
		return
	}

	// Add new bottom character
	if c.EndOffset() < c.Height {
		c.Chars[c.EndOffset()] = <-RandChar
	}

	// Update a single random character in the window
	idx := c.Offset + rand.Intn(c.WindowHeight)
	if 0 < idx && idx < c.Height {
		c.Chars[idx] = <-RandChar
	}
}

// -----------------------------------------------------------------------------
// Screen

type Screen struct {
	Columns []*Column
	Height  int
	Width   int
}

func content(ch chan<- string, height, width int) {
	screen := NewScreen(height, width)
	for range time.Tick(DrawInterval) {
		screen.Update()
		ch <- screen.String()
	}
}

func NewScreen(height, width int) *Screen {
	return &Screen{
		Columns: make([]*Column, width),
		Height:  height,
		Width:   width,
	}
}

func (s *Screen) Update() {
	idx := rand.Intn(len(s.Columns))
	if s.Columns[idx] == nil {
		s.Columns[idx] = NewColumn(s.Height)
	}

	for i, c := range s.Columns {
		if c == nil {
			continue
		}
		c.Update()
		if c.Finished {
			s.Columns[i] = nil
		}
	}
}

func (s *Screen) String() string {
	sb := strings.Builder{}
	sb.Grow(s.Height * s.Width)

	for i := 0; i < s.Height; i++ {
		if i > 0 {
			sb.WriteByte('\n')
		}
		for j := 0; j < s.Width; j++ {
			if s.Columns[j] == nil {
				sb.WriteByte(' ')
			} else {
				sb.WriteString(s.Columns[j].Char(i))
			}
		}
	}

	return sb.String()
}

// -----------------------------------------------------------------------------

func generateChars() {
	var chars []rune

	// ASCII
	for i := 33; i <= 126; i++ {
		r := rune(i)
		if unicode.IsPrint(r) {
			chars = append(chars, r)
		}
	}

	// Latin
	if !*NoLatin {
		for i := 0x00C0; i <= 0x00FF; i++ {
			r := rune(i)
			if unicode.IsLetter(r) {
				chars = append(chars, r)
			}
		}
	}

	// Greek
	if !*NoGreek {
		for i := 0x0370; i <= 0x03FF; i++ {
			r := rune(i)
			if unicode.IsLetter(r) {
				chars = append(chars, r)
			}
		}
	}

	// Cyrillic
	if !*NoCyrillic {
		for i := 0x0400; i <= 0x04FF; i++ {
			r := rune(i)
			if unicode.IsLetter(r) {
				chars = append(chars, r)
			}
		}
	}

	for {
		RandChar <- chars[rand.Intn(len(chars))]
	}
}

func main() {

	if len(os.Args) < 3 || *Help {
		fmt.Println("Usage: matrix [options] <width> <height>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	width, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	height, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	fmt.Print(EnAltScreenBuf, Hide)
	defer func() {
		fmt.Print(Show, DisAltScreenBuf)
	}()

	// Exit on Ctrl-C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		fmt.Print(Show, DisAltScreenBuf)
		os.Exit(0)
	}()

	go generateChars()
	screens := make(chan string, 1)
	go content(screens, height, width)

	for s := range screens {
		fmt.Print(Home, s)
	}
}
