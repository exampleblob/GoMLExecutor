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
		De