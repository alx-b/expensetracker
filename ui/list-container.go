package ui

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/alx-b/expensetracker/domain"
)

type ListContainer struct {
	list          material.ListStyle
	theme         *material.Theme
	nameLabel     material.LabelStyle
	dateLabel     material.LabelStyle
	categoryLabel material.LabelStyle
	amountLabel   material.LabelStyle

	deleteButtons []material.ButtonStyle
	controller    domain.API
	monthView     *domain.MonthData
}

// TODO handle this mess better.
// Update updates the list and monthView.
func (c *ListContainer) Update() {
	if len(c.deleteButtons) != len(c.monthView.Expenses) {
		delButtons := []material.ButtonStyle{}
		for _ = range c.monthView.Expenses {
			delButton := material.Button(c.theme, &widget.Clickable{}, "x")
			delButton.Background = color.NRGBA{113, 53, 53, 255}
			delButtons = append(delButtons, delButton)
		}
		c.deleteButtons = delButtons
		return
	}

	// Don't try to range through buttons if there is none.
	if len(c.deleteButtons) < 1 {
		return
	}

	for i := range c.deleteButtons {
		if c.deleteButtons[i].Button.Clicked() {
			c.controller.RemoveExpense((c.monthView.Expenses)[i].Id)

			*c.monthView = c.controller.CreateMonthData(c.monthView.Year, c.monthView.Month)
			delButtons := []material.ButtonStyle{}
			for _ = range c.monthView.Expenses {
				delButton := material.Button(c.theme, &widget.Clickable{}, "x")
				delButton.Background = color.NRGBA{113, 53, 53, 255}
				delButtons = append(delButtons, delButton)
			}
			c.deleteButtons = delButtons
			break
		}
	}
}

// Layout returns its layout.
func (c *ListContainer) Layout(gtx layout.Context) layout.Dimensions {
	margins := layout.Inset{
		Top:    unit.Dp(25),
		Bottom: unit.Dp(25),
		Right:  unit.Dp(25),
		Left:   unit.Dp(25),
	}

	insideBorderMargins := layout.Inset{
		Left:   unit.Dp(10),
		Top:    unit.Dp(10),
		Bottom: unit.Dp(10),
		// There is already some right margin, from the doc:
		// l.AnchorStrategy == Occupy
		// Increase the width to account for the space occupied by the scrollbar.
	}

	borders := widget.Border{
		Color:        color.NRGBA{R: 53, G: 53, B: 63, A: 255},
		CornerRadius: unit.Dp(3),
		Width:        unit.Dp(2),
	}

	bottomMargin := layout.Inset{Bottom: unit.Dp(10)}
	topBottomMargins := layout.Inset{Bottom: unit.Dp(6), Top: unit.Dp(12)}

	r := clip.Rect{
		Min: image.Pt(int(margins.Left)+int(borders.Width), int(margins.Top)+int(borders.Width)),
		Max: image.Pt(gtx.Constraints.Max.X-int(margins.Left)-int(borders.Width), gtx.Constraints.Max.Y-int(margins.Top)-int(borders.Width)),
	}
	blueColor := color.NRGBA{53, 53, 63, 255}
	paint.FillShape(gtx.Ops, blueColor, r.Op())

	return margins.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			return borders.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return insideBorderMargins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return c.list.Layout(gtx, len(c.monthView.Expenses), func(gtx layout.Context, i int) layout.Dimensions {
						c.nameLabel.Text = (c.monthView.Expenses)[i].Name
						c.dateLabel.Text = (c.monthView.Expenses)[i].Date
						c.categoryLabel.Text = (c.monthView.Expenses)[i].Category
						c.amountLabel.Text = fmt.Sprintf("%.2f", (c.monthView.Expenses)[i].Amount)
						return bottomMargin.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							r2 := clip.Rect{
								Min: image.Pt(0, 0),
								Max: image.Pt(gtx.Constraints.Max.X, int(gtx.Dp(24)+gtx.Sp(24))),
							}
							palerBlueColor := color.NRGBA{73, 73, 83, 255}
							paint.FillShape(gtx.Ops, palerBlueColor, r2.Op())
							return layout.Flex{
								Axis: layout.Horizontal,
							}.Layout(gtx,
								layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									return topBottomMargins.Layout(gtx, c.nameLabel.Layout)
								}),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									return topBottomMargins.Layout(gtx, c.dateLabel.Layout)
								}),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									return topBottomMargins.Layout(gtx, c.categoryLabel.Layout)
								}),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									return topBottomMargins.Layout(gtx, c.amountLabel.Layout)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									// Don't show buttons if there is a difference between
									// length of expenses vs length of buttons
									if len(c.monthView.Expenses) != len(c.deleteButtons) {
										return layout.Dimensions{}
									}
									return topBottomMargins.Layout(gtx, c.deleteButtons[i].Layout)
								}),
								layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
							)
						},
						)
					})
				})
			})
		})
}

// createListContainer returns ListContainer struct.
func createListContainer(th *material.Theme, monthData *domain.MonthData, controller domain.API) ListContainer {
	var list widget.List
	list.Axis = layout.Vertical
	listWithStyle := material.List(th, &list)

	nameLabel := material.Label(th, unit.Sp(16), "")
	dateLabel := material.Label(th, unit.Sp(16), "")
	categoryLabel := material.Label(th, unit.Sp(16), "")
	amountLabel := material.Label(th, unit.Sp(16), fmt.Sprintf("%.2f", 0.00))
	amountLabel.Alignment = text.End

	labels := []*material.LabelStyle{
		&nameLabel,
		&dateLabel,
		&categoryLabel,
		&amountLabel,
	}

	for i := range labels {
		labels[i].MaxLines = 1
	}

	delButtons := []material.ButtonStyle{}

	for _ = range monthData.Expenses {
		delButton := material.Button(th, &widget.Clickable{}, "x")
		delButton.Background = color.NRGBA{113, 53, 53, 255}
		delButtons = append(delButtons, delButton)
	}

	return ListContainer{
		list:          listWithStyle,
		theme:         th,
		nameLabel:     nameLabel,
		dateLabel:     dateLabel,
		categoryLabel: categoryLabel,
		amountLabel:   amountLabel,
		monthView:     monthData,
		deleteButtons: delButtons,
		controller:    controller,
	}
}
