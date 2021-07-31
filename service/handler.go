
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/francoispqt/gojay"
	"github.com/viant/mly/service/buffer"
	"github.com/viant/mly/service/clienterr"
	"github.com/viant/mly/service/request"
)

// Handler converts a model prediction HTTP request to its internal calls.
type Handler struct {
	maxDuration time.Duration
	service     *Service
	pool        *buffer.Pool
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, httpRequest *http.Request) {
	// use Background() since there are things to be done regardless of if the request is cancceled from the client side.
	ctx := context.Background()
	// TODO: Handle httpRequest.Context() - there are issues since this can be canceled but there should be housekeeping completed.
	ctx, cancel := context.WithTimeout(ctx, h.maxDuration)
	defer cancel()

	isDebug := h.service.config.Debug

	request := h.service.NewRequest()
	response := &Response{Status: "ok", started: time.Now()}
	if httpRequest.Method == http.MethodGet {
		if err := h.buildRequestFromQuery(httpRequest, request); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		defer httpRequest.Body.Close()
		data, size, err := buffer.Read(h.pool, httpRequest.Body)
		defer h.pool.Put(data)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		request.Body = data[:size]