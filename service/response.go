package service

import (
	"encoding/json"
	"github.com/francoispqt/gojay"
	"github.com/viant/mly/shared/common"
	"github.com/viant/xunsafe"
	"log"
	"time"
	"unsafe"
)

//Response represents service response
type Response struct {
	started        time.Time
	xSlice         *xunsafe.Slice
	sliceLen       int
	Status         string
	Error          string
	DictHash       int
	Data           in