package main

import (
	"fmt"
	"time"
)

func dryRun(stocks []Stock, d DecisionMaker) {
	fmt.Println(EvaluateGains(stocks, d, time.Unix(1714571500, 0), time.Unix(1714732405, 0)))
}

func learnParamsLWMA(stocks []Stock) {
	moea := MOEAD{
		PopulationSize: 100,
		DecisionSpace:  1,
		ObjectiveSpace: 2,
		Generations:    10,
		Crossover:      CrossoverLWMA,
		Mutate:         MutateLWMA,
		NewIndividual:  NewLWMA,
		Stocks:         stocks,
		GroupSize:      10,
	}

	moea.Initialize()
	moea.Evolve()

	for _, ind := range moea.Population {
		fmt.Printf("broker: %v, gains: %v\n", ind.broker, ind.FitFunction(stocks))
	}
}

func learnParamsSRoC(stocks []Stock) {
	moea := MOEAD{
		PopulationSize: 100,
		DecisionSpace:  1,
		ObjectiveSpace: 2,
		Generations:    10,
		Crossover:      CrossoverSRoC,
		Mutate:         MutateSRoC,
		NewIndividual:  NewSRoC,
		Stocks:         stocks,
		GroupSize:      10,
	}

	moea.Initialize()
	moea.Evolve()

	for _, ind := range moea.Population {
		fmt.Printf("broker: %v, gains: %v\n", ind.broker, ind.FitFunction(stocks))
	}
}

func learnParamsBB(stocks []Stock) {
	moea := MOEAD{
		PopulationSize: 100,
		DecisionSpace:  1,
		ObjectiveSpace: 2,
		Generations:    10,
		Crossover:      CrossoverBB,
		Mutate:         MutateBB,
		NewIndividual:  NewBB,
		Stocks:         stocks,
		GroupSize:      10,
	}

	moea.Initialize()
	moea.Evolve()

	for _, ind := range moea.Population {
		fmt.Printf("broker: %v, gains: %v\n", ind.broker, ind.FitFunction(stocks))
	}
}

func main() {
	stks, err := LoadStocks()
	if err != nil {
		panic(err)
	}
	stocks := SimplifyStocks(stks)

	// dryRun(stocks, &LWMADecisionMaker{
	// 	shortPeriod: 1,
	// 	longPeriod:  7,
	// })

	// dryRun(stocks, &BBDecisionMaker{
	// 	meanPeriod:      11,
	// 	deviationPeriod: 7,
	// 	mul:             2.2806640625,
	// })

	learnParamsSRoC(stocks)
	// learnParamsBB(stocks)
	// learnParamsLWMA(stocks)
}
