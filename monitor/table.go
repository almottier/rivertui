package monitor

import (
	"context"
	"fmt"
	"time"

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
				SetTextColor(jobTableHeaderColor).
				SetAlign(tview.AlignLeft).
				SetSelectable(false).
				SetExpansion(1))
	}
}

func (m *MonitorApp) addJobToTable(row int, job *rivertype.JobRow) {
	m.ui.jobList.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf("%d", job.ID)).SetTextColor(jobIDColor))
	m.ui.jobList.SetCell(row, 1, tview.NewTableCell(job.Kind).SetTextColor(jobKindColor))

	// Color the state cell based on job state
	stateCell := m.createStateCell(job.State)
	m.ui.jobList.SetCell(row, 2, stateCell)

	m.ui.jobList.SetCell(row, 3, tview.NewTableCell(fmt.Sprintf("%d/%d", job.Attempt, job.MaxAttempts)).SetTextColor(jobAttemptColor))

	// Error count
	if len(job.Errors) > 0 {
		m.ui.jobList.SetCell(row, 4, tview.NewTableCell(fmt.Sprintf("%d", len(job.Errors))).SetTextColor(jobErrorCountColor))
	} else {
		m.ui.jobList.SetCell(row, 4, tview.NewTableCell("").SetTextColor(jobTimestampColor))
	}

	// Duration
	m.setDurationCell(row, job)

	// Time fields
	m.ui.jobList.SetCell(row, 6, tview.NewTableCell(formatTimeAgo(job.CreatedAt)).SetTextColor(jobTimestampColor))
	m.ui.jobList.SetCell(row, 7, tview.NewTableCell(formatTimeAgo(job.ScheduledAt)).SetTextColor(jobTimestampColor))
	if job.AttemptedAt != nil {
		m.ui.jobList.SetCell(row, 8, tview.NewTableCell(formatTimeAgo(*job.AttemptedAt)).SetTextColor(jobTimestampColor))
	} else {
		m.ui.jobList.SetCell(row, 8, tview.NewTableCell("").SetTextColor(jobTimestampColor).SetBackgroundColor(jobTimestampBgColor))
	}
	if job.FinalizedAt != nil {
		m.ui.jobList.SetCell(row, 9, tview.NewTableCell(formatTimeAgo(*job.FinalizedAt)).SetTextColor(jobTimestampColor))
	} else {
		m.ui.jobList.SetCell(row, 9, tview.NewTableCell("").SetTextColor(jobTimestampColor).SetBackgroundColor(jobTimestampBgColor))
	}

	m.ui.jobList.SetCell(row, 10, tview.NewTableCell(job.Queue).SetTextColor(jobQueueColor))
}

func (m *MonitorApp) createStateCell(state rivertype.JobState) *tview.TableCell {
	stateCell := tview.NewTableCell(string(state))
	switch state {
	case rivertype.JobStateAvailable:
		stateCell.SetTextColor(jobStateAvailableColor)
	case rivertype.JobStateRunning:
		stateCell.SetTextColor(jobStateRunningColor)
	case rivertype.JobStateCompleted:
		stateCell.SetTextColor(jobStateCompletedColor)
	case rivertype.JobStateDiscarded:
		stateCell.SetTextColor(jobStateDiscardedColor)
	case rivertype.JobStateCancelled:
		stateCell.SetTextColor(jobStateCancelledColor)
	case rivertype.JobStateRetryable:
		stateCell.SetTextColor(jobStateRetryableColor)
	case rivertype.JobStateScheduled:
		stateCell.SetTextColor(jobStateScheduledColor)
	}
	return stateCell
}

func (m *MonitorApp) setDurationCell(row int, job *rivertype.JobRow) {
	if job.AttemptedAt != nil {
		if job.FinalizedAt != nil {
			duration := job.FinalizedAt.Sub(*job.AttemptedAt)
			m.ui.jobList.SetCell(row, 5, tview.NewTableCell(formatDuration(duration)).SetTextColor(jobDurationFinalizedColor))
		} else {
			duration := time.Since(*job.AttemptedAt)
			m.ui.jobList.SetCell(row, 5, tview.NewTableCell(formatDuration(duration)).SetTextColor(jobDurationActiveColor))
		}
	} else {
		m.ui.jobList.SetCell(row, 5, tview.NewTableCell("").SetTextColor(jobDurationEmptyColor).SetBackgroundColor(jobDurationEmptyBgColor))
	}
}
