package chat

import (
	"cmp"
	"encoding/json"
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/crush/internal/agent/tools"
	"github.com/charmbracelet/crush/internal/message"
	"github.com/charmbracelet/crush/internal/ui/styles"
	"github.com/charmbracelet/x/ansi"
)

type legacyBashParams struct {
	Command         string `json:"command"`
	RunInBackground bool   `json:"run_in_background,omitempty"`
}

// -----------------------------------------------------------------------------
// Bash Tool
// -----------------------------------------------------------------------------

// BashToolMessageItem is a message item that represents a bash tool call.
type BashToolMessageItem struct {
	*baseToolMessageItem
}

var _ ToolMessageItem = (*BashToolMessageItem)(nil)

// NewBashToolMessageItem creates a new [BashToolMessageItem].
func NewBashToolMessageItem(
	sty *styles.Styles,
	toolCall message.ToolCall,
	result *message.ToolResult,
	canceled bool,
) ToolMessageItem {
	return newBaseToolMessageItem(sty, toolCall, result, &BashToolRenderContext{}, canceled)
}

// BashToolRenderContext renders bash tool messages.
type BashToolRenderContext struct{}

// RenderTool implements the [ToolRenderer] interface.
func (b *BashToolRenderContext) RenderTool(sty *styles.Styles, width int, opts *ToolRenderOpts) string {
	cappedWidth := cappedMessageWidth(width)
	if opts.IsPending() {
		return pendingTool(sty, "Bash", opts.Anim, opts.Compact)
	}

	var params legacyBashParams
	if err := json.Unmarshal([]byte(opts.ToolCall.Input), &params); err != nil {
		params.Command = "failed to parse command"
	}

	// Check if this is a background job.
	var meta tools.BashResponseMetadata
	if opts.HasResult() {
		_ = json.Unmarshal([]byte(opts.Result.Metadata), &meta)
	}

	if meta.Background {
		description := cmp.Or(meta.Description, params.Command)
		content := "Command: " + params.Command + "\n" + opts.Result.Content
		return renderJobTool(sty, opts, cappedWidth, "Start", meta.ShellID, description, content)
	}

	// Regular bash command.
	cmd := params.Command
	if !opts.ExpandedContent {
		cmd = strings.ReplaceAll(cmd, "\n", " ")
	}
	cmd = strings.ReplaceAll(cmd, "\t", "    ")
	toolParams := []string{cmd}
	if params.RunInBackground {
		toolParams = append(toolParams, "background", "true")
	}

	header := toolHeader(sty, opts.Status, "Bash", cappedWidth, opts, toolParams...)
	if opts.Compact {
		return header
	}

	if earlyState, ok := toolEarlyStateContent(sty, opts, cappedWidth); ok {
		return joinToolParts(sty, header, earlyState)
	}

	if !opts.HasResult() {
		return header
	}

	output := meta.Output
	if output == "" && opts.Result.Content != tools.BashNoOutput {
		output = opts.Result.Content
	}
	if output == "" {
		return header
	}

	bodyWidth := max(1, cappedWidth-toolBodyLeftPaddingTotal)
	body := sty.Tool.Body.Render(toolOutputPlainContent(sty, output, bodyWidth, opts.ExpandedContent))
	return joinToolParts(sty, header, body)
}

// renderJobTool renders a job-related tool with the common pattern:
// header → nested check → early state → body.
func renderJobTool(sty *styles.Styles, opts *ToolRenderOpts, width int, action, shellID, description, content string) string {
	header := jobHeader(sty, opts.Status, action, shellID, description, width)
	if opts.Compact {
		return header
	}

	if earlyState, ok := toolEarlyStateContent(sty, opts, width); ok {
		return joinToolParts(sty, header, earlyState)
	}

	if content == "" {
		return header
	}

	bodyWidth := max(1, width-toolBodyLeftPaddingTotal)
	body := sty.Tool.Body.Render(toolOutputPlainContent(sty, content, bodyWidth, opts.ExpandedContent))
	return joinToolParts(sty, header, body)
}

// jobHeader builds a header for job-related tools.
// Format: "● Job (Action) PID shellID description..."
func jobHeader(sty *styles.Styles, status ToolStatus, action, shellID, description string, width int) string {
	icon := toolIcon(sty, status)
	jobPart := sty.Tool.JobToolName.Render("Job")
	actionPart := sty.Tool.JobAction.Render("(" + action + ")")
	pidPart := sty.Tool.JobPID.Render("PID " + shellID)

	prefix := fmt.Sprintf("%s %s %s %s", icon, jobPart, actionPart, pidPart)

	if description == "" {
		return prefix
	}

	prefixWidth := lipgloss.Width(prefix)
	availableWidth := width - prefixWidth - 1
	if availableWidth < 10 {
		return prefix
	}

	truncatedDesc := ansi.Truncate(description, availableWidth, "…")
	return prefix + " " + sty.Tool.JobDescription.Render(truncatedDesc)
}

// joinToolParts joins a tool header with body output, adding the tree gutter
// for each body line so tool results visually connect to their header icon.
func joinToolParts(sty *styles.Styles, header, body string) string {
	return joinToolPartsPlain(header, renderToolBodyConnector(sty, body))
}

func joinToolPartsPlain(header, body string) string {
	if body == "" {
		return header
	}
	return strings.Join([]string{header, body}, "\n")
}

func renderToolBodyConnector(sty *styles.Styles, body string) string {
	if body == "" {
		return body
	}

	lines := strings.Split(body, "\n")
	for i, line := range lines {
		connector := "├─ "
		if i == len(lines)-1 {
			connector = "╰─ "
		}
		lines[i] = "  " + sty.Tool.Connector.Render(connector) + line
	}
	return strings.Join(lines, "\n")
}
