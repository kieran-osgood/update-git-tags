package internal

const (
	Success          = iota // Indicates exiting execution successfully
	Unknown                 // Catchall for general errors
	UnknownFlag             // Flag provided doesn't exist
	InvalidFlagValue        // Flag value failed validation
)
