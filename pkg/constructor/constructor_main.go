package constructor

import (
	"riwo/wm"
	"strings"
	"syscall/js"
)

func renderAll(windows []wm.Window) {
	var sb strings.Builder
	for _, w := range windows {
		sb.WriteString(w.Render("My Window"))
	}

	js.Global().Get("document").
		Call("getElementById", "app").
		Set("innerHTML", sb.String())
}
