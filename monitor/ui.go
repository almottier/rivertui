package monitor

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func setupAppTheme() {
	tview.Styles.PrimitiveBackgroundColor = primitiveBackground
	tview.Styles.ContrastBackgroundColor = contrastBackground
	tview.Styles.MoreContrastBackgroundColor = moreContrastBackground
	tview.Styles.BorderColor = borderColor
	tview.Styles.TitleColor = titleColor
	tview.Styles.GraphicsColor = graphicsColor
	tview.Styles.PrimaryTextColor = primaryTextColor
	tview.Styles.SecondaryTextColor = secondaryTextColor
	tview.Styles.TertiaryTextColor = tertiaryTextColor
	tview.Styles.InverseTextColor = inverseTextColor
	tview.Styles.ContrastSecondaryTextColor = contrastSecondaryTextColor
}

func createJobListTable() *tview.Table {
	table := tview.NewTable()
	table.SetSelectable(true, false)
	table.SetFixed(1, 0)
	table.SetTitle(" 🚀 Jobs ↑/↓")
	table.SetBorder(true)
	table.SetBorderPadding(0, 0, 1, 1)
	table.SetBorderColor(jobListBorder)
	table.SetTitleColor(jobListTitle)
	table.SetBackgroundColor(jobListBg)
	table.SetSelectedStyle(tcell.StyleDefault.
		Background(jobListSelectedBg).
		Foreground(jobListSelectedFg))
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
	view.SetBorderColor(jobDetailsBorder)
	view.SetTitleColor(jobDetailsTitle)
	view.SetBackgroundColor(jobDetailsBg)
	return view
}

func createStatusBar() *tview.TextView {
	bar := tview.NewTextView()
	bar.SetDynamicColors(true)
	bar.SetRegions(true)
	bar.SetBorder(false)
	bar.SetBorderPadding(0, 0, 1, 1)
	bar.SetBackgroundColor(statusBarBg)
	return bar
}

func createKindFilterInput() *tview.InputField {
	input := tview.NewInputField()
	input.SetLabel("🔍 Filter by kind or job ID: ")
	input.SetFieldWidth(20)
	input.SetBorder(true)
	input.SetTitle(" 🏷️  Kind / ID Filter (Enter: Apply, Esc: Clear) ")
	input.SetLabelColor(kindFilterLabel)
	input.SetBorderColor(kindFilterBorder)
	input.SetTitleColor(kindFilterTitle)
	input.SetBackgroundColor(kindFilterBg)
	input.SetFieldBackgroundColor(kindFilterFieldBg)
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
	modal.SetBorderColor(confirmationModalBorder)
	modal.SetTitleColor(confirmationModalTitle)
	modal.SetBackgroundColor(confirmationModalBg)
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
	table.SetBorderColor(queueListBorder)
	table.SetTitleColor(queueListTitle)
	table.SetBackgroundColor(queueListBg)
	table.SetSelectedStyle(tcell.StyleDefault.
		Background(queueListSelectedBg).
		Foreground(queueListSelectedFg))
	return table
}
