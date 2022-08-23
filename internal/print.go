package internal

import "fmt"

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

// PrintInfo should be used to describe the example commands that are about to run.
func PrintInfo(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(InfoColor, fmt.Sprintf(format, args...)))
}

// PrintNotice is for upcoming changes/deprecations
func PrintNotice(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(NoticeColor, fmt.Sprintf(format, args...)))
}

// PrintWarning should be used to display a warning
func PrintWarning(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(WarningColor, fmt.Sprintf(format, args...)))
}

// PrintError should be used to display a warning
func PrintError(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(ErrorColor, fmt.Sprintf(format, args...)))
}

// PrintDebug is to be used when run in --verbose mode for collecting crash info
func PrintDebug(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(DebugColor, fmt.Sprintf(format, args...)))
}
