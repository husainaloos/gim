package main

import (
	"github.com/gdamore/tcell"
)

const (
	TabSize = 8
)

// BufferView is the view to the buffer. Buffers are not displayed completely
// in the window. Only a section of the buffer is displayed.
type BufferView struct {

	// Buf is the buffer associated with the view
	Buf *Buffer

	// StartLine is the starting line in the buffer. This is in case the
	// view is scrolled down. StartLine starts at 0.
	StartLine int

	// StartColumn is the column at which the view beings. This is needed
	// in case the view is scrolled to the right. StartColumn starts at 0.
	StartColumn int

	// Height of the view
	Height int

	// Width of the view
	Width int

	// Cursor on the buffer
	Cursor *CursorLocation
}

// Draw draws the buffer on the screen. This is done by geting the runes from
// the buffer, and manually refresh everything in the screen. This might be
// slow, but it works for now. Currently, the error returned is always nil.
func (bv *BufferView) Draw(screen tcell.Screen) {
	screen.Clear()

	cursor := &CursorLocation{0, 0}
	linecount := 0
	i := bv.StartLine
	for linecount < bv.Height && i < bv.Buf.NumOfLines() {
		runes := bv.Buf.Line(i)
		for j := bv.StartColumn; j < len(runes); j++ {
			r := runes[j]

			// if the rune is a tab, then insert tab-size spaces instead of a single character.
			if r == '\t' {
				for k := 0; k < TabSize; k++ {
					screen.SetContent(cursor.X, cursor.Y, ' ', nil, tcell.StyleDefault)
					cursor.X++
				}
				if cursor.X >= bv.Width {
					break
				}
				continue
			}
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
	screen.ShowCursor(bv.Cursor.X, bv.Cursor.Y)
}

// MoveCursorDown by one step. If the cursor exceeds the screen size downward,
// the then buffer is moved down one step. If nothing to display, then prevent
// the cursor from moving.
func (bv *BufferView) MoveCursorDown(screen tcell.Screen) {
	// TODO: this implmeentation does not support TAB movement.
	if bv.currentBufferLine() >= bv.Buf.NumOfLines() {
		return
	}
	if bv.Cursor.Y >= bv.Height-1 {
		bv.StartLine++
	} else {
		bv.Cursor.Y++
	}
	maxCol := len(bv.Buf.Line(bv.currentBufferLine()))
	if bv.currentBufferColumn() >= maxCol {
		bv.StartColumn = maxCol - bv.Width
		if bv.StartColumn < 0 {
			bv.StartColumn = 0
		}
		bv.Cursor.X = maxCol - bv.StartColumn - 1
	}
	bv.Draw(screen)
}

// MoveCursorUp by one step. If the cursor exceeds the screen size upward, the
// then buffer is moved up one step. If nothing to display, then prevent the
// cursor from moving.
func (bv *BufferView) MoveCursorUp(screen tcell.Screen) {
	// TODO: this implmeentation does not support TAB movement.
	if bv.currentBufferLine() == 0 {
		return
	}
	if bv.Cursor.Y == 0 {
		bv.StartLine--
	} else {
		bv.Cursor.Y--
	}
	maxCol := len(bv.Buf.Line(bv.currentBufferLine()))
	if bv.currentBufferColumn() >= maxCol {
		bv.StartColumn = maxCol - bv.Width
		if bv.StartColumn < 0 {
			bv.StartColumn = 0
		}
		bv.Cursor.X = maxCol - bv.StartColumn - 1
	}
	bv.Draw(screen)
}

// MoveCursorRight be one step. If the cursor at the end of a line, nothing
// happens. If the cursor is at the end of a view, then the view moves to the
// right.
func (bv *BufferView) MoveCursorRight(screen tcell.Screen) {
	// TODO: this implmeentation does not support TAB movement.
	if bv.currentBufferColumn() >= len(bv.Buf.Line(bv.currentBufferLine()))-1 {
		if bv.currentBufferLine() == bv.Buf.NumOfLines() {
			return
		}
		bv.Cursor.X = 0
		bv.Cursor.Y++
		bv.StartColumn = 0
	} else if bv.Cursor.X >= bv.Width {
		bv.Cursor.X = bv.Width - 1
		bv.StartColumn++
	} else {
		bv.Cursor.X++
	}
	bv.Draw(screen)
}

// MoveCursorLeft be one step. If the cursor at the beginning of a line,
// nothing happens. If the cursor is at the beginning of a view, then the view
// moves to the left.
func (bv *BufferView) MoveCursorLeft(screen tcell.Screen) {
	if bv.currentBufferColumn() == 0 {
		if bv.currentBufferLine() == 0 {
			return
		}
		bv.Cursor.Y--
		maxCol := len(bv.Buf.Line(bv.currentBufferLine()))
		bv.Cursor.X = maxCol - 1
		if bv.Cursor.X >= bv.Width {
			bv.Cursor.X = bv.Width - 1
			bv.StartColumn = maxCol - bv.Width
		}
	} else if bv.Cursor.X <= 0 {
		bv.Cursor.X = 0
		bv.StartColumn--
	} else {
		bv.Cursor.X--
	}
	bv.Draw(screen)
}

func (bv *BufferView) cursorLineIndexInBuffer() int {
	v := bv.StartLine + bv.Cursor.Y
	if v > bv.Buf.NumOfLines()-1 {
		v = bv.Buf.NumOfLines() - 1
	}
	return v
}

func (bv *BufferView) numberOfRunesInBufferLine() int {
	i := bv.cursorLineIndexInBuffer()
	return len(bv.Buf.Line(i))
}

func (bv *BufferView) InsertRune(screen tcell.Screen, r rune) {
	lineNum := bv.Cursor.Y + bv.StartLine
	lineCol := bv.Cursor.X + bv.StartColumn
	bv.Buf.InsertRune(r, lineNum, lineCol)
	bv.Cursor.X++
	bv.Draw(screen)
}

func (bv *BufferView) currentBufferLine() int {
	return bv.StartLine + bv.Cursor.Y
}

func (bv *BufferView) currentBufferColumn() int {
	return bv.StartColumn + bv.Cursor.X
}
