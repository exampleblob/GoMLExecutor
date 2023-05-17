package datastore

import (
	"strings"
	"time"
)

//Timeout timout setting
type Timeout struct {
	Unit       string
	Connection int
	So