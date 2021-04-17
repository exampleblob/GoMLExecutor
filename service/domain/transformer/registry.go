package transformer

import (
	"fmt"
	"github.com/viant/mly/service/domain"
)

//Register register output transformer
func Register(key string, transformer domain.Transformer) {
	Singleton().Register(key, transformer)
}

//Registry represents a registry
type Registry struct {
	registry map[string]domain.Transformer
}

//Register register transformer
func (r *Registry) Register(key string, transformer domain.Transformer) {
	r.registry[key] = transformer
}

//Lookup returns transformer or error
func (r *Registry) Lookup(key