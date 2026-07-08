package cmd

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/crush/internal/config"
	"github.com/charmbracelet/crush/internal/workspace"
	"github.com/spf13/cobra"
)

var (
	logoutHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	logoutItemStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	logoutPromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("215"))
)

var logoutCmd = &cobra.Command{
	Aliases: []string{"signout"},
	Use:     "logout [platform]",
	Short:   "Logout Crush from a platform",
	Long: `Logout Crush from a specified platform, removing stored credentials.
The platform should be provided as an argument.
If no argument is given, a list of logged-in platforms will be shown.
Available platforms are: hyper, copilot.`,
	Example: `
# Sign out from Charm Hyper
crush logout hyper

# Sign out from GitHub Copilot
crush logout copilot
  `,
	ValidArgs: []cobra.Completion{
		"hyper",
		"copilot",
		"github",
		"github-copilot",
	},
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ws, cleanup, err := setupLocalWorkspace(cmd)
		if err != nil {
			return err
		}
		defer cleanup()

		var provider string
		if len(args) == 0 {
			provider, err = pickLoggedInProvider(ws)
			if err != nil {
				return err
			}
			if provider == "" {
				return nil
			}
		} else {
			provider = args[0]
		}

		force, _ := cmd.Flags().GetBool("force")
		if !force {
			fmt.Print(logoutPromptStyle.Render(fmt.Sprintf("Are you sure you want to logout %s? (y/N) ", provider)))
			var response string
			_, err := fmt.Scanln(&response)
			if err != nil || (response != "y" && response != "Y" && response != "yes" && response != "Yes" && response != "YES") {
				fmt.Println(logoutHeaderStyle.Render("Logout cancelled."))
				return nil
			}
		}

		switch provider {
		case "hyper":
			return logoutHyper(ws)
		case "copilot", "github", "github-copilot":
			return logoutCopilot(ws)
		default:
			return fmt.Errorf("unknown platform: %s", provider)
		}
	},
}

func logoutHyper(ws workspace.Workspace) error {
	if err := ws.RemoveConfigField(config.ScopeGlobal, "providers.hyper.api_key"); err != nil {
		return err
	}
	if err := ws.RemoveConfigField(config.ScopeGlobal, "providers.hyper.oauth"); err != nil {
		return err
	}

	fmt.Println(logoutHeaderStyle.Render("Successfully logged out of Hyper."))
	return nil
}

func logoutCopilot(ws workspace.Workspace) error {
	if err := ws.RemoveConfigField(config.ScopeGlobal, "providers.copilot.api_key"); err != nil {
		return err
	}
	if err := ws.RemoveConfigField(config.ScopeGlobal, "providers.copilot.oauth"); err != nil {
		return err
	}

	fmt.Println(logoutHeaderStyle.Render("Successfully logged out of GitHub Copilot."))
	return nil
}

func pickLoggedInProvider(ws workspace.Workspace) (string, error) {
	cfg := ws.Config()

	type loggedInProvider struct {
		id   string
		name string
	}

	var loggedIn []loggedInProvider
	for p := range cfg.Providers.Seq() {
		if p.OAuthToken != nil || p.APIKey != "" {
			name := p.Name
			if name == "" {
				name = p.ID
			}
			loggedIn = append(loggedIn, loggedInProvider{id: p.ID, name: name})
		}
	}

	if len(loggedIn) == 0 {
		fmt.Println(logoutPromptStyle.Render("You are not logged in to any platform."))
		return "", nil
	}

	if len(loggedIn) == 1 {
		return loggedIn[0].id, nil
	}

	fmt.Println(logoutHeaderStyle.Render("Logged-in platforms:"))
	for i, p := range loggedIn {
		fmt.Println(logoutItemStyle.Render(fmt.Sprintf("  %d. %s", i+1, p.name)))
	}
	fmt.Print(logoutPromptStyle.Render(fmt.Sprintf("Select a platform to logout (1-%d): ", len(loggedIn))))

	var choice int
	_, err := fmt.Scanln(&choice)
	if err != nil || choice < 1 || choice > len(loggedIn) {
		fmt.Println(logoutHeaderStyle.Render("Logout cancelled."))
		return "", nil
	}

	return loggedIn[choice-1].id, nil
}

func init() {
	logoutCmd.Flags().BoolP("force", "f", false, "Skip logout confirmation prompt")
}


