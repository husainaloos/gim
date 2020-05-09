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

	// Cursor on the buffer
	Cursor *CursorLocation
}

// Draw draws the buffer on the screen. Currently, error is alwasy nil.
func (bv *BufferView) Draw(screen tcell.Screen) {
	screen.Clear()

	// do nothing if we exceeded the lines
	if bv.StartLine > len(bv.Buf.lines) {
		return
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
}

// MoveCursorDown by one step. If the cursor exceeds the screen size downward,
// the then buffer is moved down one step. If nothing to display, then prevent
// the cursor from moving.
func (bv *BufferView) MoveCursorDown(screen tcell.Screen) {
	// if the starting line being displayed is at the end of the buffer,
	// then do nothing
	if bv.StartLine >= len(bv.Buf.lines)-1 {
		return
	}

	// check if the cursor can move down since it is not at the end of the
	// buffer
	if bv.cursorAtBottomOfView() {
		// move the buffer view one step down only if there is more of
		// the buffer to view
		if !bv.atBottomOfBuffer() {
			bv.StartLine++
			bv.Draw(screen)
		}
		return
	}

	// TODO: consider not moving the cursor beyond the file
	bv.Cursor.Y++
	screen.ShowCursor(bv.Cursor.X, bv.Cursor.Y)
}

// MoveCursorUp by one step. If the cursor exceeds the screen size upward, the
// then buffer is moved up one step. If nothing to display, then prevent the
// cursor from moving.
func (bv *BufferView) MoveCursorUp(screen tcell.Screen) {
	// if the cursor is at the top of the view, then consider moving the
	// buffer view one step up in the buffer
	if bv.cursorAtTopOfView() {
		// move the buffer view one step down only if there is more of
		// the buffer to view
		if !bv.atTopOfBuffer() {
			bv.StartLine--
			bv.Draw(screen)
		}
		return
	}

	// TODO: consider hadling the case when the cursor is beyond the line
	bv.Cursor.Y--
	screen.ShowCursor(bv.Cursor.X, bv.Cursor.Y)
}

func (bv *BufferView) cursorAtBottomOfView() bool {
	return bv.Cursor.Y == bv.Height-1
}
func (bv *BufferView) cursorAtTopOfView() bool {
	return bv.Cursor.Y == 0
}

func (bv *BufferView) atBottomOfBuffer() bool {
	return bv.Height+bv.StartLine >= len(bv.Buf.lines)
}

func (bv *BufferView) atTopOfBuffer() bool {
	return bv.StartLine == 0
}
