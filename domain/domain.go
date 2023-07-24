package domain

import "time"

// STRUCTS
type Expense struct {
	Id       int
	Name     string
	Date     string
	Amount   float64
	Category string
}

type MonthData struct {
	Year           int
	Month          time.Month
	Expenses       []Expense
	Budget         float64
	TotalSpendings float64
	MoneyLeft      float64
}

// INTERFACES
type Storage interface {
	GetExpensesWithYearMonth(string) []Expense
	InsertExpense(Expense) error
	GetDefaultBudget() string
	GetBudgetWithYearMonth(string) string
	InsertBudget(string, string) error
	UpdateDefaultBudget(string) error
	DeleteExpense(int) error
}

type API interface {
	CreateMonthData(int, time.Month) MonthData
	AddExpense(Expense) error
	RemoveExpense(int) error
	InsertBudgetMonth(string, string) error
	UpdateDefaultBudget(string) error
}
