package monitor

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

var (
	themeBackground      = tcell.ColorDefault
	themeContrastPrimary = tcell.ColorDefault
)

func init() {
	if envVar := os.Getenv("RIVER_THEME_TRANSPARENT_BG"); strings.ToLower(envVar) != "true" {
		themeBackground = getEnvColor("RIVER_THEME_BACKGROUND", tcell.NewRGBColor(30, 41, 59))
		themeContrastPrimary = getEnvColor("RIVER_THEME_CONTRAST_PRIMARY", tcell.ColorBlack)
	}
}

var (
	// Theme colors
	themePrimary           = getEnvColor("RIVER_THEME_PRIMARY", tcell.ColorWhite)
	themeSecondary         = getEnvColor("RIVER_THEME_SECONDARY", tcell.NewRGBColor(120, 130, 140))
	themeTertiary          = getEnvColor("RIVER_THEME_TERTIARY", tcell.NewRGBColor(100, 116, 139))
	themeBorder            = getEnvColor("RIVER_THEME_BORDER", tcell.NewRGBColor(50, 60, 70))
	themeWarning           = getEnvColor("RIVER_THEME_WARNING", tcell.NewRGBColor(220, 38, 38))
	themeInfo              = getEnvColor("RIVER_THEME_INFO", tcell.NewRGBColor(6, 182, 212))
	themeSuccess           = getEnvColor("RIVER_THEME_SUCCESS", tcell.NewRGBColor(16, 185, 129))
	themeError             = getEnvColor("RIVER_THEME_ERROR", tcell.NewRGBColor(239, 68, 68))
	themeAvailable         = getEnvColor("RIVER_THEME_AVAILABLE", tcell.NewRGBColor(59, 130, 246))
	themeCancelled         = getEnvColor("RIVER_THEME_CANCELLED", tcell.NewRGBColor(251, 191, 36))
	themeRetryable         = getEnvColor("RIVER_THEME_RETRYABLE", tcell.NewRGBColor(139, 92, 246))
	themeScheduled         = getEnvColor("RIVER_THEME_SCHEDULED", tcell.NewRGBColor(107, 114, 128))
	themeTitle             = getEnvColor("RIVER_THEME_TITLE", tcell.NewRGBColor(200, 220, 240))
	themeContrastSecondary = getEnvColor("RIVER_THEME_CONTRAST_SECONDARY", tcell.NewRGBColor(203, 213, 225))
	themeSelectedBg        = getEnvColor("RIVER_THEME_SELECTED_BG", tcell.NewRGBColor(30, 58, 138))
)

var (
	// tview theme
	primitiveBackground        = getEnvColor("RIVER_COLOR_PRIMITIVE_BACKGROUND", themeContrastPrimary)
	contrastBackground         = getEnvColor("RIVER_COLOR_CONTRAST_BACKGROUND", themeBackground)
	moreContrastBackground     = getEnvColor("RIVER_COLOR_MORE_CONTRAST_BACKGROUND", themeBackground)
	borderColor                = getEnvColor("RIVER_COLOR_BORDER", themeBorder)
	titleColor                 = getEnvColor("RIVER_COLOR_TITLE", themeTitle)
	graphicsColor              = getEnvColor("RIVER_COLOR_GRAPHICS", themeBorder)
	primaryTextColor           = getEnvColor("RIVER_COLOR_PRIMARY_TEXT", themePrimary)
	secondaryTextColor         = getEnvColor("RIVER_COLOR_SECONDARY_TEXT", themeSecondary)
	tertiaryTextColor          = getEnvColor("RIVER_COLOR_TERTIARY_TEXT", themeTertiary)
	inverseTextColor           = getEnvColor("RIVER_COLOR_INVERSE_TEXT", themeContrastPrimary)
	contrastSecondaryTextColor = getEnvColor("RIVER_COLOR_CONTRAST_SECONDARY_TEXT", themeContrastSecondary)

	// job list
	jobListSelectedBg = getEnvColor("RIVER_COLOR_JOB_LIST_SELECTED_BG", themeSelectedBg)
	jobListSelectedFg = getEnvColor("RIVER_COLOR_JOB_LIST_SELECTED_FG", themePrimary)
	jobListBorder     = getEnvColor("RIVER_COLOR_JOB_LIST_BORDER", themeBorder)
	jobListTitle      = getEnvColor("RIVER_COLOR_JOB_LIST_TITLE", themeTitle)
	jobListBg         = getEnvColor("RIVER_COLOR_JOB_LIST_BG", themeBackground)

	// job details
	jobDetailsBorder = getEnvColor("RIVER_COLOR_JOB_DETAILS_BORDER", themeBorder)
	jobDetailsTitle  = getEnvColor("RIVER_COLOR_JOB_DETAILS_TITLE", themeTitle)
	jobDetailsBg     = getEnvColor("RIVER_COLOR_JOB_DETAILS_BG", themeBackground)

	// status bar
	statusBarBg = getEnvColor("RIVER_COLOR_STATUS_BAR_BG", themeBackground)

	// kind filter
	kindFilterLabel   = getEnvColor("RIVER_COLOR_KIND_FILTER_LABEL", themeTitle)
	kindFilterBorder  = getEnvColor("RIVER_COLOR_KIND_FILTER_BORDER", themeTitle)
	kindFilterTitle   = getEnvColor("RIVER_COLOR_KIND_FILTER_TITLE", themeTitle)
	kindFilterBg      = getEnvColor("RIVER_COLOR_KIND_FILTER_BG", themeBackground)
	kindFilterFieldBg = getEnvColor("RIVER_COLOR_KIND_FILTER_FIELD_BG", themeBackground)

	// confirmation modal
	confirmationModalBorder = getEnvColor("RIVER_COLOR_CONFIRMATION_MODAL_BORDER", themeWarning)
	confirmationModalTitle  = getEnvColor("RIVER_COLOR_CONFIRMATION_MODAL_TITLE", themeWarning)
	confirmationModalBg     = getEnvColor("RIVER_COLOR_CONFIRMATION_MODAL_BG", themeBackground)

	// queue list
	queueListSelectedBg = getEnvColor("RIVER_COLOR_QUEUE_LIST_SELECTED_BG", themeSelectedBg)
	queueListSelectedFg = getEnvColor("RIVER_COLOR_QUEUE_LIST_SELECTED_FG", themePrimary)
	queueListBorder     = getEnvColor("RIVER_COLOR_QUEUE_LIST_BORDER", themeBorder)
	queueListTitle      = getEnvColor("RIVER_COLOR_QUEUE_LIST_TITLE", themeTitle)
	queueListBg         = getEnvColor("RIVER_COLOR_QUEUE_LIST_BG", themeBackground)

	// job state
	jobStateAvailableColor = getEnvColor("RIVER_COLOR_JOB_STATE_AVAILABLE", themeAvailable)
	jobStateRunningColor   = getEnvColor("RIVER_COLOR_JOB_STATE_RUNNING", themeInfo)
	jobStateCompletedColor = getEnvColor("RIVER_COLOR_JOB_STATE_COMPLETED", themeSuccess)
	jobStateDiscardedColor = getEnvColor("RIVER_COLOR_JOB_STATE_DISCARDED", themeError)
	jobStateCancelledColor = getEnvColor("RIVER_COLOR_JOB_STATE_CANCELLED", themeCancelled)
	jobStateRetryableColor = getEnvColor("RIVER_COLOR_JOB_STATE_RETRYABLE", themeRetryable)
	jobStateScheduledColor = getEnvColor("RIVER_COLOR_JOB_STATE_SCHEDULED", themeScheduled)

	// Table header colors
	queueTableHeaderColor = getEnvColor("RIVER_COLOR_QUEUE_TABLE_HEADER", themeTitle)
	jobTableHeaderColor   = getEnvColor("RIVER_COLOR_JOB_TABLE_HEADER", themeTitle)

	// queue table cell
	queueNameColor        = getEnvColor("RIVER_COLOR_QUEUE_NAME", themePrimary)
	queueStatePausedColor = getEnvColor("RIVER_COLOR_QUEUE_STATE_PAUSED", themeCancelled)
	queueStateActiveColor = getEnvColor("RIVER_COLOR_QUEUE_STATE_ACTIVE", themeSuccess)
	queueTimestampColor   = getEnvColor("RIVER_COLOR_QUEUE_TIMESTAMP", themeSecondary)
	queueTimestampBgColor = getEnvColor("RIVER_COLOR_QUEUE_TIMESTAMP_BG", themeBackground)

	// job info
	jobIDColor           = getEnvColor("RIVER_COLOR_JOB_ID", themeSecondary)
	jobKindColor         = getEnvColor("RIVER_COLOR_JOB_KIND", themePrimary)
	jobAttemptColor      = getEnvColor("RIVER_COLOR_JOB_ATTEMPT", themeSecondary)
	jobErrorCountColor   = getEnvColor("RIVER_COLOR_JOB_ERROR_COUNT", themeError)
	jobErrorCountBgColor = getEnvColor("RIVER_COLOR_JOB_ERROR_COUNT_BG", themeBackground)
	jobTimestampColor    = getEnvColor("RIVER_COLOR_JOB_TIMESTAMP", themeSecondary)
	jobTimestampBgColor  = getEnvColor("RIVER_COLOR_JOB_TIMESTAMP_BG", themeBackground)
	jobQueueColor        = getEnvColor("RIVER_COLOR_JOB_QUEUE", themeTertiary)

	// job duration
	jobDurationFinalizedColor = getEnvColor("RIVER_COLOR_JOB_DURATION_FINALIZED", themeSecondary)
	jobDurationActiveColor    = getEnvColor("RIVER_COLOR_JOB_DURATION_ACTIVE", themeInfo)
	jobDurationEmptyColor     = getEnvColor("RIVER_COLOR_JOB_DURATION_EMPTY", themeSecondary)
	jobDurationEmptyBgColor   = getEnvColor("RIVER_COLOR_JOB_DURATION_EMPTY_BG", themeBackground)
)

func parseHexColor(s string) (tcell.Color, error) {
	s = strings.TrimPrefix(s, "#")
	if len(s) != 6 {
		return tcell.ColorDefault, strconv.ErrSyntax
	}
	r, err := strconv.ParseUint(s[0:2], 16, 8)
	if err != nil {
		return tcell.ColorDefault, err
	}
	g, err := strconv.ParseUint(s[2:4], 16, 8)
	if err != nil {
		return tcell.ColorDefault, err
	}
	b, err := strconv.ParseUint(s[4:6], 16, 8)
	if err != nil {
		return tcell.ColorDefault, err
	}
	return tcell.NewRGBColor(int32(r), int32(g), int32(b)), nil
}

func getEnvColor(envName string, fallback tcell.Color) tcell.Color {
	val := os.Getenv(envName)
	if val == "" {
		return fallback
	}
	c, err := parseHexColor(val)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to parse %s color %s, using fallback %s\n", envName, val, fallback.CSS())
		return fallback
	}
	return c
}
