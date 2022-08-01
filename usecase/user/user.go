package user

import (
	"github.com/lovoo/goka"
	"stockbit/config"
)

type userUseCase struct {
	emitters   map[string]*goka.Emitter
	viewers    map[string]*goka.View
	depositCfg config.DepositConfig
}

func NewUserUseCase(
	emitters map[string]*goka.Emitter,
	viewers map[string]*goka.View,
	cfg config.Config,
) *userUseCase {
	return &userUseCase{
		emitters:   emitters,
		viewers:    viewers,
		depositCfg: cfg.Deposit,
	}
}
