# Changelog

## Edit Tool Creates Missing Files

- Updated the `edit` tool so replacement edits against missing file paths create the file with `new_string` as the full content.
- Reused the existing edit file-creation path so parent directory creation, permissions, history tracking, and file tracking remain consistent.
- Removed the standalone `write` tool from the tool registry, agent harness, protocol aliases, permission decoding, UI renderers, and docs.
- Deleted the `write` tool implementation, description, tests, and agent test fixture now that file creation is handled by `edit`.
- Updated the edit tool description to document missing-file creation behavior.
- Added tests for creating missing files and nested parent directories through `edit`.

Validation:

- Ran `go test ./internal/agent/tools` successfully.

## Simplified Built-In Tool Set

- Removed the `fetch`, `job_kill`, and `job_output` tools.
- Removed the `agentic_fetch` sub-agent and its private `web_fetch`, `web_search`, `glob`, `grep`, and scoped `read` tool set.
- Promoted `web_search` to an independent, permission-checked top-level tool.
- Removed background and automatic background execution from `bash` so commands cannot become unmanageable after removal of the job tools.
- Removed the obsolete fetch prompts, protocol types, tests, HTML conversion dependencies, and related configuration entries.
- Retained legacy chat rendering for removed tools so existing saved sessions remain readable.
- Updated agent instructions, hook documentation, permission handling, and tool registry tests for the reduced tool set.

Deleted files:

- `internal/agent/agentic_fetch_tool.go`
- `internal/agent/templates/agentic_fetch.md`
- `internal/agent/templates/agentic_fetch_prompt.md.tpl`
- `internal/agent/tools/fetch.go`
- `internal/agent/tools/fetch.md.tpl`
- `internal/agent/tools/fetch_helpers.go`
- `internal/agent/tools/fetch_types.go`
- `internal/agent/tools/job_kill.go`
- `internal/agent/tools/job_kill.md`
- `internal/agent/tools/job_output.go`
- `internal/agent/tools/job_output.md`
- `internal/agent/tools/job_test.go`
- `internal/agent/tools/web_fetch.go`
- `internal/agent/tools/web_fetch.md.tpl`

Validation:

- Ran `go build .` successfully.
- Ran the agent, tool, protocol, and relevant UI test suites successfully.
- Ran focused tool-registry configuration tests successfully.

## Removed IPC Socket And Server Mode

- Removed `crush server` subcommand and the `--host` persistent flag.
- Removed `CRUSH_CLIENT_SERVER` environment variable gate — all commands now run in local in-process mode only.
- Removed socket auto-start, stale-socket detection, version-mismatch restart, and `spawnAndWaitReady` logic from the CLI.
- Removed the server-streaming `runNonInteractive` event loop from `crush run` — non-interactive runs now use the built-in local `App.RunNonInteractive`.
- Converted `login` and `logout` commands to use local workspace config mutations instead of the server client.
- Removed `clientserverrace` regression test package.

Changed files:

- `internal/cmd/root.go`
- `internal/cmd/run.go`
- `internal/cmd/login.go`
- `internal/cmd/logout.go`

Deleted files:

- `internal/cmd/server.go`
- `internal/cmd/server_windows.go`
- `internal/cmd/server_other.go`
- `internal/cmd/run_stream_test.go`
- `internal/cmd/clientserverrace/race_test.go`

Validation:

- Ran `go build ./...` successfully.
- Ran `go run . --help` to verify `server` subcommand and `--host` flag are gone.

## Removed Native Desktop Notification Backend

- Removed the `beeep` native desktop notification backend and its dependency (`github.com/gen2brain/beeep`).
- Removed `NativeBackend` which sent OS-level popup notifications via `beeep`.
- Removed `runtime.GOOS` auto-detection path that selected native notifications for non-macOS local sessions.
- Local sessions now use OSC terminal-based notifications instead of native popups.
- Removed `"native"` from the notification style picker dialog and schema enum.
- Updated config schema, README docs, and skill docs to reflect remaining options (`auto`, `osc`, `bell`, `disabled`).
- Removed deprecated `disable_notifications` reference in favor of `notification_style`.

Changed files:

- `internal/ui/notification/native.go` (deleted)
- `internal/ui/notification/notification.go`
- `internal/ui/notification/notification_test.go`
- `internal/ui/notification/icon_darwin.go`
- `internal/ui/model/ui.go`
- `internal/ui/dialog/notifications.go`
- `internal/config/config.go`
- `schema.json`
- `README.md`
- `internal/skills/builtin/crush-config/SKILL.md`
- `go.mod`
- `go.sum`

Deleted files:

- `internal/ui/notification/native.go`

Validation:

- Ran `go mod tidy` successfully.
- Ran `go test ./internal/ui/notification ./internal/ui/dialog ./internal/ui/model` successfully.

## Removed Built-In Tools

- Removed the `sourcegraph` tool and its Sourcegraph API client.
- Removed the `download` tool for fetching remote files to disk.
- Removed the `ls` tool for directory listing.
- Removed the `lsp_restart` tool for restarting LSP servers.
- Removed the `lsp_diagnostics` tool (diagnostic helpers retained for edit/read).
- Removed the `crush_info` tool for displaying Crush configuration.
- Removed the `crush_logs` tool for reading Crush log files.
- Removed the corresponding UI chat renderers for each removed tool.
- Removed LS and download permission dialog panels and proto unmarshal paths.
- Updated agent tool registry, config defaults, tests, docs, and schema examples.

Changed files:

- `internal/agent/coordinator.go`
- `internal/agent/agentic_fetch_tool.go`
- `internal/agent/common_test.go`
- `internal/agent/agent_test.go`
- `internal/agent/tools/diagnostics.go`
- `internal/config/config.go`
- `internal/config/load_test.go`
- `internal/ui/chat/tools.go`
- `internal/ui/chat/file.go`
- `internal/ui/chat/search.go`
- `internal/ui/dialog/permissions.go`
- `internal/proto/tools.go`
- `internal/proto/permission.go`
- `internal/proto/permission_test.go`
- `internal/skills/builtin/crush-config/SKILL.md`
- `internal/ui/AGENTS.md`
- `README.md`

Deleted files:

- `internal/agent/tools/sourcegraph.go`
- `internal/agent/tools/sourcegraph.md.tpl`
- `internal/agent/tools/download.go`
- `internal/agent/tools/download.md.tpl`
- `internal/agent/tools/ls.go`
- `internal/agent/tools/ls.md.tpl`
- `internal/agent/tools/lsp_restart.go`
- `internal/agent/tools/lsp_restart.md`
- `internal/agent/tools/diagnostics.md`
- `internal/agent/tools/crush_info.go`
- `internal/agent/tools/crush_info.md`
- `internal/agent/tools/crush_info_test.go`
- `internal/agent/tools/crush_logs.go`
- `internal/agent/tools/crush_logs.md.tpl`
- `internal/agent/tools/crush_logs_test.go`
- `internal/ui/chat/lsp_restart.go`
- `internal/ui/chat/diagnostics.go`

Validation:

- Ran `go build ./...` successfully.
- Ran `go test ./internal/config ./internal/proto ./internal/ui/chat ./internal/ui/dialog ./internal/permission` successfully.

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
