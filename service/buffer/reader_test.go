package buffer

import (
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {

	useCases := []struct {
		description string
		bufferSize  int
		reader      io.Reader
		hasError    bool
		expect      string
	}{
		{
			description: "medium buffer, small data",
			reader:      strings.NewReader("Lorem ipsum"),
			bufferSize:  1024,
			expect:      "Lorem ipsum",
		},
		{
			description: "large buffer, medium data",
			reader:      strings.NewReader(strings.Repeat("Lorem ipsum", 1024)),
			bufferSize:  1024 * 1024,
			expect:      strings.Repeat("Lorem ipsum", 1024),
		},
		{
			description: "large buffer, large data",
			reader:      strings.NewReader(strings.Repeat("Lorem ipsum", 1024*1024)),
			bufferSize:  1024 * 1024 * 26,
			expect:      strings.Repeat("Lorem ipsum", 1024*1024),
		},
		{
			description: "too small buffer",
			reader:      strings.NewReader(strings.Repeat("Lorem ipsum", 1024*1024)),
			bufferSize:  1024,
			hasError:    true,
			expect:      strings.Repeat("Lorem ipsum", 1024*1024),
		},
	}

	for _, useC