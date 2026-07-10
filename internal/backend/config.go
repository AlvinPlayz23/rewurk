package backend

import (
	"context"

	"github.com/charmbracelet/crush/internal/agent"
	"github.com/charmbracelet/crush/internal/config"
	"github.com/charmbracelet/crush/internal/oauth"
	"github.com/charmbracelet/crush/internal/proto"
	"github.com/charmbracelet/crush/internal/pubsub"
	"github.com/charmbracelet/crush/internal/skills"
)

// publishConfigChanged publishes a ConfigChanged event on the workspace's
// event broker so all subscribers (e.g. remote clients) refresh their
// cached config snapshot.
func publishConfigChanged(ws *Workspace) {
	if ws == nil || ws.App == nil {
		return
	}
	ws.SendEvent(pubsub.Event[proto.ConfigChanged]{
		Type:    pubsub.UpdatedEvent,
		Payload: proto.ConfigChanged{WorkspaceID: ws.ID},
	})
}

// SetConfigField sets a key/value pair in the config file for the
// given scope.
func (b *Backend) SetConfigField(workspaceID string, scope config.Scope, key string, value any) error {
	ws, err := b.GetWorkspace(workspaceID)
	if err != nil {
		return err
	}
	if err := ws.Cfg.SetConfigField(scope, key, value); err != nil {
		return err
	}
	publishConfigChanged(ws)
	return nil
}

// RemoveConfigField removes a key from the config file for the given
// scope.
func (b *Backend) RemoveConfigField(workspaceID string, scope config.Scope, key string) error {
	ws, err := b.GetWorkspace(workspaceID)
	if err != nil {
		return err
	}
	if err := ws.Cfg.RemoveConfigField(scope, key); err != nil {
		return err
	}
	publishConfigChanged(ws)
	return nil
}

// UpdatePreferredModel updates the preferred model for the given type
// and persists it to the config file at the given scope.
func (b *Backend) UpdatePreferredModel(workspaceID string, scope config.Scope, modelType config.SelectedModelType, model config.SelectedModel) error {
	ws, err := b.GetWorkspace(workspaceID)
	if err != nil {
		return err
	}
	if err := ws.Cfg.UpdatePreferredModel(scope, modelType, model); err != nil {
		return err
	}
	publishConfigChanged(ws)
	return nil
}

// SetCompactMode sets the compact mode setting and persists it.
func (b *Backend) SetCompactMode(workspaceID string, scope config.Scope, enabled bool) error {
	ws, err := b.GetWorkspace(workspaceID)
	if err != nil {
		return err
	}
	if err := ws.Cfg.SetCompactMode(scope, enabled); err != nil {
		return err
	}
	publishConfigChanged(ws)
	return nil
}

// SetProviderAPIKey sets the API key for a provider and persists it.
func (b *Backend) SetProviderAPIKey(workspaceID string, scope config.Scope, providerID string, apiKey any) error {
	ws, err := b.GetWorkspace(workspaceID)
	if err != nil {
		return err
	}
	if err := ws.Cfg.SetProviderAPIKey(scope, providerID, apiKey); err != nil {
		return err
	}
	publishConfigChanged(ws)
	return nil
}

// ImportCopilot attempts to import a GitHub Copilot token from disk.
func (b *Backend) ImportCopilot(workspaceID string) (*oauth.Token, bool, error) {
	ws, err := b.GetWorkspace(workspaceID)
	if err != nil {
		return nil, false, err
	}
	token, ok := ws.Cfg.ImportCopilot()
	if ok {
		publishConfigChanged(ws)
	}
	return token, ok, nil
}

// RefreshOAuthToken refreshes the OAuth token for a provider.
func (b *Backend) RefreshOAuthToken(ctx context.Context, workspaceID string, scope config.Scope, providerID string) error {
	ws, err := b.GetWorkspace(workspaceID)
	if err != nil {
		return err
	}
	if err := ws.Cfg.RefreshOAuthToken(ctx, scope, providerID); err != nil {
		return err
	}
	publishConfigChanged(ws)
	return nil
}

// ProjectNeedsInitialization checks whether the project in this
// workspace needs initialization.
func (b *Backend) ProjectNeedsInitialization(workspaceID string) (bool, error) {
	ws, err := b.GetWorkspace(workspaceID)
	if err != nil {
		return false, err
	}
	return config.ProjectNeedsInitialization(ws.Cfg)
}

// MarkProjectInitialized marks the project as initialized.
func (b *Backend) MarkProjectInitialized(workspaceID string) error {
	ws, err := b.GetWorkspace(workspaceID)
	if err != nil {
		return err
	}
	if err := config.MarkProjectInitialized(ws.Cfg); err != nil {
		return err
	}
	publishConfigChanged(ws)
	return nil
}

// InitializePrompt builds the initialization prompt for the workspace.
func (b *Backend) InitializePrompt(workspaceID string) (string, error) {
	ws, err := b.GetWorkspace(workspaceID)
	if err != nil {
		return "", err
	}
	return agent.InitializePrompt(ws.Cfg)
}

// ReadSkill reads a skill's content by ID.
func (b *Backend) ReadSkill(ctx context.Context, workspaceID, skillID string) ([]byte, proto.SkillReadResult, error) {
	ws, err := b.GetWorkspace(workspaceID)
	if err != nil {
		return nil, proto.SkillReadResult{}, err
	}

	mgr := ws.Skills
	content, result, err := skills.ReadContent(
		mgr.ActiveSkills(), mgr.ResolvedPaths(), mgr.WorkingDir(), skillID,
	)
	if err != nil {
		return nil, proto.SkillReadResult{}, err
	}
	return content, proto.SkillReadResult{
		Name:        result.Name,
		Description: result.Description,
		Source:      string(result.Source),
		Builtin:     result.Builtin,
	}, nil
}

// ListSkills returns the effective visible skills for a workspace.
func (b *Backend) ListSkills(workspaceID string) ([]proto.SkillInfo, error) {
	ws, err := b.GetWorkspace(workspaceID)
	if err != nil {
		return nil, err
	}
	mgr := ws.Skills
	entries := skills.Catalog(mgr.ActiveSkills(), mgr.ResolvedPaths(), mgr.WorkingDir())
	result := make([]proto.SkillInfo, len(entries))
	for i, entry := range entries {
		result[i] = proto.SkillInfo{
			ID:            entry.ID,
			Name:          entry.Name,
			Description:   entry.Description,
			Label:         entry.Label,
			Source:        string(entry.Source),
			UserInvocable: entry.UserInvocable,
		}
	}
	return result, nil
}

// GetWorkingDir returns the working directory for a workspace.
func (b *Backend) GetWorkingDir(workspaceID string) (string, error) {
	ws, err := b.GetWorkspace(workspaceID)
	if err != nil {
		return "", err
	}
	return ws.Cfg.WorkingDir(), nil
}
