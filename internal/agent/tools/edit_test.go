package tools

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"charm.land/fantasy"
	"github.com/charmbracelet/crush/internal/history"
	"github.com/charmbracelet/crush/internal/permission"
	"github.com/charmbracelet/crush/internal/pubsub"
	"github.com/stretchr/testify/require"
)

type mockPermissionService struct{}

func (m mockPermissionService) Subscribe(ctx context.Context) <-chan pubsub.Event[permission.PermissionRequest] {
	return make(chan pubsub.Event[permission.PermissionRequest])
}

func (m mockPermissionService) GrantPersistent(permission.PermissionRequest) bool { return true }

func (m mockPermissionService) Grant(permission.PermissionRequest) bool { return true }

func (m mockPermissionService) Deny(permission.PermissionRequest) bool { return true }

func (m mockPermissionService) Request(ctx context.Context, opts permission.CreatePermissionRequest) (bool, error) {
	return true, nil
}

func (m mockPermissionService) AutoApproveSession(sessionID string) {}

func (m mockPermissionService) SetSkipRequests(skip bool) {}

func (m mockPermissionService) SkipRequests() bool { return false }

func (m mockPermissionService) SubscribeNotifications(ctx context.Context) <-chan pubsub.Event[permission.PermissionNotification] {
	return make(chan pubsub.Event[permission.PermissionNotification])
}

type mockHistoryService struct{}

func (m mockHistoryService) Subscribe(ctx context.Context) <-chan pubsub.Event[history.File] {
	return make(chan pubsub.Event[history.File])
}

func (m mockHistoryService) Create(ctx context.Context, sessionID, path, content string) (history.File, error) {
	return history.File{SessionID: sessionID, Path: path, Content: content}, nil
}

func (m mockHistoryService) CreateVersion(ctx context.Context, sessionID, path, content string) (history.File, error) {
	return history.File{SessionID: sessionID, Path: path, Content: content}, nil
}

func (m mockHistoryService) Get(ctx context.Context, id string) (history.File, error) {
	return history.File{}, nil
}

func (m mockHistoryService) GetByPathAndSession(ctx context.Context, path, sessionID string) (history.File, error) {
	return history.File{}, nil
}

func (m mockHistoryService) ListBySession(ctx context.Context, sessionID string) ([]history.File, error) {
	return nil, nil
}

func (m mockHistoryService) ListLatestSessionFiles(ctx context.Context, sessionID string) ([]history.File, error) {
	return nil, nil
}

func (m mockHistoryService) Delete(ctx context.Context, id string) error { return nil }

func (m mockHistoryService) DeleteSessionFiles(ctx context.Context, sessionID string) error {
	return nil
}

type mockFileTrackerService struct{}

func (m mockFileTrackerService) RecordRead(ctx context.Context, sessionID, path string) {}

func (m mockFileTrackerService) LastReadTime(ctx context.Context, sessionID, path string) time.Time {
	return time.Now()
}

func (m mockFileTrackerService) ListReadFiles(ctx context.Context, sessionID string) ([]string, error) {
	return nil, nil
}

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

func TestEditToolAppliesMultipleEditsAtomically(t *testing.T) {
	t.Parallel()

	workingDir := t.TempDir()
	filePath := filepath.Join(workingDir, "main.go")
	err := os.WriteFile(filePath, []byte("package main\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n"), 0o644)
	require.NoError(t, err)

	ctx := context.WithValue(context.Background(), SessionIDContextKey, "test-session")
	tool := NewEditTool(&mockPermissionService{}, &mockHistoryService{}, mockFileTrackerService{}, workingDir)

	input, err := json.Marshal(EditParams{
		FilePath: "main.go",
		Edits: []EditOperation{
			{
				OldString: "func main() {\n",
				NewString: "func main() {\n\t// Greeting.\n",
			},
			{
				OldString: "Hello, World!",
				NewString: "Hello, Crush!",
			},
		},
	})
	require.NoError(t, err)

	resp, err := tool.Run(ctx, fantasy.ToolCall{
		ID:    "test-call",
		Name:  EditToolName,
		Input: string(input),
	})
	require.NoError(t, err)
	require.False(t, resp.IsError)

	b, err := os.ReadFile(filePath)
	require.NoError(t, err)
	require.Equal(t, "package main\n\nfunc main() {\n\t// Greeting.\n\tfmt.Println(\"Hello, Crush!\")\n}\n", string(b))
}

func TestEditToolEditsMatchOriginalContent(t *testing.T) {
	t.Parallel()

	newContent, err := applyEditOperations("foo\nbar\nbaz\n", []EditOperation{
		{OldString: "foo\n", NewString: "foo bar\n"},
		{OldString: "bar\n", NewString: "BAR\n"},
	})
	require.NoError(t, err)
	require.Equal(t, "foo bar\nBAR\nbaz\n", newContent)
}

func TestEditToolEditsRejectOverlap(t *testing.T) {
	t.Parallel()

	_, err := applyEditOperations("one\ntwo\nthree\n", []EditOperation{
		{OldString: "one\ntwo\n", NewString: "ONE\nTWO\n"},
		{OldString: "two\nthree\n", NewString: "TWO\nTHREE\n"},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "overlap")
}

func TestEditToolEditsDoNotPartiallyApply(t *testing.T) {
	t.Parallel()

	workingDir := t.TempDir()
	filePath := filepath.Join(workingDir, "test.txt")
	original := "alpha\nbeta\ngamma\n"
	err := os.WriteFile(filePath, []byte(original), 0o644)
	require.NoError(t, err)

	ctx := context.WithValue(context.Background(), SessionIDContextKey, "test-session")
	tool := NewEditTool(&mockPermissionService{}, &mockHistoryService{}, mockFileTrackerService{}, workingDir)

	input, err := json.Marshal(EditParams{
		FilePath: "test.txt",
		Edits: []EditOperation{
			{OldString: "alpha\n", NewString: "ALPHA\n"},
			{OldString: "missing\n", NewString: "MISSING\n"},
		},
	})
	require.NoError(t, err)

	resp, err := tool.Run(ctx, fantasy.ToolCall{
		ID:    "test-call",
		Name:  EditToolName,
		Input: string(input),
	})
	require.NoError(t, err)
	require.True(t, resp.IsError)

	b, err := os.ReadFile(filePath)
	require.NoError(t, err)
	require.Equal(t, original, string(b))
}

func TestEditToolEditsPreserveCRLF(t *testing.T) {
	t.Parallel()

	workingDir := t.TempDir()
	filePath := filepath.Join(workingDir, "test.txt")
	err := os.WriteFile(filePath, []byte("first\r\nsecond\r\nthird\r\n"), 0o644)
	require.NoError(t, err)

	ctx := context.WithValue(context.Background(), SessionIDContextKey, "test-session")
	tool := NewEditTool(&mockPermissionService{}, &mockHistoryService{}, mockFileTrackerService{}, workingDir)

	input, err := json.Marshal(EditParams{
		FilePath: "test.txt",
		Edits: []EditOperation{
			{OldString: "second\n", NewString: "SECOND\n"},
		},
	})
	require.NoError(t, err)

	resp, err := tool.Run(ctx, fantasy.ToolCall{
		ID:    "test-call",
		Name:  EditToolName,
		Input: string(input),
	})
	require.NoError(t, err)
	require.False(t, resp.IsError)

	b, err := os.ReadFile(filePath)
	require.NoError(t, err)
	require.Equal(t, "first\r\nSECOND\r\nthird\r\n", string(b))
}
