package upload

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ConfigManager struct {
	
}

func (d *ConfigManager) Init() tea.Cmd {
	return tea.Batch(
		tea.ClearScreen,
		tea.SetWindowTitle("Setting Config"),
		textinput.Blink,
	)
}

func (d *ConfigManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return d, tea.Quit
		default:

		}
	}

	return d, nil
}

func (d *ConfigManager) View() string {
	return "Config Manager"
}