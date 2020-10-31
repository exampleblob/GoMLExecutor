package server

import (
	"context"

	"github.com/viant/mly/example/transformer/slf"
	slfmodel "github.com/viant/mly/example/transformer/slf/model"
	"github.com/viant/mly/service/domain/transformer"
	"github.com/viant/mly/service/endpoint"
	"github.com/viant/mly/shared/common"
	"github.com/viant/mly/shared/common/storable"
)

func RunApp(Version string