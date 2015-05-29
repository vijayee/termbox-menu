package main

import (
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
	keyEventService  chan termbox.Event
	isFocused        bool
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

var subscribers []chan termbox.Event
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
		titleStart := 3
		var c rune
		for x := 0; x < w; x++ {
			switch {
			case x > titleStart && titleIndex < len(title):
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
	m.isFocused = false
	m.items[m.currentSelection].Invoke()
	m.isFocused = true
}
func (m *Menu) ListenToKeys() {
	m.keyEventService = make(chan termbox.Event)
	Subscribe(m.keyEventService)
	m.isFocused = true
	//defer UnSubscribe(m.keyEventService)
	for {
		select {
		case keyEvent := <-m.keyEventService:
			switch keyEvent.Type {
			case termbox.EventKey:
				switch keyEvent.Key {
				case termbox.KeyEsc:
					if m.isFocused == true {
						return
					}
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

}
func (m *Menu) StopListeningToKeys() {
	UnSubscribe(m.keyEventService)
	close(m.keyEventService)
}
func NewMenu(title string, items []Item) Menu {
	return Menu{"More Options", items, 0, termbox.ColorWhite, termbox.ColorBlue, nil, false}
}
func ListenToKeys() {
	isListening = true
	for isListening == true {
		Emit(termbox.PollEvent())
	}

}
func StopListeningToKeys() {
	isListening = false
}
func Emit(event termbox.Event) {
	for _, listener := range subscribers {
		select {
		case listener <- event:
			continue
		default:
			continue
		}
	}
}
func Subscribe(listener chan termbox.Event) {
	if subscribers == nil {
		subscribers = make([]chan termbox.Event, 1)
	} else {
		oneUp := make([]chan termbox.Event, cap(subscribers)+1)
		for i := range subscribers {
			oneUp[i] = subscribers[i]
		}
		subscribers = oneUp
	}
	subscribers[len(subscribers)-1] = listener
}

func UnSubscribe(listener chan termbox.Event) {
	if subscribers == nil {
		return
	} else {
		if len(subscribers) == 1 {
			subscribers = nil
			return
		}
		oneDown := make([]chan termbox.Event, cap(subscribers)-1)
		for i := range subscribers {
			if listener == subscribers[i] {
				continue
			}
			oneDown[i] = subscribers[i]
		}
		subscribers = oneDown
	}
	subscribers[len(subscribers)-1] = listener
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
	menu3 := NewMenu("Even More Options", []Item{item4, item1})
	menu2 = NewMenu("More Options", []Item{item4, item5, item6, &menu3})
	menu1 = NewMenu("Tour of IPFS", []Item{item1, item2, item3, &menu2})
	go ListenToKeys()
	menu1.Invoke()

}
