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
		if len(tp.Single) > 0 {
			testData = tp.Single

			for _, field := range inputs {
				n := field.Name
				sv, ok := tp.Single[n]
				switch field.DataType {
				case "int", "int32", "int64":
					if !ok {
						testData[n] = rand.Int31()
					} else {
						switch tsv := sv.(type) {
						case string:
							testData[n], err = strconv.Atoi(tsv)
							if err != nil {
								return err
							}
						}
					}
				case "float", "float32", "float64":
					testData[n] = rand.Float32()
				default:
					if !ok {
						testData[n] = fmt.Sprintf("test-%d", rand.Int31())
					} else {
						testData[n] = toolbox.AsString(sv)
					}
				}
			}
		} else {
			testData = make(map[string]interface{})
			for _, field := range inputs {
				n := field.Name
				switch field.DataType {
				case "int", "int32", "int64":
					testData[n] = rand.Int31()
				case "float", "float32", "float64":
					testData[n] = rand.Float32()
				default:
					testData[n] = fmt.Sprintf("test-%d", rand.Int31())
				}
			}
		}

		if tp.SingleBatch {
			for _, field := range inputs {
				fn := field.Name
				tv := testData[fn]
				switch field.DataType {
				case "int", "int32", "int64":
					var v int
					switch atv := tv.(type) {
					case int:
						v = atv
					case int32:
					case int64:
						v = int(atv)
					default:
						return fmt.Errorf("test data malformed: %s expected int-like, found %T", fn, tv)
					}

					b := [1]int{v}
					testData[fn] = b[:]
				case "float", "float32", "float64":
					var v float32
					switch atv := tv.(type) {
					case float32:
						v = atv
					case float64:
						v = float32(atv)
					default:
						return fmt.Errorf("test data malformed: %s expected float32-like, found %T", fn, tv)
					}

					b := [1]float32{v}
					testData[fn] = b[:]
				default:
					switch atv := tv.(type) {
					case string:
						b