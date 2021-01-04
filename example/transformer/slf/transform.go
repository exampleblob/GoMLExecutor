package slf

import (
	"context"
	"fmt"

	"github.com/viant/gtly"
	"github.com/viant/mly/example/transformer/slf/model"
	"github.com/viant/mly/service/domain"
	"github.com/viant/mly/shared"
	"github.com/viant/mly/shared/common"
)

// conforms to service/domain.Transformer
// only supports single-item requests
func Transform(ctx context.Context, signature *domain.Signature, input *gtly.Object, output interface{}) (common.Storable, error) {
	actual, err := extract(output, 0)
	if err != nil {
		return nil, err
	}

	segment := "other"
	if actual < 1 {
		segment = "one"
	} else if actual < 2 {
		segment = "two"
	} else if actual < 5 {
		segment = "five"
	}

	s := new(model.Segmented)
	s.Class = segment
	return s, nil
}

func extract(o interface{}, i int) (float32, error) {
	switch typed := o.(