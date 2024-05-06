package main

import (
	"fmt"
	"math/rand"
)

func getLWMA(stocks []Stock, periods int) (float64, error) {
	if periods > len(stocks) {
		return 0, fmt.Errorf("too many periods requested")
	}

	weightSum := 0
	weightedPriceSum := 0.0

	for i := 0; i < periods; i++ {
		weight := periods - i
		weightSum += weight
		weightedPriceSum += stocks[len(stocks)-1-i].Price * float64(weight)
	}

	return weightedPriceSum / float64(weightSum), nil
}

type LWMADecisionMaker struct {
	shortPeriod int
	longPeriod  int
}

func (d *LWMADecisionMaker) SellOrBuy(stocks []Stock) (Decision, error) {
	long, err := getLWMA(stocks, d.longPeriod)
	if err != nil {
		return "", err
	}

	short, err := getLWMA(stocks, d.shortPeriod)
	if err != nil {
		return "", err
	}

	if short >= long {
		return BUY, nil
	}

	return SELL, nil
}

func CrossoverLWMA(i1, i2 Individual) Individual {
	d1, ok := i1.broker.(*LWMADecisionMaker)
	if !ok {
		panic("not convertable")
	}
	d2, ok := i2.broker.(*LWMADecisionMaker)
	if !ok {
		panic("not convertable")
	}
	return Individual{
		broker: &LWMADecisionMaker{
			shortPeriod: int(0.5*float64(d1.shortPeriod) + 0.5*float64(d2.shortPeriod)),
			longPeriod:  int(0.5*float64(d1.longPeriod) + 0.5*float64(d2.longPeriod)),
		},
	}
}

func MutateLWMA(i Individual) Individual {
	d, ok := i.broker.(*LWMADecisionMaker)
	if !ok {
		panic("not convertable")
	}
	res := &LWMADecisionMaker{
		shortPeriod: d.shortPeriod,
		longPeriod:  d.longPeriod,
	}

	if rand.Float64() < MutationRate {
		d.shortPeriod += rand.Intn(3) - 1
		if d.shortPeriod <= 0 {
			d.shortPeriod = 1
		}
	}
	if rand.Float64() < MutationRate {
		d.longPeriod += rand.Intn(3) - 1
		if d.longPeriod <= 0 {
			d.longPeriod = 1
		}
	}
	return Individual{
		broker: res,
	}
}

func NewLWMA() Individual {
	short := rand.Intn(10)
	return Individual{
		broker: &LWMADecisionMaker{
			shortPeriod: short,
			longPeriod:  short + 1 + rand.Intn(10),
		},
	}
}
