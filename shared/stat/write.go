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
		L1Write,
		L2Write,
	}
}

func (p write) Map(value interface{}) int {
	if value == nil {
		return -1
	}

	switch val := value.(type) {
	case error:
		return 0
	case string:
		switch val {
		c