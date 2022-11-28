package cache

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/pestanko/miniscrape/pkg/config"
)

// DefaultContentFile contains the name of the file where to store processed
// content
const DefaultContentFile = "content.txt"

// NamespacePath defines a generic interface for each type to have method
// to return the namespace path
type NamespacePath interface {
	// Path return the path of the item
	Path() string
}

// Item represents an unique item stored in the cache
type Item struct {
	// Namespace of the item
	Namespace ItemNamespace
	// FileName of the file that will be/is stored in the cache
	FileName string
	// CachePolicy for the item - it can be nocache
	CachePolicy string
}

// ItemNamespace contains tuple Category/Page
type ItemNamespace struct {
	// Page codename
	Page string
	// Category codename
	Category string
}

// Path returns the path representation of the namespace (OS specific)
func (n ItemNamespace) Path() string {
	return path.Join(n.Category, n.Page)
}

// String returns the string representation of the namespace
func (n ItemNamespace) String() string {
	return fmt.Sprintf("%s/%s", n.Category, n.Page)
}

// NewNamespace creates a instance of the new ItemNamespace
func NewNamespace(cat string, page string) ItemNamespace {
	return ItemNamespace{Category: cat, Page: page}
}

// Cache interface
type Cache interface {
	// Store the item to the cache with provided content
	Store(item Item, content []byte) error
	// IsPageCached checks whether the page is in the cache
	IsPageCached(nm ItemNamespace) bool
	// IsItemCached checks whether the item is in the cache
	// Deprecated: use IsPageCached instead
	IsItemCached(item Item) bool
	// GetContent returns the content for the item
	GetContent(item Item) []byte
	// Invalidate the cache content
	Invalidate(sel config.RunSelector)
}

// NewCache creates an instance of the new cache
func NewCache(cacheCfg config.CacheCfg, date time.Time) Cache {
	if cacheCfg.Enabled {
		root := cacheCfg.Root
		if root == "" {
			root = path.Join(os.TempDir(), "mini-scrape")
		}
		log.Info().Str("cache_root", root).Msg("cache is enabled")
		return &cacheFs{
			rootDir:     root,
			forceUpdate: cacheCfg.Update,
			date:        date,
		}
	} else {
		log.Info().Msg("cache is disabled")
	}
	return nil
}

// cacheFs cache implemented over the filesystem
type cacheFs struct {
	rootDir     string
	forceUpdate bool
	date        time.Time
}

func (c *cacheFs) Invalidate(sel config.RunSelector) {
	nm := NewNamespace(sel.Category, sel.Page)
	pth := c.getNamespaceDir(nm.Path())
	removeDir(pth)
}

func (c *cacheFs) IsItemCached(item Item) bool {
	return !c.forceUpdate && isPathExists(c.getFileForItem(item))
}

func (c *cacheFs) GetContent(item Item) []byte {
	fp := c.getFileForItem(item)
	content, err := os.ReadFile(filepath.Clean(fp))

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
	return !c.forceUpdate && isPathExists(c.getNamespaceDir(nm.Path()))
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

	if err := os.WriteFile(fp, content, 0600); err != nil {
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
	return path.Join(c.rootDir, c.getDateDirName())
}

func (c *cacheFs) getDateDirName() string {
	return c.date.Format("2006-01-02")
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

func isPathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func removeDir(pth string) {
	if isPathExists(pth) {
		if err := os.RemoveAll(pth); err != nil {
			log.Error().
				Err(err).
				Str("path", pth).
				Str("type", "cache").
				Msg("CACHE: Unable to remove directory")
		}
	}
}
