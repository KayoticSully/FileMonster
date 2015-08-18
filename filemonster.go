package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/KayoticSully/gocui"
	"github.com/dustin/go-humanize"
)

var GUI *gocui.Gui

// main is the point of execution entry for FileMonster
func main() {
	flag.Parse()

	// Get any options set via cmdline
	source := flag.Arg(0)
	target := flag.Arg(1)

	// Todo: If no source or target prompt user

	// Gather system specific info
	var numWorkers = runtime.NumCPU()

	gui := SetupGUI()
	GUI = gui
	defer gui.Close()

	gocui.Edit = gocui.EditorFunc(SettingsEditor)

	SetupKeyEvents(gui)
	SetupLayout(gui, numWorkers)

	gui.ShowCursor = true

	// Start Processing
	/*inFiles*/ _, filesFound := GoWalk(source, numWorkers)
	//_, filesProcessed := StartWorkers(inFiles, target, gui, numWorkers)

	// Update all other views
	go func() {
		gui.SetCurrentView("settings-source")
		view, _ := gui.View("settings-labels")
		fmt.Fprintf(view, "Number of Workers: %d\n", numWorkers)
		fmt.Fprintf(view, "Source Directory: %s\n", source)
		fmt.Fprintf(view, "Target Directory: %s\n", target)

		view, _ = gui.View("start")
		fmt.Fprintf(view, "Start\n")

		view, _ = gui.View("stats")

		for {
			view.Clear()
			total := int64(0) //sum(filesProcessed)

			fmt.Fprintf(view, "Total Files Processed: %s\n", humanize.Comma(total))
			fmt.Fprintf(view, "Total Files Found:     %s", humanize.Comma(*filesFound))

			time.Sleep(FPSDelay(60))
		}
	}()

	gui.MainLoop()
}

// GoWalk traverses the source directory tree and returns a FileData chan
// that every found file gets pushed into
func GoWalk(source string, n int) (<-chan FileData, *int64) {
	// Buffer channel with entries equal to the number of workers.  This ensures
	// maximum throughput efficiency
	out := make(chan FileData, n)
	var filesFound int64

	go func() {
		filepath.Walk(source, func(path string, f os.FileInfo, _ error) error {
			filesFound++
			out <- FileData{f, path}
			return nil
		})

		defer close(out)
	}()

	return out, &filesFound
}

// StartWorkers sets up and starts n number of worker threads to process the queue of files
func StartWorkers(files <-chan FileData, target string, gui *gocui.Gui, n int) (*sync.WaitGroup, []int64) {
	var wg sync.WaitGroup

	var filesProcessed []int64
	filesProcessed = make([]int64, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go Worker(i, files, target, &wg, filesProcessed, gui)
	}

	return &wg, filesProcessed
}

func FPSDelay(fps int) time.Duration {
	timeout := (1.0 / float64(fps)) * 1000.0
	return time.Duration(timeout) * time.Millisecond
}
