package service

import (
	info "github.com/Dimadetected/my-bank-service/interface"
	"github.com/Dimadetected/my-bank-service/internal/repository"
)

type Service struct {
	Account info.AccountInterface
}

func NewService(r *repository.Repository) *Service {
	return &Service{
		Account: NewAccount(r.Account),
	}
}
