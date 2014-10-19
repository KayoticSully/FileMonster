package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/KayoticSully/gocui"
	"github.com/dustin/go-humanize"
)

func Worker(id int, files <-chan FileData, target string, wg *sync.WaitGroup, filesProcessed []int64, gui *gocui.Gui) {
	view := gui.View("worker" + strconv.Itoa(id))

	for file := range files {
		view.Clear()
		filesProcessed[id] = filesProcessed[id] + 1
		fmt.Fprintf(view, "%-5s%-15s%s", strconv.Itoa(id), humanize.Comma(filesProcessed[id]), file.Name())

		var err error
		// Filter out directories here so new features are easier to add
		// in the future dealing with processing directories
		//if !file.IsDir() {
			err = processFile(file, target)
		//}

		if err != nil {
			log.Fatal(err)
		}
	}

	view.Clear()
	fmt.Fprintf(view, "%-5s%-15s%s", strconv.Itoa(id), humanize.Comma(filesProcessed[id]), "Done!")

	wg.Done()
}

func processFile(file FileData, target string) error {
	fileParts := strings.Split(file.Name(), ".")

	baseName := strings.Join(fileParts[0:len(fileParts)-1], ".")
	ext := fileParts[len(fileParts)-1]

	targetFolder := target + "/" + ext
	targetPath := targetFolder + "/" + file.Name()

	// Make sure folder exists
	os.MkdirAll(targetFolder, os.ModeDir)

	// Logic It
	// Need to make sure the file does not exist yet
	var found = false
	var i int
	for !found {
		if _, err := os.Stat(targetPath); err == nil {
			i++
			if strings.TrimSpace(baseName) == "" {
				targetPath = targetFolder + "/" + "." + ext + "(" + strconv.Itoa(i) + ")"
			} else {
				targetPath = targetFolder + "/" + baseName + "(" + strconv.Itoa(i) + ")." + ext
			}
		} else {
			found = true
		}
	}

	return os.Rename(file.Path(), targetPath)
}
