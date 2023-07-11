package controller

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/alx-b/expensetracker/domain"
	"github.com/alx-b/expensetracker/logger"
)

type Controller struct {
	db domain.Storage
}

// CreateController returns pointer to Controller struct
// which contains Storage interface.
func CreateController(db domain.Storage) *Controller {
	return &Controller{db: db}
}

// splitDate takes in a date string (YYYY-MM-DD | YYYY-MM) and
// returns a slice separating the year, month, day.
func splitDate(date string) []string {
	splittedDate := strings.Split(date, "-")

	if len(splittedDate) == 1 {
		splittedDate = strings.Split(date, ".")
	}

	if len(splittedDate) == 1 {
		splittedDate = strings.Split(date, "/")
	}

	return splittedDate
}

// isNumberBetween returns a bool "if number is between a and b inclusively".
func isNumberBetween(number, a, b int) bool {
	return number >= a && number <= b
}

// formatDate validate and format the date and returns a new date string.
func formatDate(dateString string) (string, error) {
	splittedDate := splitDate(dateString)
	splittedNumberDate := []int{}

	for i := range splittedDate {
		digit, err := strconv.Atoi(splittedDate[i])
		if err != nil {
			logger.Error("Should be integer")
			return "", err
		}
		splittedNumberDate = append(splittedNumberDate, digit)
	}

	year := splittedNumberDate[0]
	month := splittedNumberDate[1]
	day := func() int {
		if len(splittedNumberDate) == 3 {
			return splittedNumberDate[2]
		}
		return 0
	}()

	if !isNumberBetween(month, 1, 12) {
		return "", errors.New("Month should be from 1 to 12.")
	}

	if len(splittedNumberDate) == 3 && !isNumberBetween(day, 1, 31) {
		return "", errors.New("Day should be from 1 to 31.")
	}

	date := func() string {
		if day == 0 {
			return fmt.Sprintf("%d-%02d", year, month)
		}
		return fmt.Sprintf("%d-%02d-%02d", year, month, day)
	}()

	return date, nil
}

// AddExpense adds Expense to database if valid.
func (c *Controller) AddExpense(expense domain.Expense) error {
	date, err := formatDate(expense.Date)
	if err != nil {
		return err
	}

	expense.Date = date

	return c.db.InsertExpense(expense)
}

// InsertBudgetMonth adds budget amount to database if valid.
func (c *Controller) InsertBudgetMonth(amount, date string) error {
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return err
	}

	formattedAmount := fmt.Sprintf("%.2f", amountFloat)

	formattedDate, err := formatDate(date)
	if err != nil {
		return err
	}

	return c.db.InsertBudget(formattedAmount, formattedDate)
}

// UpdateDefaultBudget updates the amount of default budget.
func (c *Controller) UpdateDefaultBudget(amount string) error {
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return err
	}

	formattedAmount := fmt.Sprintf("%.2f", amountFloat)

	return c.db.UpdateDefaultBudget(formattedAmount)
}

// TODO
func (c *Controller) CreateMonthData(year int, monthNumber time.Month) domain.MonthData {
	expenses := c.getExpensesForYearMonth(year, monthNumber)
	totalSpendings := calculateTotalExpenses(expenses)

	yearMonth := fmt.Sprintf("%d-%02d", year, int(monthNumber))

	budgetMonth := c.db.GetBudgetWithYearMonth(yearMonth)

	if budgetMonth == "" {
		budgetMonth = c.db.GetDefaultBudget()
	}

	budget, err := strconv.ParseFloat(budgetMonth, 64)

	if err != nil {
		logger.Error("Could not parse string to float: " + err.Error())
		budget = 0.00
	}

	return domain.MonthData{
		Year:           year,
		Month:          monthNumber,
		Expenses:       expenses,
		Budget:         budget,
		TotalSpendings: totalSpendings,
		MoneyLeft:      budget - totalSpendings,
	}
}

// getExpensesForYearMonth format year and month into a string and
// calls the database to fetch expenses of that month and returns it.
func (c *Controller) getExpensesForYearMonth(year int, month time.Month) []domain.Expense {
	yearMonth := fmt.Sprintf("%%%d-%02d%%", year, int(month))
	spendings := c.db.GetExpensesWithYearMonth(yearMonth)

	return spendings
}

// calculateTotalSpending returns the total amount for all expenses.
func calculateTotalExpenses(expenses []domain.Expense) float64 {
	total := 0.00

	for _, s := range expenses {
		total += s.Amount
	}

	return total
}
