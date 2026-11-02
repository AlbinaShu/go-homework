package ledger

import "testing"

func TestTransactionValidate(t *testing.T) {
	tests := []struct {
		name    string
		tx      Transaction
		wantErr bool
	}{
		{
			name:    "valid transaction",
			tx:      Transaction{Amount: 100, Category: "Food"},
			wantErr: false,
		},
		{
			name:    "zero amount",
			tx:      Transaction{Amount: 0, Category: "Food"},
			wantErr: true,
		},
		{
			name:    "negative amount",
			tx:      Transaction{Amount: -100, Category: "Food"},
			wantErr: true,
		},
		{
			name:    "empty category",
			tx:      Transaction{Amount: 100, Category: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.tx.Validate()
			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected nil, got %v", err)
			}
		})
	}
}

func TestBudgetValidate(t *testing.T) {
	tests := []struct {
		name    string
		budget  Budget
		wantErr bool
	}{
		{
			name:    "valid budget",
			budget:  Budget{Category: "Food", Limit: 5000},
			wantErr: false,
		},
		{
			name:    "zero limit",
			budget:  Budget{Category: "Food", Limit: 0},
			wantErr: true,
		},
		{
			name:    "negative limit",
			budget:  Budget{Category: "Food", Limit: -100},
			wantErr: true,
		},
		{
			name:    "empty category",
			budget:  Budget{Category: "", Limit: 5000},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.budget.Validate()
			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected nil, got %v", err)
			}
		})
	}
}

func TestAddTransactionBudgetExceeded(t *testing.T) {
	Reset()
	t.Cleanup(Reset)

	err := SetBudget(Budget{Category: "Food", Limit: 5000})
	if err != nil {
		t.Fatalf("failed to set budget: %v", err)
	}

	err = AddTransaction(Transaction{Amount: 1000, Category: "Food"})
	if err != nil {
		t.Fatalf("first transaction should be added: %v", err)
	}

	err = AddTransaction(Transaction{Amount: 4500, Category: "Food"})
	if err == nil {
		t.Fatal("expected budget exceeded error")
	}
	if err.Error() != "budget exceeded" {
		t.Fatalf("expected budget exceeded, got %v", err)
	}

	transactions := ListTransactions()
	if len(transactions) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(transactions))
	}
}
