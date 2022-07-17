package filters

type PageFilter interface {
	Filter(content string) (string, error)
	IsEnabled() bool
	Name() string
}
