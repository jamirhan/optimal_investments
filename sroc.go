package main

import (
	"math/rand"
)

func EXMA(stocks []Stock, period int, alpha float64) float64 {
	total := 0.0
	total += stocks[len(stocks)-1].Price
	for i := 1; i < period; i++ {
		total += stocks[len(stocks)-i-1].Price
	}
	return total / (float64(period))
}

func calculateSRoC(stocks []Stock, period int, alpha float64) float64 {
	roc := (stocks[len(stocks)-1].Price - stocks[len(stocks)-1-period].Price) / stocks[len(stocks)-1-period].Price
	// fmt.Println(roc)
	sma := EXMA(stocks, period, alpha)
	return (roc - sma) / sma
}

type SRoCDecisionMaker struct {
	period int
	alpha  float64
}

func (d *SRoCDecisionMaker) SellOrBuy(stocks []Stock) (Decision, error) {
	sroc := calculateSRoC(stocks, d.period, d.alpha)
	// fmt.Println(sroc)
	if sroc > 0.2 {
		return BUY, nil
	}
	if sroc < -0.2 {
		return SELL, nil
	}

	return HOLD, nil
}

func CrossoverSRoC(i1, i2 Individual) Individual {
	d1, ok := i1.broker.(*SRoCDecisionMaker)
	if !ok {
		panic("not convertable")
	}
	d2, ok := i2.broker.(*SRoCDecisionMaker)
	if !ok {
		panic("not convertable")
	}
	return Individual{
		broker: &SRoCDecisionMaker{
			period: int(0.5*float64(d1.period) + 0.5*float64(d2.period)),
			alpha:  0.5*float64(d1.alpha) + 0.5*float64(d2.alpha),
		},
	}
}

func MutateSRoC(i Individual) Individual {
	d, ok := i.broker.(*SRoCDecisionMaker)
	if !ok {
		panic("not convertable")
	}
	res := &SRoCDecisionMaker{
		period: d.period,
		alpha:  d.alpha,
	}

	if rand.Float64() < MutationRate {
		d.period += rand.Intn(5) - 3
		if d.period <= 0 {
			d.period = 1
		}
	}
	if rand.Float64() < MutationRate {
		d.alpha += rand.Float64()
	}
	return Individual{
		broker: res,
	}
}

func NewSRoC() Individual {
	return Individual{
		broker: &SRoCDecisionMaker{
			period: rand.Intn(20),
			alpha:  rand.Float64(),
		},
	}
}
