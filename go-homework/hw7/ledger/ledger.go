package ledger

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var ErrBudgetExceeded = errors.New("budget exceeded")

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

func (tx Transaction) Validate() error {
	if tx.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if tx.Category == "" {
		return errors.New("category cannot be empty")
	}
	if _, err := time.Parse("2006-01-02", tx.Date); err != nil {
		return errors.New("date must be in YYYY-MM-DD format")
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
	if db == nil {
		return errors.New("database is not initialized")
	}

	_, err := db.Exec(`
        INSERT INTO budgets(category, limit_amount)
        VALUES($1, $2)
        ON CONFLICT(category)
        DO UPDATE SET limit_amount = EXCLUDED.limit_amount
    `, b.Category, b.Limit)
	return err
}

func ListBudgets() ([]Budget, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	rows, err := db.Query(`SELECT category, limit_amount FROM budgets ORDER BY category`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []Budget{}
	for rows.Next() {
		var b Budget
		if err := rows.Scan(&b.Category, &b.Limit); err != nil {
			return nil, err
		}
		result = append(result, b)
	}
	return result, rows.Err()
}

func AddTransaction(tx Transaction) error {
	if err := tx.Validate(); err != nil {
		return err
	}
	if db == nil {
		return errors.New("database is not initialized")
	}

	var limit float64
	err := db.QueryRow(`SELECT limit_amount FROM budgets WHERE category = $1`, tx.Category).Scan(&limit)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if err == nil {
		var spent float64
		if err := db.QueryRow(`
            SELECT COALESCE(SUM(amount), 0)
            FROM expenses
            WHERE category = $1
        `, tx.Category).Scan(&spent); err != nil {
			return err
		}

		if spent+tx.Amount > limit {
			return ErrBudgetExceeded
		}
	}

	date, _ := time.Parse("2006-01-02", tx.Date)
	if err := db.QueryRow(`
        INSERT INTO expenses(amount, category, description, date)
        VALUES($1, $2, $3, $4)
        RETURNING id
    `, tx.Amount, tx.Category, tx.Description, date).Scan(&tx.ID); err != nil {
		return fmt.Errorf("insert transaction: %w", err)
	}

	return nil
}

func ListTransactions() ([]Transaction, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	rows, err := db.Query(`
        SELECT id, amount, category, COALESCE(description, ''), date
        FROM expenses
        ORDER BY date DESC, id DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []Transaction{}
	for rows.Next() {
		var tx Transaction
		var date time.Time
		if err := rows.Scan(&tx.ID, &tx.Amount, &tx.Category, &tx.Description, &date); err != nil {
			return nil, err
		}
		tx.Date = date.Format("2006-01-02")
		result = append(result, tx)
	}
	return result, rows.Err()
}

func GetLastTransaction() (Transaction, error) {
	if db == nil {
		return Transaction{}, errors.New("database is not initialized")
	}

	var tx Transaction
	var date time.Time
	err := db.QueryRow(`
		SELECT id, amount, category, COALESCE(description, ''), date
		FROM expenses
		ORDER BY id DESC
		LIMIT 1
	`).Scan(&tx.ID, &tx.Amount, &tx.Category, &tx.Description, &date)
	if err != nil {
		return Transaction{}, err
	}
	tx.Date = date.Format("2006-01-02")
	return tx, nil
}
