# Changelog

## UI Message And Footer Cleanup

- Removed the visible `C` and `U` markers from assistant and user chat messages.
- Kept the chat messages slightly indented by replacing the removed markers with blank two-column prefixes.
- Removed the visible `:::` prompt continuation markers from normal, yolo, and bang prompt modes while preserving prompt alignment.
- Replaced the bottom-left `Ready` status text with the active model name and provider.
- Removed the `ctrl+g more` and `ctrl+l models` hints from the bottom help/footer display while leaving the shortcuts active.
- Removed the top compact-header model/provider badge and token percentage because that information now appears in the bottom status area.
- Removed the `ctrl+d` compact session-details shortcut and the compact session-details overlay screen.
- Replaced the assistant random/scrambled loading animation with rotating plain status phrases such as `Thinking...`, `Streaming...`, `Ruminating...`, and `Hold on, let me think more...`.
- Added a subtle rounded box (`UserBox`) around user messages so they are visually grouped and separated from assistant/tool output.
- Aligned the left edge of the user-message box with assistant messages by removing the extra user prefix indentation.
- Removed the completed-assistant footer line that showed model, provider, duration, and a horizontal separator after each request.
- Removed the user-message box background so it no longer looks like the thinking-block highlight on wrapped or multi-line user messages.
- Limited live thinking/reasoning previews to roughly two lines plus `...` while the assistant is streaming.
- Added a `/` command, `Toggle Thinking Blocks`, to show or hide completed thinking blocks after a response finishes.
- Hidden completed thinking blocks now render a compact placeholder instead of the full reasoning body.

Changed files:

- `internal/ui/styles/quickstyle.go`
- `internal/ui/styles/styles.go`
- `internal/ui/model/status.go`
- `internal/ui/model/ui.go`
- `internal/ui/model/header.go`
- `internal/ui/model/keys.go`
- `internal/ui/chat/assistant.go`
- `internal/ui/chat/user.go`
- `internal/ui/chat/messages.go`
- `internal/ui/model/chat.go`
- `internal/ui/dialog/actions.go`
- `internal/ui/dialog/commands.go`

Validation:

- Ran `gofmt` on the touched Go files.
- Ran `go test ./internal/ui/chat ./internal/ui/model ./internal/ui/styles` successfully.
