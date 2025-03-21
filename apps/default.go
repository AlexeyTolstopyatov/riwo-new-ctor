package apps

import (
	"riwo/wm"
	"syscall/js"
)

// init
// Registers Application in Window Manager
func init() {
	AppRegistry["Default"] = AppDefault
}

// AppDefault
// Like Windows API logic - registers the Window of
// this application and all interaction logic too.
func AppDefault(window *wm.Window) {
	document := js.Global().Get("document")
	themeColor := "green"
	colorFG := wm.GetColor[themeColor]["vivid"]
	colorBG := wm.GetColor[themeColor]["faded"]
	colorMG := wm.GetColor[themeColor]["normal"]

	// Create a container div for the grid
	container := document.Call("createElement", "div")
	container.Get("style").Set("display", "grid")
	container.Get("style").Set("gridTemplateColumns", "repeat(auto-fit, minmax(120px, 1fr))")
	container.Get("style").Set("gap", "5%")
	container.Get("style").Set("padding", "5%")
	container.Get("style").Set("height", "100%")
	container.Get("style").Set("background", colorBG)

	// Create a title
	title := document.Call("createElement", "div")
	title.Set("innerHTML", "Applications")
	title.Get("style").Set("gridColumn", "1 / -1")
	title.Get("style").Set("fontSize", "24px")
	title.Get("style").Set("textAlign", "center")
	title.Get("style").Set("marginBottom", "20px")
	title.Get("style").Set("color", colorFG)
	container.Call("appendChild", title)

	// Iterate over AppRegistry
	// and create a button for each app
	for appName, appFunc := range AppRegistry {
		/*if appName == "Default" {
		    continue
		}*/ // This can be used to skip default app itself

		btnContainer := document.Call("createElement", "div")
		btnContainer.Get("style").Set("textAlign", "center")

		btn := document.Call("createElement", "div")
		btn.Set("innerText", appName)
		btn.Get("style").Set("cursor", "url(assets/cursor-inverted.svg), auto")
		btn.Get("style").Set("padding", "15px")
		btn.Get("style").Set("background", colorBG)
		btn.Get("style").Set("color", "black")
		btn.Get("style").Set("borderRadius", "0")
		btn.Get("style").Set("border", "solid "+colorMG)
		btn.Get("style").Set("transition", "all 0.2s ease")
		btn.Get("style").Set("userSelect", "none")

		// Hover effects
		btn.Call("addEventListener", "mouseover", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			btn.Get("style").Set("background", colorFG)
			btn.Get("style").Set("color", colorBG)
			return nil
		}))
		btn.Call("addEventListener", "mouseout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			btn.Get("style").Set("background", colorBG)
			btn.Get("style").Set("color", "black")
			return nil
		}))

		// Click handler
		btn.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if wm.Verbose {
				wm.Print("App " + appName + " selected")
			}
			appFunc(window)
			return nil
		}))

		btnContainer.Call("appendChild", btn)
		container.Call("appendChild", btnContainer)
	}

	window.Element.Set("innerHTML", "")
	window.Element.Call("appendChild", container)
}
