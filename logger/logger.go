// Package logger implements an application log settings
package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	Error   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
)

type Color string

const (
	ColorRed    Color = "\u001b[31m"
	ColorGreen  Color = "\u001b[32m"
	ColorYellow Color = "\u001b[33m"
	ColorBlue   Color = "\u001b[34m"
	ColorReset  Color = "\u001b[0m"
)

// Init setup application logger settings
func Init() {
	Error = log.New(
		os.Stdout,
		colorize(ColorRed, "ERROR: "),
		log.Ldate|log.Ltime,
	)
	Info = log.New(
		os.Stdout,
		colorize(ColorBlue, "INFO: "),
		log.Ldate|log.Ltime,
	)
	Warning = log.New(
		os.Stdout,
		colorize(ColorYellow, "WARNING: "),
		log.Ldate|log.Ltime,
	)
}

// colorize return string with selected color
func colorize(color Color, message string) string {
	return fmt.Sprintf("%s %s %s", string(color), message, string(ColorReset))
}
