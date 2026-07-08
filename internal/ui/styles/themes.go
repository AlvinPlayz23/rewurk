package styles

import (
	"image/color"
	"strings"

	"github.com/charmbracelet/x/exp/charmtone"
)

func colorRGBA(r, g, b byte) color.Color {
	return color.RGBA{r, g, b, 255}
}

// ThemeKeyForProvider returns a stable identifier for the theme associated
// with the given provider ID.
func ThemeKeyForProvider(providerID string) string {
	switch providerID {
	case "hyper":
		return "hyper"
	default:
		return "purple"
	}
}

// ThemeKey returns the normalized key for a configured theme name. Empty,
// "auto", and unknown names resolve to the provider-derived theme key.
func ThemeKey(themeName, providerID string) string {
	switch strings.ToLower(strings.TrimSpace(themeName)) {
	case "", "auto", "default":
		return ThemeKeyForProvider(providerID)
	case "purple", "pantera":
		return "purple"
	case "midnight":
		return "midnight"
	case "forest":
		return "forest"
	case "amber":
		return "amber"
	case "light", "paper":
		return "light"
	default:
		return ThemeKeyForProvider(providerID)
	}
}

// Theme returns the Styles for a configured theme name. Empty and "auto"
// preserve provider-based defaults.
func Theme(themeName, providerID string) Styles {
	switch ThemeKey(themeName, providerID) {
	case "hyper":
		return HypercrushObsidiana()
	case "midnight":
		return Midnight()
	case "forest":
		return Forest()
	case "amber":
		return Amber()
	case "light":
		return Paper()
	default:
		return CharmtonePantera()
	}
}

// ThemeForProvider returns the Styles associated with the given provider ID.
func ThemeForProvider(providerID string) Styles {
	return Theme("auto", providerID)
}

// CharmtonePantera returns the default purple dark theme.
func CharmtonePantera() Styles {
	return applyPanteraOverrides(quickStyle(purplePalette()))
}

// HypercrushObsidiana returns the Hypercrush dark theme.
func HypercrushObsidiana() Styles {
	return CharmtonePantera()
}

// Midnight returns a blue dark theme.
func Midnight() Styles {
	return quickStyle(midnightPalette())
}

// Forest returns a green dark theme.
func Forest() Styles {
	return quickStyle(forestPalette())
}

// Amber returns a warm dark theme.
func Amber() Styles {
	return quickStyle(amberPalette())
}

// Paper returns a light theme.
func Paper() Styles {
	return quickStyle(paperPalette())
}

func applyPanteraOverrides(s Styles) Styles {
	// Bang ! prompt overrides - use Salt/Hazy/Larple colors.
	s.Editor.PromptBangIconFocused = s.Editor.PromptBangIconFocused.
		Foreground(charmtone.Salt).
		Background(charmtone.Hazy)
	s.Editor.PromptBangDotsFocused = s.Editor.PromptBangDotsFocused.
		Foreground(charmtone.Hazy)
	s.Editor.PromptBangDotsBlurred = s.Editor.PromptBangDotsBlurred.
		Foreground(charmtone.Larple)

	// Shell bar/prompt overrides - use Charple/Iron/Hazy colors.
	s.Messages.ShellBarFocused = s.Messages.ShellBarFocused.
		BorderForeground(charmtone.Charple)
	s.Messages.ShellBarBlurred = s.Messages.ShellBarBlurred.
		BorderForeground(charmtone.Iron)
	s.Messages.ShellPrompt = s.Messages.ShellPrompt.
		Foreground(charmtone.Hazy)
	s.Messages.ShellPromptBlurred = s.Messages.ShellPromptBlurred.
		Foreground(charmtone.Hazy)

	return s
}

func purplePalette() quickStyleOpts {
	return quickStyleOpts{
		primary:   colorRGBA(0x6c, 0x5c, 0xe7),
		secondary: colorRGBA(0xa2, 0x9b, 0xfe),
		accent:    colorRGBA(0x00, 0xce, 0xc9),
		keyword:   colorRGBA(0xff, 0x6e, 0xc7),

		fgBase:       colorRGBA(0xc8, 0xc8, 0xd0),
		fgBright:     colorRGBA(0xe8, 0xe8, 0xff),
		fgMoreSubtle: colorRGBA(0x88, 0x88, 0xbb),
		fgSubtle:     colorRGBA(0x88, 0x88, 0xbb),
		fgMostSubtle: colorRGBA(0x5a, 0x5a, 0x7a),
		onPrimary:    colorRGBA(0xff, 0xff, 0xff),

		bgBase:         colorRGBA(0x0a, 0x0a, 0x0f),
		bgPanel:        colorRGBA(0x0d, 0x0d, 0x1a),
		bgChat:         colorRGBA(0x08, 0x08, 0x10),
		bgLeastVisible: colorRGBA(0x14, 0x14, 0x2a),
		bgLessVisible:  colorRGBA(0x1a, 0x1a, 0x2e),
		bgMostVisible:  colorRGBA(0x2a, 0x2a, 0x4a),
		separator:      colorRGBA(0x2a, 0x2a, 0x4a),

		destructive:       colorRGBA(0xe1, 0x70, 0x55),
		error:             colorRGBA(0xe1, 0x70, 0x55),
		warningSubtle:     colorRGBA(0xfd, 0xcb, 0x6e),
		warning:           colorRGBA(0xfd, 0xcb, 0x6e),
		denied:            colorRGBA(0xe1, 0x70, 0x55),
		busy:              colorRGBA(0xfd, 0xcb, 0x6e),
		info:              colorRGBA(0x74, 0xb9, 0xff),
		infoMoreSubtle:    colorRGBA(0x74, 0xb9, 0xff),
		infoMostSubtle:    colorRGBA(0x5a, 0x8a, 0xcc),
		success:           colorRGBA(0x00, 0xb8, 0x94),
		successMoreSubtle: colorRGBA(0x00, 0xb8, 0x94),
		successMostSubtle: colorRGBA(0x00, 0xce, 0xc9),

		ansiBlack: colorRGBA(0x14, 0x14, 0x2a), ansiRed: colorRGBA(0xe1, 0x70, 0x55), ansiGreen: colorRGBA(0x00, 0xb8, 0x94), ansiYellow: colorRGBA(0xfd, 0xcb, 0x6e), ansiBlue: colorRGBA(0x6c, 0x5c, 0xe7), ansiMagenta: colorRGBA(0xa2, 0x9b, 0xfe), ansiCyan: colorRGBA(0x74, 0xb9, 0xff), ansiWhite: colorRGBA(0xc8, 0xc8, 0xd0),
		ansiBrightBlack: colorRGBA(0x2a, 0x2a, 0x4a), ansiBrightRed: colorRGBA(0xe1, 0x70, 0x55), ansiBrightGreen: colorRGBA(0x00, 0xb8, 0x94), ansiBrightYellow: colorRGBA(0xfd, 0xcb, 0x6e), ansiBrightBlue: colorRGBA(0xa2, 0x9b, 0xfe), ansiBrightMagenta: colorRGBA(0xff, 0x6e, 0xc7), ansiBrightCyan: colorRGBA(0x74, 0xb9, 0xff), ansiBrightWhite: colorRGBA(0xe8, 0xe8, 0xff),
	}
}

func midnightPalette() quickStyleOpts {
	o := purplePalette()
	o.primary = colorRGBA(0x4f, 0x8c, 0xff)
	o.secondary = colorRGBA(0x8a, 0xb4, 0xff)
	o.accent = colorRGBA(0x2d, 0xd4, 0xbf)
	o.keyword = colorRGBA(0xc0, 0x84, 0xfc)
	o.fgBase = colorRGBA(0xc9, 0xd7, 0xee)
	o.fgBright = colorRGBA(0xf0, 0xf6, 0xff)
	o.fgMoreSubtle = colorRGBA(0x7d, 0x91, 0xb3)
	o.fgSubtle = colorRGBA(0x8b, 0xa3, 0xc7)
	o.fgMostSubtle = colorRGBA(0x4d, 0x5f, 0x7c)
	o.bgBase = colorRGBA(0x06, 0x0a, 0x12)
	o.bgPanel = colorRGBA(0x08, 0x10, 0x1e)
	o.bgChat = colorRGBA(0x04, 0x08, 0x10)
	o.bgLeastVisible = colorRGBA(0x0e, 0x1a, 0x2f)
	o.bgLessVisible = colorRGBA(0x13, 0x22, 0x3a)
	o.bgMostVisible = colorRGBA(0x1e, 0x35, 0x58)
	o.separator = colorRGBA(0x1e, 0x35, 0x58)
	o.info = colorRGBA(0x60, 0xa5, 0xfa)
	o.infoMoreSubtle = colorRGBA(0x60, 0xa5, 0xfa)
	o.infoMostSubtle = colorRGBA(0x3b, 0x82, 0xf6)
	o.success = colorRGBA(0x34, 0xd3, 0x99)
	o.successMoreSubtle = colorRGBA(0x34, 0xd3, 0x99)
	o.successMostSubtle = colorRGBA(0x2d, 0xd4, 0xbf)
	o.ansiBlue = o.primary
	o.ansiBrightBlue = o.secondary
	o.ansiCyan = o.accent
	o.ansiBrightCyan = colorRGBA(0x5e, 0xe0, 0xd1)
	o.ansiWhite = o.fgBase
	o.ansiBrightWhite = o.fgBright
	return o
}

func forestPalette() quickStyleOpts {
	o := midnightPalette()
	o.primary = colorRGBA(0x22, 0xc5, 0x5e)
	o.secondary = colorRGBA(0x86, 0xef, 0xac)
	o.accent = colorRGBA(0x14, 0xb8, 0xa6)
	o.keyword = colorRGBA(0xf4, 0x72, 0xb6)
	o.bgBase = colorRGBA(0x05, 0x0f, 0x0a)
	o.bgPanel = colorRGBA(0x08, 0x18, 0x10)
	o.bgChat = colorRGBA(0x03, 0x0c, 0x08)
	o.bgLeastVisible = colorRGBA(0x0d, 0x21, 0x15)
	o.bgLessVisible = colorRGBA(0x13, 0x2d, 0x1d)
	o.bgMostVisible = colorRGBA(0x1e, 0x46, 0x2c)
	o.separator = colorRGBA(0x1e, 0x46, 0x2c)
	o.successMostSubtle = o.secondary
	o.ansiGreen = o.primary
	o.ansiBrightGreen = o.secondary
	return o
}

func amberPalette() quickStyleOpts {
	o := midnightPalette()
	o.primary = colorRGBA(0xf5, 0x9e, 0x0b)
	o.secondary = colorRGBA(0xfd, 0xba, 0x74)
	o.accent = colorRGBA(0xfb, 0x71, 0x85)
	o.keyword = colorRGBA(0xf4, 0x72, 0xb6)
	o.bgBase = colorRGBA(0x12, 0x0b, 0x05)
	o.bgPanel = colorRGBA(0x1a, 0x10, 0x08)
	o.bgChat = colorRGBA(0x10, 0x08, 0x04)
	o.bgLeastVisible = colorRGBA(0x2a, 0x1a, 0x0d)
	o.bgLessVisible = colorRGBA(0x36, 0x22, 0x12)
	o.bgMostVisible = colorRGBA(0x54, 0x35, 0x1a)
	o.separator = colorRGBA(0x54, 0x35, 0x1a)
	o.warning = o.primary
	o.warningSubtle = o.secondary
	o.busy = o.primary
	o.ansiYellow = o.primary
	o.ansiBrightYellow = o.secondary
	return o
}

func paperPalette() quickStyleOpts {
	return quickStyleOpts{
		primary:   colorRGBA(0x5b, 0x5b, 0xd6),
		secondary: colorRGBA(0x7c, 0x3a, 0xed),
		accent:    colorRGBA(0x0f, 0x76, 0x66),
		keyword:   colorRGBA(0xbe, 0x18, 0x5d),

		fgBase:       colorRGBA(0x2f, 0x32, 0x3a),
		fgBright:     colorRGBA(0x12, 0x16, 0x1f),
		fgMoreSubtle: colorRGBA(0x6b, 0x72, 0x80),
		fgSubtle:     colorRGBA(0x4b, 0x55, 0x63),
		fgMostSubtle: colorRGBA(0x9c, 0xa3, 0xaf),
		onPrimary:    colorRGBA(0xff, 0xff, 0xff),

		bgBase:         colorRGBA(0xf4, 0xef, 0xe5),
		bgPanel:        colorRGBA(0xec, 0xe4, 0xd6),
		bgChat:         colorRGBA(0xfb, 0xf7, 0xef),
		bgLeastVisible: colorRGBA(0xee, 0xe7, 0xda),
		bgLessVisible:  colorRGBA(0xe2, 0xd8, 0xc9),
		bgMostVisible:  colorRGBA(0xc8, 0xb8, 0xa5),
		separator:      colorRGBA(0xd3, 0xc5, 0xb4),

		destructive:       colorRGBA(0xdc, 0x26, 0x26),
		error:             colorRGBA(0xdc, 0x26, 0x26),
		warningSubtle:     colorRGBA(0xd9, 0x77, 0x06),
		warning:           colorRGBA(0xb4, 0x53, 0x09),
		denied:            colorRGBA(0xdc, 0x26, 0x26),
		busy:              colorRGBA(0xd9, 0x77, 0x06),
		info:              colorRGBA(0x25, 0x63, 0xeb),
		infoMoreSubtle:    colorRGBA(0x25, 0x63, 0xeb),
		infoMostSubtle:    colorRGBA(0x60, 0x75, 0x8a),
		success:           colorRGBA(0x16, 0xa3, 0x4a),
		successMoreSubtle: colorRGBA(0x16, 0xa3, 0x4a),
		successMostSubtle: colorRGBA(0x0f, 0x76, 0x66),

		ansiBlack: colorRGBA(0x2f, 0x32, 0x3a), ansiRed: colorRGBA(0xdc, 0x26, 0x26), ansiGreen: colorRGBA(0x16, 0xa3, 0x4a), ansiYellow: colorRGBA(0xd9, 0x77, 0x06), ansiBlue: colorRGBA(0x25, 0x63, 0xeb), ansiMagenta: colorRGBA(0xbe, 0x18, 0x5d), ansiCyan: colorRGBA(0x0f, 0x76, 0x66), ansiWhite: colorRGBA(0x4b, 0x55, 0x63),
		ansiBrightBlack: colorRGBA(0x9c, 0xa3, 0xaf), ansiBrightRed: colorRGBA(0xef, 0x44, 0x44), ansiBrightGreen: colorRGBA(0x22, 0xc5, 0x5e), ansiBrightYellow: colorRGBA(0xf5, 0x9e, 0x0b), ansiBrightBlue: colorRGBA(0x60, 0xa5, 0xfa), ansiBrightMagenta: colorRGBA(0xdb, 0x27, 0x77), ansiBrightCyan: colorRGBA(0x14, 0xb8, 0xa6), ansiBrightWhite: colorRGBA(0x12, 0x16, 0x1f),
	}
}
