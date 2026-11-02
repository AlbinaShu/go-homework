package main

import (
	"context"
	"fmt"
	"net/http"

	"go-homework/hw7/gateway/internal/api"
	"go-homework/hw7/ledger"
)

func main() {
	if err := ledger.InitDB(); err != nil {
		fmt.Println("Database connection error:", err)
		return
	}
	defer ledger.CloseDB()
	fmt.Println("Ledger connected to PostgreSQL")

	if err := ledger.InitRedis(context.Background()); err != nil {
		fmt.Println("Redis connection error:", err)
		return
	}
	defer ledger.CloseRedis()
	fmt.Println("Ledger connected to Redis")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/transactions", api.TransactionsHandler)
	mux.HandleFunc("/api/budgets", api.BudgetsHandler)
	mux.HandleFunc("/api/reports/summary", api.ReportSummaryHandler)

	handler := api.LoggingMiddleware(mux)

	fmt.Println("Gateway service started on port 8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		fmt.Println("Server error:", err)
	}
}
