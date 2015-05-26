package menu

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"unicode/utf8"
)

type Item interface {
	Title() string
	Invoke() error
}
type Menu struct {
	title            string
	items            []Item
	currentSelection int
	foreground       termbox.Attribute
	background       termbox.Attribute
}
type MenuItem struct {
	title string
}

func (m *MenuItem) Title() string {
	return m.title
}
func (m *MenuItem) Invoke() error {
	m.title = "selected"
	return nil
}

var keyInput chan termbox.Event
var isListening bool

func (m *Menu) drawTitle() {
	w, _ := termbox.Size()
	titleStart := (w / 2) - (len(m.title) / 2)
	titleRow := 2
	titleIndex := 0
	for x := 0; x < w; x++ {
		if x > titleStart && titleIndex < len(m.title) {
			c, _ := utf8.DecodeRuneInString(m.title[titleIndex:])
			titleIndex++
			termbox.SetCell(x, titleRow, c, m.foreground, m.background)
			c = '_'
		}
		termbox.SetCell(x, titleRow+1, '_', m.foreground, m.background)
	}

}
func (m *Menu) drawItems() {
	w, h := termbox.Size()
	currrentRow := 5
	for index, item := range m.items {
		if currrentRow > h {
			break
		}
		titleIndex := 0
		title := item.Title()
		titleStart := 5
		itemNumber := fmt.Sprintf("%03d", index+1)
		numIndex := 0
		var c rune
		for x := 0; x < w; x++ {
			switch {
			case x >= titleStart && x <= titleStart+3 && numIndex < len(itemNumber):
				c, _ = utf8.DecodeRuneInString(itemNumber[numIndex:])
				numIndex++
			case x > titleStart+4 && titleIndex < len(title):
				c, _ = utf8.DecodeRuneInString(title[titleIndex:])
				titleIndex++
			default:
				c = ' '

			}
			if m.currentSelection == index {
				termbox.SetCell(x, currrentRow, c, m.background, m.foreground)
			} else {

				termbox.SetCell(x, currrentRow, c, m.foreground, m.background)
			}
		}
		currrentRow += 2
	}

}
func (m *Menu) draw() error {
	termbox.Clear(m.background, m.background)
	m.drawTitle()
	m.drawItems()
	termbox.Flush()
	return nil
}
func (m *Menu) Invoke() error {
	m.draw()
	m.ListenToKeys()
	return nil
}

func (m *Menu) Title() string {
	return m.title
}
func (m *Menu) Up() {
	m.currentSelection--
	if m.currentSelection < 0 {
		m.currentSelection = 0
	}
}

func (m *Menu) Down() {
	m.currentSelection++
	if m.currentSelection >= len(m.items) {
		m.currentSelection = len(m.items) - 1
	}
}
func (m *Menu) Select() {
	m.ignoreKeys()
	m.items[m.currentSelection].Invoke()
	m.draw()
	m.ListenToKeys()
}
func (m *Menu) listenToKeys() <-chan termbox.Event {
	if keyInput == nil {
		keyInput = make(chan termbox.Event)

		go func() {
			for {
				select {
				case keyInput <- termbox.PollEvent():
					continue
				default:
					continue
				}

			}
		}()
	}
	return keyInput
}
func (m *Menu) ListenToKeys() {
	isListening = true
	keyEventService := m.listenToKeys()
listenerLoop:
	for {
		select {
		case keyEvent := <-keyEventService:
			switch keyEvent.Type {
			case termbox.EventKey:
				switch keyEvent.Key {
				case termbox.KeyEsc:
					break listenerLoop
				case termbox.KeyArrowUp:
					go func() {
						m.Up()
						m.draw()
					}()
				case termbox.KeyArrowDown:
					go func() {
						m.Down()
						m.draw()
					}()
				case termbox.KeyEnter:
					go func() {
						m.Select()
						m.draw()
					}()
				}
			case termbox.EventError:
				panic(keyEvent.Err)
			}
		}
		m.draw()
	}
	isListening = false

}

func (m *Menu) ignoreKeys() {
	isListening = false

}

var menu1 Menu
var menu2 Menu

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
	item1 := &MenuItem{"first"}
	item2 := &MenuItem{"second"}
	item3 := &MenuItem{"third"}
	item4 := &MenuItem{"four"}
	item5 := &MenuItem{"five"}
	item6 := &MenuItem{"six"}
	menu2 = Menu{"More Options", []Item{item4, item5, item6}, 0, termbox.ColorWhite, termbox.ColorRed}
	menu1 = Menu{"Tour of IPFS", []Item{item1, item2, item3, &menu2}, 0, termbox.ColorWhite, termbox.ColorRed}
	menu1.Invoke()

}
