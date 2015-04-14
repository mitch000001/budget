package main

func NewBudget() *Budget {
	return &Budget{
		Einnahmen: make(map[string]float64),
		Ausgaben:  make(map[string]float64),
	}
}

type Budget struct {
	Einnahmen map[string]float64
	Ausgaben  map[string]float64
}

func (b *Budget) Balance() float64 {
	var balance float64
	return balance
}
