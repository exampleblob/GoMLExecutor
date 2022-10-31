package client

import (
	"encoding/json"
	"fmt"

	"github.com/francoispqt/gojay"
)

func Marshal(data interface{}, id string) ([]byte, error) {
	if 