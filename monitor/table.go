package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/rivo/tview"
)

// updateJobList refreshes the job list table
func (m *MonitorApp) updateJobList() error {
	ctx := context.Background()

	// Build ListJobsOpts
	opts := river.NewJobListParams().
		First(m.pagination.pageSize).
		OrderBy(river.JobListOrderByField("id"), river.SortOrderDesc)

	// Apply filters
	opts = m.filter.ApplyToParams(opts)

	// Apply cursor for pagination if we're not on first page
	if cursor := m.pagination.GetCurrentCursor(); cursor != nil {
		opts = opts.After(cursor)
	}

	// Fetch jobs
	result, err := m.client.RiverClient.JobList(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to list jobs: %w", err)
	}

	// Update pagination state
	m.pagination.totalJobsOnPage = len(result.Jobs)
	m.pagination.lastCursor = result.LastCursor
	m.pagination.hasNextPage = result.LastCursor != nil && len(result.Jobs) == m.pagination.pageSize

	// Collect unique kinds for modal
	kindSet := make(map[string]struct{})
	for _, job := range result.Jobs {
		kindSet[job.Kind] = struct{}{}
	}
	m.lastJobKinds = m.lastJobKinds[:0]
	for kind := range kindSet {
		m.lastJobKinds = append(m.lastJobKinds, kind)
	}

	// Clear existing table
	m.ui.jobList.Clear()

	// Update the title with pagination info
	paginationInfo := ""
	if m.pagination.hasNextPage || m.pagination.currentPage > 1 {
		paginationInfo = fmt.Sprintf(" (Page %d)", m.pagination.currentPage)
	}
	m.ui.jobList.SetTitle(fmt.Sprintf(" 🚀 Jobs%s ", paginationInfo))

	// Set headers
	m.setTableHeaders()

	// Add jobs to table
	for i, job := range result.Jobs {
		m.addJobToTable(i+1, job)
	}

	if m.scrollToBeginning && len(result.Jobs) > 0 {
		m.ui.jobList.ScrollToBeginning()
		m.ui.jobList.Select(1, 0)
		m.scrollToBeginning = false
	}

	// Update status bar
	m.setListModeStatus()

	return nil
}

func (m *MonitorApp) setTableHeaders() {
	headers := []string{"ID", "KIND", "STATE", "ATTEMPT", "ERRORS", "DURATION", "CREATED", "SCHEDULED", "LAST_ATTEMPT", "FINALIZED", "QUEUE"}
	for i, header := range headers {
		m.ui.jobList.SetCell(0, i,
			tview.NewTableCell(header).
				SetTextColor(ColorPrimary).
				SetAlign(tview.AlignLeft).
				SetSelectable(false).
				SetExpansion(1))
	}
}

func (m *MonitorApp) addJobToTable(row int, job *rivertype.JobRow) {
	m.ui.jobList.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf("%d", job.ID)).SetTextColor(ColorSecondary))
	m.ui.jobList.SetCell(row, 1, tview.NewTableCell(job.Kind).SetTextColor(tcell.ColorWhite))

	// Color the state cell based on job state
	stateCell := m.createStateCell(job.State)
	m.ui.jobList.SetCell(row, 2, stateCell)

	m.ui.jobList.SetCell(row, 3, tview.NewTableCell(fmt.Sprintf("%d/%d", job.Attempt, job.MaxAttempts)).SetTextColor(ColorSecondary))

	// Error count
	if len(job.Errors) > 0 {
		m.ui.jobList.SetCell(row, 4, tview.NewTableCell(fmt.Sprintf("%d", len(job.Errors))).SetTextColor(ColorError))
	} else {
		m.ui.jobList.SetCell(row, 4, tview.NewTableCell("").SetTextColor(ColorSecondary).SetBackgroundColor(ColorBackground))
	}

	// Duration
	m.setDurationCell(row, job)

	// Time fields
	m.ui.jobList.SetCell(row, 6, tview.NewTableCell(formatTimeAgo(job.CreatedAt)).SetTextColor(ColorSecondary))
	m.ui.jobList.SetCell(row, 7, tview.NewTableCell(formatTimeAgo(job.ScheduledAt)).SetTextColor(ColorSecondary))
	if job.AttemptedAt != nil {
		m.ui.jobList.SetCell(row, 8, tview.NewTableCell(formatTimeAgo(*job.AttemptedAt)).SetTextColor(ColorSecondary))
	} else {
		m.ui.jobList.SetCell(row, 8, tview.NewTableCell("").SetTextColor(ColorSecondary).SetBackgroundColor(ColorBackground))
	}
	if job.FinalizedAt != nil {
		m.ui.jobList.SetCell(row, 9, tview.NewTableCell(formatTimeAgo(*job.FinalizedAt)).SetTextColor(ColorSecondary))
	} else {
		m.ui.jobList.SetCell(row, 9, tview.NewTableCell("").SetTextColor(ColorSecondary).SetBackgroundColor(ColorBackground))
	}

	m.ui.jobList.SetCell(row, 10, tview.NewTableCell(job.Queue).SetTextColor(tcell.NewRGBColor(100, 116, 139)))
}

func (m *MonitorApp) createStateCell(state rivertype.JobState) *tview.TableCell {
	stateCell := tview.NewTableCell(string(state))
	switch state {
	case rivertype.JobStateAvailable:
		stateCell.SetTextColor(tcell.NewRGBColor(59, 130, 246)) // Blue-500
	case rivertype.JobStateRunning:
		stateCell.SetTextColor(ColorInfo) // Cyan-500
	case rivertype.JobStateCompleted:
		stateCell.SetTextColor(ColorSuccess) // Emerald-500
	case rivertype.JobStateDiscarded:
		stateCell.SetTextColor(ColorError) // Red-500
	case rivertype.JobStateCancelled:
		stateCell.SetTextColor(ColorWarning) // Amber-500
	case rivertype.JobStateRetryable:
		stateCell.SetTextColor(tcell.NewRGBColor(139, 92, 246)) // Violet-500
	case rivertype.JobStateScheduled:
		stateCell.SetTextColor(tcell.NewRGBColor(107, 114, 128)) // Gray-500
	}
	return stateCell
}

func (m *MonitorApp) setDurationCell(row int, job *rivertype.JobRow) {
	if job.AttemptedAt != nil {
		if job.FinalizedAt != nil {
			duration := job.FinalizedAt.Sub(*job.AttemptedAt)
			m.ui.jobList.SetCell(row, 5, tview.NewTableCell(formatDuration(duration)).SetTextColor(ColorSecondary))
		} else {
			duration := time.Since(*job.AttemptedAt)
			m.ui.jobList.SetCell(row, 5, tview.NewTableCell(formatDuration(duration)).SetTextColor(ColorInfo))
		}
	} else {
		m.ui.jobList.SetCell(row, 5, tview.NewTableCell("").SetTextColor(ColorSecondary).SetBackgroundColor(ColorBackground))
	}
}
