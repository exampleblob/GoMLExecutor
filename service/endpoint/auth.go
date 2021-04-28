package endpoint

import (
	"net/http"

	"github.com/viant/mly/shared/common"
)

type AuthHandler struct {
	*Config
	handler http.Handler
}

func NewAuthHandler(config *Config, handler http.Handler) *AuthHandler {
	h := new(AuthHandler)
	h.Config = config
	h.handler