package client

import (
	"github.com/viant/mly/shared/client/config"
)

//Config represents a client config
type Config struct {
	Hosts              []*Host
	Model              string
	CacheSizeMb        int