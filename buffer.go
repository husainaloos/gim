package main

import (
	"bufio"
	"io"
	"os"
	"unicode/utf8"

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
		// if the line ends with '\n', then remove it
		if line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}
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

	// StartLine is the starting line in the buffer. This is in case the
	// view is scrolled down.
	StartLine int

	// StartColumn is the column at which the view beings. This is needed
	// in case the view is scrolled to the right.
	StartColumn int

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

	cursor := &CursorLocation{0, 0}
	linecount := 0
	i := bv.StartLine
	for linecount < bv.Height && i < len(bv.Buf.lines) {
		line := bv.Buf.lines[i]
		runes := []rune(line)
		for j := bv.StartColumn; j < len(runes); j++ {
			r := runes[j]
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

	bv.Cursor.Y++
	if bv.Cursor.X > bv.numberOfRunesInBufferLine()-1 {
		bv.Cursor.X = bv.numberOfRunesInBufferLine() - 1
	}
	if bv.Cursor.X < 0 {
		bv.Cursor.X = 0
	}
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

	bv.Cursor.Y--
	if bv.Cursor.X > bv.numberOfRunesInBufferLine()-1 {
		bv.Cursor.X = bv.numberOfRunesInBufferLine() - 1
	}
	if bv.Cursor.X < 0 {
		bv.Cursor.X = 0
	}
	screen.ShowCursor(bv.Cursor.X, bv.Cursor.Y)
}

// MoveCursorRight be one step. If the cursor at the end of a line, nothing
// happens. If the cursor is at the end of a view, then the view moves to the
// right.
func (bv *BufferView) MoveCursorRight(screen tcell.Screen) {
	if bv.Cursor.X >= bv.numberOfRunesInBufferLine()-1 {
		return
	}

	if bv.Cursor.X < bv.Width {
		bv.Cursor.X++
		screen.ShowCursor(bv.Cursor.X, bv.Cursor.Y)
		return
	}

	bv.StartColumn++
	bv.Draw(screen)
}

// MoveCursorLeft be one step. If the cursor at the beginning of a line,
// nothing happens. If the cursor is at the beginning of a view, then the view
// moves to the left.
func (bv *BufferView) MoveCursorLeft(screen tcell.Screen) {
	if bv.Cursor.X == 0 && bv.StartColumn == 0 {
		return
	}

	if bv.Cursor.X > 0 {
		bv.Cursor.X--
		screen.ShowCursor(bv.Cursor.X, bv.Cursor.Y)
		return
	}

	bv.StartColumn--
	bv.Draw(screen)
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

func (bv *BufferView) cursorLineIndexInBuffer() int {
	v := bv.StartLine + bv.Cursor.Y
	if v > len(bv.Buf.lines)-1 {
		v = len(bv.Buf.lines) - 1
	}
	return v
}

func (bv *BufferView) numberOfRunesInBufferLine() int {
	i := bv.cursorLineIndexInBuffer()
	return utf8.RuneCountInString(bv.Buf.lines[i])
}
