package monitor

import (
	"os"
	"strings"

	"github.com/almottier/rivertui/config"
	"github.com/almottier/rivertui/internal/client"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/rivo/tview"
)

// Page names
const (
	PageList         = "list"
	PageDetails      = "details"
	PageKindFilter   = "kindFilter"
	PageConfirmation = "confirmation"
	PageQueues       = "queues"
)

// State filter configuration
type StateFilterConfig struct {
	States []rivertype.JobState
	Labels []string
}

func newStateFilterConfig() *StateFilterConfig {
	return &StateFilterConfig{
		States: []rivertype.JobState{
			rivertype.JobStateCompleted,
			rivertype.JobStateAvailable,
			rivertype.JobStateDiscarded,
			rivertype.JobStateRunning,
			rivertype.JobStateRetryable,
			rivertype.JobStateCancelled,
			rivertype.JobStateScheduled,
		},
		Labels: []string{
			"All",
			string(rivertype.JobStateCompleted),
			string(rivertype.JobStateAvailable),
			string(rivertype.JobStateDiscarded),
			string(rivertype.JobStateRunning),
			string(rivertype.JobStateRetryable),
			string(rivertype.JobStateCancelled),
			string(rivertype.JobStateScheduled),
		},
	}
}

func (sfc *StateFilterConfig) GetStateByNumber(num int) rivertype.JobState {
	if num == 0 || num > len(sfc.States) {
		return ""
	}
	return sfc.States[num-1]
}

// Pagination handles pagination logic
type Pagination struct {
	currentPage     int
	pageSize        int
	cursors         []*river.JobListCursor
	lastCursor      *river.JobListCursor
	hasNextPage     bool
	totalJobsOnPage int
}

func newPagination() *Pagination {
	return &Pagination{
		currentPage: 1,
		pageSize:    50,
		cursors:     make([]*river.JobListCursor, 0),
	}
}

func (p *Pagination) Reset() {
	p.currentPage = 1
	p.cursors = make([]*river.JobListCursor, 0)
	p.lastCursor = nil
	p.hasNextPage = false
	p.totalJobsOnPage = 0
}

func (p *Pagination) NextPage() bool {
	if !p.hasNextPage {
		return false
	}
	if p.lastCursor != nil && len(p.cursors) == p.currentPage-1 {
		p.cursors = append(p.cursors, p.lastCursor)
	}
	p.currentPage++
	return true
}

func (p *Pagination) PreviousPage() bool {
	if p.currentPage <= 1 {
		return false
	}
	p.currentPage--
	if len(p.cursors) > 0 {
		p.cursors = p.cursors[:len(p.cursors)-1]
	}
	return true
}

func (p *Pagination) GetCurrentCursor() *river.JobListCursor {
	if p.currentPage > 1 && len(p.cursors) >= p.currentPage-1 {
		return p.cursors[p.currentPage-2]
	}
	return nil
}

// JobFilter handles job filtering logic
type JobFilter struct {
	stateFilter      []rivertype.JobState
	kindFilter       []string
	selectedStateNum int
	stateConfig      *StateFilterConfig
}

func newJobFilter() *JobFilter {
	return &JobFilter{
		selectedStateNum: 0,
		stateConfig:      newStateFilterConfig(),
	}
}

func (jf *JobFilter) SetStateFilter(stateNum int) {
	jf.selectedStateNum = stateNum
	if stateNum == 0 {
		jf.stateFilter = nil
	} else {
		state := jf.stateConfig.GetStateByNumber(stateNum)
		if state != "" {
			jf.stateFilter = []rivertype.JobState{state}
		}
	}
}

func (jf *JobFilter) SetKindFilter(kinds []string) {
	jf.kindFilter = kinds
}

func (jf *JobFilter) ApplyToParams(opts *river.JobListParams) *river.JobListParams {
	if len(jf.stateFilter) > 0 {
		opts = opts.States(jf.stateFilter...)
	}
	if len(jf.kindFilter) > 0 {
		opts = opts.Kinds(jf.kindFilter...)
	}
	return opts
}

// UIComponents holds all UI components
type UIComponents struct {
	app               *tview.Application
	pages             *tview.Pages
	jobList           *tview.Table
	jobDetails        *tview.TextView
	queueList         *tview.Table
	filterStatusBar   *tview.TextView
	statusBar         *tview.TextView
	kindFilterInput   *tview.InputField
	confirmationModal *tview.TextView
}

func newUIComponents() *UIComponents {
	app := tview.NewApplication()
	setupAppTheme()

	return &UIComponents{
		app:               app,
		pages:             tview.NewPages(),
		jobList:           createJobListTable(),
		jobDetails:        createJobDetailsView(),
		queueList:         createQueueListTable(),
		filterStatusBar:   createStatusBar(),
		statusBar:         createStatusBar(),
		kindFilterInput:   createKindFilterInput(),
		confirmationModal: createConfirmationModal(),
	}
}

// ModalState handles confirmation modal state
type ModalState struct {
	onYes func()
	onNo  func()
}

func newModalState() *ModalState {
	return &ModalState{}
}

func (ms *ModalState) Set(onYes, onNo func()) {
	ms.onYes = onYes
	ms.onNo = onNo
}

func (ms *ModalState) Clear() {
	ms.onYes = nil
	ms.onNo = nil
}

func (ms *ModalState) ExecuteYes() {
	if ms.onYes != nil {
		ms.onYes()
	}
}

func (ms *ModalState) ExecuteNo() {
	if ms.onNo != nil {
		ms.onNo()
	}
}

// MonitorApp represents the main TUI application
type MonitorApp struct {
	ui                *UIComponents
	client            *client.Client
	config            *config.Config
	pagination        *Pagination
	filter            *JobFilter
	modalState        *ModalState
	currentJobID      string
	initialJobID      int64
	lastJobKinds      []string
	scrollToBeginning bool
	lastActivePage    string
}

// NewMonitorApp creates a new monitor application
func NewMonitorApp(cli *client.Client, cfg *config.Config, jobID int64, kindFilter string) *MonitorApp {
	// Set COLORTERM and TERM if not already set
	if os.Getenv("COLORTERM") == "" {
		os.Setenv("COLORTERM", "truecolor")
	}
	term := os.Getenv("TERM")
	if term == "" || !strings.Contains(term, "256color") {
		os.Setenv("TERM", "xterm-256color")
	}

	ui := newUIComponents()

	monitor := &MonitorApp{
		ui:                ui,
		client:            cli,
		config:            cfg,
		pagination:        newPagination(),
		filter:            newJobFilter(),
		modalState:        newModalState(),
		initialJobID:      jobID,
		lastJobKinds:      make([]string, 0),
		scrollToBeginning: true,
	}

	// Set initial kind filter if provided
	if kindFilter != "" {
		monitor.filter.SetKindFilter([]string{kindFilter})
	}

	monitor.setupUI()
	monitor.setupKeyBindings()
	return monitor
}
