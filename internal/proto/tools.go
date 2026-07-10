package proto

// The wire schema for per-tool permission parameters is owned by the
// tool itself, not duplicated here. We alias the canonical types so
// there is exactly one source of truth and so values survive a
// round-trip across the client/server boundary as the same Go type
// the UI asserts on.
import "github.com/charmbracelet/crush/internal/agent/tools"

// ToolResponseType represents the type of tool response.
type ToolResponseType string

const (
	ToolResponseTypeText  ToolResponseType = "text"
	ToolResponseTypeImage ToolResponseType = "image"
)

// ToolResponse represents a response from a tool.
type ToolResponse struct {
	Type     ToolResponseType `json:"type"`
	Content  string           `json:"content"`
	Metadata string           `json:"metadata,omitempty"`
	IsError  bool             `json:"is_error"`
}

const BashToolName = "bash"

// BashParams represents the parameters for the bash tool.
type BashParams struct {
	Command string `json:"command"`
	Timeout int    `json:"timeout"`
}

// BashPermissionsParams represents the permission parameters for the bash tool.
type BashPermissionsParams = tools.BashPermissionsParams

// BashResponseMetadata represents the metadata for a bash tool response.
type BashResponseMetadata struct {
	StartTime        int64  `json:"start_time"`
	EndTime          int64  `json:"end_time"`
	Output           string `json:"output"`
	WorkingDirectory string `json:"working_directory"`
}

const EditToolName = "edit"

// EditParams represents the parameters for the edit tool.
type EditParams struct {
	FilePath   string `json:"file_path"`
	OldString  string `json:"old_string"`
	NewString  string `json:"new_string"`
	ReplaceAll bool   `json:"replace_all,omitempty"`
}

// EditPermissionsParams represents the permission parameters for the edit tool.
type EditPermissionsParams = tools.EditPermissionsParams

// EditResponseMetadata represents the metadata for an edit tool response.
type EditResponseMetadata struct {
	Additions  int    `json:"additions"`
	Removals   int    `json:"removals"`
	OldContent string `json:"old_content,omitempty"`
	NewContent string `json:"new_content,omitempty"`
}

const GlobToolName = "glob"

// GlobParams represents the parameters for the glob tool.
type GlobParams struct {
	Pattern string `json:"pattern"`
	Path    string `json:"path"`
}

// GlobResponseMetadata represents the metadata for a glob tool response.
type GlobResponseMetadata struct {
	NumberOfFiles int  `json:"number_of_files"`
	Truncated     bool `json:"truncated"`
}

const GrepToolName = "grep"

// GrepParams represents the parameters for the grep tool.
type GrepParams struct {
	Pattern     string `json:"pattern"`
	Path        string `json:"path"`
	Include     string `json:"include"`
	LiteralText bool   `json:"literal_text"`
}

// GrepResponseMetadata represents the metadata for a grep tool response.
type GrepResponseMetadata struct {
	NumberOfMatches int  `json:"number_of_matches"`
	Truncated       bool `json:"truncated"`
}

const MultiEditToolName = "multiedit"

// MultiEditOperation represents a single edit operation in a multi-edit.
type MultiEditOperation struct {
	OldString  string `json:"old_string"`
	NewString  string `json:"new_string"`
	ReplaceAll bool   `json:"replace_all,omitempty"`
}

// MultiEditParams represents the parameters for the multi-edit tool.
type MultiEditParams struct {
	FilePath string               `json:"file_path"`
	Edits    []MultiEditOperation `json:"edits"`
}

// MultiEditPermissionsParams represents the permission parameters for the multi-edit tool.
type MultiEditPermissionsParams = tools.MultiEditPermissionsParams

// MultiEditResponseMetadata represents the metadata for a multi-edit tool response.
type MultiEditResponseMetadata struct {
	Additions    int    `json:"additions"`
	Removals     int    `json:"removals"`
	OldContent   string `json:"old_content,omitempty"`
	NewContent   string `json:"new_content,omitempty"`
	EditsApplied int    `json:"edits_applied"`
}

const ReadToolName = "read"

// ReadParams represents the parameters for the read tool.
type ReadParams struct {
	FilePath string `json:"file_path"`
	Offset   int    `json:"offset"`
	Limit    int    `json:"limit"`
}

// ReadPermissionsParams represents the permission parameters for the read tool.
type ReadPermissionsParams = tools.ReadPermissionsParams

// ReadResponseMetadata represents the metadata for a read tool response.
type ReadResponseMetadata struct {
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
}

const WriteToolName = "write"

// WriteParams represents the parameters for the write tool.
type WriteParams struct {
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
}

// WritePermissionsParams represents the permission parameters for the write tool.
type WritePermissionsParams = tools.WritePermissionsParams

// WriteResponseMetadata represents the metadata for a write tool response.
type WriteResponseMetadata struct {
	Diff      string `json:"diff"`
	Additions int    `json:"additions"`
	Removals  int    `json:"removals"`
}
