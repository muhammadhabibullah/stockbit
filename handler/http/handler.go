package http

import (
	"stockbit/domain"
)

type httpHandler struct {
	userUseCase domain.UserUseCase
}

func NewHTTPHandler(
	userUseCase domain.UserUseCase,
) *httpHandler {
	return &httpHandler{
		userUseCase: userUseCase,
	}
}
