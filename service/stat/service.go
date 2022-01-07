package stat

import (
	"github.com/viant/gmetric/counter"
	"github.com/viant/mly/shared/stat"
)

const (
	Evaluate = "eval"
	Invalid  = "invalid"
)

type service struct{}

// implements github.com/viant/gmetric/counter.Provider
func (e service) Keys() []string {
	return []string{
		stat.ErrorKey,
		Evaluate,
		Pending,
		stat.Timeout,
		Invalid,
		stat.Canceled,
		stat.DeadlineExceeded,
	}
}

// implements github.com/viant/gmetric/counter.Provider
func (e service) Map(value interface{}) int {
	if value == nil {
		return -1
	}
	switch val := value.(type) {
	case error:
		ret