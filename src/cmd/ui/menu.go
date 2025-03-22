package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type menu struct {
	manifests []string
	Cursor int
	Selected map[int]string
}

func MenuPage(manifests []string) *menu {
	return &menu{
		manifests: manifests,
		Cursor: 0,
		Selected: make(map[int]string),
	}
}

func (*menu) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("Chunker"),
		tea.ClearScreen,
	)
}

func (p *menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return p, tea.Quit
		case "up", "k", "w":
			if p.Cursor > 0 {
				p.Cursor--
			}
		case "down", "j", "s":
			
			if len(p.Selected) > 0 {
				if p.Cursor < len(p.manifests) {
					p.Cursor ++
				}
			} else {
				if p.Cursor < len(p.manifests) - 1 {
					p.Cursor++
				}
			}
		case "enter", " ":
			if len(p.Selected) > 0 && p.Cursor == len(p.manifests) {
				for i := range p.Selected {
					p.Cursor = i
				}

				return p, tea.Quit
			}

			p.Selected = make(map[int]string)
			p.Selected[p.Cursor] = p.manifests[p.Cursor]
		}
	}

	return p, nil
}

func (p *menu) View() string {
	render := "Please select a Fortnite version to download: \n\n"

	if len(p.manifests) == 0 {
		render += "No manifests found, please come back later\n"
		return render
	}

	manifests := make([]string, len(p.manifests))
	copy(manifests, p.manifests)

	if len(p.Selected) > 0 {
		manifests = append(manifests, "Download Selected Version")
	}

	for i, manifest := range manifests {
		cursor := " "
		if i == p.Cursor {
			cursor = ">"
		}

		selected := " "
		if _, ok := p.Selected[i]; ok {
			selected = "*"
		}

		if manifest == "Download selected version" {
			render += fmt.Sprintf("%s [ %s ]\n", cursor, manifest)
			continue
		}

		render += fmt.Sprintf("%s [%s] %s\n", cursor, selected, manifest)
	}

	render += fmt.Sprintf(
		"\n%s\n", 
		gloss.NewStyle().Foreground(gloss.Color("#626262")).Render("q/esc/ctrl+c to quit"),
	)
	
	return render
}