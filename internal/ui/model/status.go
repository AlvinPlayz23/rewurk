package model

import (
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/crush/internal/ui/common"
	"github.com/charmbracelet/crush/internal/ui/util"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/ansi"
)

const DefaultStatusTTL = 5 * time.Second

type StatusData struct {
	LSPCount     int
	MCPCount     int
	TokenUsed    int
	TokenMax     int
	ModelName    string
	ProviderName string
	IsBusy       bool
}

type Status struct {
	com      *common.Common
	hideHelp bool
	help     help.Model
	helpKm   help.KeyMap
	msg      util.InfoMsg
	data     StatusData
}

func NewStatus(com *common.Common, km help.KeyMap) *Status {
	s := new(Status)
	s.com = com
	s.help = help.New()
	s.help.Styles = com.Styles.Help
	s.helpKm = km
	return s
}

func (s *Status) SetInfoMsg(msg util.InfoMsg) {
	s.msg = msg
}

func (s *Status) SetData(data StatusData) {
	s.data = data
}

func (s *Status) ClearInfoMsg() {
	s.msg = util.InfoMsg{}
}

func (s *Status) SetWidth(width int) {
	helpStyle := s.com.Styles.Status.Help
	horizontalPadding := helpStyle.GetPaddingLeft() + helpStyle.GetPaddingRight()
	s.help.SetWidth(width - horizontalPadding)
}

func (s *Status) ShowingAll() bool {
	return s.help.ShowAll
}

func (s *Status) ToggleHelp() {
	s.help.ShowAll = !s.help.ShowAll
}

func (s *Status) SetHideHelp(hideHelp bool) {
	s.hideHelp = hideHelp
}

func (s *Status) Draw(scr uv.Screen, area uv.Rectangle) {
	t := s.com.Styles
	width := area.Dx()
	if width <= 0 {
		return
	}

	// Build left section: status dot + info items.
	var leftParts []string

	modelText := s.data.ModelName
	if modelText == "" {
		modelText = "No model"
	}
	if s.data.ProviderName != "" {
		modelText += " · " + s.data.ProviderName
	}
	leftParts = append(leftParts, modelText)

	// Token usage.
	if s.data.TokenMax > 0 {
		tokenStr := fmt.Sprintf("Tokens: %s / %s", formatTokenCount(s.data.TokenUsed), formatTokenCount(s.data.TokenMax))
		leftParts = append(leftParts, tokenStr)
	}

	// LSP count.
	if s.data.LSPCount > 0 {
		leftParts = append(leftParts, fmt.Sprintf("LSP: %d active", s.data.LSPCount))
	}

	// MCP count.
	if s.data.MCPCount > 0 {
		leftParts = append(leftParts, fmt.Sprintf("MCP: %d tools", s.data.MCPCount))
	}

	sep := "  │  "
	leftStr := strings.Join(leftParts, sep)
	baseStyle := lipgloss.NewStyle().Foreground(t.Status.Help.GetForeground())

	// Build right section: keyboard hints (from help model).
	var rightStr string
	if !s.hideHelp {
		rightStr = s.help.View(s.helpKm)
	}

	// Check if we should show notification instead.
	if !s.msg.IsEmpty() {
		var indStyle, msgStyle lipgloss.Style
		switch s.msg.Type {
		case util.InfoTypeError:
			indStyle = t.Status.ErrorIndicator
			msgStyle = t.Status.ErrorMessage
		case util.InfoTypeWarn:
			indStyle = t.Status.WarnIndicator
			msgStyle = t.Status.WarnMessage
		case util.InfoTypeUpdate:
			indStyle = t.Status.UpdateIndicator
			msgStyle = t.Status.UpdateMessage
		case util.InfoTypeInfo:
			indStyle = t.Status.InfoIndicator
			msgStyle = t.Status.InfoMessage
		case util.InfoTypeSuccess:
			indStyle = t.Status.SuccessIndicator
			msgStyle = t.Status.SuccessMessage
		}

		ind := indStyle.String()
		msg := strings.Join(strings.Split(s.msg.Msg, "\n"), " ")
		notifStr := ind + msgStyle.Render(msg)
		notifWidth := lipgloss.Width(notifStr)
		if notifWidth < width {
			notifStr += strings.Repeat(" ", width-notifWidth)
		}
		uv.NewStyledString(baseStyle.Render(notifStr)).Draw(scr, area)
		return
	}

	// Combine left and right sections.
	leftWidth := lipgloss.Width(leftStr)
	rightWidth := lipgloss.Width(rightStr)
	totalWidth := leftWidth + rightWidth

	var full string
	if totalWidth+4 <= width {
		padding := width - totalWidth
		full = leftStr + strings.Repeat(" ", padding) + rightStr
	} else if leftWidth+10 <= width {
		availRight := width - leftWidth - 2
		rightStr = ansi.Truncate(rightStr, max(0, availRight), "…")
		padding := width - leftWidth - lipgloss.Width(rightStr)
		full = leftStr + strings.Repeat(" ", max(0, padding)) + rightStr
	} else {
		full = ansi.Truncate(leftStr, width, "…")
	}

	uv.NewStyledString(baseStyle.Render(full)).Draw(scr, area)
}

func formatTokenCount(n int) string {
	switch {
	case n >= 1000000:
		return fmt.Sprintf("%.1fM", float64(n)/1000000)
	case n >= 1000:
		return fmt.Sprintf("%.1fk", float64(n)/1000)
	default:
		return fmt.Sprintf("%d", n)
	}
}

func clearInfoMsgCmd(ttl time.Duration) tea.Cmd {
	return tea.Tick(ttl, func(time.Time) tea.Msg {
		return util.ClearStatusMsg{}
	})
}
