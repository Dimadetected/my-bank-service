package service

import (
	"github.com/Dimadetected/my-bank-service/internal/repository"
	info "github.com/Dimadetected/my-bank-service/interface"
)

type Service struct {
	repo    *repository.Repository
	Account  info.AccountInterface
}

func NewService(r *repository.Repository) *Service {
	return &Service{
		repo:    r,
		Account: NewAccount(),
	}
}
