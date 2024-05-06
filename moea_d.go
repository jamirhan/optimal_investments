package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type CrossoverFuncton func(Individual, Individual) Individual

type MutateFunction func(Individual) Individual

type NewIndividualFunction func() Individual

type Individual struct {
	broker    DecisionMaker
	counted   bool
	fitResult float64
}

func (i *Individual) Fit_(stocks []Stock) float64 {
	startDate := time.Unix(1711967605, 0) // 01.04.2024
	gains := 0.0
	for startDate.Before(time.Unix(1714050742, 0)) { // 25.04.2024
		gain, _, err := EvaluateGains(stocks, i.broker, startDate, startDate.Add(1*24*time.Hour))
		if err != nil {
			panic(err)
		}
		gains += gain
		startDate = startDate.Add(5 * 24 * time.Hour)
	}
	return gains
}

func (i *Individual) FitFunction(stocks []Stock) float64 {
	if !i.counted {
		i.counted = true
		i.fitResult = i.Fit_(stocks)
	}
	return i.fitResult
}

type MOEAD struct {
	PopulationSize int
	DecisionSpace  int
	ObjectiveSpace int
	GroupSize      int
	Generations    int
	Population     []Individual
	Crossover      CrossoverFuncton
	Mutate         MutateFunction
	NewIndividual  NewIndividualFunction
	Stocks         []Stock
}

func (m *MOEAD) Initialize() {
	m.Population = make([]Individual, m.PopulationSize)

	for i := range m.Population {
		ind := m.NewIndividual()
		m.Population[i] = ind
	}
}

func (m *MOEAD) Breed(i, j int) {
	parent1 := m.Population[i]
	parent2 := m.Population[j]
	child := m.Crossover(parent1, parent2)
	child = m.Mutate(child)
	// fmt.Printf("breed %d and %d. fit function: %v", i, j, child.FitFunction(m.Stocks))
	if child.FitFunction(m.Stocks) >= m.Population[i].FitFunction(m.Stocks) {
		m.Population[i] = child
	}
}

func (m *MOEAD) Inbreed(start int) {
	for i := 0; i < m.GroupSize; i++ {
		another := min(start+rand.Intn(m.GroupSize), len(m.Population))
		m.Breed(start+i, another)
		another = rand.Intn(len(m.Population))
		m.Breed(i, another)
	}
}

func (m *MOEAD) Evolve() {
	for gen := 0; gen < m.Generations; gen++ {
		sort.Slice(m.Population, func(i, j int) bool {
			return m.Population[i].FitFunction(m.Stocks) < m.Population[j].FitFunction(m.Stocks)
		})
		fmt.Printf("generation %d: highest gain is %v\n", gen, m.Population[len(m.Population)-1].fitResult)
		for i := 0; i < len(m.Population); i += m.GroupSize {
			m.Inbreed(i)
		}
	}
}
