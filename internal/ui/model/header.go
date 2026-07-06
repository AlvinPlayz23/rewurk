package model

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/crush/internal/fsext"
	"github.com/charmbracelet/crush/internal/session"
	"github.com/charmbracelet/crush/internal/ui/common"
	"github.com/charmbracelet/crush/internal/ui/styles"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/ansi"
)

const (
	leftPadding   = 1
	rightPadding  = 1
	headerSpacing = 1
)

type header struct {
	logo        string
	compactLogo string

	com     *common.Common
	width   int
	compact bool
}

func newHeader(com *common.Common) *header {
	h := &header{
		com: com,
	}
	h.refresh()
	return h
}

func (h *header) refresh() {
	t := h.com.Styles
	name := "crush"
	h.compactLogo = styles.ApplyBoldForegroundGrad(
		t.Header.LogoGradCanvas,
		name,
		t.Header.LogoGradFromColor,
		t.Header.LogoGradToColor,
	) + " "
	h.width = 0
	h.logo = ""
}

func (h *header) drawHeader(
	scr uv.Screen,
	area uv.Rectangle,
	session *session.Session,
	compact bool,
	width int,
	hyperCredits *int,
) {
	t := h.com.Styles
	if width != h.width || compact != h.compact {
		h.logo = renderLogo(h.com.Styles, compact, h.com.IsHyper(), width)
	}

	h.width = width
	h.compact = compact

	if !compact || session == nil {
		uv.NewStyledString(h.logo).Draw(scr, area)
		return
	}

	if session.ID == "" {
		return
	}

	var b strings.Builder
	b.WriteString(h.compactLogo)

	availDetailWidth := width - leftPadding - rightPadding - lipgloss.Width(b.String()) - headerSpacing
	details := renderHeaderDetails(
		h.com,
		availDetailWidth,
		hyperCredits,
	)

	cwd := fsext.DirTrim(fsext.PrettyPath(h.com.Workspace.WorkingDir()), 4)
	cwdStyle := lipgloss.NewStyle().
		Foreground(t.Header.Charm.GetForeground()).
		Background(t.Dialog.ContentPanel.GetBackground()).
		Padding(0, 1).
		SetString(cwd)
	b.WriteString(cwdStyle.String())
	b.WriteString(" ")

	if remaining := availDetailWidth - lipgloss.Width(cwd) - 1; remaining > 0 {
		details = ansi.Truncate(details, remaining, "…")
	}
	b.WriteString(details)

	view := uv.NewStyledString(
		t.Header.Wrapper.Padding(0, rightPadding, 0, leftPadding).Render(b.String()),
	)
	view.Draw(scr, area)
}

func renderHeaderDetails(
	com *common.Common,
	availWidth int,
	hyperCredits *int,
) string {
	t := com.Styles

	var parts []string

	// Hypercredits.
	if com.IsHyper() && hyperCredits != nil {
		hc := t.Header.HypercreditIcon.Render(styles.HypercreditIcon) + " " + t.Header.Percentage.Render(common.FormatCredits(*hyperCredits))
		parts = append(parts, hc)
	}

	if len(parts) == 0 {
		return ""
	}

	sep := "  "
	result := sep + strings.Join(parts, sep)
	return ansi.Truncate(result, max(0, availWidth), "…")
}
