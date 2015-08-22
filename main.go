package main

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/KayoticSully/gocui"
)

const (
	NumOfSettings = 2
)

var (
	GUI         *gocui.Gui
	NUM_WORKERS int
)

// main is the point of execution entry for FileMonster
func main() {
	flag.Parse()

	// Get any options set via cmdline
	source := flag.Arg(0)
	target := flag.Arg(1)

	// Gather system specific info
	NUM_WORKERS = runtime.NumCPU()
	NUM_WORKERS = 1

	gui := SetupGUI(source, target)
	GUI = gui
	defer gui.Close()

	// Give initial control over to settings
	// This needs to be a goroutine so the MainLoop can get setup
	// before this executes
	go gui.SetCurrentView("settings")

	// Start the main application
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
