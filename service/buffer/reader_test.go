package buffer

import (
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {

	useCases := []struc