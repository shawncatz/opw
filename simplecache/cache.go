package simplecache

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Cache struct {
	File string
	data CacheData
}

type CacheData map[string]*CacheItem

type CacheItem struct {
	Value   string
	Expires time.Time
}

func New(file string) (*Cache, error) {
	cache := &Cache{
		File: file,
		data: CacheData{},
	}

	if _, err := os.Stat(cache.File); os.IsNotExist(err) {
		if err := cache.write(); err != nil {
			return nil, err
		}
	}

	if err := cache.read(); err != nil {
		return nil, err
	}

	return cache, nil
}

func (c *Cache) Get(key string) string {
	v := c.data[key]
	if v == nil {
		logrus.Debugf("cache miss: %s", key)
		return ""
	}

	t := time.Now()
	if v.Expires.Sub(t) < 0 {
		logrus.Debugf("cache expired: %s", key)
		return ""
	}

	logrus.Debugf("cache hit: %s", key)
	return v.Value
}

func (c *Cache) Set(key, value string, expires time.Duration) error {
	c.data[key] = &CacheItem{Value: value, Expires: time.Now().Add(expires)}
	return c.write()
}

func (c *Cache) Fetch(key string, expires time.Duration, force bool, f func() (string, error)) (string, error) {
	var err error

	v := c.Get(key)
	if v == "" || force {
		v, err = f()
		if err != nil {
			return v, err
		}

		if err = c.Set(key, v, expires); err != nil {
			return v, err
		}
	}

	return v, nil
}

func (c *Cache) read() error {
	data, err := ioutil.ReadFile(c.File)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &c.data)
}

func (c *Cache) write() error {
	data, err := c.dataAsJSON()
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(c.File, data, 0600); err != nil {
		return err
	}

	return nil
}

func (c *Cache) dataAsJSON() ([]byte, error) {
	return json.Marshal(c.data)
}
