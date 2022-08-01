package http

import (
	"github.com/lovoo/goka"
)

type httpHandler struct {
	emitter *goka.Emitter
	view    *goka.View
}

func NewHTTPHandler(
	emitter *goka.Emitter,
	view *goka.View,
) *httpHandler {
	return &httpHandler{
		emitter: emitter,
		view:    view,
	}
}
