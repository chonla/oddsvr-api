package logger

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

// Info prints info log
func Info(msg string) {
	log(color.New(color.FgYellow, color.Bold).SprintFunc()("INFO"), msg)
}

// Debug prints debug log
func Debug(msg string) {
	log(color.New(color.FgMagenta).SprintFunc()("DEBUG"), msg)
}

// Error prints error log
func Error(msg string) {
	log(color.New(color.FgRed, color.Bold).SprintFunc()("ERROR"), msg)
}

func log(level, msg string) {
	now := time.Now()
	stamp := now.Format("02-01-2006 03:04:05PM")
	colorizedStamp := color.New(color.FgCyan).SprintFunc()(stamp)

	fmt.Printf("[%s] %s: %s\n", level, colorizedStamp, msg)
}
