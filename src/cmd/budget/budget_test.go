package main

import (
	"testing"
	"time"
)

func TestNewBudget(t *testing.T) {
	budget := NewBudget()

	if budget == nil {
		t.Logf("Expected budget not to be nil")
		t.FailNow()
	}

	if budget.Einnahmen == nil {
		t.Logf("Expected Einnahmen not to be nil")
		t.Fail()
	}

	if budget.Ausgaben == nil {
		t.Logf("Expected Ausgaben not to be nil")
		t.Fail()
	}
}

func TestBudgetBalance(t *testing.T) {
	budget := NewBudget()

	budget.Einnahmen["foo"] = NewBudgetColumnEntry(400, ShortDate{})
	budget.Einnahmen["bar"] = NewBudgetColumnEntry(300, ShortDate{})
	budget.Ausgaben["baz"] = NewBudgetColumnEntry(100, ShortDate{})
	budget.Ausgaben["qux"] = NewBudgetColumnEntry(900, ShortDate{})

	balance := budget.Balance()

	expectedBalance := -300.0

	if expectedBalance != balance {
		t.Logf("Expected balance to equal %f, got %f\n", expectedBalance, balance)
		t.Fail()
	}
}

func TestNewBudgetColumn(t *testing.T) {
	entry := NewBudgetColumn()

	if entry == nil {
		t.Logf("Expected entry not to be nil")
		t.FailNow()
	}
}

func TestBudgetColumnSum(t *testing.T) {
	column := NewBudgetColumn()

	column["foo"] = NewBudgetColumnEntry(23, ShortDate{})
	column["bar"] = NewBudgetColumnEntry(55.0, ShortDate{})

	sum := column.Sum()

	if sum != 78 {
		t.Logf("Expected sum to equal %f, got %f\n", 78, sum)
		t.Fail()
	}
}

func TestNewBudgetColumnEntry(t *testing.T) {
	entry := NewBudgetColumnEntry(22, Date(2015, 01, 01, time.Local))

	if entry == nil {
		t.Logf("Expected entry not to be nil")
		t.FailNow()
	}

	if entry.value != 22 {
		t.Logf("Expected value to equal %f, got %f\n", 22, entry.value)
		t.Fail()
	}

	expectedDate := Date(2015, 01, 01, time.Local)
	if entry.bookingDate != expectedDate {
		t.Logf("Expected value to equal %q, got %q\n", expectedDate, entry.bookingDate)
		t.Fail()
	}
}

func TestBudgetColumnEntryValue(t *testing.T) {
	entry := NewBudgetColumnEntry(22, ShortDate{})

	if entry.Value() != 22 {
		t.Logf("Expected Value() to equal %f, got %f\n", 22, entry.Value())
		t.Fail()
	}
}

func TestBudgetColumnEntryBookingDate(t *testing.T) {
	entry := NewBudgetColumnEntry(22, Date(2015, 01, 01, time.Local))

	expectedDate := Date(2015, 01, 01, time.Local)
	if entry.BookingDate() != expectedDate {
		t.Logf("Expected value to equal %q, got %q\n", expectedDate, entry.BookingDate())
		t.Fail()
	}
}

func TestBudgetColumnEntryString(t *testing.T) {
	entry := NewBudgetColumnEntry(2.2, ShortDate{})

	if entry.String() != "2.20" {
		t.Logf("Expected Value() to equal %q, got %q\n", "2.20", entry.String())
		t.Fail()
	}
}
