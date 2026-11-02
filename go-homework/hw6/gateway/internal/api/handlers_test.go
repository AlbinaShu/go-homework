package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-homework/hw6/ledger"
)

func TestBudgetsHandler(t *testing.T) {
	ledger.Reset()
	t.Cleanup(ledger.Reset)

	t.Run("post_ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/budgets", strings.NewReader(`{"category":"Food","limit":5000}`))
		rec := httptest.NewRecorder()

		BudgetsHandler(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d", rec.Code)
		}
		if rec.Header().Get("Content-Type") != "application/json; charset=utf-8" {
			t.Fatalf("unexpected Content-Type: %s", rec.Header().Get("Content-Type"))
		}

		var response BudgetResponse
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if response.Category != "Food" || response.Limit != 5000 {
			t.Fatalf("unexpected response: %+v", response)
		}
	})

	t.Run("get_ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/budgets", nil)
		rec := httptest.NewRecorder()

		BudgetsHandler(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rec.Code)
		}
		if rec.Header().Get("Content-Type") != "application/json; charset=utf-8" {
			t.Fatalf("unexpected Content-Type: %s", rec.Header().Get("Content-Type"))
		}

		var response []BudgetResponse
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if len(response) != 1 || response[0].Category != "Food" {
			t.Fatalf("unexpected response: %+v", response)
		}
	})

	t.Run("invalid_limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/budgets", strings.NewReader(`{"category":"Transport","limit":0}`))
		rec := httptest.NewRecorder()

		BudgetsHandler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rec.Code)
		}
		if rec.Header().Get("Content-Type") != "application/json; charset=utf-8" {
			t.Fatalf("unexpected Content-Type: %s", rec.Header().Get("Content-Type"))
		}

		var response map[string]string
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if response["error"] == "" {
			t.Fatal("expected error message")
		}
	})
}

func TestTransactionsHandler(t *testing.T) {
	ledger.Reset()
	t.Cleanup(ledger.Reset)

	budgetReq := httptest.NewRequest(http.MethodPost, "/api/budgets", strings.NewReader(`{"category":"Food","limit":5000}`))
	budgetRec := httptest.NewRecorder()
	BudgetsHandler(budgetRec, budgetReq)
	if budgetRec.Code != http.StatusCreated {
		t.Fatalf("failed to create budget, status %d", budgetRec.Code)
	}

	t.Run("ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/transactions", strings.NewReader(`{"amount":1000,"category":"Food","description":"Groceries","date":"2026-11-02"}`))
		rec := httptest.NewRecorder()

		TransactionsHandler(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d", rec.Code)
		}
		if rec.Header().Get("Content-Type") != "application/json; charset=utf-8" {
			t.Fatalf("unexpected Content-Type: %s", rec.Header().Get("Content-Type"))
		}

		var response TransactionResponse
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if response.ID != 1 || response.Amount != 1000 || response.Category != "Food" {
			t.Fatalf("unexpected response: %+v", response)
		}
	})

	t.Run("get_ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
		rec := httptest.NewRecorder()

		TransactionsHandler(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rec.Code)
		}

		var response []TransactionResponse
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if len(response) != 1 || response[0].Description != "Groceries" {
			t.Fatalf("unexpected response: %+v", response)
		}
	})

	t.Run("exceeded", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/transactions", strings.NewReader(`{"amount":4500,"category":"Food","description":"Restaurant","date":"2026-11-02"}`))
		rec := httptest.NewRecorder()

		TransactionsHandler(rec, req)

		if rec.Code != http.StatusConflict {
			t.Fatalf("expected status 409, got %d", rec.Code)
		}

		var response map[string]string
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if response["error"] != "budget exceeded" {
			t.Fatalf("unexpected error: %s", response["error"])
		}
	})

	t.Run("bad_json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/transactions", strings.NewReader(`{"amount":`))
		rec := httptest.NewRecorder()

		TransactionsHandler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rec.Code)
		}
		if rec.Header().Get("Content-Type") != "application/json; charset=utf-8" {
			t.Fatalf("unexpected Content-Type: %s", rec.Header().Get("Content-Type"))
		}

		var response map[string]string
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if response["error"] == "" {
			t.Fatal("expected error message")
		}
	})
}
