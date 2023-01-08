package apprun

// EmptyDeps represents an empty dependency instance
type EmptyDeps struct{}

// Close the empty dependencies
func (e EmptyDeps) Close() error {
	return nil
}

// NewEmptyDeps return a new instance of the empty dependencies
func NewEmptyDeps() *EmptyDeps {
	return new(EmptyDeps)
}
