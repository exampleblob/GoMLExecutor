package client

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/bintly"
	"github.com/viant/mly/shared"
	cconfig "github.com/viant/mly/shared/client/config"
	"github.com/viant/mly/shared/client/faker"
	"github.com/viant/mly/shared/common"
	"github.com/viant/mly/shared/config"
	"github.com/viant/mly/shared/datastore/mock"
	"github.com/viant/scache"
	"github.com/viant/toolbox"
	"path"
	"reflect"
	"testing"
	"time"
)

type TestOutput struct {
	Prediction float32
}

func (t *TestOutput) EncodeBinary(stream *bintly.Writer) error {
	stream.Float32(