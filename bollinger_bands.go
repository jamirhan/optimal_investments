package main

import (
	"math"
	"math/rand"
)

func meanPrice(stocks []Stock) float64 {
	total := 0.0
	for _, number := range stocks {
		total += number.Price
	}
	return total / float64(len(stocks))
}

func stdDev(stocks []Stock, period int) float64 {
	stocks = stocks[len(stocks)-period:]
	m := meanPrice(stocks)
	total := 0.0
	for _, stock := range stocks {
		total += math.Pow(stock.Price-m, 2)
	}
	variance := total / float64(len(stocks)-1)
	return math.Sqrt(variance)
}

func bollingerBands(stocks []Stock, meanPeriod, deviationPeriod int, mul float64) (float64, float64, error) {
	m, err := getLWMA(stocks, meanPeriod)
	if err != nil {
		return 0, 0, err
	}
	sd := stdDev(stocks, deviationPeriod)
	return m - mul*sd, m + mul*sd, nil
}

type BBDecisionMaker struct {
	meanPeriod      int
	deviationPeriod int
	mul             float64
}

func (d *BBDecisionMaker) SellOrBuy(stocks []Stock) (Decision, error) {
	lowerBand, higherBand, err := bollingerBands(stocks, d.meanPeriod, d.deviationPeriod, d.mul)
	if err != nil {
		return "", err
	}
	lastPrice := stocks[len(stocks)-1].Price
	if lastPrice >= higherBand {
		return BUY, nil
	}
	if lastPrice <= lowerBand {
		return SELL, nil
	}
	return HOLD, nil
}

func CrossoverBB(i1, i2 Individual) Individual {
	d1, ok := i1.broker.(*BBDecisionMaker)
	if !ok {
		panic("not convertable")
	}
	d2, ok := i2.broker.(*BBDecisionMaker)
	if !ok {
		panic("not convertable")
	}
	return Individual{
		broker: &BBDecisionMaker{
			meanPeriod:      int(0.5*float64(d1.meanPeriod) + 0.5*float64(d2.meanPeriod)),
			deviationPeriod: int(0.5*float64(d1.deviationPeriod) + 0.5*float64(d2.deviationPeriod)),
			mul:             0.5*float64(d1.mul) + 0.5*float64(d2.mul),
		},
	}
}

func MutateBB(i Individual) Individual {
	d, ok := i.broker.(*BBDecisionMaker)
	if !ok {
		panic("not convertable")
	}
	res := &BBDecisionMaker{
		meanPeriod:      d.meanPeriod,
		deviationPeriod: d.deviationPeriod,
		mul:             d.mul,
	}

	if rand.Float64() < MutationRate {
		d.meanPeriod += rand.Intn(5) - 3
		if d.meanPeriod <= 0 {
			d.meanPeriod = 1
		}
	}
	if rand.Float64() < MutationRate {
		d.deviationPeriod += rand.Intn(5) - 3
		if d.deviationPeriod <= 0 {
			d.deviationPeriod = 1
		}
	}
	if rand.Float64() < MutationRate {
		d.mul += 3*rand.Float64() - 1
	}
	return Individual{
		broker: res,
	}
}

func NewBB() Individual {
	return Individual{
		broker: &BBDecisionMaker{
			meanPeriod:      rand.Intn(20),
			deviationPeriod: rand.Intn(20),
			mul:             1 + float64(rand.Intn(3)),
		},
	}
}
