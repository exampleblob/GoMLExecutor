package config

import (
	"fmt"
	"github.com/viant/mly/shared/config/datastore"
)

//DatastoreList represents datastore list
type DatastoreList struct {
	Connections []*datastore.Connection
	Datastores  []*D