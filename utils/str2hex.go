package main

import (
	"fmt"
	"os"
	"runtime"
)

var (
	putf = fmt.Printf
)

func main() {

	if len(os.Args) > 1 {
		str := os.Args[1]
		putf("0x%x\n", str)

		if runtime.GOOS == "windows" {
			fmt.Scanln() // Enter Key to terminate the console screen
		}
	} else {
		os.Exit(0)
	}

}
