package main

import (
	"fmt"

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
