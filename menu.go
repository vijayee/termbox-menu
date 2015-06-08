package menu

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

var subscribers []chan termbox.Event
var isListening bool

func (m *Menu) drawTitle() {
	w, _ := termbox.Size()
	titleStart := (w / 2) - (len(m.title) / 2)
	titleRow := 2
	titleIndex := 0
	for x := 0; x < w; x++ {
		if x >= titleStart && titleIndex < len(m.title) {
			c, rw := utf8.DecodeRuneInString(m.title[titleIndex:])
			titleIndex += rw
			titleStart += rw
			termbox.SetCell(x, titleRow, c, m.foreground, m.background)
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
		var rw int
		for x := 0; x < w; x++ {
			switch {
			case x >= titleStart && titleIndex < len(title):
				c, rw = utf8.DecodeRuneInString(title[titleIndex:])
				titleIndex += rw
				titleStart += rw
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
	for {
		select {
		case keyEvent := <-m.keyEventService:
			switch keyEvent.Type {
			case termbox.EventKey:
				if m.isFocused == true {
					switch keyEvent.Key {
					case termbox.KeyEsc:
						if m.isFocused == true {
							return
						}
					case termbox.KeyArrowUp:
						if m.isFocused == true {
							go func() {
								m.Up()
								m.draw()
							}()
						}
					case termbox.KeyArrowDown:
						if m.isFocused == true {
							go func() {
								m.Down()
								m.draw()
							}()
						}
					case termbox.KeyEnter:
						if m.isFocused == true {
							go func() {
								m.Select()
								m.draw()
							}()
						}
					}
				}

			case termbox.EventError:
				panic(keyEvent.Err)
			}
		}
		if m.isFocused == true {
			m.draw()
		}
	}

}
func (m *Menu) StopListeningToKeys() {
	UnSubscribe(m.keyEventService)
	close(m.keyEventService)
}
func NewMenu(title string, items []Item, foreground termbox.Attribute, background termbox.Attribute) Menu {
	return Menu{title, items, 0, foreground, background, nil, false}
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
