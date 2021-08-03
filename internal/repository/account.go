package repository

import (
	"database/sql"
	"fmt"

	"github.com/Dimadetected/my-bank-service/internal/models"
)

type Account struct {
	db *sql.DB
}

func NewAccount(db *sql.DB) *Account {
	return &Account{
		db: db,
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
	var acc *models.Account

	//Получаем все данные аккаунта
	if err := a.db.QueryRow(`select * from $1`, accountTable).Scan(&acc.ID, &acc.Currency, &acc.Sum); err != nil {
		return nil, err
	}
	fmt.Println(acc)
	return acc, nil
}

// GetAccount возвращает информацию о счете
func (a *Account) CreatePayment(sum float64) error {
	//Создаем транзакцию
	fmt.Println(0)
	tx, err := a.db.Begin()
	defer tx.Commit()
	if err != nil {
		return nil
	}
	fmt.Println(1)
	//Записываем в бд начисление
	if _, err := tx.Exec(`insert into $1 (sum,is_checked,type) values ($2,$3,$4)`, paymentsTable, sum, false, paymentTypePayment); err != nil {
		tx.Rollback()
		return err
	}
	fmt.Println(2)

	//Увеличиваем сумму счета
	if _, err := tx.Exec(`update $1 set sum = sum + 2`, accountTable, sum); err != nil {
		tx.Rollback()
		return err
	}
	fmt.Println(3)

	return nil
}
func (a *Account) PercentCalculate() error {
	//Получаем из бд все начисления, на которые не были начислены проценты
	rows, err := a.db.Query(`select id, sum from $1 where type = $2 and is_checked = $3`, accountTable, paymentTypePayment, false)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var sum float64
		//Получаем из бд id и сумму начисления
		if err := rows.Scan(&id, &sum); err != nil {
			return err
		}

		//Начинаем транзакцию
		tx, err := a.db.Begin()

		if err != nil {
			return err
		}
		//Считаем процент начисления
		sum *= 0.06

		//Добавляем в таблицу payments начисление процентов
		if _, err := tx.Exec(`insert into $1 (sum,is_checked,type) values ($2,$3,$4)`, paymentsTable, sum, true, paymentTypePercent); err != nil {
			tx.Rollback()
			return err
		}

		//Изменяем сумму аккаунта
		if _, err := tx.Exec(`update $1 set sum = sum + 2`, accountTable, sum); err != nil {
			tx.Rollback()
			return err
		}

		//Изменяем статус платежа на проверенный
		if _, err := tx.Exec(`update $1 set is_checked = $2`, paymentsTable, true); err != nil {
			tx.Rollback()
			return err
		}

		tx.Commit()
	}
	return nil
}
