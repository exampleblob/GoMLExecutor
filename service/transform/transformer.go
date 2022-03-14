package transform

import (
	"github.com/viant/mly/service/domain"
	"github.com/viant/mly/service/domain/transformer"
)

// Get from transformer singleton
func Get(name string) (domain.Transformer, error) {
	result, err := transformer.Singleton().Lookup(name)
	if err == nil && result !