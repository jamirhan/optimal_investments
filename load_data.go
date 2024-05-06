package main

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

const filename = "data.json"

type StockRaw struct {
	timestamp int
	high      float64
	low       float64
	open      float64
	close     float64
	volume    float64
}

func LoadStocks() ([]StockRaw, error) {
	var raw [][]float64
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bytes, &raw); err != nil {
		return nil, err
	}
	var res []StockRaw
	for _, rawStock := range raw {
		res = append(res, StockRaw{
			timestamp: int(rawStock[0]),
			open:      rawStock[1],
			high:      rawStock[2],
			low:       rawStock[3],
			close:     rawStock[4],
			volume:    rawStock[5],
		})
	}
	return res, nil
}

func SimplifyStocks(stocks []StockRaw) []Stock {
	var res []Stock
	for _, raw := range stocks {
		res = append(res, Stock{
			Date:  time.Unix(int64(raw.timestamp)/1000, 0),
			Price: raw.close,
		})
	}
	return res
}

func LoadTillTimestamp(stocks []Stock, timestamp time.Time) []Stock {
	for ind, stock := range stocks {
		if stock.Date.After(timestamp) {
			return stocks[:ind]
		}
	}
	return stocks
}
