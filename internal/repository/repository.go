package repository

import (
	"database/sql"
	"github.com/Dimadetected/my-bank-service/internal/models"
)

type Repository struct {
	db      *sql.DB
	Account AccountInterface
}


type AccountInterface interface {
	GetAccount() (*models.Account, error)
}

func NewRepositry(db *sql.DB) *Repository {
	return &Repository{
		db:      db,
		Account: NewAccount(db),
	}
}
