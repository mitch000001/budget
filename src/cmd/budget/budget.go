package main

import "fmt"

func NewBudget() *Budget {
	return &Budget{
		Einnahmen: NewBudgetColumn(),
		Ausgaben:  NewBudgetColumn(),
	}
}

type Budget struct {
	Einnahmen BudgetColumn
	Ausgaben  BudgetColumn
}

func (b *Budget) Balance() float64 {
	earnings := b.Einnahmen.Sum()
	expenses := b.Ausgaben.Sum()
	return earnings - expenses
}

func NewBudgetColumn() BudgetColumn {
	entry := make(BudgetColumn)
	return entry
}

type BudgetColumn map[string]*BudgetColumnEntry

func (b BudgetColumn) Sum() float64 {
	sum := 0.0
	for _, columnEntry := range b {
		sum += columnEntry.Value()
	}
	return sum
}

func NewBudgetColumnEntry(val float64, bookingDate ShortDate) *BudgetColumnEntry {
	return &BudgetColumnEntry{
		value:       val,
		bookingDate: bookingDate,
	}
}

type BudgetColumnEntry struct {
	value       float64
	bookingDate ShortDate
}

func (b *BudgetColumnEntry) Value() float64 {
	return b.value
}

func (b *BudgetColumnEntry) BookingDate() ShortDate {
	return b.bookingDate
}

func (b *BudgetColumnEntry) String() string {
	return fmt.Sprintf("%.2f", b.value)
}
