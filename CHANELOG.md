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

Changed files:

- `internal/ui/styles/quickstyle.go`
- `internal/ui/styles/styles.go`
- `internal/ui/model/status.go`
- `internal/ui/model/ui.go`
- `internal/ui/model/header.go`
- `internal/ui/model/keys.go`
- `internal/ui/chat/assistant.go`

Validation:

- Ran `gofmt` on the touched Go files.
- Ran `go test ./internal/ui/chat ./internal/ui/model ./internal/ui/styles` successfully.
