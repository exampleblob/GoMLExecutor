
package test

import (
	"context"
	"fmt"

	"github.com/francoispqt/gojay"
	"github.com/viant/bintly"
	"github.com/viant/gtly"
	"github.com/viant/mly/service/domain"
	"github.com/viant/mly/shared"
	"github.com/viant/mly/shared/common"
	"github.com/viant/toolbox"
)

type VO struct {
	Value float32
}

func (p *VO) EncodeBinary(enc *bintly.Writer) error {
	enc.Float32(p.Value)
	return nil
}

func (p *VO) DecodeBinary(dec *bintly.Reader) error {
	dec.Float32(&p.Value)
	return nil
}

func (s *VO) MarshalJSONObject(enc *gojay.Encoder) {