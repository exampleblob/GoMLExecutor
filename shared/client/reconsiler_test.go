
package client

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viant/mly/shared"
	"github.com/viant/mly/shared/common"
)

func TestReconcileData(t *testing.T) {

	messages := NewMessages(func() *Dictionary {
		return NewDictionary(&common.Dictionary{}, []*shared.Field{
			{