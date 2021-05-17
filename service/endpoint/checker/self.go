package checker

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/viant/mly/service/config"
	"github.com/viant/mly/shared"
	"github.com/viant/mly/shared/client"
	"github.com/viant/toolbox"
)

func SelfTest(host []*client.Host, timeout time.Duration, modelID string, usesTransformer bool, inputs_ []*shared.Field, tp config.TestPayload, outputs []*shared.Field, debug bool) error {
	cli, err := client.New(modelID, host, client.WithDebug(true))
	if err != nil {
		return fmt.Errorf("%s:%w", modelID, err)
	}

	inputs := cli.Config.Datastore.MetaInput.Inputs

	// generate payload

	var testData map[string]interface{}
	var batchSize int
	if len(tp.Batch) > 0 {
		for k, v := range tp.Batch {
			testData[k] = v
			batchSize = len(v)
		}
	} else {
		if len(tp.S