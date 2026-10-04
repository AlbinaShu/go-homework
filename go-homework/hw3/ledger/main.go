package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

type Transaction struct {
	ID          int
	Amount      float64
	Category    string
	Description string
	Date        string
}

type Budget struct {
	Category string
	Limit    float64
}

var transactions = []Transaction{}
var budgets = make(map[string]Budget)

func SetBudget(b Budget) error {
	if b.Category == "" {
		return errors.New("category cannot be empty")
	}

	if b.Limit <= 0 {
		return errors.New("budget limit must be greater than zero")
	}

	budgets[b.Category] = b

	return nil
}

func GetCategoryTotal(category string) float64 {
	var total float64

	for _, tx := range transactions {
		if tx.Category == category {
			total += tx.Amount
		}
	}

	return total
}

func AddTransaction(tx Transaction) error {
	if tx.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	budget, exists := budgets[tx.Category]

	if exists {
		total := GetCategoryTotal(tx.Category)

		if total+tx.Amount > budget.Limit {
			return errors.New("budget exceeded")
		}
	}

	tx.ID = len(transactions) + 1
	transactions = append(transactions, tx)

	return nil
}

func ListTransactions() []Transaction {
	return transactions
}

func LoadBudgets(r io.Reader) error {
	var loadedBudgets []Budget

	err := json.NewDecoder(r).Decode(&loadedBudgets)
	if err != nil {
		return fmt.Errorf("failed to read budgets: %w", err)
	}

	for _, budget := range loadedBudgets {
		err := SetBudget(budget)
		if err != nil {
			return fmt.Errorf("failed to set budget: %w", err)
		}
	}

	return nil
}

func main() {
	fmt.Println("Ledger service started")

	file, err := os.Open("budgets.json")
	if err != nil {
		fmt.Println("Error opening budgets file:", err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	err = LoadBudgets(reader)
	if err != nil {
		fmt.Println("Error loading budgets:", err)
		return
	}

	fmt.Println("Budgets loaded successfully")

	tx1 := Transaction{
		Amount:      1000,
		Category:    "Food",
		Description: "Groceries",
		Date:        "2026-10-04",
	}

	tx2 := Transaction{
		Amount:      4500,
		Category:    "Food",
		Description: "Restaurant",
		Date:        "2026-10-04",
	}

	tx3 := Transaction{
		Amount:      500,
		Category:    "Transport",
		Description: "Taxi",
		Date:        "2026-10-04",
	}

	err = AddTransaction(tx1)
	if err != nil {
		fmt.Println("Transaction rejected:", err)
	} else {
		fmt.Println("Transaction added:", tx1.Description)
	}

	err = AddTransaction(tx2)
	if err != nil {
		fmt.Println("Transaction rejected:", err)
	} else {
		fmt.Println("Transaction added:", tx2.Description)
	}

	err = AddTransaction(tx3)
	if err != nil {
		fmt.Println("Transaction rejected:", err)
	} else {
		fmt.Println("Transaction added:", tx3.Description)
	}

	fmt.Println("Transactions:")

	for _, tx := range ListTransactions() {
		fmt.Println(tx)
	}
}