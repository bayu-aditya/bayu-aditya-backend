package rest

import (
	"github.com/bayu-aditya/bayu-aditya-backend/lib/core/usecase/monopolybanking"
)

type Handler struct {
	prefix                 string
	usecaseMonopolyBanking usecasemonopolybanking.IMonopolyBankingUsecase
}

func NewHttpHandler(ucMonopoly usecasemonopolybanking.IMonopolyBankingUsecase) *Handler {
	return &Handler{
		prefix:                 "rest",
		usecaseMonopolyBanking: ucMonopoly,
	}
}
