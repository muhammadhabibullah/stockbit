package http

import (
	"github.com/lovoo/goka"
)

type httpHandler struct {
	emitter            *goka.Emitter
	balanceView        *goka.View
	aboveThresholdView *goka.View
}

func NewHTTPHandler(
	emitter *goka.Emitter,
	balanceView *goka.View,
	aboveThresholdView *goka.View,
) *httpHandler {
	return &httpHandler{
		emitter:            emitter,
		balanceView:        balanceView,
		aboveThresholdView: aboveThresholdView,
	}
}
