package health

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/viant/mly/service"
	"github.com/viant/mly/service/config"
)

type HealthHandler struct {
	healths map[string]*int32
	mu   