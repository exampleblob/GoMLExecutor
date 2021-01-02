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
func Transform(ctx context.Context, signature *domain.Signature, input *gtly.Object, output interface{}) (co