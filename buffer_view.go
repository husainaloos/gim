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
	// if the starting line being displayed is at the end of the buffer,
	// then do nothing
	if bv.StartLine >= bv.Buf.NumOfLines()-1 {
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
	return bv.Height+bv.StartLine >= bv.Buf.NumOfLines()
}

func (bv *BufferView) atTopOfBuffer() bool {
	return bv.StartLine == 0
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

// Adjust adjusts the current view of the buffer. This means scrolling the
// BufferView up or down, right of left, if the cursor is out-of-bound. This
// also adjust the cursor if it is going beyond the Buffer. Also, this adjust
// the position of the cursor to make sure it is at a valid character.
func (bv *BufferView) Adjust() {
	//TODO(husainaloos): need to finish the implementation of this function

	// find out what line of the buffer are we one.
}
