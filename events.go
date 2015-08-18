package main

import (
	"fmt"

	"github.com/KayoticSully/gocui"
)

func quit(gui *gocui.Gui, view *gocui.View) error {
	return gocui.Quit
}

func selectSettingTarget(g *gocui.Gui, v *gocui.View) error {
	resetView(v)

	g.ShowCursor = true
	setSettingLabelIndex(g, 2)
	return selectView(g, "settings-target")
}

func selectSettingSource(g *gocui.Gui, v *gocui.View) error {
	resetView(v)
	g.ShowCursor = true
	setSettingLabelIndex(g, 1)
	return selectView(g, "settings-source")
}

func selectStart(g *gocui.Gui, v *gocui.View) error {
	var err error

	g.ShowCursor = false
	resetView(v)

	// reset labels view
	if v, err = g.View("settings-labels"); err != nil {
		return err
	}
	resetView(v)

	return selectView(g, "start")
}

func setSettingLabelIndex(g *gocui.Gui, i int) error {
	var v *gocui.View
	var err error

	if v, err = g.View("settings-labels"); err != nil {
		return err
	}

	// Make sure hilight is enabled so we can see that it is selected
	v.Highlight = true

	if err = v.SetCursor(0, i); err != nil {
		return err
	}

	if err = v.SetOrigin(0, 0); err != nil {
		return err
	}

	return nil
}

func resetView(v *gocui.View) error {
	v.Highlight = false

	cx, cy := v.Cursor()
	ox, oy := v.Origin()

	gLog(fmt.Sprintf("cx: %d, cy: %d, ox: %d, oy: %d", cx, cy, ox, oy))

	if err := v.SetCursor(0, 0); err != nil {
		return err
	}

	if err := v.SetOrigin(0, 0); err != nil {
		return err
	}

	return nil
}

func stats(g *gocui.Gui, v *gocui.View) error {

	cx, cy := v.Cursor()
	ox, oy := v.Origin()

	gLog(fmt.Sprintf("cx: %d, cy: %d, ox: %d, oy: %d", cx, cy, ox, oy))

	return nil
}

func selectView(g *gocui.Gui, viewName string) error {
	var view *gocui.View
	var err error

	// get view
	if view, err = g.View(viewName); err != nil {
		return err
	}

	// set active settings
	view.Highlight = true

	// set as active
	return g.SetCurrentView(viewName)
}

func validateSetting(g *gocui.Gui, v *gocui.View) error {
	return nil
}

func selectLineDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		ox, oy := v.Origin()

		if err := v.SetCursor(cx, cy+1); err != nil {
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}

	return nil
}

func selectLineUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()

		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}

	return nil
}

// DefaultEditor is used as the default gocui editor.
func SettingsEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case key == gocui.KeyEnter:
		v.EditNewLine()
	}
}
