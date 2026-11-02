package main

import (
	"fmt"
	"net/http"

	"go-homework/hw6/gateway/internal/api"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/transactions", api.TransactionsHandler)
	mux.HandleFunc("/api/budgets", api.BudgetsHandler)

	handler := api.LoggingMiddleware(mux)

	fmt.Println("Gateway service started on port 8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		fmt.Println("Server error:", err)
	}
}
