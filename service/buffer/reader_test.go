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
			reader:      strings.NewReader(strings.Repeat("Lorem ipsum"