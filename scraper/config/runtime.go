package config

// RunResultStatus enum of runtime statuses
type RunResultStatus string

const (
	// RunSuccess status OK
	RunSuccess RunResultStatus = "ok"
	// RunError status ERROR
	RunError RunResultStatus = "error"
	// RunEmpty status EMPTY
	RunEmpty RunResultStatus = "empty"
)

// RunResult representation of the run result
type RunResult struct {
	// Page instance for which the result is
	Page Page
	// Content of the page
	Content string
	// Status of run
	Status RunResultStatus
	// Kind of the result
	Kind string
}
