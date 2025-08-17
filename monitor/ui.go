package monitor

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func setupAppTheme() {
	tview.Styles.PrimitiveBackgroundColor = ColorPrimativeBackground
	tview.Styles.ContrastBackgroundColor = ColorContrastBackground
	tview.Styles.MoreContrastBackgroundColor = ColorMoreContrastBackground
	tview.Styles.BorderColor = ColorBorder
	tview.Styles.TitleColor = ColorTitle
	tview.Styles.GraphicsColor = ColorBorder
	tview.Styles.PrimaryTextColor = ColorPrimary
	tview.Styles.SecondaryTextColor = ColorSecondary
	tview.Styles.TertiaryTextColor = ColorTertiary
	tview.Styles.InverseTextColor = ColorPrimativeBackground
	tview.Styles.ContrastSecondaryTextColor = ColorContrastSecondary
}

func createJobListTable() *tview.Table {
	table := tview.NewTable()
	table.SetSelectable(true, false)
	table.SetFixed(1, 0)
	table.SetTitle(" 🚀 Jobs ↑/↓")
	table.SetBorder(true)
	table.SetBorderPadding(0, 0, 1, 1)
	table.SetBorderColor(ColorBorder)
	table.SetTitleColor(ColorTitle)
	table.SetBackgroundColor(ColorContrastBackground)
	table.SetSelectedStyle(tcell.StyleDefault.
		Background(ColorSelectedBg).
		Foreground(ColorSelectedFg))
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
	view.SetTitleColor(ColorTitle)
	view.SetBackgroundColor(ColorContrastBackground)
	return view
}

func createStatusBar() *tview.TextView {
	bar := tview.NewTextView()
	bar.SetDynamicColors(true)
	bar.SetRegions(true)
	bar.SetBorder(false)
	bar.SetBorderPadding(0, 0, 1, 1)
	bar.SetBackgroundColor(ColorContrastBackground)
	return bar
}

func createKindFilterInput() *tview.InputField {
	input := tview.NewInputField()
	input.SetLabel("🔍 Filter by kind or job ID: ")
	input.SetFieldWidth(20)
	input.SetBorder(true)
	input.SetTitle(" 🏷️  Kind / ID Filter (Enter: Apply, Esc: Clear) ")
	input.SetLabelColor(ColorTitle)
	input.SetBorderColor(ColorTitle)
	input.SetTitleColor(ColorTitle)
	input.SetBackgroundColor(ColorContrastBackground)
	input.SetFieldBackgroundColor(ColorContrastBackground)
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
	modal.SetBackgroundColor(ColorContrastBackground)
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

func createQueueListTable() *tview.Table {
	table := tview.NewTable()
	table.SetSelectable(true, false)
	table.SetFixed(1, 0)
	table.SetTitle(" 🔀 Queues ↑/↓")
	table.SetBorder(true)
	table.SetBorderPadding(0, 0, 1, 1)
	table.SetBorderColor(ColorBorder)
	table.SetTitleColor(ColorTitle)
	table.SetBackgroundColor(ColorContrastBackground)
	table.SetSelectedStyle(tcell.StyleDefault.
		Background(ColorSelectedBg).
		Foreground(ColorSelectedFg))
	return table
}
