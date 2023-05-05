package storable

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/mly/shared/common"
	"testing"
)

func TestGeneric_Iterator(t *testing.T) {
	afoo := &foo{A: 1, 