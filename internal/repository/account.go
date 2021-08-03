package repository

import (
	"database/sql"
	"fmt"
	"log"

	info "github.com/Dimadetected/my-bank-service/interface"
	"github.com/Dimadetected/my-bank-service/internal/models"
)

type Account struct {
	DB *sql.DB
}

func NewAccount(db *sql.DB) *Account {
	return &Account{
		DB: db,
	}
}

const (
	accountTable  = "accounts"
	paymentsTable = "payments"

	paymentTypePayment = "payment"
	paymentTypePercent = "percent"
)

// GetAccount возвращает информацию о счете
func (a *Account) GetAccount() (*models.Account, error) {
	acc := models.Account{}
	var curr string
	//Получаем все данные аккаунта
	if err := a.DB.QueryRow(`select * from accounts`).Scan(&acc.ID, &curr, &acc.Sum); err != nil {
		return nil, err
	}
	acc.Currency = info.Currency(curr)
	return &acc, nil
}

// GetAccount возвращает информацию о счете
func (a *Account) CreatePayment(sum float64) error {
	//Создаем транзакцию
	tx, err := a.DB.Begin()
	defer tx.Commit()
	if err != nil {
		return nil
	}
	//Записываем в бд начисление
	q := fmt.Sprintf(`insert into %s (sum,is_checked,type) values (%f,%t,'%s')`, paymentsTable, sum, false, paymentTypePayment)
	if _, err := tx.Exec(q); err != nil {
		tx.Rollback()
		return err
	}

	//Увеличиваем сумму счета
	q = fmt.Sprintf(`update %s set sum = sum + %f`, accountTable, sum)
	if _, err := tx.Exec(q); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

type Percent struct {
	ID  int
	Sum float64
}

func (a *Account) PercentCalculate() error {
	//Получаем из бд все начисления, на которые не были начислены проценты
	acc, err := a.GetAccount()
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`select id, sum from %s where type = '%s' and is_checked = %d`, paymentsTable, paymentTypePayment, 0)
	rows, err := a.DB.Query(q)
	if err != nil {
		return err
	}
	var perc []Percent
	for rows.Next() {
		var p Percent
		//Получаем из бд id и сумму начисления
		if err := rows.Scan(&p.ID, &p.Sum); err != nil {
			return err
		}
		perc = append(perc, p)
	}
	rows.Close()
	for _, p := range perc {
		//Начинаем транзакцию
		tx, err := a.DB.Begin()

		if err != nil {
			return err
		}
		//Считаем процент начисления
		sum := acc.Sum * 0.06

		//Добавляем в таблицу payments начисление процентов
		q = fmt.Sprintf(`insert into %s (sum,is_checked,type) values (%f,%d,'%s')`, paymentsTable, sum, 1, paymentTypePercent)
		if _, err := tx.Exec(q); err != nil {
			log.Print(err)
			tx.Rollback()
			return err
		}

		//Изменяем сумму аккаунта
		q = fmt.Sprintf(`update %s set sum = sum + %f`, accountTable, sum)
		if _, err := tx.Exec(q); err != nil {
			log.Print(err)
			tx.Rollback()
			return err
		}

		//Изменяем статус платежа на проверенный
		q = fmt.Sprintf(`update %s set is_checked = %d where id = %d`, paymentsTable, 1, p.ID)
		if _, err := tx.Exec(q); err != nil {
			log.Print(err)
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			log.Print(err)
		}
	}
	return nil
}
