package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/KayoticSully/gocui"
)

const (
	logoViewHeight     = 7 // x + 1
	settingsViewHeight = 4
)

// SetupGUI creates and initalizes a gocui gui object
func SetupGUI() *gocui.Gui {
	var err error

	// Create GUI
	gui := gocui.NewGui()

	// Init
	if err = gui.Init(); err != nil {
		log.Panicln(err)
	}

	return gui
}

// SetupKeyEvents ties keyboard events to actions within the GUI
func SetupKeyEvents(gui *gocui.Gui) {
	var err error

	// Quit
	if err = gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	// Move down a setting
	if err := gui.SetKeybinding("settings-source", gocui.KeyArrowDown, gocui.ModNone, selectSettingTarget); err != nil {
		log.Panicln(err)
	}

	// Print Stats
	if err := gui.SetKeybinding("settings-source", gocui.KeyEnd, gocui.ModNone, stats); err != nil {
		log.Panicln(err)
	}

	// Move up a settings
	if err := gui.SetKeybinding("settings-target", gocui.KeyArrowUp, gocui.ModNone, selectSettingSource); err != nil {
		log.Panicln(err)
	}

	// Print Stats
	if err := gui.SetKeybinding("settings-target", gocui.KeyEnd, gocui.ModNone, stats); err != nil {
		log.Panicln(err)
	}

	// Move down a setting
	if err := gui.SetKeybinding("settings-target", gocui.KeyArrowDown, gocui.ModNone, selectStart); err != nil {
		log.Panicln(err)
	}

	// Validate setting
	if err := gui.SetKeybinding("settings", gocui.KeyEnter, gocui.ModNone, validateSetting); err != nil {
		log.Panicln(err)
	}

	// Move up to the settings view
	if err := gui.SetKeybinding("start", gocui.KeyArrowUp, gocui.ModNone, selectSettingTarget); err != nil {
		log.Panicln(err)
	}
}

// SetupLayout creates the interface objects used within the GUI
// Note: May be able to remove the SetupLayout function call, and just
// use a function called Layout() and call gui.SetLayout from mainline
func SetupLayout(gui *gocui.Gui, n int) {
	gui.SetLayout(func(g *gocui.Gui) error {
		maxX, maxY := g.Size()

		// Use half as base to center everything
		halfX := maxX / 2
		halfY := maxY / 2

		// DevLog Heading
		if v, err := gui.SetView("log-label", halfX-70, halfY-15, halfX-30, halfY-13); err != nil {
			if err != gocui.ErrorUnkView {
				return err
			}

			v.Frame = true
			fmt.Fprintln(v, "Dev Log")
		}

		// DevLog
		if v, err := gui.SetView("log", halfX-70, halfY-13, halfX-30, halfY+12); err != nil {
			if err != gocui.ErrorUnkView {
				return err
			}

			v.Frame = true
		}

		// Logo
		if v, err := gui.SetView("logo", halfX-26, halfY-16, halfX+28, halfY-9); err != nil {
			if err != gocui.ErrorUnkView {
				return err
			}

			v.Frame = false
			fmt.Fprintln(v, ` ______ _ _      __  __                 _            `)
			fmt.Fprintln(v, `|  ____(_) |    |  \/  |               | |           `)
			fmt.Fprintln(v, `| |__   _| | ___| \  / | ___  _ __  ___| |_ ___ _ __ `)
			fmt.Fprintln(v, `|  __| | | |/ _ \ |\/| |/ _ \| '_ \/ __| __/ _ \ '__|`)
			fmt.Fprintln(v, `| |    | | |  __/ |  | | (_) | | | \__ \ ||  __/ |   `)
			fmt.Fprintln(v, `|_|    |_|_|\___|_|  |_|\___/|_| |_|___/\__\___|_|   `)
		}

		// Settings Labels
		if v, err := gui.SetView("settings-labels", halfX-25, halfY-9, halfX-7, halfY-5); err != nil {
			if err != gocui.ErrorUnkView {
				return err
			}

			v.Highlight = true
			v.Frame = false

			// Hilight the initial setting that is active upon start
			if err := v.SetCursor(0, 1); err != nil {
				return err
			}

			if err := v.SetOrigin(0, 0); err != nil {
				return err
			}
		}

		// Settings Fields
		// Note/Todo: Create seperate view per setting field.  One view for the labels is fine, but
		// the fields need to be editable halfY....  this also "simplifies" the logic for
		// switching between settings.... many functions, but very straightforward.
		// Also... SLEEP on this... might come up with a way to do this dynamically :)
		// Settings Labels
		if v, err := gui.SetView("settings-source", halfX-6, halfY-8, halfX+25, halfY-6); err != nil {
			if err != gocui.ErrorUnkView {
				return err
			}

			v.Highlight = true
			v.Frame = false
			v.Editable = true
			v.Wrap = false
			v.Autoscroll = false

			if err := v.SetCursor(18, 1); err != nil {
				if err := v.SetOrigin(18, 1); err != nil {
					return err
				}
			}
		}

		if v, err := gui.SetView("settings-target", halfX-6, halfY-7, halfX+25, halfY-5); err != nil {
			if err != gocui.ErrorUnkView {
				return err
			}

			v.Highlight = true
			v.Frame = false
			v.Editable = true
			v.Wrap = false
			v.Autoscroll = false

			if err := v.SetCursor(18, 1); err != nil {
				if err := v.SetOrigin(18, 1); err != nil {
					return err
				}
			}
		}

		// Start
		if v, err := gui.SetView("start", halfX-25, halfY-6, halfX+25, halfY-2); err != nil {
			if err != gocui.ErrorUnkView {
				return err
			}

			v.Frame = false
		}

		// Stats
		if v, err := gui.SetView("stats", halfX-25, halfY-4, halfX+25, halfY-2); err != nil {
			if err != gocui.ErrorUnkView {
				return err
			}

			v.Frame = false
		}

		// Workers
		if v, err := gui.SetView("workers", halfX-25, halfY-2, halfX+25, halfY); err != nil {
			if err != gocui.ErrorUnkView {
				return err
			}

			v.Frame = false
			fmt.Fprintln(v, "Workers")
		}

		// Headings
		if v, err := gui.SetView("headings", halfX-25, halfY, halfX+25, halfY+2); err != nil {
			if err != gocui.ErrorUnkView {
				return err
			}

			fmt.Fprintf(v, "%-5s%-15s%s", "#", "Processed", "Current File")
		}

		for i := 0; i < n; i++ {
			if v, err := gui.SetView("worker"+strconv.Itoa(i), halfX-25, (halfY+2)+(i*2), halfX+25, (halfY+4)+(i*2)); err != nil {
				if err != gocui.ErrorUnkView {
					return err
				}

				v.Frame = false
			}
		}

		return nil
	})
}
