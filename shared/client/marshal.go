package client

import (
	"encoding/json"
	"fmt"

	"github.com/francoispqt/gojay"
)

func Marshal(data interface{}, id string) ([]byte, error) {
	if data == nil {
		return nil, fmt.Errorf("data was nil")
	}
	switch val := data.(type) {
	case *Message:
		if !val.isValid() {
			return nil, fmt.Errorf("invalid message: has been already sent befo