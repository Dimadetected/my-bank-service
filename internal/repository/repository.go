package repository

import (
	"database/sql"
	"fmt"

	"github.com/Dimadetected/my-bank-service/internal/models"
)

type Repository struct {
	Account AccountInterface
}

type AccountInterface interface {
	GetAccount() (*models.Account, error)
}

func NewRepositry(db *sql.DB) *Repository {
	fmt.Println("=====")
	fmt.Println(db)
	return &Repository{
		Account: NewAccount(db),
	}
}
