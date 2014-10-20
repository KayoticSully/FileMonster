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

func main() {
	flag.Parse()
	source := flag.Arg(0)
	target := flag.Arg(1)

	// Todo: If no source or target prompt user

	var numWorkers = runtime.NumCPU()

	gui := SetupGUI()
	defer gui.Close()

	SetupKeyEvents(gui)
	SetupLayout(gui, numWorkers)

	// Start Processing
	inFiles, filesFound := GoWalk(source, numWorkers)
	_, filesProcessed := StartWorkers(inFiles, target, gui, numWorkers)

	// Update all other views
	go func() {
		view := gui.View("stats")

		for {
			view.Clear()
			total := sum(filesProcessed)

			fmt.Fprintf(view, "Total Files Processed: %s\n", humanize.Comma(total))
			fmt.Fprintf(view, "Total Files Found:     %s", humanize.Comma(*filesFound))

			time.Sleep(gui.FPSDelay(60))
		}
	}()

	gui.RefreshLoop(60)
}

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
