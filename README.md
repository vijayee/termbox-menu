# termbox-menu
UI Menu Using termbox
This library provides very basic terminal menu functionality
### To use:
import the following libraries

```go
import (
	"github.com/nsf/termbox-go"
	"github.com/vijayee/termbox-menu"
)
```

initialize termbox

``` go
err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
```

Now to create a menu. Menus takes an array of Items. An item is an interface type of the following structure.

```go
	type Item interface {
	Title() string
	Invoke() error
}
```
So long as it conforms it can be used as an item in the menu. Menus themselves conform to this interface and can thus be nested.
Menus are created by a call to the new menu function.

```go
	menu3 := menu.NewMenu("Even More Options", []Item{item7, item8}, termbox.ColorWhite, termbox.ColorBlue)
	menu2 := menu.NewMenu("More Options", []Item{item4, item5, item6, &menu3}, termbox.ColorWhite, termbox.ColorBlue)
	menu1 := menu.NewMenu("Main Menu", []Item{item1, item2, item3, &menu2}, termbox.ColorWhite, termbox.ColorBlue)
```

Finally, kick off the key service and Invoke the top level menu.

```go	
	go menu.ListenToKeys()
	menu1.Invoke()
```