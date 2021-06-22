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
	mu      *sync.Mutex
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		mu:      new(sync.Mutex),
		healths: make(map[string]*int32),
	}
}

func (h *HealthHandler) RegisterHealthPoint(name string, isOkPtr *int32) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.healths[name] = isOkPtr
}

// implements Hook
func (h *HealthHandler) Hook(model *config.Model, modelSr