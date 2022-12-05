package main

import "fmt"

type Color string

const (
	Reset   Color = "\033[0m"
	Black   Color = "\033[1;30m"
	Red     Color = "\033[1;31m"
	Green   Color = "\033[1;32m"
	Yellow  Color = "\033[1;33m"
	Purple  Color = "\033[1;34m"
	Magenta Color = "\033[1;35m"
	Teal    Color = "\033[1;36m"
	White   Color = "\033[1;37m"
)

func PrintColor(msg string, color Color) {
	fmt.Print(color, msg, Reset)
}
