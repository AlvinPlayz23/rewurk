package dialog

import (
	"slices"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/crush/internal/config"
	"github.com/charmbracelet/crush/internal/ui/common"
	"github.com/charmbracelet/crush/internal/ui/list"
	"github.com/charmbracelet/crush/internal/ui/styles"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/sahilm/fuzzy"
)

const (
	// ToolsID is the identifier for the extra tools dialog.
	ToolsID              = "tools"
	toolsDialogMaxWidth  = 54
	toolsDialogMaxHeight = 9
)

type toolOption struct {
	Name        string
	Title       string
	Description string
}

var extraToolOptions = []toolOption{
	{Name: "glob", Title: "Glob", Description: "Find files by name pattern"},
	{Name: "grep", Title: "Grep", Description: "Search file contents"},
}

// Tools represents a dialog for enabling and disabling extra tools.
type Tools struct {
	com  *common.Common
	help help.Model
	list *list.List

	keyMap struct {
		Toggle   key.Binding
		Next     key.Binding
		Previous key.Binding
		UpDown   key.Binding
		Close    key.Binding
	}
}

// ToolItem represents an extra tool list item.
type ToolItem struct {
	*list.Versioned
	tool    toolOption
	enabled bool
	t       *styles.Styles
	cache   map[int]string
	focused bool
}

var (
	_ Dialog   = (*Tools)(nil)
	_ ListItem = (*ToolItem)(nil)
)

// NewTools creates an extra tools dialog.
func NewTools(com *common.Common) *Tools {
	t := &Tools{com: com}

	h := help.New()
	h.Styles = com.Styles.DialogHelpStyles()
	t.help = h

	t.list = list.NewList()
	t.list.Focus()
	t.list.RegisterRenderCallback(list.FocusedRenderCallback(t.list))

	t.keyMap.Toggle = key.NewBinding(
		key.WithKeys("enter", "ctrl+y", "space"),
		key.WithHelp("enter", "toggle"),
	)
	t.keyMap.Next = key.NewBinding(
		key.WithKeys("down", "ctrl+n"),
		key.WithHelp("down", "next item"),
	)
	t.keyMap.Previous = key.NewBinding(
		key.WithKeys("up", "ctrl+p"),
		key.WithHelp("up", "previous item"),
	)
	t.keyMap.UpDown = key.NewBinding(
		key.WithKeys("up", "down"),
		key.WithHelp("up/down", "choose"),
	)
	t.keyMap.Close = CloseKey

	t.Refresh()
	return t
}

// ID implements Dialog.
func (*Tools) ID() string { return ToolsID }

// HandleMsg implements Dialog.
func (t *Tools) HandleMsg(msg tea.Msg) Action {
	keyMsg, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return nil
	}

	switch {
	case key.Matches(keyMsg, t.keyMap.Close):
		return ActionClose{}
	case key.Matches(keyMsg, t.keyMap.Previous):
		if t.list.IsSelectedFirst() {
			t.list.SelectLast()
		} else {
			t.list.SelectPrev()
		}
	case key.Matches(keyMsg, t.keyMap.Next):
		if t.list.IsSelectedLast() {
			t.list.SelectFirst()
		} else {
			t.list.SelectNext()
		}
	case key.Matches(keyMsg, t.keyMap.Toggle):
		item, ok := t.list.SelectedItem().(*ToolItem)
		if ok {
			return ActionToggleTool{Name: item.tool.Name}
		}
	}
	return nil
}

// Draw implements Dialog.
func (t *Tools) Draw(scr uv.Screen, area uv.Rectangle) *tea.Cursor {
	sty := t.com.Styles
	width := max(0, min(toolsDialogMaxWidth, area.Dx()))
	height := max(0, min(toolsDialogMaxHeight, area.Dy()))
	innerWidth := width - sty.Dialog.View.GetHorizontalFrameSize()
	heightOffset := sty.Dialog.Title.GetVerticalFrameSize() + titleContentHeight +
		sty.Dialog.HelpView.GetVerticalFrameSize() + sty.Dialog.View.GetVerticalFrameSize()

	t.list.SetSize(innerWidth, height-heightOffset)
	t.help.SetWidth(innerWidth)

	rc := NewRenderContext(sty, width)
	rc.Title = "Extra Tools"
	rc.AddPart(sty.Dialog.List.Height(t.list.Height()).Render(t.list.Render()))
	rc.Help = t.help.View(t)
	DrawCenterCursor(scr, area, rc.Render(), nil)
	return nil
}

// ShortHelp implements help.KeyMap.
func (t *Tools) ShortHelp() []key.Binding {
	return []key.Binding{t.keyMap.UpDown, t.keyMap.Toggle, t.keyMap.Close}
}

// FullHelp implements help.KeyMap.
func (t *Tools) FullHelp() [][]key.Binding {
	return [][]key.Binding{{t.keyMap.Toggle, t.keyMap.Next, t.keyMap.Previous, t.keyMap.Close}}
}

// Refresh rebuilds the tool list from the current configuration.
func (t *Tools) Refresh() {
	disabled := config.ExtraToolNames
	if cfg := t.com.Config(); cfg != nil && cfg.Options != nil && cfg.Options.DisabledTools != nil {
		disabled = cfg.Options.DisabledTools
	}

	selected := t.list.Selected()
	items := make([]list.Item, 0, len(extraToolOptions))
	for _, tool := range extraToolOptions {
		items = append(items, &ToolItem{
			Versioned: list.NewVersioned(),
			tool:      tool,
			enabled:   !slices.Contains(disabled, tool.Name),
			t:         t.com.Styles,
		})
	}
	t.list.SetItems(items...)
	if selected < 0 {
		selected = 0
	}
	t.list.SetSelected(min(selected, len(items)-1))
}

// Finished implements list.Item.
func (*ToolItem) Finished() bool { return true }

// Filter implements ListItem.
func (t *ToolItem) Filter() string { return t.tool.Title + " " + t.tool.Description }

// ID implements ListItem.
func (t *ToolItem) ID() string { return t.tool.Name }

// SetFocused implements ListItem.
func (t *ToolItem) SetFocused(focused bool) {
	if t.focused == focused {
		return
	}
	t.focused = focused
	t.cache = nil
	t.Bump()
}

// SetMatch implements ListItem.
func (*ToolItem) SetMatch(fuzzy.Match) {}

// Render implements list.Item.
func (t *ToolItem) Render(width int) string {
	status := "disabled"
	if t.enabled {
		status = "enabled"
	}
	itemStyles := ListItemStyles{
		ItemBlurred:     t.t.Dialog.NormalItem,
		ItemFocused:     t.t.Dialog.SelectedItem,
		InfoTextBlurred: t.t.Dialog.ListItem.InfoBlurred,
		InfoTextFocused: t.t.Dialog.ListItem.InfoFocused,
	}
	title := t.tool.Title + " - " + t.tool.Description
	return renderItem(itemStyles, title, status, t.focused, width, t.cache, nil)
}
