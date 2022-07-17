package filters

// PageFilter main interface for all the filters
type PageFilter interface {
	// Filter would filter the content for page
	Filter(content string) (string, error)
	// IsEnabled whether the filter is enabled or not
	IsEnabled() bool
	// Name of the filter
	Name() string
}
