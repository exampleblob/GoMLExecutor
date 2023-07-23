package stat

import "github.com/viant/gmetric/counter"

const (
	L1Write = "L1Write"
	L2Write = "L2Write"
)

type write struct{}

func (p write) Keys() []string {
	return []string{
		ErrorKey,
		