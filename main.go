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

	Ascii      = flag.Bool("ascii", false, "Use only ASCII characters")
	NoLatin    = flag.Bool("no-latin", false, "Disable latin extended characters")
	NoGreek    = flag.Bool("no-greek", false, "Disable greek characters")
	NoCyrillic = flag.Bool("no-cyrillic", false, "Disable cyrillic characters")
	Help       = flag.Bool("help", false, "Show help")
)

type Screen struct {
	Columns []Column
	Height  int
	Width   int
}

func NewScreen(height, width int) *Screen {
	return &Screen{
		Columns: make([]Column, width),
		Height:  height,
		Width:   width,
	}
}

func (s *Screen) Update() {
	// Randomly add a new column. This seems to create a nice distribution
	// similar to the original matrix, with minimal logic.
	idx := rand.Intn(len(s.Columns))
	if s.Columns[idx] == nil {
		s.Columns[idx] = NewColumn(s.Height)
	}

	for i, c := range s.Columns {
		if c == nil {
			continue
		}
		c.Update()
		if c.Finished() {
			s.Columns[i] = nil
		}
	}
}

func (s *Screen) String() string {
	sb := strings.Builder{}
	sb.Grow(s.Height * s.Width)

	// Build the entire screen as a single string, line by line.
	for i := 0; i < s.Height; i++ {
		if i > 0 {
			sb.WriteByte('\n')
		}
		for _, c := range s.Columns {
			if c == nil {
				sb.WriteByte(' ')
			} else {
				sb.WriteString(c.Char(i))
			}
		}
	}

	return sb.String()
}

func generateContent(ch chan<- string, height, width int) {
	screen := NewScreen(height, width)
	for range time.Tick(DrawInterval) {
		screen.Update()
		ch <- screen.String()
	}
}

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
	if !(*NoLatin || *Ascii) {
		for i := 0x00C0; i <= 0x00FF; i++ {
			r := rune(i)
			if unicode.IsLetter(r) {
				chars = append(chars, r)
			}
		}
	}

	// Greek
	if !(*NoGreek || *Ascii) {
		for i := 0x0370; i <= 0x03FF; i++ {
			r := rune(i)
			if unicode.IsLetter(r) {
				chars = append(chars, r)
			}
		}
	}

	// Cyrillic
	if !(*NoCyrillic || *Ascii) {
		for i := 0x0400; i <= 0x04FF; i++ {
			r := rune(i)
			if unicode.IsLetter(r) {
				chars = append(chars, r)
			}
		}
	}

	/*
		// Jog: this can even be evaluated at compile time
		chars := [ rune(i) for i := 33; i <= 126; i++ if unicode.IsPrint(rune(i)) ]

		if !(*NoLatin || *Ascii) {
			chars += [ rune(i) for i := 0x00C0; i <= 0x00FF; i++ if unicode.IsLetter(rune(i)) ]
		}

		if !(*NoGreek || *Ascii) {
			chars += [ rune(i) for i := 0x0370; i <= 0x03FF; i++ if unicode.IsLetter(rune(i)) ]
		}

		if !(*NoCyrillic || *Ascii) {
			chars += [ rune(i) for i := 0x0400; i <= 0x04FF; i++ if unicode.IsLetter(rune(i)) ]
		}
	*/

	for {
		RandChar <- chars[rand.Intn(len(chars))]
	}
}

func main() {
	//@Error panic(err)

	flag.Parse()
	if len(flag.Args()) != 2 || *Help {
		fmt.Println("Usage: matrix [options] <width> <height>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	//width := strconv.Atoi(flag.Arg(0))?
	width, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	height, err := strconv.Atoi(flag.Arg(1))
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

	ch := make(chan string, 1)
	go generateChars()
	go generateContent(ch, height, width)

	for text := range ch {
		fmt.Print(Home, text)
	}
}
