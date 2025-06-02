package monitor

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/rivo/tview"
)

// updateQueueList refreshes the queue list table
func (m *MonitorApp) updateQueueList() error {
	ctx := context.Background()

	// Fetch queues using River's QueueList API
	result, err := m.client.RiverClient.QueueList(ctx, river.NewQueueListParams().First(100))
	if err != nil {
		return fmt.Errorf("failed to list queues: %w", err)
	}

	if len(result.Queues) == 0 {
		m.ui.statusBar.SetText("No queues found")
		return nil
	}

	m.ui.queueList.Clear()
	m.setQueueTableHeaders()

	// Add queues to table
	for i, queue := range result.Queues {
		m.addQueueToTable(i+1, queue)
	}

	return nil
}

func (m *MonitorApp) setQueueTableHeaders() {
	headers := []string{"NAME", "STATE", "PAUSED", "CREATED", "UPDATED"}
	for i, header := range headers {
		m.ui.queueList.SetCell(0, i,
			tview.NewTableCell(header).
				SetTextColor(ColorPrimary).
				SetAlign(tview.AlignLeft).
				SetSelectable(false).
				SetExpansion(1))
	}
}

func (m *MonitorApp) addQueueToTable(row int, queue *rivertype.Queue) {
	// Queue name
	m.ui.queueList.SetCell(row, 0, tview.NewTableCell(queue.Name).SetTextColor(tcell.ColorWhite))

	// Queue state (paused or active)
	var stateCell *tview.TableCell
	if queue.PausedAt != nil {
		stateCell = tview.NewTableCell("PAUSED").SetTextColor(ColorWarning)
	} else {
		stateCell = tview.NewTableCell("ACTIVE").SetTextColor(ColorSuccess)
	}
	m.ui.queueList.SetCell(row, 1, stateCell)

	// Paused at timestamp
	if queue.PausedAt != nil {
		m.ui.queueList.SetCell(row, 2, tview.NewTableCell(formatTimeAgo(*queue.PausedAt)).SetTextColor(ColorSecondary))
	} else {
		m.ui.queueList.SetCell(row, 2, tview.NewTableCell("").SetTextColor(ColorSecondary).SetBackgroundColor(ColorBackground))
	}

	// Created at timestamp
	m.ui.queueList.SetCell(row, 3, tview.NewTableCell(formatTimeAgo(queue.CreatedAt)).SetTextColor(ColorSecondary))

	// Updated at timestamp
	m.ui.queueList.SetCell(row, 4, tview.NewTableCell(formatTimeAgo(queue.UpdatedAt)).SetTextColor(ColorSecondary))
}

// showQueues switches to the queue view
func (m *MonitorApp) showQueues() {
	m.ui.pages.SwitchToPage(PageQueues)
	m.ui.app.SetFocus(m.ui.queueList)
	m.setQueueModeStatus()

	// Update queue list when switching to queue view
	if err := m.updateQueueList(); err != nil {
		m.ui.statusBar.SetText(fmt.Sprintf("Error: %v", err))
	}
}

func (m *MonitorApp) setQueueModeStatus() {
	m.ui.statusBar.SetText("[#60A5FA]Mode:[white] Queues | Esc: Back to jobs | p: Pause queue | r: Resume queue | q: Quit")
}

// handleQueuePause pauses the selected queue
func (m *MonitorApp) handleQueuePause() {
	row, _ := m.ui.queueList.GetSelection()
	if row > 0 {
		queueName := m.ui.queueList.GetCell(row, 0).Text
		m.showQueuePauseConfirmation(queueName)
	}
}

// handleQueueResume resumes the selected queue
func (m *MonitorApp) handleQueueResume() {
	row, _ := m.ui.queueList.GetSelection()
	if row > 0 {
		queueName := m.ui.queueList.GetCell(row, 0).Text
		m.showQueueResumeConfirmation(queueName)
	}
}

func (m *MonitorApp) showQueuePauseConfirmation(queueName string) {
	m.showConfirmationModal(
		"Pause Queue",
		fmt.Sprintf("Are you sure you want to pause queue '%s'?\n\n[#60A5FA]Y[white]: Yes, pause the queue\n[#60A5FA]N[white]: No, cancel", queueName),
		func() { m.pauseQueue(queueName) },
		func() {},
	)
}

func (m *MonitorApp) showQueueResumeConfirmation(queueName string) {
	m.showConfirmationModal(
		"Resume Queue",
		fmt.Sprintf("Are you sure you want to resume queue '%s'?\n\n[#60A5FA]Y[white]: Yes, resume the queue\n[#60A5FA]N[white]: No, cancel", queueName),
		func() { m.resumeQueue(queueName) },
		func() {},
	)
}

func (m *MonitorApp) pauseQueue(queueName string) {
	ctx := context.Background()
	err := m.client.RiverClient.QueuePause(ctx, queueName, nil)
	if err != nil {
		m.ui.statusBar.SetText(fmt.Sprintf("[red]Error pausing queue: %v[white]", err))
	} else {
		m.ui.statusBar.SetText(fmt.Sprintf("[green]Queue '%s' paused[white]", queueName))
		// Refresh queue list
		if err := m.updateQueueList(); err != nil {
			m.ui.statusBar.SetText(fmt.Sprintf("Error: %v", err))
		}
	}
}

func (m *MonitorApp) resumeQueue(queueName string) {
	ctx := context.Background()
	err := m.client.RiverClient.QueueResume(ctx, queueName, nil)
	if err != nil {
		m.ui.statusBar.SetText(fmt.Sprintf("[red]Error resuming queue: %v[white]", err))
	} else {
		m.ui.statusBar.SetText(fmt.Sprintf("[green]Queue '%s' resumed[white]", queueName))
		// Refresh queue list
		if err := m.updateQueueList(); err != nil {
			m.ui.statusBar.SetText(fmt.Sprintf("Error: %v", err))
		}
	}
}
