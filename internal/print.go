package internal

import "fmt"

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf(InfoColor, fmt.Sprintf(format, args...) + "\n")
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	fmt.Printf(WarningColor, fmt.Sprintf(format, args...) + "\n")
}

// Error should be used to display a warning
func Error(format string, args ...interface{}) {
	fmt.Printf(ErrorColor, fmt.Sprintf(format, args...) + "\n")
}