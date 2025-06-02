package monitor

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

func (m *MonitorApp) setupKeyBindings() {
	m.setupJobListKeyBindings()
	m.setupKindFilterKeyBindings()
	m.setupConfirmationKeyBindings()
	m.setupJobDetailsKeyBindings()
	m.setupQueueKeyBindings()
}

func (m *MonitorApp) setupJobListKeyBindings() {
	m.ui.jobList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			row, _ := m.ui.jobList.GetSelection()
			if row > 0 {
				jobID := m.ui.jobList.GetCell(row, 0).Text
				m.showJobDetails(jobID)
				m.ui.pages.SwitchToPage(PageDetails)
				m.setDetailsModeStatus()
			}
			return nil
		case tcell.KeyCtrlQ:
			m.showQueues()
			return nil
		case tcell.KeyRune:
			if event.Rune() == 'q' {
				m.ui.app.Stop()
				return nil
			}
			if event.Rune() == '/' {
				m.openKindFilter()
				return nil
			}
			if event.Rune() == 'r' {
				m.handleJobRetry()
				return nil
			}
			if event.Rune() == 'c' {
				m.handleJobCancel()
				return nil
			}
			if event.Rune() == 'n' {
				m.nextPage()
				return nil
			}
			if event.Rune() == 'p' {
				m.previousPage()
				return nil
			}
			if event.Rune() >= '0' && event.Rune() <= '7' {
				stateNum := int(event.Rune() - '0')
				m.setStateFilter(stateNum)
				return nil
			}
		}
		return event
	})
}

func (m *MonitorApp) setupKindFilterKeyBindings() {
	m.ui.kindFilterInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			text := strings.TrimSpace(m.ui.kindFilterInput.GetText())
			if text == "" {
				m.filter.SetKindFilter(nil)
				m.pagination.Reset()
				m.scrollToBeginning = true
				m.updateFilterStatusBar()
				m.closeKindFilter()
				return nil
			}

			// Check if the input is a valid integer (job ID)
			if _, err := strconv.ParseInt(text, 10, 64); err == nil {
				// Input is a valid integer, treat as job ID and show details
				m.closeKindFilter()
				m.showJobDetails(text)
				m.ui.pages.SwitchToPage(PageDetails)
				m.setDetailsModeStatus()
				return nil
			}

			// Input is not an integer, treat as kind filter
			m.filter.SetKindFilter([]string{text})
			m.pagination.Reset()
			m.scrollToBeginning = true
			m.updateFilterStatusBar()
			m.closeKindFilter()
			return nil
		case tcell.KeyEsc:
			m.filter.SetKindFilter(nil)
			m.pagination.Reset()
			m.scrollToBeginning = true
			m.updateFilterStatusBar()
			m.closeKindFilter()
			return nil
		}
		return event
	})
}

func (m *MonitorApp) setupConfirmationKeyBindings() {
	m.ui.confirmationModal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'Y', 'y':
				m.modalState.ExecuteYes()
				m.closeConfirmationModal()
				return nil
			case 'N', 'n':
				m.modalState.ExecuteNo()
				m.closeConfirmationModal()
				return nil
			case 'q':
				m.ui.app.Stop()
				return nil
			}
		case tcell.KeyEsc:
			m.modalState.ExecuteNo()
			m.closeConfirmationModal()
			return nil
		}
		return event
	})
}

func (m *MonitorApp) setupJobDetailsKeyBindings() {
	m.ui.jobDetails.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter, tcell.KeyEsc:
			m.currentJobID = ""
			m.ui.pages.SwitchToPage(PageList)
			m.ui.app.SetFocus(m.ui.jobList)
			m.setListModeStatus()
			return nil
		case tcell.KeyRune:
			if event.Rune() == 'q' {
				m.ui.app.Stop()
				return nil
			}
			if event.Rune() == 'r' {
				m.handleJobRetryInDetails()
				return nil
			}
			if event.Rune() == 'c' {
				m.handleJobCancelInDetails()
				return nil
			}
		}
		return event
	})
}

func (m *MonitorApp) setupQueueKeyBindings() {
	m.ui.queueList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter, tcell.KeyEsc:
			m.ui.pages.SwitchToPage(PageList)
			m.ui.app.SetFocus(m.ui.jobList)
			m.setListModeStatus()
			return nil
		case tcell.KeyRune:
			if event.Rune() == 'q' {
				m.ui.app.Stop()
				return nil
			}
			if event.Rune() == 'p' {
				m.handleQueuePause()
				return nil
			}
			if event.Rune() == 'r' {
				m.handleQueueResume()
				return nil
			}
		}
		return event
	})
}
