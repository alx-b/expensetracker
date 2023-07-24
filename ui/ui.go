package ui

import (
	"image/color"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"github.com/alx-b/expensetracker/domain"
)

type Page int

const (
	List Page = iota
	Add
)

func Run(w *app.Window, controller domain.API) error {
	// Defines material design style
	th := material.NewTheme(gofont.Collection())
	th.Bg = color.NRGBA{100, 100, 100, 255}
	th.Fg = color.NRGBA{200, 200, 200, 255}
	// Operations from the UI
	var ops op.Ops

	currentPage := List
	year, month, _ := time.Now().Date()
	monthView := controller.CreateMonthData(year, month)

	// Create UI parts
	topBar := createTopBar(th, &currentPage, &monthView, controller)
	list := createListContainer(th, &monthView, controller)
	dataDisplay := createDataDisplay(th, controller, &monthView)
	addFormPage := createFormPage(th, controller)

	for {
		e := <-w.Events()

		switch e := e.(type) {
		case system.FrameEvent:
			// CONTEXT = Reference to operation + event
			gtx := layout.NewContext(&ops, e)

			// paint background
			paint.Fill(&ops, color.NRGBA{43, 43, 53, 255})

			// UPDATE
			topBar.Update()
			dataDisplay.Update()
			addFormPage.Update()
			list.Update()

			// LAYOUT
			if currentPage == List {
				layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(topBar.Layout),
					layout.Flexed(1, list.Layout),
					layout.Rigid(dataDisplay.Layout),
					layout.Rigid(layout.Spacer{Height: unit.Dp(25)}.Layout),
				)
			} else if currentPage == Add {
				layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(topBar.Layout),
					layout.Flexed(1, layout.Spacer{Height: unit.Dp(25)}.Layout),
					layout.Rigid(addFormPage.Layout),
					layout.Flexed(1, layout.Spacer{Height: unit.Dp(25)}.Layout),
				)
			}
			// Send context operation to event frame
			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			return e.Err
		}
	}
}
