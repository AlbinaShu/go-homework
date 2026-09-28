package main

import (
	"errors"
	"fmt"
)

type Transaction struct {
	ID          int
	Amount      float64
	Category    string
	Description string
	Date        string
}

var transactions = []Transaction{}

func AddTransaction(tx Transaction) error {
	if tx.Amount == 0 {
		return errors.New("amount cannot be zero")
	}

	tx.ID = len(transactions) + 1

	transactions = append(transactions, tx)

	return nil
}

func ListTransactions() []Transaction {
	return transactions
}

func main() {
	fmt.Println("Ledger service started")

	tx1 := Transaction{
		Amount:      150.50,
		Category:    "Food",
		Description: "Dinner",
		Date:        "2026-09-27",
	}

	tx2 := Transaction{
		Amount:      50,
		Category:    "Transport",
		Description: "Taxi",
		Date:        "2026-09-27",
	}

	tx3 := Transaction{
		Amount:      1000,
		Category:    "Shopping",
		Description: "Shoes",
		Date:        "2026-09-27",
	}

	AddTransaction(tx1)
	AddTransaction(tx2)
	AddTransaction(tx3)

	fmt.Println("Transactions:")

	for _, tx := range ListTransactions() {
		fmt.Println(tx)
	}
}