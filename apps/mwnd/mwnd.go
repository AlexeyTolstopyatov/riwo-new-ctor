package mwnd

import (
	"riwo/apps"
	"riwo/wm"
)

func init() {
	apps.AppRegistry["AppMyWindow"] = AppMyWindow
}

// AppMyWindow
// Describes handlers and other backend of application
func AppMyWindow(window *wm.Window) {
	// Window registration process.
	window.Render("My Window")

	// callbacks
	window.SetPosition("200", "200", "400", "200")
}
