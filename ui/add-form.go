package ui

import (
	"fmt"
	"image"
	"image/color"
	"strconv"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/alx-b/expensetracker/domain"
)

type FormPage struct {
	nameInput     material.EditorStyle
	dateInput     material.EditorStyle
	categoryInput material.EditorStyle
	amountInput   material.EditorStyle
	submitButton  material.ButtonStyle
	cancelButton  material.ButtonStyle
	controller    domain.API
	allInputs     []*material.EditorStyle
	refreshData   bool
}

// createFormPage returns FormPage struct.
func createFormPage(th *material.Theme, controller domain.API) FormPage {
	nameInput := material.Editor(th, &widget.Editor{}, "name")
	dateInput := material.Editor(th, &widget.Editor{}, "date (YYYY-MM-DD or YYYY-MM)")
	categoryInput := material.Editor(th, &widget.Editor{}, "category")
	amountInput := material.Editor(th, &widget.Editor{}, "amount")

	inputs := []*material.EditorStyle{
		&nameInput,
		&dateInput,
		&categoryInput,
		&amountInput,
	}

	for i := range inputs {
		inputs[i].Editor.Alignment = text.Middle
		inputs[i].Editor.SingleLine = true
		inputs[i].Color = color.NRGBA{235, 235, 235, 255}
		inputs[i].HintColor = color.NRGBA{255, 255, 255, 40}
	}

	submitButton := material.Button(th, &widget.Clickable{}, "Submit")
	submitButton.Background = color.NRGBA{53, 53, 113, 255}
	cancelButton := material.Button(th, &widget.Clickable{}, "Cancel")
	cancelButton.Background = color.NRGBA{113, 53, 53, 255}

	return FormPage{
		nameInput:     nameInput,
		dateInput:     dateInput,
		categoryInput: categoryInput,
		amountInput:   amountInput,
		submitButton:  submitButton,
		cancelButton:  cancelButton,
		controller:    controller,
		allInputs:     inputs,
	}
}

// clearInputs clear its inputs.
func (fp *FormPage) clearInputs() {
	for i := range fp.allInputs {
		fp.allInputs[i].Editor.SetText("")
	}
}

// Update updates data based on button clicks.
func (fp *FormPage) Update() {
	if fp.cancelButton.Button.Clicked() {
		fp.clearInputs()
	}

	if fp.submitButton.Button.Clicked() {
		amount, err := strconv.ParseFloat(fp.amountInput.Editor.Text(), 64)
		if err != nil {
			fmt.Println(err)
		}
		fp.controller.AddExpense(domain.Expense{
			Name:     fp.nameInput.Editor.Text(),
			Date:     fp.dateInput.Editor.Text(),
			Category: fp.categoryInput.Editor.Text(),
			Amount:   amount,
		})
		fp.clearInputs()
	}
}

// Layout returns its layout.
func (fp *FormPage) Layout(gtx layout.Context) layout.Dimensions {
	margins := layout.Inset{
		Top:    unit.Dp(25),
		Bottom: unit.Dp(25),
		Right:  unit.Dp(25),
		Left:   unit.Dp(25),
	}

	marginTop := layout.Inset{
		Top: unit.Dp(25),
	}

	insideBorderMargins := layout.UniformInset(unit.Dp(10))

	borders := widget.Border{
		Color:        color.NRGBA{R: 53, G: 53, B: 63, A: 255},
		CornerRadius: unit.Dp(3),
		Width:        unit.Dp(2),
	}
	color := color.NRGBA{53, 53, 63, 255}

	return margins.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return marginTop.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								r := clip.Rect{
									Min: image.Pt(gtx.Dp(borders.Width), gtx.Dp(borders.Width)),
									Max: image.Pt(fp.nameInput.Layout(gtx).Size.X, fp.nameInput.Layout(gtx).Size.Y+gtx.Dp(insideBorderMargins.Top*2)-gtx.Dp(borders.Width)),
								}
								paint.FillShape(gtx.Ops, color, r.Op())
								return borders.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return insideBorderMargins.Layout(gtx, fp.nameInput.Layout)
								})
							},
						)
					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return marginTop.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								r := clip.Rect{
									Min: image.Pt(gtx.Dp(borders.Width), gtx.Dp(borders.Width)),
									Max: image.Pt(fp.dateInput.Layout(gtx).Size.X, fp.dateInput.Layout(gtx).Size.Y+gtx.Dp(insideBorderMargins.Top*2)-gtx.Dp(borders.Width)),
								}
								paint.FillShape(gtx.Ops, color, r.Op())
								return borders.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return insideBorderMargins.Layout(gtx, fp.dateInput.Layout)
								})
							},
						)
					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return marginTop.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								r := clip.Rect{
									Min: image.Pt(gtx.Dp(borders.Width), gtx.Dp(borders.Width)),
									Max: image.Pt(fp.categoryInput.Layout(gtx).Size.X, fp.categoryInput.Layout(gtx).Size.Y+gtx.Dp(insideBorderMargins.Top*2)-gtx.Dp(borders.Width)),
								}
								paint.FillShape(gtx.Ops, color, r.Op())
								return borders.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return insideBorderMargins.Layout(gtx, fp.categoryInput.Layout)
								})
							},
						)
					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return marginTop.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								r := clip.Rect{
									Min: image.Pt(gtx.Dp(borders.Width), gtx.Dp(borders.Width)),
									Max: image.Pt(fp.amountInput.Layout(gtx).Size.X, fp.amountInput.Layout(gtx).Size.Y+gtx.Dp(insideBorderMargins.Top*2)-gtx.Dp(borders.Width)),
								}
								paint.FillShape(gtx.Ops, color, r.Op())
								return borders.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return insideBorderMargins.Layout(gtx, fp.amountInput.Layout)
								})
							},
						)
					},
				),
				layout.Rigid(
					layout.Spacer{Height: unit.Dp(80)}.Layout,
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return marginTop.Layout(gtx,
							fp.submitButton.Layout,
						)
					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return marginTop.Layout(gtx,
							fp.cancelButton.Layout,
						)
					},
				),
			)
		},
	)
}
