package config

import (
	"github.com/viant/mly/shared"
	"github.com/viant/mly/shared/config"
	"github.com/viant/mly/shared/config/datastore"
)

//Remote represents client datastore
type Remote struct {
	Connections [