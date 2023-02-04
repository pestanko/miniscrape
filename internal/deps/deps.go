package deps

import (
	"github.com/pestanko/miniscrape/internal/config"
)

// Deps Represents an application dependencies
type Deps struct {
	Cfg *config.AppConfig
}

// Close the dependencies
func (d Deps) Close() error {
	return nil
}

// InitAppDeps init the application dependencies
func InitAppDeps() (*Deps, error) {
	return &Deps{
		Cfg: config.GetAppConfig(),
	}, nil
}
