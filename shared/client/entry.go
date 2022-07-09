
package client

import (
	"math"
	"math/big"
)

type entry struct {
	ints     map[int]bool
	float32s map[float32]bool // Deprecated: floats are never used for lookup
	strings  map[string]bool

	prec uint
	fm64 float64
}

func (e *entry) hasString(val string) bool {
	if len(e.strings) == 0 {
		return false
	}
	_, ok := e.strings[val]
	return ok
}

func (e *entry) hasInt(val int) bool {