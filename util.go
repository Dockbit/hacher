package main

import (
	"fmt"
	"os"

	"github.com/mgutz/ansi"
)

func printInfo(message string, args ...interface{}) {
	logger.Println(colorizeMessage("green", "info:", message, args...))
}

func printFatal(message string, args ...interface{}) {
	logger.Println(colorizeMessage("red", "error:", message, args...))
	os.Exit(1)
}

func checkError(err error) {
	if err != nil {
		printFatal(err.Error())
	}
}

func colorizeMessage(color, prefix, message string, args ...interface{}) string {
	prefResult := ""
	if prefix != "" {
		prefResult = ansi.Color(prefix, color+"+b") + " " + ansi.ColorCode("reset")
	}
	return prefResult + ansi.Color(fmt.Sprintf(message, args...), color) + ansi.ColorCode("reset")
}
