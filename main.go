package main

import (
	"os"

	"gioui.org/app"
	"gioui.org/unit"

	"github.com/alx-b/expensetracker/controller"
	"github.com/alx-b/expensetracker/database"
	"github.com/alx-b/expensetracker/logger"
	"github.com/alx-b/expensetracker/ui"
)

func main() {
	defer logger.CloseFile()

	db := database.CreateDB()
	defer db.Close()

	controller := controller.CreateController(db)

	go func() {
		w := app.NewWindow(app.Title("Simple Expense Tracker"), app.Size(unit.Dp(500), unit.Dp(700)))
		if err := ui.Run(w, controller); err != nil {
			logger.Error(err.Error())
		}
		os.Exit(0)
	}()
	app.Main()
}
