package main

import (
	"bufio"
	"io"
	"os"

	"github.com/gdamore/tcell"
)

// Buffer is a represetnation of the buffer in the editor. It is structured to
// have an array of lines. This is done in order to help with editing the
// files. It is much easiser to keep track of the changes this way.
type Buffer struct {
	filepath string
	lines    []string
}

// NewBuffer generates a new buffer from a file. The file is closed after
// reading. If an error is encountered while reading the file, it is returned.
func NewBuffer(filepath string) (*Buffer, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	br := bufio.NewReader(f)

	line, err := br.ReadString('\n')
	lines := make([]string, 0)
	for err != io.EOF {
		lines = append(lines, line)
		line, err = br.ReadString('\n')
	}

	if err != io.EOF {
		return nil, err
	}

	return &Buffer{
		filepath: filepath,
		lines:    lines,
	}, nil
}

// BufferView is the view to the buffer. Buffers are not displayed completely
// in the window. Only a section of the buffer is displayed.
type BufferView struct {

	// Buf is the buffer associated with the view
	Buf *Buffer

	// StartLine is the starting line in the buffer
	StartLine int

	// Height of the view
	Height int

	// Width of the view
	Width int
}

// Draw draws the buffer on the screen. Currently, error is alwasy nil.
func (bv *BufferView) Draw(screen tcell.Screen) error {
	screen.Clear()
	if bv.StartLine > len(bv.Buf.lines) {
		return nil
	}

	cursor := CursorLocation{0, 0}
	linecount := 0
	i := bv.StartLine
	for linecount < bv.Height && i < len(bv.Buf.lines) {
		line := bv.Buf.lines[i]
		for _, r := range line {
			screen.SetContent(cursor.X, cursor.Y, r, nil, tcell.StyleDefault)
			cursor.X++
			if cursor.X >= bv.Width {
				break
			}
		}
		cursor.X = 0
		cursor.Y++
		i++
		linecount++
	}
	return nil
}
