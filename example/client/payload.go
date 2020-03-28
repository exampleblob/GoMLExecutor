
package client

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"github.com/viant/mly/shared/client"
	"github.com/viant/mly/shared/common"
)

type CliPayload struct {
	Data  map[string]interface{}
	Batch int
}

func (c *CliPayload) Iterator(pair common.Pair) error {
	for field, values := range c.Data {
		err := pair(field, values)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CliPayload) Pair(key string, value interface{}) error {
	return nil
}

func (c *CliPayload) SetBatch(msg *client.Message) {