package styles

import (
	"image/color"

	"github.com/charmbracelet/x/exp/charmtone"
)

func colorRGBA(r, g, b byte) color.Color {
	return color.RGBA{r, g, b, 255}
}

// ThemeKeyForProvider returns a stable identifier for the theme
// associated with the given provider ID. Providers that share a theme
// yield the same key, so callers can cheaply detect when switching
// providers would not actually change the active theme and skip the
// expensive style rebuild. This is the single source of truth for the
// provider-to-theme mapping; [ThemeForProvider] builds on it.
func ThemeKeyForProvider(providerID string) string {
	switch providerID {
	case "hyper":
		return "hyper"
	default:
		return "default"
	}
}

// ThemeForProvider returns the Styles associated with the given provider
// ID. Unknown or empty provider IDs yield the default Charmtone Pantera
// theme.
func ThemeForProvider(providerID string) Styles {
	switch ThemeKeyForProvider(providerID) {
	case "hyper":
		return HypercrushObsidiana()
	default:
		return CharmtonePantera()
	}
}

// CharmtonePantera returns the Charmtone dark theme. It's the default style
// for the UI.
func CharmtonePantera() Styles {
	s := quickStyle(quickStyleOpts{
		primary:   colorRGBA(0x6c, 0x5c, 0xe7),
		secondary: colorRGBA(0xa2, 0x9b, 0xfe),
		accent:    colorRGBA(0x00, 0xce, 0xc9),
		keyword:   colorRGBA(0xff, 0x6e, 0xc7),

		fgBase:       colorRGBA(0xc8, 0xc8, 0xd0),
		fgBright:     colorRGBA(0xe8, 0xe8, 0xff),
		fgMoreSubtle: colorRGBA(0x88, 0x88, 0xbb),
		fgSubtle:     colorRGBA(0x88, 0x88, 0xbb),
		fgMostSubtle: colorRGBA(0x5a, 0x5a, 0x7a),

		onPrimary: colorRGBA(0xff, 0xff, 0xff),

		bgBase:         colorRGBA(0x0a, 0x0a, 0x0f),
		bgPanel:        colorRGBA(0x0d, 0x0d, 0x1a),
		bgChat:         colorRGBA(0x08, 0x08, 0x10),
		bgLeastVisible: colorRGBA(0x14, 0x14, 0x2a),
		bgLessVisible:  colorRGBA(0x1a, 0x1a, 0x2e),
		bgMostVisible:  colorRGBA(0x2a, 0x2a, 0x4a),

		separator: colorRGBA(0x2a, 0x2a, 0x4a),

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

		// ANSI 16-color palette for remapping raw terminal output
		// (e.g. bang-mode shell commands) onto legible colors.
		ansiBlack:   colorRGBA(0x14, 0x14, 0x2a),
		ansiRed:     colorRGBA(0xe1, 0x70, 0x55),
		ansiGreen:   colorRGBA(0x00, 0xb8, 0x94),
		ansiYellow:  colorRGBA(0xfd, 0xcb, 0x6e),
		ansiBlue:    colorRGBA(0x6c, 0x5c, 0xe7),
		ansiMagenta: colorRGBA(0xa2, 0x9b, 0xfe),
		ansiCyan:    colorRGBA(0x74, 0xb9, 0xff),
		ansiWhite:   colorRGBA(0xc8, 0xc8, 0xd0),

		ansiBrightBlack:   colorRGBA(0x2a, 0x2a, 0x4a),
		ansiBrightRed:     colorRGBA(0xe1, 0x70, 0x55),
		ansiBrightGreen:   colorRGBA(0x00, 0xb8, 0x94),
		ansiBrightYellow:  colorRGBA(0xfd, 0xcb, 0x6e),
		ansiBrightBlue:    colorRGBA(0xa2, 0x9b, 0xfe),
		ansiBrightMagenta: colorRGBA(0xff, 0x6e, 0xc7),
		ansiBrightCyan:    colorRGBA(0x74, 0xb9, 0xff),
		ansiBrightWhite:   colorRGBA(0xe8, 0xe8, 0xff),
	})

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

// HypercrushObsidiana returns the Hypercrush dark theme.
func HypercrushObsidiana() Styles {
	return CharmtonePantera()
}
