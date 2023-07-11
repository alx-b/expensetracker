package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"

	"github.com/alx-b/expensetracker/domain"
	"github.com/alx-b/expensetracker/logger"
)

type DB struct {
	db *sql.DB
}

// CreateDB opens sqlite database connection and returns pointer to DB struct.
func CreateDB() *DB {
	db, err := sql.Open("sqlite", "./db.sqlite3")
	if err != nil {
		logger.Error(fmt.Errorf("Could not open database: %w", err).Error())
	}

	if createExpensesTable(db); err != nil {
		logger.Error(err.Error())
	}

	if createBudgetTable(db); err != nil {
		logger.Error(err.Error())
	}

	if addDefaultMonthlyBudget(db); err != nil {
		logger.Error(err.Error())
	}

	return &DB{db: db}
}

// createExpensesTable takes in a database connection and
// creates the expenses table if it doesn't exist.
func createExpensesTable(db *sql.DB) error {
	result, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS expenses (
id INTEGER PRIMARY KEY,
name TEXT,
category TEXT,
date TEXT,
amount TEXT
)`,
	)
	if err != nil {
		return fmt.Errorf("Could not create table: %w", err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("Could not retrieve last inserted id: %w", err)
	}

	return nil
}

// createBudgetTable takes in a database connection and
// creates the budget table if it doesn't exist.
func createBudgetTable(db *sql.DB) error {
	result, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS budget (
id INTEGER PRIMARY KEY,
date TEXT UNIQUE,
amount TEXT
)`,
	)
	if err != nil {
		return fmt.Errorf("Could not create table: %w", err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("Could not retrieve last inserted id: %w", err)
	}

	return nil
}

// addDefaultMonthlyBudget takes in a database connection and inserts a
// monthly default budget amount of 0.00 using "default" instead of a date.
func addDefaultMonthlyBudget(db *sql.DB) error {
	result, err := db.Exec(
		"INSERT OR IGNORE INTO budget (date, amount) VALUES (?,?)",
		"default",
		"0.00",
	)
	if err != nil {
		return fmt.Errorf("Could not insert into table: %w", err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("Could not retrieve last inserted id: %w", err)
	}

	return nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	return db.db.Close()
}

// GetWithMonthYear returns expenses of a specific month and year (YYYY-MM).
func (db *DB) GetExpensesWithYearMonth(yearMonth string) []domain.Expense {
	rows, err := db.db.Query("SELECT name, date, amount, category FROM expenses WHERE date LIKE ?", yearMonth)
	if err != nil {
		logger.Error(fmt.Errorf("Could not query database: %w", err).Error())
	}

	defer rows.Close()

	list := []domain.Expense{}

	for rows.Next() {
		expense := domain.Expense{}
		rows.Scan(
			&expense.Name,
			&expense.Date,
			&expense.Amount,
			&expense.Category,
		)
		list = append(list, expense)
	}

	return list
}

// GetDefaultBudget returns the default monthly budget amount.
func (db *DB) GetDefaultBudget() string {
	rows, err := db.db.Query("SELECT amount FROM budget WHERE date='default'")
	if err != nil {
		logger.Error(fmt.Errorf("Could not query database: %w", err).Error())
	}

	defer rows.Close()

	amount := ""

	for rows.Next() {
		rows.Scan(
			&amount,
		)
	}

	return amount
}

// UpdateDefaultBudget updates the default monthly budget amount.
func (db DB) UpdateDefaultBudget(amount string) error {
	result, err := db.db.Exec("UPDATE budget SET amount=? WHERE date='default'", amount)
	if err != nil {
		return fmt.Errorf("Could not update table: %w", err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("Could not retrieve last inserted id: %w", err)
	}

	return nil
}

// InsertBudget inserts budget amount for a specific month and year (YYYY-MM).
func (db DB) InsertBudget(amount, date string) error {
	result, err := db.db.Exec("INSERT OR REPLACE INTO budget (amount, date) VALUES (?,?)", amount, date)
	if err != nil {
		return fmt.Errorf("Could not insert into table: %w", err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("Could not retrieve last inserted id: %w", err)
	}

	return nil
}

// GetBudgetWithYearMonth returns budget amount
// of a specific month and year (YYYY-MM).
func (db *DB) GetBudgetWithYearMonth(date string) string {
	rows, err := db.db.Query("SELECT amount FROM budget WHERE date=?", date)
	if err != nil {
		logger.Error(fmt.Errorf("Could not query database: %w", err).Error())
	}

	defer rows.Close()

	amount := ""

	for rows.Next() {
		rows.Scan(
			&amount,
		)
	}

	return amount
}

// InsertExpense inserts a given expense into expenses table.
func (db DB) InsertExpense(expense domain.Expense) error {
	result, err := db.db.Exec(
		"INSERT INTO expenses (name, date, amount, category) VALUES (?,?,?,?)",
		expense.Name,
		expense.Date,
		expense.Amount,
		expense.Category,
	)
	if err != nil {
		return fmt.Errorf("Could not insert into table: %w", err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("Could not retrieve last inserted id: %w", err)
	}

	return nil
}
