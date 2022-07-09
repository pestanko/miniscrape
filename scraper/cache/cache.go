package cache

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/pestanko/miniscrape/scraper/config"
)

const DefaultContentFile = "content.txt"

type NamespacePath interface {
	Path() string
}

type Item struct {
	Namespace   ItemNamespace
	FileName    string
	CachePolicy string
}

type ItemNamespace struct {
	Page     string
	Category string
}

func (n *ItemNamespace) Path() string {
	return path.Join(n.Category, n.Page)
}

func (n ItemNamespace) String() string {
	return fmt.Sprintf("%s/%s", n.Category, n.Page)
}

func NewNamespace(cat string, page string) ItemNamespace {
	return ItemNamespace{Category: cat, Page: page}
}

type Cache interface {
	Store(item Item, content []byte) error
	IsPageCached(nm ItemNamespace) bool
	IsItemCached(item Item) bool
	GetContent(item Item) []byte
	Invalidate(sel config.RunSelector)
}

func NewCache(cacheCfg config.CacheCfg, date time.Time) Cache {
	if cacheCfg.Enabled {
		root := cacheCfg.Root
		if root == "" {
			root = path.Join(os.TempDir(), "mini-scrape")
		}
		return &cacheFs{
			RootDir:     root,
			ForceUpdate: cacheCfg.Update,
			Date:        date,
		}
	}
	return nil
}

type cacheFs struct {
	BlockList   []string
	RootDir     string
	ForceUpdate bool
	Date        time.Time
}

func (c *cacheFs) Invalidate(sel config.RunSelector) {
	nm := NewNamespace(sel.Category, sel.Page)
	pth := c.getNamespaceDir(nm.Path())
	RemoveDir(pth)
}

func (c *cacheFs) IsItemCached(item Item) bool {
	return !c.ForceUpdate && IsPathExists(c.getFileForItem(item))
}

func (c *cacheFs) GetContent(item Item) []byte {
	fp := c.getFileForItem(item)
	content, err := ioutil.ReadFile(fp)

	if err != nil {
		log.Warn().
			Err(err).
			Str("file", fp).
			Str("type", "cache").
			Msg("CACHE: Unable to load content")

		return []byte{}
	}

	log.Trace().
		Str("file", fp).
		Str("type", "cache").
		Msg("CACHE: Loading cached content")

	return content
}

func (c *cacheFs) IsPageCached(nm ItemNamespace) bool {
	return !c.ForceUpdate && IsPathExists(c.getNamespaceDir(nm.Path()))
}

func (c *cacheFs) Store(item Item, content []byte) error {
	if item.CachePolicy == "no-cache" || item.CachePolicy == "no" {
		return nil
	}

	fp := c.getFileForItem(item)

	if !c.IsPageCached(item.Namespace) {
		if err := os.MkdirAll(c.getNamespaceDir(item.Namespace.Path()), 0700); err != nil {
			log.Error().
				Err(err).
				Str("file", fp).
				Str("type", "cache").
				Msg("CACHE: Unable to create directory file")

			return err
		}
	}

	if c.IsItemCached(item) {
		log.Trace().
			Str("item", item.Namespace.String()).
			Str("type", "cache").
			Msg("CACHE: Item already cache")
		return nil
	}

	log.Trace().
		Str("item", item.Namespace.String()).
		Str("file", fp).
		Str("type", "cache").
		Msg("CACHE: Writing cache file")

	if err := ioutil.WriteFile(fp, content, 0600); err != nil {
		log.Error().
			Err(err).
			Str("file", fp).
			Str("type", "cache").
			Msg("CACHE: Unable to write file")
		return err
	}

	return nil
}

func (c *cacheFs) GetDateDir() string {
	return path.Join(c.RootDir, c.getDateDirName())
}

func (c *cacheFs) getDateDirName() string {
	return c.Date.Format("2006-01-02")
}

func (c *cacheFs) getNamespaceDir(namespace string) string {
	return path.Join(c.GetDateDir(), namespace)
}

func (c *cacheFs) getFileForItem(item Item) string {
	fileName := item.FileName
	if fileName == "" {
		fileName = DefaultContentFile
	}

	return filepath.Join(c.getNamespaceDir(item.Namespace.Path()), fileName)
}

func IsPathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func RemoveDir(pth string) {
	if IsPathExists(pth) {
		if err := os.RemoveAll(pth); err != nil {
			log.Error().
				Err(err).
				Str("path", pth).
				Str("type", "cache").
				Msg("CACHE: Unable to remove directory")
		}
	}
}
