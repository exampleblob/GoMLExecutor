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
	Data           interface{}
	ServiceTimeMcs int
}

//SetError sets errors
func (r *Response) SetError(err error) {
	if err == nil {
		return
	}
	r.Error = err.Error()
	r.Status = common.StatusError
}

//MarshalJSONObject marshal response
func (r *Response) MarshalJSONObject(enc *gojay.Encoder) {
	enc.StringKeyOmitEmpty("status", r.Status)
	enc.StringKeyOmitEmpty("error", r.Error)
	enc.IntKeyOmitEmpty("dictHash", r.DictHash