package internal

const (
	Success          = iota // 0 -> Indicates exiting execution successfully
	Unknown                 // 1 -> Catchall for general errors
	UnknownFlag             // 2 -> Flag provided doesn't exist
	InvalidFlagValue        // 3 -> Flag value failed validation
)
