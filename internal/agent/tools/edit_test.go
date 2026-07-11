package tools

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"charm.land/fantasy"
	"github.com/stretchr/testify/require"
)

func TestEditToolCreatesMissingFileWithNewString(t *testing.T) {
	t.Parallel()

	workingDir := t.TempDir()
	ctx := context.WithValue(context.Background(), SessionIDContextKey, "test-session")

	tool := NewEditTool(&mockPermissionService{}, &mockHistoryService{}, mockFileTrackerService{}, workingDir)

	input, err := json.Marshal(EditParams{
		FilePath:  "new.txt",
		OldString: "does not exist yet",
		NewString: "created content\n",
	})
	require.NoError(t, err)

	resp, err := tool.Run(ctx, fantasy.ToolCall{
		ID:    "test-call",
		Name:  EditToolName,
		Input: string(input),
	})
	require.NoError(t, err)
	require.False(t, resp.IsError)

	b, err := os.ReadFile(filepath.Join(workingDir, "new.txt"))
	require.NoError(t, err)
	require.Equal(t, "created content\n", string(b))
}

func TestEditToolCreatesMissingFileParentDirs(t *testing.T) {
	t.Parallel()

	workingDir := t.TempDir()
	ctx := context.WithValue(context.Background(), SessionIDContextKey, "test-session")

	tool := NewEditTool(&mockPermissionService{}, &mockHistoryService{}, mockFileTrackerService{}, workingDir)

	input, err := json.Marshal(EditParams{
		FilePath:  filepath.Join("nested", "dir", "new.txt"),
		OldString: "does not exist yet",
		NewString: "created content\n",
	})
	require.NoError(t, err)

	resp, err := tool.Run(ctx, fantasy.ToolCall{
		ID:    "test-call",
		Name:  EditToolName,
		Input: string(input),
	})
	require.NoError(t, err)
	require.False(t, resp.IsError)

	b, err := os.ReadFile(filepath.Join(workingDir, "nested", "dir", "new.txt"))
	require.NoError(t, err)
	require.Equal(t, "created content\n", string(b))
}
