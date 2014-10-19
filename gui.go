package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/KayoticSully/gocui"
)

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

func SetupKeyEvents(gui *gocui.Gui) {
	var err error
	if err = gui.SetKeybinding("", gocui.KeyCtrlC, 0, quit); err != nil {
		log.Panicln(err)
	}
}

func SetupLayout(gui *gocui.Gui, n int) {
	gui.SetLayout(func(g *gocui.Gui) error {
		maxX, maxY := g.Size()

		// Logo
		if v, err := gui.SetView("logo", maxX/2-26, maxY/2-16, maxX/2+28, maxY/2-9); err != nil {
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

		// Stats
		if v, err := gui.SetView("stats", maxX/2-25, maxY/2-9, maxX/2+25, maxY/2-6); err != nil {
			if err != gocui.ErrorUnkView {
				return err
			}

			v.Frame = false
		}

		// Workers
		if v, err := gui.SetView("workers", maxX/2-25, maxY/2-6, maxX/2+25, maxY/2-4); err != nil {
			if err != gocui.ErrorUnkView {
				return err
			}

			v.Frame = false
			fmt.Fprintln(v, "Workers")
		}

		// Headings
		if v, err := gui.SetView("headings", maxX/2-25, maxY/2-4, maxX/2+25, maxY/2-2); err != nil {
			if err != gocui.ErrorUnkView {
				return err
			}

			fmt.Fprintf(v, "%-5s%-15s%s", "#", "Processed", "Current File")
		}

		for i := 0; i < n; i++ {
			if v, err := gui.SetView("worker"+strconv.Itoa(i), maxX/2-25, (maxY/2-2)+(i*2), maxX/2+25, (maxY/2)+(i*2)); err != nil {
				if err != gocui.ErrorUnkView {
					return err
				}

				v.Frame = false
			}
		}

		return nil
	})
}

/* Event Functions */
func quit(gui *gocui.Gui, view *gocui.View) error {
	return gocui.ErrorQuit
}
