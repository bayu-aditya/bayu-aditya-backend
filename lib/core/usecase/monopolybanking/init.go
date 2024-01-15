package usecasemonopolybanking

import reponats "github.com/bayu-aditya/bayu-aditya-backend/lib/core/repository/nats"

func NewMonopolyBankingUsecase(repoNats reponats.INatsRepository) IMonopolyBankingUsecase {
	return &usecase{
		prefix:         "ucmonopolybanking",
		repositoryNats: repoNats,
	}
}
