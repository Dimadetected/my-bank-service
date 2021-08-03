package repository

import (
	"database/sql"

	"github.com/Dimadetected/my-bank-service/internal/models"
)

type Repository struct {
	Account AccountInterface
}

type AccountInterface interface {
	GetAccount() (*models.Account, error)
	CreatePayment(sum float64) error
	PercentCalculate() error
}

func NewRepositry(db *sql.DB) *Repository {
	return &Repository{
		Account: NewAccount(db),
	}
}
