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
	cli, err := client.