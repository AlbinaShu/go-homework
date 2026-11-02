package api

import (
	"encoding/json"
	"net/http"

	"go-homework/hw6/ledger"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func TransactionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createTransaction(w, r)
	case http.MethodGet:
		listTransactions(w)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func createTransaction(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	tx := requestToTransaction(req)
	if err := tx.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := ledger.AddTransaction(tx); err != nil {
		if err.Error() == "budget exceeded" {
			writeError(w, http.StatusConflict, "budget exceeded")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	transactions := ledger.ListTransactions()
	created := transactions[len(transactions)-1]
	writeJSON(w, http.StatusCreated, transactionToResponse(created))
}

func listTransactions(w http.ResponseWriter) {
	transactions := ledger.ListTransactions()
	response := make([]TransactionResponse, 0, len(transactions))
	for _, tx := range transactions {
		response = append(response, transactionToResponse(tx))
	}
	writeJSON(w, http.StatusOK, response)
}

func BudgetsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createBudget(w, r)
	case http.MethodGet:
		listBudgets(w)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func createBudget(w http.ResponseWriter, r *http.Request) {
	var req CreateBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	budget := requestToBudget(req)
	if err := budget.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := ledger.SetBudget(budget); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, budgetToResponse(budget))
}

func listBudgets(w http.ResponseWriter) {
	budgets := ledger.ListBudgets()
	response := make([]BudgetResponse, 0, len(budgets))
	for _, budget := range budgets {
		response = append(response, budgetToResponse(budget))
	}
	writeJSON(w, http.StatusOK, response)
}
