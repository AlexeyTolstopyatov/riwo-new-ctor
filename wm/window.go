/*
 * Window struct for WM
 * and its functions + listeners
 */

package wm

import (
	"bytes"
	"html/template"
	"strconv"
	"syscall/js"
)

// ContextEntry
// Defines a named type for the context entry.
type ContextEntry struct {
	Name     string
	Callback func()
}

// Window
// manages single abstract window’s properties.
type Window struct {
	ID             int            // For the most part unites DOM object and Go object
	Element        js.Value       // Connected DOM element.
	ContextEntries []ContextEntry // Tho "Move", "Resize", "Delete" and "Hide" are basic ones
}

// WindowCreate
// creates a new Window,
// sets up its DOM element,
// and returns a pointer to it.
func WindowCreate(x, y, width, height, content string) *Window {
	WindowCount++
	id := WindowCount

	document := js.Global().Get("document")
	body := document.Get("body")

	// Create the DOM element for the window.
	winElem := document.Call("createElement", "div")
	style := "overflow: hidden; position: absolute; z-index: " + strconv.Itoa(HighestZIndex) + "; left: " + x +
		"; top: " + y + "; width: " + width + "; height: " + height +
		"; background-color: #f0f0f0; border: solid #55AAAA; padding: 0;"
	winElem.Set("style", style)
	winElem.Set("innerHTML", content)
	winElem.Set("id", strconv.Itoa(id)) // Assing shared ID

	body.Call("appendChild", winElem)

	// Logging
	if Verbose {
		Print("Generated window's ID (wid) is \"" +
			strconv.Itoa(id) + "\"")
	}

	newWindow := &Window{
		ID:      id,
		Element: winElem,
		// No custom ContextEntries
	}
	CurrentWindow = newWindow
	ActiveWindow = winElem
	AllWindows[strconv.Itoa(newWindow.ID)] = newWindow

	// Bring to front when clicked
	winElem.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if !IsResizingInit {
			CurrentWindow = newWindow
			ActiveWindow = winElem
		}
		// Right-click (RMB) on the window to select it
		// for resizing, second right-click activates resizing
		if IsResizingMode && !IsResizingInit && args[0].Get("button").Int() == 2 {
			// First RMB hold - Select the window for resizing
			args[0].Call("preventDefault")
			args[0].Call("stopPropagation")
			JustSelected = true
			if Verbose {
				Print("First right-click: Window selected for resizing.")
			}
			HighestZIndex++
			winElem.Get("style").Set("z-index", strconv.Itoa(HighestZIndex))
			IsResizingInit = true
			js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor-selection.svg) 12 12, auto")
		}
		// Mouse down event for selecting and
		// dragging the window (click brings it to front)
		if !IsResizingInit {
			HighestZIndex++
			winElem.Get("style").Set("z-index", strconv.Itoa(HighestZIndex))
			if Verbose {
				Print("Window brought to front.")
			}

			if IsMovingMode && args[0].Get("button").Int() == 2 {
				args[0].Call("preventDefault")
				args[0].Call("stopPropagation")

				//JustSelected = true
				StartX = args[0].Get("clientX").Float() - winElem.Get("offsetLeft").Float()
				StartY = args[0].Get("clientY").Float() - winElem.Get("offsetTop").Float()
				IsDragging = true

				// Create ghost window
				GhostWindow = document.Call("createElement", "div")
				rect := winElem.Call("getBoundingClientRect")
				width := rect.Get("width").Float()
				height := rect.Get("height").Float()

				// Ensure ghost window is above
				// everything during drag
				GhostWindow.Set("style", "position: absolute; z-index: "+strconv.Itoa(HighestZIndex+1)+"; width: "+Ftoa(width)+"px; height: "+Ftoa(height)+"px; border: solid 2px #FF0000; cursor: url(assets/cursor-drag.svg) 12 12, auto;")
				GhostWindow.Get("style").Set("left", Ftoa(winElem.Get("offsetLeft").Float())+"px")
				GhostWindow.Get("style").Set("top", Ftoa(winElem.Get("offsetTop").Float())+"px")
				body.Call("appendChild", GhostWindow)
				JustSelected = true

				if Verbose {
					Print("Dragging initiated with ghost window.")
				}
			}
			if IsHiding && args[0].Get("button").Int() == 2 {
				// Hide window
				args[0].Call("preventDefault")
				args[0].Call("stopPropagation")
				JustSelected = true
				IsHiding = false

				hiddenWindowOption := CreateMenuOption("wid " + strconv.Itoa(newWindow.ID))
				if winElem.Get("title").String() != "" {
					hiddenWindowOption = CreateMenuOption(winElem.Get("title").String())
				}

				hiddenWindowOption.Set("id", "menuopt"+strconv.Itoa(newWindow.ID))
				// Show option activation
				hiddenWindowOption.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					if args[0].Get("button").Int() == 2 {
						args[0].Call("preventDefault")
						args[0].Call("stopPropagation")
						JustSelected = true
						RemoveMenuOption(hiddenWindowOption)
						winElem.Get("style").Set("display", "block")
						ContextMenu.Get("style").Set("display", "none")
						// Delete by value
						for index, value := range ContextMenuHides {
							if value.Get("id").String() == hiddenWindowOption.Get("id").String() {
								ContextMenuHides = append(ContextMenuHides[:index], ContextMenuHides[index+1:]...)
							}
						}
						JustSelected = false
						if Verbose {
							Print("activated.")
						}
					}
					return nil
				}))
				ContextMenuHides = append(ContextMenuHides, hiddenWindowOption)
				winElem.Get("style").Set("display", "none")
				js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor.svg), auto")
				JustSelected = false
				if Verbose {
					Print("WID " + strconv.Itoa(newWindow.ID) + " hidden")
				}

			}
		}
		// Right-click (RMB) deletes the window in delete mode
		if IsDeleteMode && args[0].Get("button").Int() == 2 {
			args[0].Call("preventDefault")
			args[0].Call("stopPropagation")
			JustSelected = true
			WindowRemove(newWindow)
			IsDeleteMode = false
			js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor.svg), auto")
			JustSelected = false
			if Verbose {
				Print("Window deleted.")
			}
		}
		return nil
	}))

	return newWindow
}

// SetPosition
// Sets positions and dimensions for the window. (Actually useless)
func (w *Window) SetPosition(newX, newY, newWidth, newHeight string) {
	w.Element.Get("style").Set("left", newX)
	w.Element.Get("style").Set("top", newY)
	w.Element.Get("style").Set("width", newWidth)
	w.Element.Get("style").Set("height", newHeight)
}

// WindowRemove
// Deletes the window from DOM and Go.
func WindowRemove(w *Window) {
	w.ID = -1                              // Remove reference for apps
	w.Element.Call("remove")               // Remove html part
	delete(AllWindows, strconv.Itoa(w.ID)) // Remove from list
	return
}

// Render
// Returns
func (w *Window) Render(title string) string {
	tmpl := `
    <div class="window" id="{{.ID}}" 
         style="position: absolute; left: {{.X}}px; top: {{.Y}}px">
        <div class="title-bar">{{.Title}}</div>
        <div class="content">{{.Content}}</div>
        <button data-action="close" data-target="{{.ID}}">×</button>
    </div>
    `
	t := template.Must(template.New(title).Parse(tmpl))
	var buf bytes.Buffer
	err := t.Execute(&buf, w)

	if err != nil {
		return "errors at renderer"
	}
	return buf.String()
}
