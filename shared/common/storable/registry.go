package storable

import (
	"fmt"
	"github.com/viant/mly/shared/common"
)

//Registry represents storable registry
type Registry struct {
	registry map[string]func() common.Storable
}

//Register represents storable registry
func (r *Registry) Register(key string, fn func() common.Storable) {
	r.reg