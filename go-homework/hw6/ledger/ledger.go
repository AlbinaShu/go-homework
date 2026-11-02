package ledger

import "errors"

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

func SetBudget(b Budget) error {
	if err := b.Validate(); err != nil {
		return err
	}
	budgets[b.Category] = b
	return nil
}

func ListBudgets() []Budget {
	result := []Budget{}
	for _, budget := range budgets {
		result = append(result, budget)
	}
	return result
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
	if err := tx.Validate(); err != nil {
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

func Reset() {
	transactions = []Transaction{}
	budgets = make(map[string]Budget)
}
