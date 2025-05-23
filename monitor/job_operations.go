package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/riverqueue/river/rivertype"
)

// showJobDetails displays detailed information about a selected job
func (m *MonitorApp) showJobDetails(jobID string) {
	m.currentJobID = jobID // Store the current job ID
	ctx := context.Background()

	// Parse job ID
	id, err := strconv.ParseInt(jobID, 10, 64)
	if err != nil {
		m.ui.jobDetails.SetText(fmt.Sprintf("Error: Invalid job ID: %v", err))
		return
	}

	// Fetch job details
	job, err := m.client.RiverClient.JobGet(ctx, id)
	if err != nil {
		m.ui.jobDetails.SetText(fmt.Sprintf("Error: Failed to get job: %v", err))
		return
	}

	// Format job details with spacing, alignment, and color
	var details strings.Builder
	// details.WriteString("\n") // Top blank line
	details.WriteString("[#60A5FA]Job Details[white]\n") // Blue-400

	// Helper for alignment
	pad := func(label string) string { return fmt.Sprintf("%-13s", label) }

	details.WriteString(fmt.Sprintf("%s %v\n", pad("ID:"), job.ID))
	details.WriteString(fmt.Sprintf("%s %s\n", pad("Kind:"), job.Kind))
	details.WriteString(fmt.Sprintf("%s %s\n", pad("Queue:"), job.Queue))

	jobState := fmt.Sprintf("%s", job.State)
	switch job.State {
	case rivertype.JobStateAvailable:
		jobState = "[#3B82F6]" + jobState + "[white]" // Blue-500
	case rivertype.JobStateRunning:
		jobState = "[#06B6D4]" + jobState + "[white]" // Cyan-500
	case rivertype.JobStateCompleted:
		jobState = "[#10B981]" + jobState + "[white]" // Emerald-500
	case rivertype.JobStateDiscarded:
		jobState = "[#EF4444]" + jobState + "[white]" // Red-500
	case rivertype.JobStateCancelled:
		jobState = "[#F59E0B]" + jobState + "[white]" // Amber-500
	case rivertype.JobStateRetryable:
		jobState = "[#8B5CF6]" + jobState + "[white]" // Violet-500
	case rivertype.JobStateScheduled:
		jobState = "[#6B7280]" + jobState + "[white]" // Gray-500
	}
	details.WriteString(fmt.Sprintf("%s %s\n", pad("State:"), jobState))

	details.WriteString(fmt.Sprintf("%s %d/%d\n", pad("Attempt:"), job.Attempt, job.MaxAttempts))

	// Add duration calculation similar to the job list
	if job.AttemptedAt != nil {
		if job.FinalizedAt != nil {
			duration := job.FinalizedAt.Sub(*job.AttemptedAt)
			details.WriteString(fmt.Sprintf("%s %s\n", pad("Duration:"), formatDuration(duration)))
		} else {
			duration := time.Since(*job.AttemptedAt)
			details.WriteString(fmt.Sprintf("%s %s[#06B6D4] (running)[white]\n", pad("Duration:"), formatDuration(duration))) // Cyan-500 for running
		}
	}

	details.WriteString(fmt.Sprintf("%s %s (%s)\n", pad("Created:"), job.CreatedAt.Format(time.RFC3339), formatTimeAgo(job.CreatedAt)))
	details.WriteString(fmt.Sprintf("%s %s (%s)\n", pad("Scheduled:"), job.ScheduledAt.Format(time.RFC3339), formatTimeAgo(job.ScheduledAt)))
	if job.AttemptedAt != nil {
		details.WriteString(fmt.Sprintf("%s %s (%s)\n", pad("Last Attempt:"), job.AttemptedAt.Format(time.RFC3339), formatTimeAgo(*job.AttemptedAt)))
	}
	if job.FinalizedAt != nil {
		details.WriteString(fmt.Sprintf("%s %s (%s)\n", pad("Finalized:"), job.FinalizedAt.Format(time.RFC3339), formatTimeAgo(*job.FinalizedAt)))
	}
	details.WriteString(fmt.Sprintf("%s %s\n", pad("Attempted By:"), strings.Join(job.AttemptedBy, ",")))
	details.WriteString("\n")

	// Add arguments if present
	if len(job.EncodedArgs) > 0 {
		details.WriteString("[#60A5FA]Arguments[white]\n") // Blue-400
		var prettyArgs map[string]interface{}
		if err := json.Unmarshal(job.EncodedArgs, &prettyArgs); err == nil {
			argsJSON, _ := json.MarshalIndent(prettyArgs, "", "  ")
			for _, line := range strings.Split(string(argsJSON), "\n") {
				details.WriteString("  " + line + "\n")
			}
		} else {
			for _, line := range strings.Split(string(job.EncodedArgs), "\n") {
				details.WriteString("  " + line + "\n")
			}
		}
		details.WriteString("\n")
	}

	// Add metadata if present
	if len(job.Metadata) > 0 {
		details.WriteString("[#60A5FA]Metadata[white]\n") // Blue-400
		var prettyMeta map[string]interface{}
		if err := json.Unmarshal(job.Metadata, &prettyMeta); err == nil {
			metaJSON, _ := json.MarshalIndent(prettyMeta, "", "  ")
			for _, line := range strings.Split(string(metaJSON), "\n") {
				details.WriteString("  " + line + "\n")
			}
		} else {
			for _, line := range strings.Split(string(job.Metadata), "\n") {
				details.WriteString("  " + line + "\n")
			}
		}
		details.WriteString("\n")
	}

	// Add errors if present
	if len(job.Errors) > 0 {
		details.WriteString("[#60A5FA]Errors[#EF4444]\n") // Blue-400 for header, Red-500 for content
		errorsJSON, _ := json.MarshalIndent(job.Errors, "", "  ")
		for _, line := range strings.Split(string(errorsJSON), "\n") {
			details.WriteString("  " + line + "\n")
		}
		details.WriteString("[white]\n")
	}

	// Add tags if present
	if len(job.Tags) > 0 {
		details.WriteString("[#60A5FA]Tags[white]\n") // Blue-400
		tagsJSON, _ := json.MarshalIndent(job.Tags, "", "  ")
		for _, line := range strings.Split(string(tagsJSON), "\n") {
			details.WriteString("  " + line + "\n")
		}
		details.WriteString("\n")
	}

	m.ui.jobDetails.SetText(details.String())
	m.ui.app.SetFocus(m.ui.jobDetails)

	// Update status bar with retry shortcut
	m.ui.statusBar.SetText("[#60A5FA]Mode:[white] Details | Enter/Esc: Back to list | r: Retry job | c: Cancel job | q: Quit") // Blue-400
}

func (m *MonitorApp) handleJobRetry() {
	row, _ := m.ui.jobList.GetSelection()
	if row > 0 {
		jobID := m.ui.jobList.GetCell(row, 0).Text
		m.showJobRetryConfirmation(jobID)
	}
}

func (m *MonitorApp) handleJobCancel() {
	row, _ := m.ui.jobList.GetSelection()
	if row > 0 {
		jobID := m.ui.jobList.GetCell(row, 0).Text
		m.showJobCancelConfirmation(jobID)
	}
}

func (m *MonitorApp) handleJobRetryInDetails() {
	if m.currentJobID != "" {
		m.showJobRetryConfirmation(m.currentJobID)
	}
}

func (m *MonitorApp) handleJobCancelInDetails() {
	if m.currentJobID != "" {
		m.showJobCancelConfirmation(m.currentJobID)
	}
}

func (m *MonitorApp) showJobRetryConfirmation(jobID string) {
	m.showConfirmationModal(
		"Retry Job",
		fmt.Sprintf("Are you sure you want to retry job %s?\n\n[#60A5FA]Y[white]: Yes, retry the job\n[#60A5FA]N[white]: No, cancel", jobID),
		func() { m.retryJob(jobID) },
		func() {},
	)
}

func (m *MonitorApp) showJobCancelConfirmation(jobID string) {
	m.showConfirmationModal(
		"Cancel Job",
		fmt.Sprintf("Are you sure you want to cancel job %s?\n\n[#60A5FA]Y[white]: Yes, cancel the job\n[#60A5FA]N[white]: No, go back", jobID),
		func() { m.cancelJob(jobID) },
		func() {},
	)
}

func (m *MonitorApp) retryJob(jobID string) {
	id, err := strconv.ParseInt(jobID, 10, 64)
	if err != nil {
		m.ui.statusBar.SetText(fmt.Sprintf("[red]Error: Invalid job ID: %v[white]", err))
		return
	}

	ctx := context.Background()
	_, err = m.client.RiverClient.JobRetry(ctx, id)
	if err != nil {
		m.ui.statusBar.SetText(fmt.Sprintf("[red]Error retrying job: %v[white]", err))
	} else {
		m.ui.statusBar.SetText("[green]Job retry initiated[white]")
		if m.currentJobID == jobID {
			m.showJobDetails(m.currentJobID)
		}
	}
}

func (m *MonitorApp) cancelJob(jobID string) {
	id, err := strconv.ParseInt(jobID, 10, 64)
	if err != nil {
		m.ui.statusBar.SetText(fmt.Sprintf("[red]Error: Invalid job ID: %v[white]", err))
		return
	}

	ctx := context.Background()
	_, err = m.client.RiverClient.JobCancel(ctx, id)
	if err != nil {
		m.ui.statusBar.SetText(fmt.Sprintf("[red]Error cancelling job: %v[white]", err))
	} else {
		m.ui.statusBar.SetText("[green]Job cancellation initiated[white]")
		if m.currentJobID == jobID {
			m.showJobDetails(m.currentJobID)
		}
	}
}
