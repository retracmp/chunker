package upload

import (
	tea "github.com/charmbracelet/bubbletea"
)

func Start() error {
	if menu := tea.NewProgram(&ConfigManager{}); menu != nil {
		if _, err := menu.Run(); err != nil {
			return err
		}
	}

	return nil
}