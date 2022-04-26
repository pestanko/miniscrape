package cache

import (
	"github.com/pestanko/miniscrape/scraper/config"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
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
		log.Printf("CACHE: Unable to load content '%s': %v", fp, err)
		return []byte{}
	}
	log.Printf("CACHE: Loading cached content '%s'", fp)
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
			log.Printf("Unable to create directory file '%s': %v", fp, err)
			return err
		}
	}

	if c.IsItemCached(item) {
		log.Printf("Item already cached: %v", item)
		return nil
	}

	log.Printf("Writing cache file '%s': %v", fp, item)
	if err := ioutil.WriteFile(fp, content, 0600); err != nil {
		log.Printf("Unable to write file '%s': %v", fp, err)
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
			log.Printf("Unable to remove directory '%s': %v", pth, err)
		}
	}
}
