package http

import (
	"github.com/lovoo/goka"
)

type httpHandler struct {
	emitters map[string]*goka.Emitter
	viewers  map[string]*goka.View
}

func NewHTTPHandler(
	emitters map[string]*goka.Emitter,
	viewers map[string]*goka.View,
) *httpHandler {
	return &httpHandler{
		emitters: emitters,
		viewers:  viewers,
	}
}
