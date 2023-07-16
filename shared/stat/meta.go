package stat

import "github.com/viant/gmetric/counter"

// shared error-only multi operation tracker provider
type errorOnly struct{}

func (p errorOnly) Keys() []string {
	return []string{
		ErrorKey,
	}
}

fun