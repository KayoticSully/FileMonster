package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/KayoticSully/gocui"
)

func SettingsEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	labelview, _ := GUI.View("settings-labels")
	cx, cy := v.Cursor()
	ox, oy := v.Origin()
	line, _ := v.Line(cy)

	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)

	case key == gocui.KeySpace:
		v.EditWrite(' ')

	case key == gocui.KeyTab:
		searchDir := filepath.Dir(line)
		partialDir := filepath.Base(line)

		// get a list of all files at the current search path
		files, err := ioutil.ReadDir(searchDir)
		if err != nil {
			break
		}

		// Find all partially matching directories
		matchingFiles := []string{}
		for _, f := range files {
			if strings.HasPrefix(f.Name(), partialDir) {
				matchingFiles = append(matchingFiles, f.Name())
			}
		}

		// Pick the replacement directory or partial directory based
		// on what was found
		var foundDir string
		if len(matchingFiles) == 1 {
			foundDir = matchingFiles[0]
		} else if len(matchingFiles) > 1 {
			foundDir = FindLongestCommonPrefix(matchingFiles)
		}

		// If something was found replace the current line
		if foundDir != "" {
			newline := strings.Trim(filepath.Join(searchDir, foundDir), " ")
			v.SetLine(cy, newline)
			// TODO: take origin into account for cursor placement
			lineLength := utf8.RuneCountInString(newline)

			if lineLength > 31 {
				err = v.SetOrigin(lineLength-30, oy)
				if err != nil {
					gLog(err.Error())
				}

				err = v.SetCursor(30, cy)
				if err != nil {
					gLog(err.Error())
				}
			} else {
				err = v.SetCursor(lineLength, cy)
				if err != nil {
					gLog(err.Error())
				}
			}
		}

	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		// Keep the cursor on the line
		if cx+ox > 0 {
			v.EditDelete(true)
		}

	case key == gocui.KeyDelete:
		// Keep the cursor on the line
		if cx+ox < utf8.RuneCountInString(line) {
			v.EditDelete(false)
		}

	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite

	case key == gocui.KeyArrowDown:
		// Only move cursor if there is another setting
		if cy+oy < NumOfSettings-1 {
			v.MoveCursor(0, 1, false)

			// make sure the correct label is hilighted
			labelview.MoveCursor(0, 1, false)
		} else if cy+oy == NumOfSettings-1 {
			selectStart(GUI, v)
			//gLog("HERE")
		}

	case key == gocui.KeyArrowUp:
		// Only move cursor if there is another setting
		if cy+oy >= NumOfSettings-1 {
			v.MoveCursor(0, -1, false)

			// make sure the correct label is hilighted
			labelview.MoveCursor(0, -1, false)
		}

	case key == gocui.KeyArrowLeft:
		// Keep the cursor on the line
		if cx+ox > 0 {
			v.MoveCursor(-1, 0, false)
		}

	case key == gocui.KeyArrowRight:
		// Keep the cursor on the line
		if cx+ox < utf8.RuneCountInString(line) {
			v.MoveCursor(1, 0, false)
		}
	}
}
