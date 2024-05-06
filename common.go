package main

import "time"

type Stock struct {
	Price float64
	Date  time.Time
}

type Decision string

const (
	SELL Decision = "SELL"
	BUY           = "BUY"
	HOLD          = "HOLD"
)

const MutationRate float64 = 0.3

type DecisionMaker interface {
	SellOrBuy([]Stock) (Decision, error)
}
