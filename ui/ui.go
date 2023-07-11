package ui

import (
	"image/color"

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

func Run(w *app.Window, controller domain.API) error {
	// Defines the material design style
	th := material.NewTheme(gofont.Collection())
	th.Bg = color.NRGBA{100, 100, 100, 255}
	th.Fg = color.NRGBA{200, 200, 200, 255}
	// Operations from the UI
	var ops op.Ops

	for {
		// listen for events in the window.
		e := <-w.Events()
		switch e := e.(type) {
		case system.FrameEvent:
			// CONTEXT = REFERENCE TO OPERATION + EVENT
			gtx := layout.NewContext(&ops, e)
			// background color
			paint.Fill(&ops, color.NRGBA{43, 43, 53, 255})

			// UPDATE

			// LAYOUT FLEX
			layout.Flex{
				Axis: layout.Vertical,
				//Spacing: layout.SpaceStart,
			}.Layout(gtx,
				layout.Rigid(layout.Spacer{Height: unit.Dp(25)}.Layout),
			)
			// SEND CONTEXT OPERATION TO EVENT FRAME
			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			return e.Err
		}
	}
}
