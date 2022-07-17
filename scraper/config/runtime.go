package config

type RunResultStatus string

const (
	RunSuccess RunResultStatus = "ok"
	RunError   RunResultStatus = "error"
	RunEmpty   RunResultStatus = "empty"
)

type RunResult struct {
	Page    Page
	Content string
	Status  RunResultStatus
	Kind    string
}
