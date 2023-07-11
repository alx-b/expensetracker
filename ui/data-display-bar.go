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

type DataDisplay struct {
	checkBox      material.CheckBoxStyle
	cancelBudget  material.ButtonStyle
	submitBudget  material.ButtonStyle
	inputBudget   material.EditorStyle
	editBudget    material.ButtonStyle
	budgetLabel   material.LabelStyle
	totalLabel    material.LabelStyle
	leftoverLabel material.LabelStyle
	state         State
	controller    domain.API
	monthData     *domain.MonthData
}

type State int

const (
	Editing State = iota
	Visual
)

// Update updates data based on button clicks
func (d *DataDisplay) Update() {
	if d.editBudget.Button.Clicked() {
		d.state = Editing
	}

	if d.cancelBudget.Button.Clicked() {
		d.inputBudget.Editor.SetText("")
		d.state = Visual
	}

	if d.submitBudget.Button.Clicked() {
		money := d.inputBudget.Editor.Text()

		_, err := strconv.ParseFloat(money, 64)
		if err != nil {
			d.inputBudget.Editor.SetText("")
			return
		}

		date := fmt.Sprintf("%d-%02d", d.monthData.Year, d.monthData.Month)
		d.controller.InsertBudgetMonth(money, date)

		if d.checkBox.CheckBox.Value == true {
			d.controller.UpdateDefaultBudget(money)
			d.checkBox.CheckBox.Value = false
		}

		d.inputBudget.Editor.SetText("")
		*d.monthData = d.controller.CreateMonthData(d.monthData.Year, d.monthData.Month)
		d.state = Visual
	}

	d.budgetLabel.Text = fmt.Sprintf("Budget: %.2f", d.monthData.Budget)
	d.totalLabel.Text = fmt.Sprintf("Total: %.2f", d.monthData.TotalSpendings)
	d.leftoverLabel.Text = fmt.Sprintf("Leftover: %.2f", d.monthData.MoneyLeft)
}

// Layout returns its layout.
func (d *DataDisplay) Layout(gtx layout.Context) layout.Dimensions {
	if d.state == Editing {
		return layout.Flex{
			Axis:    layout.Horizontal,
			Spacing: layout.SpaceEvenly,
		}.Layout(gtx,
			layout.Rigid(
				func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis: layout.Horizontal,
					}.Layout(gtx,
						layout.Rigid(d.cancelBudget.Layout),
						layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
						layout.Rigid(d.submitBudget.Layout),
						layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
						layout.Rigid(d.checkBox.Layout),
						layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
						layout.Rigid(
							func(gtx layout.Context) layout.Dimensions {
								r := clip.Rect{
									Min: image.Pt(gtx.Dp(0), gtx.Dp(0)),
									Max: image.Pt(d.inputBudget.Layout(gtx).Size.X, d.inputBudget.Layout(gtx).Size.Y),
								}
								paint.FillShape(gtx.Ops, color.NRGBA{53, 53, 63, 255}, r.Op())
								return d.inputBudget.Layout(gtx)
							},
						),
					)
				},
			),
			layout.Rigid(
				d.totalLabel.Layout,
			),
			layout.Rigid(
				d.leftoverLabel.Layout,
			),
		)

	}

	return layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceEvenly,
	}.Layout(gtx,
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Horizontal,
				}.Layout(gtx,
					layout.Rigid(d.editBudget.Layout),
					layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
					layout.Rigid(d.budgetLabel.Layout),
				)
			},
		),
		layout.Rigid(
			d.totalLabel.Layout,
		),
		layout.Rigid(
			d.leftoverLabel.Layout,
		),
	)
}

// createDataDisplay returns DataDisplay struct.
func createDataDisplay(th *material.Theme, controller domain.API, monthData *domain.MonthData) DataDisplay {
	inputBudget := material.Editor(th, &widget.Editor{}, "0.00")
	checkBox := material.CheckBox(th, &widget.Bool{}, "Default")
	cancelBudget := material.Button(th, &widget.Clickable{}, "Cancel")
	submitBudget := material.Button(th, &widget.Clickable{}, "Submit")
	editBudget := material.Button(th, &widget.Clickable{}, "Edit")
	budgetLabel := material.Label(th, unit.Sp(16), fmt.Sprintf("Budget: %.2f", 0.00))
	totalLabel := material.Label(th, unit.Sp(16), fmt.Sprintf("Total: %.2f", 0.00))
	leftoverLabel := material.Label(th, unit.Sp(16), fmt.Sprintf("Leftover: %.2f", 0.00))
	state := Visual

	submitBudget.Background = color.NRGBA{53, 53, 113, 255}
	cancelBudget.Background = color.NRGBA{113, 53, 53, 255}
	editBudget.Background = color.NRGBA{53, 53, 113, 255}

	inputBudget.Editor.Alignment = text.Middle
	inputBudget.Editor.SingleLine = true
	inputBudget.Color = color.NRGBA{235, 235, 235, 255}
	inputBudget.HintColor = color.NRGBA{255, 255, 255, 40}

	budgetLabel.MaxLines = 1
	totalLabel.MaxLines = 1
	leftoverLabel.MaxLines = 1

	return DataDisplay{
		inputBudget:   inputBudget,
		checkBox:      checkBox,
		cancelBudget:  cancelBudget,
		submitBudget:  submitBudget,
		editBudget:    editBudget,
		budgetLabel:   budgetLabel,
		totalLabel:    totalLabel,
		leftoverLabel: leftoverLabel,
		state:         state,
		controller:    controller,
		monthData:     monthData,
	}
}
