package main

import (
	"fmt"
	"os"

	"cli/display"
)

func main() {
	// ensure correct usage
	if len(os.Args) < 2 {
		errAndExit("Usage: " + display.Blue + os.Args[0] + " [command] [options]" + display.Normal + "\n\nPossible commands are:" + display.Yellow + "\n   create")

	}

	switch cmd := os.Args[1]; cmd {
	case "create":

	default:
		errAndExit("Unrecognized command '" + cmd + "'" + display.Normal + "\n\nPossible commands are: " + display.Yellow + "\n   create")
	}

}

func check(err error) {
	if err != nil {
		errAndExit(err.Error())
	}
}

func errAndExit(msg string) {
	fmt.Println(display.ZephyrError(msg) + "\n\nRun 'zephyr help [cmd]' for more info.\n")
	os.Exit(1)
}
