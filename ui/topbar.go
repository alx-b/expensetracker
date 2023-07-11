package ui

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"time"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/alx-b/expensetracker/domain"
)

type TopBar struct {
	prevMonthButton material.ButtonStyle
	nextMonthButton material.ButtonStyle
	listPageButton  material.ButtonStyle
	addPageButton   material.ButtonStyle
	closeButton     material.ButtonStyle
	labelMonth      material.LabelStyle
	margins         layout.Inset
	labelMarginTop  layout.Inset
	currentMonth    string
	currentPage     *Page
	monthView       *domain.MonthData
	controller      domain.API
}

// createTopBar returns TopBar struct
func createTopBar(th *material.Theme, currentPage *Page, monthData *domain.MonthData, controller domain.API) TopBar {
	currentMonth := fmt.Sprintf("%s %d", monthData.Month.String(), monthData.Year)

	prevMonthButton := material.Button(th, &widget.Clickable{}, "<")
	labelMonth := material.Label(th, unit.Sp(16), currentMonth)
	nextMonthButton := material.Button(th, &widget.Clickable{}, ">")
	listPageButton := material.Button(th, &widget.Clickable{}, "MAIN")
	addPageButton := material.Button(th, &widget.Clickable{}, "ADD")
	closeButton := material.Button(th, &widget.Clickable{}, "X")

	labelMonth.MaxLines = 1

	buttons := []*material.ButtonStyle{
		&prevMonthButton,
		&nextMonthButton,
		&listPageButton,
		&addPageButton,
		&closeButton,
	}

	for i := range buttons {
		buttons[i].Background = color.NRGBA{3, 106, 102, 255}
	}

	margins := layout.UniformInset(unit.Dp(6))

	labelMarginTop := layout.Inset{Top: unit.Dp(8)}

	return TopBar{
		prevMonthButton: prevMonthButton,
		nextMonthButton: nextMonthButton,
		listPageButton:  listPageButton,
		addPageButton:   addPageButton,
		closeButton:     closeButton,
		labelMonth:      labelMonth,
		margins:         margins,
		labelMarginTop:  labelMarginTop,
		currentMonth:    currentMonth,
		currentPage:     currentPage,
		monthView:       monthData,
		controller:      controller,
	}
}

// Update updates data based on button clicks.
func (t *TopBar) Update() {
	if t.prevMonthButton.Button.Clicked() {
		if t.monthView.Month == time.January {
			t.monthView.Month = time.December
			t.monthView.Year--
		} else {
			t.monthView.Month--
		}
		*t.monthView = t.controller.CreateMonthData(t.monthView.Year, t.monthView.Month)
	} else if t.nextMonthButton.Button.Clicked() {
		if t.monthView.Month == time.December {
			t.monthView.Month = time.January
			t.monthView.Year++
		} else {
			t.monthView.Month++
		}
		*t.monthView = t.controller.CreateMonthData(t.monthView.Year, t.monthView.Month)
	} else if t.listPageButton.Button.Clicked() {
		*t.currentPage = List
		*t.monthView = t.controller.CreateMonthData(t.monthView.Year, t.monthView.Month)
	} else if t.addPageButton.Button.Clicked() {
		*t.currentPage = Add
		*t.monthView = t.controller.CreateMonthData(t.monthView.Year, t.monthView.Month)
	} else if t.closeButton.Button.Clicked() {
		os.Exit(0)
	}
	t.currentMonth = fmt.Sprintf("%s %d", t.monthView.Month.String(), t.monthView.Year)
	t.labelMonth.Text = t.currentMonth
}

// Layout returns its layout
func (t *TopBar) Layout(gtx layout.Context) layout.Dimensions {
	r := clip.Rect{
		Min: image.Pt(0, 0),
		Max: image.Pt(gtx.Constraints.Max.X, gtx.Metric.Dp(50)),
	}

	color := color.NRGBA{3, 106, 102, 255}
	paint.FillShape(gtx.Ops, color, r.Op())

	if *t.currentPage == Add {
		return t.margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx,
				layout.Flexed(1, layout.Spacer{}.Layout),
				layout.Rigid(t.listPageButton.Layout),
				layout.Rigid(layout.Spacer{Width: unit.Dp(18)}.Layout),
				layout.Rigid(t.addPageButton.Layout),
				layout.Rigid(layout.Spacer{Width: unit.Dp(18)}.Layout),
				layout.Rigid(t.closeButton.Layout),
			)
		})
	}

	return t.margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis: layout.Horizontal,
		}.Layout(gtx,
			layout.Rigid(t.prevMonthButton.Layout),
			layout.Rigid(layout.Spacer{Width: unit.Dp(18)}.Layout),
			layout.Rigid(
				func(gtx layout.Context) layout.Dimensions {
					return t.labelMarginTop.Layout(gtx, t.labelMonth.Layout)
				},
			),
			layout.Rigid(layout.Spacer{Width: unit.Dp(18)}.Layout),
			layout.Rigid(t.nextMonthButton.Layout),
			layout.Flexed(1, layout.Spacer{}.Layout),
			layout.Rigid(t.listPageButton.Layout),
			layout.Rigid(layout.Spacer{Width: unit.Dp(18)}.Layout),
			layout.Rigid(t.addPageButton.Layout),
			layout.Rigid(layout.Spacer{Width: unit.Dp(18)}.Layout),
			layout.Rigid(t.closeButton.Layout),
		)
	})
}
