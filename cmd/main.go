package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Dimadetected/my-bank-service/internal/handler"
	"github.com/Dimadetected/my-bank-service/internal/repository"
	"github.com/Dimadetected/my-bank-service/internal/service"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db := dbInit()

	// defer db.Close()

	repo := repository.NewRepositry(db)
	service := service.NewService(repo)
	h := handler.NewHandler(service)

	http.HandleFunc("/AddFunds", h.AddFunds)
	http.HandleFunc("/Withdraw", h.Withdraw)
	http.HandleFunc("/GetCurrency", h.GetCurrency)
	http.HandleFunc("/SumProfit", h.SumProfit)
	http.HandleFunc("/GetAccountCurrencyRate", h.GetAccountCurrencyRate)
	http.HandleFunc("/GetBalance", h.GetBalance)
	fmt.Println(db)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
	db.Close()
}

func dbInit() *sql.DB {
	db, err := sql.Open("sqlite3", "db.db")
	if err != nil {
		panic(err)
	}

	var count int
	if err := db.QueryRow(`select count(*) from accounts`).Scan(&count); err != nil {
		panic(err)
	}

	if count == 0 {
		if _, err := db.Exec(`INSERT INTO accounts (currency, sum) VALUES ('SBP', 0);`); err != nil {
			panic(err)
		}
	}
	return db
}
