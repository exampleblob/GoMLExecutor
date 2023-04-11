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
	stream.Float32(t.Prediction)
	return nil
}

func (t *TestOutput) DecodeBinary(stream *bintly.Reader) error {
	stream.Float32(&t.Prediction)
	return nil
}

func TestService_Run(t *testing.T) {
	baseURL := toolbox.CallerDirectory(3)

	server := faker.Server{URL: path.Join(baseURL, "testdata"), Port: 8087, Debug: true}
	go server.Start()
	time.Sleep(time.Second)
	defer server.Stop()
	var metaInput = shared.MetaInput{
		Inputs: []*shared.Field{
			{
				Name: "i1",
			},
			{
				Name:     "i2",
				Wildcard: true,
			},
		},
	}

	var dictionary = NewDictionary(&common.Dictionary{
		Layers: []common.Layer{
			{
				Name: "i1",
				Strings: []string{
					"v1", "v2",
				},
			},
		},
		Hash: 123,
	}, metaInput.Inputs)

	var testCases = []struct {
		description string
		model       string
		opt