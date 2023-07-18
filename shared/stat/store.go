package stat

import "github.com/viant/gmetric/counter"

type store struct{}

func (p store) Keys() []string {
	return []string{
		ErrorKey,
		NoSuchKey,
		Timeout,
		Down,
		Canceled,
		DeadlineExceeded,
	}
}

func (p store) Map(value interface{}) int {
	if value == nil {
		return -1
	}
	switch val := value.(type) {
	case error:
		return 0
	case string:
		switch