package service

import (
	"errors"
	"fmt"

	info "github.com/Dimadetected/my-bank-service/interface"
	"github.com/Dimadetected/my-bank-service/internal/repository"
)

const (
	SBP2RUB = 0.7523
	RUB2SBP = 1.3333
)

type Account struct {
	r repository.Account
}

func NewAccount() *Account {
	return &Account{}
}

func (a *Account) AddFunds(sum float64) {
	fmt.Println(sum)
	if err := a.r.CreatePayment(sum); err != nil {
		panic(err)
	}
	a.SumProfit()
}

// SumProfit Рассчитывает процент по вкладу и полученные деньги вносит на счёт
func (a *Account) SumProfit() {
	if err := a.r.PercentCalculate(); err != nil {
		panic(err)
	}
}

// Withdraw Производит списание со счёта по указанным правилам. Если списание выходит за рамки правил, выдаёт ошибку
func (a *Account) Withdraw(f float64) error {
	acc, err := a.r.GetAccount()
	if err != nil {
		return err
	}
	if acc.Sum*0.3 < f {
		return errors.New("сумма для списания превышает 30% от суммы вклада")
	}
	if err := a.r.CreatePayment(-f); err != nil {
		return err
	}
	return nil
}

// GetCurrency Выдаёт валюту счёта
func (a *Account) GetCurrency() info.Currency {
	acc, err := a.r.GetAccount()
	if err != nil {
		panic(err)
	}
	return acc.Currency
}

// GetAccountCurrencyRate Выдаёт курс валюты счёта к передаваемой валюте cur
func (a *Account) GetAccountCurrencyRate(cur info.Currency) float64 {
	switch cur {
	case info.CurrencyRUB:
		return SBP2RUB
	case info.CurrencySBP:
		return RUB2SBP
	}
	return 0
}

// GetBalance Выдаёт баланс счёта в указанной валюте
func (a *Account) GetBalance(cur info.Currency) float64 {
	acc, err := a.r.GetAccount()
	if err != nil {
		panic(err)
	}

	var balance float64
	switch cur {
	case info.CurrencyRUB:
		balance = acc.Sum * SBP2RUB
	case info.CurrencySBP:
		balance = acc.Sum * RUB2SBP
	}

	return balance
}
