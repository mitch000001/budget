package main

func NewBudget() *Budget {
	return &Budget{
		Einnahmen: make(map[string]float64),
		Ausgaben:  make(map[string]float64),
	}
}

type Budget struct {
	Einnahmen Earnings
	Ausgaben  Expenses
}

func (b *Budget) Balance() float64 {
	earnings := b.Einnahmen.Sum()
	expenses := b.Ausgaben.Sum()
	return earnings - expenses
}

type Earnings map[string]float64

func (e *Earnings) Sum() float64 {
	sum := 0.0
	for _, val := range *e {
		sum += val
	}
	return sum
}

type Expenses map[string]float64

func (e *Expenses) Sum() float64 {
	sum := 0.0
	for _, val := range *e {
		sum += val
	}
	return sum
}
