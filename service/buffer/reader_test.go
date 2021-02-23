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
			reader:      string