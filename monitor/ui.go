package monitor

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func setupAppTheme() {
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorBlack
	tview.Styles.ContrastBackgroundColor = tcell.NewRGBColor(30, 41, 59)
	tview.Styles.MoreContrastBackgroundColor = ColorBackground
	tview.Styles.BorderColor = ColorBorder
	tview.Styles.TitleColor = ColorPrimary
	tview.Styles.GraphicsColor = ColorBorder
	tview.Styles.PrimaryTextColor = tcell.ColorWhite
	tview.Styles.SecondaryTextColor = ColorSecondary
	tview.Styles.TertiaryTextColor = tcell.NewRGBColor(100, 116, 139)
	tview.Styles.InverseTextColor = tcell.ColorBlack
	tview.Styles.ContrastSecondaryTextColor = tcell.NewRGBColor(203, 213, 225)
}

func createJobListTable() *tview.Table {
	table := tview.NewTable()
	table.SetSelectable(true, false)
	table.SetFixed(1, 0)
	table.SetTitle(" 🚀 Jobs ↑/↓")
	table.SetBorder(true)
	table.SetBorderPadding(0, 0, 1, 1)
	table.SetBorderColor(ColorBorder)
	table.SetTitleColor(ColorPrimary)
	table.SetBackgroundColor(ColorBackground)
	table.SetSelectedStyle(tcell.StyleDefault.
		Background(tcell.NewRGBColor(30, 58, 138)).
		Foreground(tcell.ColorWhite))
	return table
}

func createJobDetailsView() *tview.TextView {
	view := tview.NewTextView()
	view.SetDynamicColors(true)
	view.SetRegions(true)
	view.SetWordWrap(true)
	view.SetTitle(" 📋 Job Details (Enter/Esc to return) ")
	view.SetBorder(true)
	view.SetBorderPadding(0, 0, 1, 1)
	view.SetBorderColor(ColorBorder)
	view.SetTitleColor(ColorPrimary)
	view.SetBackgroundColor(ColorBackground)
	return view
}

func createStatusBar() *tview.TextView {
	bar := tview.NewTextView()
	bar.SetDynamicColors(true)
	bar.SetRegions(true)
	bar.SetBorder(false)
	bar.SetBorderPadding(0, 0, 1, 1)
	bar.SetBackgroundColor(tcell.NewRGBColor(30, 41, 59))
	return bar
}

func createKindFilterInput() *tview.InputField {
	input := tview.NewInputField()
	input.SetLabel("🔍 Filter by kind or job ID: ")
	input.SetFieldWidth(20)
	input.SetBorder(true)
	input.SetTitle(" 🏷️  Kind / ID Filter (Enter: Apply, Esc: Clear) ")
	input.SetBorderColor(ColorPrimary)
	input.SetTitleColor(ColorPrimary)
	input.SetBackgroundColor(ColorBackground)
	input.SetFieldBackgroundColor(tcell.NewRGBColor(30, 41, 59))
	return input
}

func createConfirmationModal() *tview.TextView {
	modal := tview.NewTextView()
	modal.SetDynamicColors(true)
	modal.SetRegions(true)
	modal.SetWordWrap(true)
	modal.SetTextAlign(tview.AlignCenter)
	modal.SetBorder(true)
	modal.SetTitle(" ⚠️  Confirmation ")
	modal.SetBorderColor(ColorWarning)
	modal.SetTitleColor(ColorWarning)
	modal.SetBackgroundColor(ColorBackground)
	return modal
}

func createCenteredModal(component tview.Primitive, width, height int) *tview.Flex {
	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(nil, 0, 1, false).
			AddItem(component, width, 0, true).
			AddItem(nil, 0, 1, false), height, 0, true).
		AddItem(nil, 0, 1, false)
}
