package main

import (
	"bufio"
	"io"
	"os"
)

// Buffer is a represetnation of the buffer in the editor. It is structured to
// have an array of lines. This is done in order to help with editing the
// files. It is much easiser to keep track of the changes this way.
type Buffer struct {
	filepath string
	lines    [][]rune
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
	lines := make([][]rune, 0)
	for err != io.EOF {
		lines = append(lines, []rune(line))
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

func (buf *Buffer) InsertRune(r rune, linenum, linecol int) {
	line := buf.lines[linenum]
	line = append(line[:linecol], append([]rune{r}, line[linecol:]...)...)
	buf.lines[linenum] = line
}

// NumOfLines returns the number of lines in the buffer.
func (buf *Buffer) NumOfLines() int {
	return len(buf.lines)
}

// Line gets the ith line.
func (buf *Buffer) Line(i int) []rune {
	if i < len(buf.lines) {
		return buf.lines[i]
	}
	return []rune{}
}
