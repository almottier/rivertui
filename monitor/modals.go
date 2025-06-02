package monitor

import (
	"fmt"
)

// openKindFilter opens the kind filter input modal
func (m *MonitorApp) openKindFilter() {
	// Clear previous input and show current filter if any
	if len(m.filter.kindFilter) > 0 {
		m.ui.kindFilterInput.SetText(m.filter.kindFilter[0])
	} else {
		m.ui.kindFilterInput.SetText("")
	}
	m.ui.pages.ShowPage(PageKindFilter)
	m.ui.app.SetFocus(m.ui.kindFilterInput)
}

// closeKindFilter closes the kind filter input modal and returns to the list
func (m *MonitorApp) closeKindFilter() {
	m.ui.pages.HidePage(PageKindFilter)
	m.ui.pages.SwitchToPage(PageList)
	m.ui.app.SetFocus(m.ui.jobList)
}

// showConfirmationModal displays a confirmation modal with the given title, message, and callbacks
func (m *MonitorApp) showConfirmationModal(title, message string, onYes, onNo func()) {
	// Track the current page before showing the modal
	m.lastActivePage, _ = m.ui.pages.GetFrontPage()

	// Set the callbacks
	m.modalState.Set(onYes, onNo)

	// Update the modal title and content
	m.ui.confirmationModal.SetTitle(fmt.Sprintf(" ⚠️  %s ", title))
	m.ui.confirmationModal.SetText(message)

	// Show the modal and focus on it
	m.ui.pages.ShowPage(PageConfirmation)
	m.ui.app.SetFocus(m.ui.confirmationModal)
}

// closeConfirmationModal closes the confirmation modal and returns to the previous page
func (m *MonitorApp) closeConfirmationModal() {
	// Clear the callbacks
	m.modalState.Clear()

	// Hide the modal and return to the appropriate page
	m.ui.pages.HidePage(PageConfirmation)

	// Return to the page that was active before the confirmation modal
	switch m.lastActivePage {
	case PageQueues:
		m.ui.pages.SwitchToPage(PageQueues)
		m.ui.app.SetFocus(m.ui.queueList)
	case PageDetails:
		m.ui.pages.SwitchToPage(PageDetails)
		m.ui.app.SetFocus(m.ui.jobDetails)
	default:
		m.ui.pages.SwitchToPage(PageList)
		m.ui.app.SetFocus(m.ui.jobList)
	}
}
