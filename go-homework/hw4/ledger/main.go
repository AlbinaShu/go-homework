package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

type Validatable interface {
	Validate() error
}

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

func (tx Transaction) Validate() error {
	if tx.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	if tx.Category == "" {
		return errors.New("category cannot be empty")
	}

	return nil
}

func (b Budget) Validate() error {
	if b.Category == "" {
		return errors.New("category cannot be empty")
	}

	if b.Limit <= 0 {
		return errors.New("budget limit must be greater than zero")
	}

	return nil
}

func CheckValid(v Validatable) error {
	return v.Validate()
}

func SetBudget(b Budget) error {
	err := b.Validate()
	if err != nil {
		return err
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
	err := tx.Validate()
	if err != nil {
		return err
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

	testTransaction := Transaction{
		Amount:      100,
		Category:    "Food",
		Description: "Test transaction",
		Date:        "2026-10-31",
	}

	testBudget := Budget{
		Category: "Food",
		Limit:    5000,
	}

	err = CheckValid(testTransaction)
	if err != nil {
		fmt.Println("Transaction validation error:", err)
	} else {
		fmt.Println("Transaction is valid")
	}

	err = CheckValid(testBudget)
	if err != nil {
		fmt.Println("Budget validation error:", err)
	} else {
		fmt.Println("Budget is valid")
	}

	tx1 := Transaction{Amount: 1000, Category: "Food", Description: "Groceries", Date: "2026-10-31"}
	tx2 := Transaction{Amount: 4500, Category: "Food", Description: "Restaurant", Date: "2026-10-31"}
	tx3 := Transaction{Amount: 500, Category: "Transport", Description: "Taxi", Date: "2026-10-31"}

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
