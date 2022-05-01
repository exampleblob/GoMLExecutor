package client

import (
	"github.com/viant/mly/shared/client/config"
)

//Config represents a client config
type Config struct {
	Hosts              []*Host
	Model              string
	CacheSizeMb        int
	CacheScope         *CacheScope
	Datastore          *config.Remote
	MaxRetry           int
	Debug              bool
	DictHashValidation bool
}

//CacheSize returns cache size
func (c *Config) CacheSize