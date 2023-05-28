package main

import (
	"fmt"
	"math/rand"
)

type Column interface {
	// Char return the "character" at the given index.
	// The character is a string, because it may be multiple runes
	// to allow for ANSI escape sequences.
	Char(i int) string

	Update()
	Finished() bool
}

func NewColumn(height int) Column {
	wh := ColumnMinHeight + rand.Intn(height-ColumnMinHeight)

	chars := make([]rune, height)
	for i := range chars {
		chars[i] = ' '
	}

	return &NormalColumn{
		windowHeight: wh,
		offset:       -wh,
		chars:        chars,
		height:       height,
	}
}

// -----------------------------------------------------------------------------
// Normal Column

type NormalColumn struct {
	windowHeight int
	offset       int
	chars        []rune
	height       int
	finished     bool
}

func (c *NormalColumn) EndOffset() int {
	return c.offset + c.windowHeight
}

func (c *NormalColumn) Finished() bool {
	return c.finished
}

func (c *NormalColumn) Char(i int) string {
	ch := c.chars[i]

	if i == c.EndOffset() {
		return fmt.Sprintf("%s%s%c%s%s", FgWhite, Bold, ch, Reset, FgDefault)
	}

	if c.EndOffset()-2 <= i && i < c.EndOffset() {
		return fmt.Sprintf("%s%c%s", FgWhite, ch, FgDefault)
	}

	if c.offset <= i && i <= c.offset+2 {
		return fmt.Sprintf("%s%s%c%s%s", FgGreen, Dim, ch, Reset, FgDefault)
	}

	return fmt.Sprintf("%s%c%s", FgGreen, ch, FgDefault)
}

func (c *NormalColumn) Update() {
	if c.finished {
		return
	}

	// Remove old top character
	if c.offset >= 0 {
		c.chars[c.offset] = ' '
	}

	// Update offset
	c.offset++
	if c.offset >= c.height {
		c.finished = true
		return
	}

	// Add new bottom character
	if c.EndOffset() < c.height {
		c.chars[c.EndOffset()] = <-RandChar
	}

	// Update a single random character in the window
	idx := c.offset + rand.Intn(c.windowHeight)
	if 0 < idx && idx < c.height {
		c.chars[idx] = <-RandChar
	}
}
