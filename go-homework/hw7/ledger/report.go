package ledger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type ReportItem struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
}

type ReportSummary struct {
	From  string       `json:"from"`
	To    string       `json:"to"`
	Items []ReportItem `json:"items"`
}

func GetReportSummary(ctx context.Context, from, to string) (ReportSummary, error) {
	if db == nil {
		return ReportSummary{}, errors.New("database is not initialized")
	}
	if redisClient == nil {
		return ReportSummary{}, errors.New("redis is not initialized")
	}

	fromDate, err := time.Parse("2006-01-02", from)
	if err != nil {
		return ReportSummary{}, errors.New("from must be in YYYY-MM-DD format")
	}
	toDate, err := time.Parse("2006-01-02", to)
	if err != nil {
		return ReportSummary{}, errors.New("to must be in YYYY-MM-DD format")
	}
	if fromDate.After(toDate) {
		return ReportSummary{}, errors.New("from must not be after to")
	}

	key := fmt.Sprintf("report:summary:%s:%s", from, to)
	cached, err := redisClient.Get(ctx, key).Bytes()
	if err == nil {
		var result ReportSummary
		if err := json.Unmarshal(cached, &result); err == nil {
			log.Printf("report cache hit: %s", key)
			return result, nil
		}
	} else if !errors.Is(err, redis.Nil) {
		log.Printf("redis get error: %v", err)
	}

	log.Printf("report cache miss: %s", key)
	rows, err := db.QueryContext(ctx, `
        SELECT category, COALESCE(SUM(amount), 0)
        FROM expenses
        WHERE date BETWEEN $1 AND $2
        GROUP BY category
        ORDER BY category
    `, fromDate, toDate)
	if err != nil {
		return ReportSummary{}, err
	}
	defer rows.Close()

	result := ReportSummary{From: from, To: to, Items: []ReportItem{}}
	for rows.Next() {
		var item ReportItem
		if err := rows.Scan(&item.Category, &item.Total); err != nil {
			return ReportSummary{}, err
		}
		result.Items = append(result.Items, item)
	}
	if err := rows.Err(); err != nil {
		return ReportSummary{}, err
	}

	data, err := json.Marshal(result)
	if err == nil {
		if err := redisClient.Set(ctx, key, data, 30*time.Second).Err(); err != nil {
			log.Printf("redis set error: %v", err)
		}
	}

	return result, nil
}
