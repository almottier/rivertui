package monitor

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
)

func (m *MonitorApp) setupUI() {
	// Create layouts
	listFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(m.ui.jobList, 0, 1, true).
		AddItem(m.ui.filterStatusBar, 1, 0, false).
		AddItem(m.ui.statusBar, 1, 0, false)

	detailsFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(m.ui.jobDetails, 0, 1, true).
		AddItem(m.ui.statusBar, 1, 0, false)

	kindFilterModal := createCenteredModal(m.ui.kindFilterInput, 60, 3)
	confirmationModalLayout := createCenteredModal(m.ui.confirmationModal, 60, 8)

	// Add pages
	m.ui.pages.AddPage(PageList, listFlex, true, true)
	m.ui.pages.AddPage(PageDetails, detailsFlex, true, false)
	m.ui.pages.AddPage(PageKindFilter, kindFilterModal, true, false)
	m.ui.pages.AddPage(PageConfirmation, confirmationModalLayout, true, false)

	// Initialize filter status bar and help text
	m.updateFilterStatusBar()

	// Check if we should start with job details
	if m.initialJobID > 0 {
		m.showJobDetails(fmt.Sprintf("%d", m.initialJobID))
		m.ui.pages.SwitchToPage(PageDetails)
		m.setDetailsModeStatus()
	} else {
		m.setListModeStatus()
	}
}

// StartRefreshLoop begins the background refresh loop
func (m *MonitorApp) StartRefreshLoop() {
	go func() {
		for {
			m.ui.app.QueueUpdateDraw(func() {
				if err := m.updateJobList(); err != nil {
					m.ui.statusBar.SetText(fmt.Sprintf("Error: %v", err))
				}

				// If we're on the details page and have a current job ID, refresh the details
				if m.ui.pages.HasPage(PageDetails) {
					if pageName, _ := m.ui.pages.GetFrontPage(); pageName == PageDetails && m.currentJobID != "" {
						m.showJobDetails(m.currentJobID)
					}
				}
			})
			time.Sleep(1 * time.Second)
		}
	}()
}

// Run starts the monitor application
func (m *MonitorApp) Run() error {
	m.StartRefreshLoop()
	return m.ui.app.SetRoot(m.ui.pages, true).EnableMouse(false).Run()
}
