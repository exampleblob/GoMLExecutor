
package client

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/francoispqt/gojay"
)

//Response represents a response
type Response struct {
	Status      string        `json:"status"`
	Error       string        `json:"error,omitempty"`
	ServiceTime time.Duration `json:"serviceTime"`
	DictHash    int           `json:"dictHash"`
	Data        interface{}   `json:"data"`
}

//UnmarshalJSONObject unmsrhal JSON (gojay API)
func (r *Response) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "status":
		if err := dec.String(&r.Status); err != nil {
			return err
		}
	case "error":
		if err := dec.String(&r.Error); err != nil {
			return err
		}
	case "dictHash":
		if err := dec.Int(&r.DictHash); err != nil {
			return err
		}
	case "serviceTimeMcs":
		serviceTime := 0
		if err := dec.Int(&serviceTime); err != nil {
			return err