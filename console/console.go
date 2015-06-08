/*
	This package includes some usefull methods for displaying data on the terminal.
	This package depends on termbox which is a crossplatform library written by Georg Reinke.
	Termbox is available at github.com/nsf/termbox-go
*/

package console

import (
	//"fmt"
	"github.com/indidev/vocable-o/util/mathutil"
	"github.com/indidev/vocable-o/util/stringutil"
	"github.com/nsf/termbox-go"
	//"strings"
	//"math"
)

type Rect struct {
	x, y, w, h int
}

const noItemInfo = "Sorry, no item to choose."
var width, height int
var event = make(chan termbox.Event)
var info = make(chan bool)
var infoLine string

/*
	initiates the terminal
*/
func Init() {
	go run()

	<-info //wait for termbox to open
}

func run() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	width, height = termbox.Size()

	Clear()
	info <- true

	for repeate := true; repeate; {
		ev := termbox.PollEvent()

		if ev.Type == termbox.EventResize {
			update()
		}

		if ev.Type != termbox.EventInterrupt {
			event <- ev
		} else {
			repeate = false
			close(event)
		}
	}
	termbox.Close()
	//fmt.Println("Terminal closed")
	info <- true
}

/*
	updates the terminal
*/
func update() {
	width, height = termbox.Size()
	updateInfoBottom()
}

/*
	quits the terminal, highly recommended to use!
*/
func Quit() {
	termbox.Interrupt()
	<-info // wait for the terminal
}

/*
	returns an event channel
*/
func Event() chan termbox.Event {
	return event
}

/*
	clears the terminal
*/
func Clear() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	updateInfoBottom()
}

/*
	writes a string to the terminal at given x, y position
*/
func Write(x, y int, value string) {
	writeNoFlush(x, y, value)
	termbox.Flush()
}

func writeNoFlush(x, y int, value string) {
	if height > y && y >= 0 {
		i := 0
		for _, c := range value {
			if width > (x + i) {
				termbox.SetCell(x+i, y, c, termbox.ColorDefault, termbox.ColorDefault)
			}
			i++
		}
	}
}

/*
	Creates a vertical and horizontal centered menu with the given info text and items to choose from.
	returns -1 if Esc has been pressed otherwise the index of the selected item.
*/
func Menu(items []string, info string) int {

	selection, _ := ExtendedMenu(items, info, []rune{}, 0)

	return selection
}

/*
	Like Menu, just with preselection of an item as well as characters which can trigger an invent.
	returns -1 if Esc has been pressed.
	returns index of the selected item and rune(0) if enter has been pressed.
	If an accepted character has been pressed, the menu returns the index as well as the pressed character.
*/
func ExtendedMenu(items []string, info string, acceptedChars []rune, selectedID int) (int, rune) {
	selection := selectedID

	key := rune(0)

	drawMenu(items, selection, info)

loop:
	for ev := range Event() {
		switch ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				selection = -1
				break loop
			case termbox.KeyEnter:
				break loop
			case termbox.KeyArrowUp:
				selection--
				selection = mathutil.MaxInt(selection, 0)
				drawMenu(items, selection, info)

			case termbox.KeyArrowDown:
				selection++
				selection = mathutil.MinInt(selection, len(items)-1)
				drawMenu(items, selection, info)

			default:
				if ev.Ch != rune(0) {
					for _, x := range acceptedChars {
						if x == ev.Ch {
							key = ev.Ch
							break loop
						}
					}
				}
			}
		case termbox.EventResize:
			drawMenu(items, selection, info)
		}
	}

  //avoid out of bounds exeption (no items means always escape, no enter)
	if len(items) == 0 {
		selection = -1
	}

	return selection, key
}

func drawMenu(items []string, selectedId int, info string) {
	empty := "[ ] "
	selected := "[#] "
	above := "ᐃ"
	below := "ᐁ"

	w := 0
	h := len(items)

	h = mathutil.MinInt(h, height-8)

	//determine width of longest item and therefor width of the box spacing the items
	for _, elem := range items {
		w = mathutil.MaxInt(stringutil.Size(elem), w)
	}

	if len(items) == 0 {
		w = stringutil.Size(noItemInfo) - 4
	}

	//add width of selection item
	w += 4

	x := width - w - 4
	y := height - h - 4

	x /= 2
	y /= 2

	//clear frame
	Clear()

	//write info
	writeNoFlush((width-stringutil.Size(info))/2, y-2, info)

	//draw box
	termbox.SetCell(x, y, '┌', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(x, y+h+3, '└', termbox.ColorDefault, termbox.ColorDefault)

	termbox.SetCell(x+w+3, y, '┐', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(x+w+3, y+h+3, '┘', termbox.ColorDefault, termbox.ColorDefault)

	for i := 1; i < w+3; i++ {
		termbox.SetCell(x+i, y, '─', termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCell(x+i, y+h+3, '─', termbox.ColorDefault, termbox.ColorDefault)
	}

	for i := 1; i <= h+2; i++ {
		termbox.SetCell(x, y+i, '│', termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCell(x+w+3, y+i, '│', termbox.ColorDefault, termbox.ColorDefault)
	}

	x += 2
	y += 2

	//fill box‍

	startindex := 0

	if h < len(items) {
		startindex = mathutil.MinInt(selectedId-1, len(items)-h)
		startindex = mathutil.MaxInt(0, startindex)
	}

	for i := 0; i < h; i++ {
		symbol := empty
		if (i + startindex) == selectedId {
			symbol = selected
		}

		writeNoFlush(x, y+i, stringutil.Join(symbol, items[i+startindex]))
	}

	if len(items) == 0 {
		writeNoFlush(x, y, noItemInfo)
	}

	if startindex > 0 {
		writeNoFlush(width/2-1, y-1, above)
	}
	if (startindex + h) < len(items) {
		writeNoFlush(width/2-1, y+h, below)
	}
	termbox.Flush()
}

/*
	Creates an input including cursor.
	a replacement map can be used to support fast usage of special characters.
*/
func Input(x, y int, replacements map[string]string, text string) (string, bool) {
	input := text
	//cursor := 0
	maxsize := stringutil.Size(input)

	writeNoFlush(x, y, input)
	termbox.SetCursor(x+stringutil.Size(input), y)
	termbox.Flush()
	char, key := getKeyEvent()

	for key != termbox.KeyEnter && key != termbox.KeyEsc {
		if char != 0 {
			input = stringutil.Join(input, string(char))

		} else {
			switch key {
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				input = stringutil.RemoveTail(input, 1)

			case termbox.KeySpace:
				input = stringutil.Join(input, " ")
			}
		}

		input = stringutil.ReplaceMap(input, replacements)
		maxsize = mathutil.MaxInt(maxsize, stringutil.Size(input))
		clearLine(x, y, maxsize)
		Write(x, y, input)

		termbox.SetCursor(x+stringutil.Size(input), y)
		termbox.Flush()
		char, key = getKeyEvent()
	}

	termbox.HideCursor()

	return input, key == termbox.KeyEnter
}

/*
	clears a line from x to x + wdith
*/
func clearLine(x, y, width int) {
	for i := 0; i < width; i++ {
		termbox.SetCell(x+i, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}
}

/*
	listens to the event chanel and waits for a key to be pressed.
	returns the character as well as the key
*/
func getKeyEvent() (rune, termbox.Key) {
	eventChan := Event()

	event := <-eventChan

	for event.Type != termbox.EventKey {
		event = <-eventChan
	}

	return event.Ch, event.Key
}

/*
	displays multiple lines vertical and horizontal centered on the termianl
	The text inside the virtual box is left aligned
*/
func DisplayCentered(lines []string) Rect {
	ySize := len(lines)
	xSize := 0

	for _, line := range lines {
		xSize = mathutil.MaxInt(xSize, stringutil.Size(line))
	}

	deltaX := (width - xSize) / 2
	deltaY := (height - ySize) / 2

	for i, line := range lines {
		Write(deltaX, deltaY+i, line)
	}
	termbox.Flush()

	return Rect{deltaX, deltaY, xSize, ySize}
}

/*
	like DisplayCentered just with an input line below the output
*/
func DisplayCenteredWithInput(lines []string, replacements map[string]string, inputText string) (string, bool) {
	lines = append(lines, "")
	lines = append(lines, "")

	rect := DisplayCentered(lines)

	return Input(rect.x, rect.y+rect.h-1, replacements, inputText)
}

/*
	sets the info text at the bottom of the terminal
	e.q. good for displaying key shortcuts
*/
func SetInfoBottom(infoString string) {
	infoLine = infoString
	updateInfoBottom()
}

func updateInfoBottom() {
	y := height - 1
	x := (width - stringutil.Size(infoLine)) / 2

	Write(x, y, infoLine)
	termbox.Flush()
}

/*
	Blocks until the user presses escape
*/
func WaitForEsc() {
wait:
	for ev := range Event() {
		if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
			break wait
		}
	}
}

/*
	Blocks until the user presses enter
*/
func WaitForEnter() {
wait:
	for ev := range Event() {
		if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEnter {
			break wait
		}
	}
}

/*
	Blocks until the user presses any key
*/
func WaitForAnyInput() rune {
	key, _ := getKeyEvent()
	return key
}

/*
	returns the size of the terminal
*/
func Size() (int, int) {
	return width, height
}
