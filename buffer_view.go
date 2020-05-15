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
	// view is scrolled down. StartLine starts at 1.
	StartLine int

	// StartColumn is the column at which the view beings. This is needed
	// in case the view is scrolled to the right. StartColumn starts at 1.
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
	bv.Cursor.Y++
	bv.Adjust()
	bv.Draw(screen)
}

// MoveCursorUp by one step. If the cursor exceeds the screen size upward, the
// then buffer is moved up one step. If nothing to display, then prevent the
// cursor from moving.
func (bv *BufferView) MoveCursorUp(screen tcell.Screen) {
	bv.Cursor.Y--
	bv.Adjust()
	bv.Draw(screen)
}

// MoveCursorRight be one step. If the cursor at the end of a line, nothing
// happens. If the cursor is at the end of a view, then the view moves to the
// right.
func (bv *BufferView) MoveCursorRight(screen tcell.Screen) {
	bv.Cursor.X++
	bv.Adjust()
	bv.Draw(screen)
}

// MoveCursorLeft be one step. If the cursor at the beginning of a line,
// nothing happens. If the cursor is at the beginning of a view, then the view
// moves to the left.
func (bv *BufferView) MoveCursorLeft(screen tcell.Screen) {
	bv.Cursor.X--
	bv.Adjust()
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
	bv.Adjust()
	bv.Draw(screen)
}

// Adjust adjusts the current view of the buffer. This means scrolling the
// BufferView up or down, right of left, if the cursor is out-of-bound. This
// also adjust the cursor if it is going beyond the Buffer. Also, this adjust
// the position of the cursor to make sure it is at a valid character.
func (bv *BufferView) Adjust() {
	// If the cursor is too far down, then move the view to where the
	// cursor is at the bottom line of the view.
	if bv.Cursor.Y > bv.Height {
		bv.StartLine += bv.Cursor.Y - bv.Height
	}

	// If the cursor is too far up, then move the view up to where the
	// curosr is at the top line of the view.
	if bv.Cursor.Y < 0 {
		bv.StartLine += bv.Cursor.Y
	}

	// If the view moved down too much, then reset the buffer to display up
	// to the last line of the file.
	if bv.StartLine > bv.Buf.NumOfLines() {
		bv.StartLine = bv.Buf.NumOfLines() - bv.Height
	}

	// If the view moved up too much, then reset the buffer to display the
	// first line of the buffer.
	if bv.StartLine < 0 {
		bv.StartLine = 0
	}

	// Do a similar logic for X
	if bv.Cursor.X > bv.Width {
		bv.StartColumn += bv.Cursor.X - bv.Width
	}

	if bv.Cursor.X < 0 {
		bv.StartColumn -= bv.Cursor.X
	}

	lineLength := bv.numberOfRunesInBufferLine()
	if bv.StartColumn > lineLength {
		bv.StartColumn = lineLength - bv.Width
	}

	if bv.StartColumn < 0 {
		bv.StartColumn = 0
	}

	// Sanity checks that the cursor is not beyond the view.
	if bv.Cursor.X < 0 {
		bv.Cursor.X = 0
	}
	if bv.Cursor.Y < 0 {
		bv.Cursor.Y = 0
	}
	if bv.Cursor.X >= bv.Width {
		bv.Cursor.X = bv.Width - 1
	}
	if bv.Cursor.Y >= bv.Height {
		bv.Cursor.Y = bv.Height - 1
	}

	if bv.Cursor.Y >= bv.Buf.NumOfLines() {
		bv.Cursor.Y = bv.Buf.NumOfLines() - 1
	}
	lineLength = bv.numberOfRunesInBufferLine()
	if bv.Cursor.X >= lineLength {
		bv.Cursor.X = lineLength - 1
	}
}
