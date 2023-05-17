package config

import (
	"fmt"
	"github.com/viant/mly/shared/config/datastore"
)

//DatastoreList represents datastore list
type DatastoreList struct {
	Connections []*datastore.Connection
	Datastores  []*Datastore
}

//Init initialises list
func (d *DatastoreList) Init() {
	if len(d.Connections) > 0 {
		for i := range d.Connections {
			d.Connections[i].Init()
		}
	}
	if len(d.D