package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/gdamore/tcell"
)

type CursorLocation struct {
	X int
	Y int
}

func main() {
	cursor := &CursorLocation{0, 0}
	logf, err := os.OpenFile("./tmp/gim.log", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("failed to create log file: %v", err)
	}
	defer logf.Close()

	log.SetOutput(logf)
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("failed to create new terminal screen: %v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("failed to init screen: %v", err)
	}

	if len(os.Args) > 1 {
		filename := os.Args[1]
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("failed to open file: %v", err)
		}
		defer file.Close()
		if err := loadReader(screen, file); err != nil {
			log.Fatalf("failed to load file: %v", err)
		}
	}

	cursor = &CursorLocation{0, 0}
	screen.ShowCursor(cursor.X, cursor.Y)
	screen.Sync()

	for {
		event := screen.PollEvent()
		log.Printf("received key press: %#v", event)

		switch e := event.(type) {
		case *tcell.EventError:
			log.Printf("event_error: %v", e)
		case *tcell.EventInterrupt:
			log.Printf("event_interrupt: %v", e)
			break
		case *tcell.EventKey:
			key := e.Key()
			if key == tcell.KeyCtrlC {
				log.Println("the key is control c")
				goto exit
			}
			if key == tcell.KeyUp {
				if cursor.Y > 0 {
					cursor.Y--
					screen.ShowCursor(cursor.X, cursor.Y)
					screen.Show()
				}
			}
			if key == tcell.KeyDown {
				_, y := screen.Size()
				if cursor.Y < y-1 {
					cursor.Y++
					screen.ShowCursor(cursor.X, cursor.Y)
					screen.Show()
				}
			}
			if key == tcell.KeyRight {
				x, _ := screen.Size()
				if cursor.X < x-1 {
					cursor.X++
					screen.ShowCursor(cursor.X, cursor.Y)
					screen.Show()
				}
			}
			if key == tcell.KeyLeft {
				if cursor.X > 0 {
					cursor.X--
					screen.ShowCursor(cursor.X, cursor.Y)
					screen.Show()
				}
			}
			if key == tcell.KeyRune {
				printRune(screen, cursor, e.Rune())
				screen.ShowCursor(cursor.X, cursor.Y)
				screen.Show()
				break
			}
		case *tcell.EventMouse:
			log.Printf("event_mouse: %v", e)
		case *tcell.EventResize:
			log.Printf("event_resize: %v", e)
		case *tcell.EventTime:
			log.Printf("event_time: %v", e)
		}
	}

exit:
	log.Printf("application ends")
}

func loadReader(screen tcell.Screen, reader io.Reader) error {
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	tmpCursor := &CursorLocation{0, 0}
	str := string(content)
	for _, r := range str {
		_, y := screen.Size()
		if tmpCursor.Y >= y {
			break
		}

		printRune(screen, tmpCursor, r)
	}

	return nil
}

func printRune(screen tcell.Screen, cursor *CursorLocation, r rune) {
	log.Printf("printing rune %q", r)
	if r == '\n' {
		cursor.Y++
		cursor.X = 0
		return
	}
	if r == '\t' {
		cursor.X += 8
		return
	}
	screen.SetContent(cursor.X, cursor.Y, r, nil, tcell.StyleDefault)
	cursor.X++
}

func moveCursor(screen tcell.Screen, cursor *CursorLocation) {
	cursor.X++
	x, _ := screen.Size()
	if cursor.X >= x {
		cursor.X = 0
		cursor.Y++
	}
}
