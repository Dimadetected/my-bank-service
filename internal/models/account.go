package models

import (
	info "github.com/Dimadetected/my-bank-service/interface"
)

type Account struct {
	ID       int
	Currency info.Currency
	Sum      float64
}
