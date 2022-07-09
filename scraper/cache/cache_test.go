package cache

import (
	"testing"
	"time"

	"github.com/pestanko/miniscrape/scraper/config"
	"github.com/stretchr/testify/assert"
)

func TestCreateACacheInstanceForDisabledIsNil(t *testing.T) {
	assert := assert.New(t)

	now := time.Now()
	cache := NewCache(config.CacheCfg{
		Enabled: false,
		Update:  false,
		Root:    "",
	}, now)

	assert.Nil(cache)
}
