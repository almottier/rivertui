package monitor

import (
	"fmt"
	"strings"

	"github.com/riverqueue/river/rivertype"
)

// updateFilterStatusBar updates the filter status bar with numbered options
func (m *MonitorApp) updateFilterStatusBar() {
	var text strings.Builder

	// Search kind/id filter information
	text.WriteString("[#60A5FA]Search:[white] ")
	if len(m.filter.kindFilter) > 0 {
		text.WriteString(fmt.Sprintf("[#3B82F6]%s[white]", m.filter.kindFilter[0]))
	} else {
		text.WriteString("All")
	}
	text.WriteString(" ([#60A5FA]/[white])")

	text.WriteString(" | [#60A5FA]State:[white] ")

	// State filter information
	for i, state := range m.filter.stateConfig.Labels {
		if i > 0 {
			text.WriteString(" | ")
		}
		if i == m.filter.selectedStateNum {
			text.WriteString(fmt.Sprintf("[#3B82F6][[%d:%s]][white]", i, state))
		} else {
			text.WriteString(fmt.Sprintf("[#94A3B8]%d:%s[white]", i, state))
		}
	}

	m.ui.filterStatusBar.SetText(text.String())
}

// setStateFilter sets the state filter based on the selected number
func (m *MonitorApp) setStateFilter(stateNum int) {
	m.filter.SetStateFilter(stateNum)
	m.pagination.Reset()
	m.scrollToBeginning = true
	m.updateFilterStatusBar()
}

// getStateByNumber returns the job state for a given number (0=All, 1=completed, etc.)
func (m *MonitorApp) getStateByNumber(num int) rivertype.JobState {
	return m.filter.stateConfig.GetStateByNumber(num)
}

// nextPage navigates to the next page if available
func (m *MonitorApp) nextPage() {
	if !m.pagination.NextPage() {
		m.ui.statusBar.SetText("[yellow]No more pages available[white]")
	}
}

// previousPage navigates to the previous page if available
func (m *MonitorApp) previousPage() {
	if !m.pagination.PreviousPage() {
		m.ui.statusBar.SetText("[yellow]Already on first page[white]")
	}
}

func (m *MonitorApp) setListModeStatus() {
	m.ui.statusBar.SetText("[#60A5FA]Mode:[white] List | Enter: View details | Ctrl+Q: View queues | n: Next page | p: Prev page | r: Retry job | c: Cancel job | q: Quit")
}

func (m *MonitorApp) setDetailsModeStatus() {
	m.ui.statusBar.SetText("[#60A5FA]Mode:[white] Details | Enter/Esc: Back to list | r: Retry job | c: Cancel job | q: Quit")
}
