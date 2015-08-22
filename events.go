package main

import (
	"fmt"
	"time"

	"github.com/KayoticSully/gocui"
	"github.com/dustin/go-humanize"
)

func quit(gui *gocui.Gui, view *gocui.View) error {
	return gocui.Quit
}

func selectSettings(g *gocui.Gui, v *gocui.View) error {
	resetView(v)

	g.ShowCursor = true
	setSettingLabelIndex(g, 1)

	return selectView(g, "settings")
}

func selectStart(g *gocui.Gui, v *gocui.View) error {
	var err error

	g.ShowCursor = false
	v.Highlight = false

	// reset labels view
	if v, err = g.View("settings-labels"); err != nil {
		return err
	}
	v.Highlight = false

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

	if err := v.SetCursor(0, 0); err != nil {
		return err
	}

	if err := v.SetOrigin(0, 0); err != nil {
		return err
	}

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
	err = g.SetCurrentView(viewName)
	if err != nil {
		return err
	}

	return err
}

func start(g *gocui.Gui, v *gocui.View) error {
	// Disable all input
	resetView(v)
	g.SetCurrentView("logo")

	// Get Final Settings
	var settings *gocui.View
	var source, target string
	var err error

	if settings, err = g.View("settings"); err != nil {
		return err
	}

	if source, err = settings.Line(0); err != nil {
		return err
	}

	if target, err = settings.Line(1); err != nil {
		return err
	}

	// Start Processing
	inFiles, filesFound := GoWalk(source, NUM_WORKERS)
	_, filesProcessed := StartWorkers(inFiles, target, g, NUM_WORKERS)

	view, err := g.View("stats")
	if err != nil {
		return err
	}

	go func() {
		for {
			view.Clear()
			total := sum(filesProcessed)

			fmt.Fprintf(view, "Number of Workers: %d\n", NUM_WORKERS)
			fmt.Fprintf(view, "Total Files Processed: %s\n", humanize.Comma(total))
			fmt.Fprintf(view, "Total Files Found:     %s", humanize.Comma(*filesFound))

			time.Sleep(FPSDelay(60))
		}
	}()

	return nil
}
