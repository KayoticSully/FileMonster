package main

import (
	"fmt"
	"unicode/utf8"

	"github.com/KayoticSully/gocui"
)

func sum(arr []int64) int64 {
	var sum int64

	for num := range arr {
		sum = sum + arr[num]
	}

	return sum
}

func gLog(str string) error {
	var v *gocui.View
	var err error

	if v, err = GUI.View("log"); err != nil {
		return err
	}

	fmt.Fprintln(v, str)

	return nil
}

func FindLongestCommonPrefix(data []string) string {
	var found = false
	var length = 0
	var prefix = ""

	// Loop until the longest common prefix has been found
	for !found {
		// Try the next longest string
		length = length + 1
		// TODO: Fix issue when first string is shortest
		for _, str := range data {
			// if any string's full length has been met, the
			// longest prefix has been found
			if utf8.RuneCountInString(str) < length {
				prefix = prefix[:length-1]
				found = true
				break
			}

			prefix = data[0][:length]
			if str[:length] != prefix {
				prefix = prefix[:length-1]
				found = true
				break
			}
		}
	}

	return prefix
}
