package main

import (
	"fmt"
	"time"
)

const (
	bought = iota
	sold
	notStarted
)

func EvaluateGains(stocks []Stock, m DecisionMaker, from, to time.Time) (float64, int, error) {
	buyCount := 0
	gains := 0.0
	state := notStarted
	boughtAt := -1.0

	fromInd := -1
	toInd := len(stocks) - 1
	for i, stock := range stocks {
		if stock.Date.After(from) && fromInd == -1 {
			fromInd = i
		}
		if stock.Date.Before(to) {
			toInd = i
		}
	}
	if fromInd == -1 {
		return 0, 0, fmt.Errorf("wrong FROM date")
	}

	for i := fromInd; i <= toInd; i += 1 {
		decision, err := m.SellOrBuy(stocks[:i+1])
		if err != nil {
			return 0, 0, err
		}
		if decision == BUY {
			if state == notStarted || state == sold {
				boughtAt = stocks[i+1].Price
				state = bought
			}
		}
		if decision == SELL {
			if state == bought {
				buyCount += 1
				gains += stocks[i+1].Price - boughtAt
				state = sold
			}
		}
	}

	return gains, buyCount, nil
}
