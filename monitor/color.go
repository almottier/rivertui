package monitor

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

var (
	ColorPrimary                = getEnvColor("RIVER_COLOR_PRIMARY", tcell.ColorWhite)
	ColorSecondary              = getEnvColor("RIVER_COLOR_SECONDARY", tcell.NewRGBColor(120, 130, 140))
	ColorTertiary               = getEnvColor("RIVER_COLOR_TERTIARY", tcell.NewRGBColor(100, 116, 139))
	ColorBorder                 = getEnvColor("RIVER_COLOR_BORDER", tcell.NewRGBColor(50, 60, 70))
	ColorWarning                = getEnvColor("RIVER_COLOR_WARNING", tcell.NewRGBColor(245, 158, 11))
	ColorInfo                   = getEnvColor("RIVER_COLOR_INFO", tcell.NewRGBColor(6, 182, 212))
	ColorSuccess                = getEnvColor("RIVER_COLOR_SUCCESS", tcell.NewRGBColor(16, 185, 129))
	ColorError                  = getEnvColor("RIVER_COLOR_ERROR", tcell.NewRGBColor(239, 68, 68))
	ColorAvailable              = getEnvColor("RIVER_COLOR_AVAILABLE", tcell.NewRGBColor(59, 130, 246))
	ColorCancelled              = getEnvColor("RIVER_COLOR_CANCELLED", tcell.NewRGBColor(251, 191, 36))
	ColorRetryable              = getEnvColor("RIVER_COLOR_RETRYABLE", tcell.NewRGBColor(139, 92, 246))
	ColorScheduled              = getEnvColor("RIVER_COLOR_SCHEDULED", tcell.NewRGBColor(107, 114, 128))
	ColorTitle                  = getEnvColor("RIVER_COLOR_TITLE", tcell.NewRGBColor(200, 220, 240))
	ColorContrastSecondary      = getEnvColor("RIVER_COLOR_CONTRAST_SECONDARY", tcell.NewRGBColor(203, 213, 225))
	ColorSelectedFg             = getEnvColor("RIVER_COLOR_SELECTED_FG", tcell.ColorWhite)
	ColorSelectedBg             = getEnvColor("RIVER_COLOR_SELECTED_BG", tcell.NewRGBColor(30, 58, 138))
	ColorPrimativeBackground    = getEnvColor("RIVER_COLOR_PRIMATIVE_BACKGROUND", tcell.ColorBlack)
	ColorContrastBackground     = getEnvColor("RIVER_COLOR_CONTRAST_BACKGROUND", tcell.NewRGBColor(30, 41, 59))
	ColorMoreContrastBackground = getEnvColor("RIVER_COLOR_MORE_CONTRAST_BACKGROUND", tcell.NewRGBColor(30, 41, 59))
)

func init() {
	if envVar := os.Getenv("RIVER_COLOR_TRANSPARENT_BG"); strings.ToLower(envVar) == "true" {
		ColorContrastBackground = tcell.ColorDefault
		ColorPrimativeBackground = tcell.ColorDefault
	}
}

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
