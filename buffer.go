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

// Load loads the buffer to the screen. This clears the already existing
// content on the screen and replaces it with the content of the buffer. The
// content is alwasy from the beginning of the buffer
func (buf *Buffer) Load(screen tcell.Screen) error {
	maxx, maxy := screen.Size()
	cursor := &CursorLocation{0, 0}
	screen.Clear()
	for _, line := range buf.lines {
		cursor.X = 0
		for _, r := range line {
			//TODO: implement wrapping
			if cursor.X >= maxx {
				break
			}

			screen.SetContent(cursor.X, cursor.Y, r, nil, tcell.StyleDefault)
			cursor.X++
		}
		cursor.Y++
		if cursor.Y >= maxy {
			break
		}
	}

	return nil
}
