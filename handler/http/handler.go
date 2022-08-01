package http

import (
	"github.com/lovoo/goka"
	"stockbit/config"
)

type httpHandler struct {
	emitters   map[string]*goka.Emitter
	viewers    map[string]*goka.View
	depositCfg config.DepositConfig
}

func NewHTTPHandler(
	emitters map[string]*goka.Emitter,
	viewers map[string]*goka.View,
	cfg config.Config,
) *httpHandler {
	return &httpHandler{
		emitters:   emitters,
		viewers:    viewers,
		depositCfg: cfg.Deposit,
	}
}
