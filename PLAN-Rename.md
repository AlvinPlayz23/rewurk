plan for renaming all View-prefixed internal types to Read-prefixed types:

## Renaming Map

| Old Name | New Name |
|----------|----------|
| `ViewParams` | `ReadParams` |
| `ViewPermissionsParams` | `ReadPermissionsParams` |
| `ViewResourceType` | `ReadResourceType` |
| `ViewResourceUnset` | `ReadResourceUnset` |
| `ViewResourceSkill` | `ReadResourceSkill` |
| `ViewResponseMetadata` | `ReadResponseMetadata` |
| `ViewToolName` | `ReadToolName` |
| `NewViewTool` | `NewReadTool` |
| `NewViewToolMessageItem` | `NewReadToolMessageItem` |
| `ViewToolMessageItem` | `ReadToolMessageItem` |
| `ViewToolRenderContext` | `ReadToolRenderContext` |
| `viewDescription` | `readDescription` |
| `viewDescriptionTmpl` | `readDescriptionTmpl` |
| `viewDescriptionTpl` | `readDescriptionTpl` |
| `viewDescriptionData` | `readDescriptionData` |

## Files Affected (12 files)

| File | Changes |
|------|---------|
| `internal/agent/tools/view.go` | Rename all View* types/functions, rename file to `read.go` |
| `internal/agent/tools/view_test.go` | Rename all View* references, rename file to `read_test.go` |
| `internal/proto/tools.go` | Rename `ViewParams`, `ViewPermissionsParams`, `ViewResponseMetadata`, `ViewToolName` |
| `internal/proto/permission.go` | Rename `ViewPermissionsParams`, `ViewToolName` |
| `internal/proto/permission_test.go` | Rename `ViewPermissionsParams`, `ViewToolName` |
| `internal/ui/chat/file.go` | Rename `ViewParams`, `ViewResourceSkill`, `ViewResponseMetadata`, `ViewToolMessageItem`, `NewViewToolMessageItem`, `ViewToolRenderContext` |
| `internal/ui/chat/tools.go` | Rename `ViewParams`, `ViewToolName`, `ViewResponseMetadata`, `NewViewToolMessageItem` |
| `internal/ui/dialog/permissions.go` | Rename `ViewPermissionsParams`, `ViewToolName` |
| `internal/agent/coordinator.go` | Rename `NewViewTool` |
| `internal/agent/common_test.go` | Rename `NewViewTool` |
| `internal/agent/agentic_fetch_tool.go` | Rename `NewViewTool` |
| `internal/agent/agent_test.go` | Rename `ViewToolName` |
| `internal/cmd/session.go` | Rename `ViewResponseMetadata`, `ViewResourceSkill` |

## Steps

1. Rename `internal/agent/tools/view.go` → `read.go`
2. Rename `internal/agent/tools/view_test.go` → `read_test.go`
3. Apply all type renames across all 12 files
4. Run `go build ./...` and `go test ./internal/config/... ./internal/proto/... ./internal/ui/chat/... ./internal/ui/dialog/... ./internal/permission/...` to verify