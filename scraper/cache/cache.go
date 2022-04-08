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

type Cache interface {
	Store(item Item, content []byte) error
	BlockListPage(pageName string) error
	IsPageCached(pageName string) bool
	IsItemCached(item Item) bool
	GetContent(item Item) []byte
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

func (c *cacheFs) IsItemCached(item Item) bool {
	return !c.ForceUpdate && IsPathExists(c.getFileForItem(item))
}

type Item struct {
	PageName     string
	FileName     string
	CategoryName string
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

func (c *cacheFs) IsPageCached(pageName string) bool {
	return !c.ForceUpdate && IsPathExists(c.getPageDir(pageName))
}

func (c *cacheFs) Store(item Item, content []byte) error {
	fp := c.getFileForItem(item)

	if !c.IsPageCached(item.PageName) {
		if err := os.MkdirAll(c.getPageDir(item.PageName), 0700); err != nil {
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

func (c *cacheFs) BlockListPage(name string) error {
	// TODO: Not implemented
	return nil
}

func (c *cacheFs) GetDateDir() string {
	return path.Join(c.RootDir, c.getDateDirName())
}

func (c *cacheFs) getDateDirName() string {
	return c.Date.Format("2006-01-02")
}

func (c *cacheFs) getPageDir(pageName string) string {
	return path.Join(c.GetDateDir(), pageName)
}

func (c *cacheFs) getFileForItem(item Item) string {
	fileName := item.FileName
	if fileName == "" {
		fileName = DefaultContentFile
	}

	return filepath.Join(c.getPageDir(item.PageName), fileName)
}

func IsPathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
