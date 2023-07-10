package domain

// STRUCTS
type Expense struct {
	Name     string
	Date     string
	Amount   float64
	Category string
}

// INTERFACES
type Storage interface {
	GetExpensesWithYearMonth(string) []Expense
	InsertExpense(Expense) error
	GetDefaultBudget() string
	GetBudgetWithYearMonth(string) string
	InsertBudget(string, string) error
	UpdateDefaultBudget(string) error
}
