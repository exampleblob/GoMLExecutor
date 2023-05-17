package datastore

import (
	"strings"
	"time"
)

//Timeout timout setting
type Timeout struct {
	Unit       string
	Connection int
	Socket     int
	Total      int
}

func (t *Timeout) DurationUnit() time.Duration {
	if t.Unit == "" {
		t.Unit = "ms"
	}
	switch strings.ToLower(t.Unit) {
	case "ms":
		return time.Millise