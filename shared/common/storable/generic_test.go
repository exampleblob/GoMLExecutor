package storable

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/mly/shared/common"
	"testing"
)

func TestGeneric_Iterator(t *testing.T) {
	afoo := &foo{A: 1, B: "aer", C: []int{2, 4}}
	g := NewGeneric(afoo)
	aMap := map[string]interface{}{}
	iter := g.Iterator()
	err := iter(func(key string, value interface{}) error {
		aMap[key] = value
		r