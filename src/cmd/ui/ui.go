package ui

import (
	"acid/chunker/src/helpers/http"

	tea "github.com/charmbracelet/bubbletea"

	"fmt"
)

var (
	retracHttpClient *http.HttpClient = http.NewClient("https://retrac.site/retrac")
)

func Start() error {
	manifests, err := http.Get[[]string](retracHttpClient, "manifests")
	if err != nil {
		return err
	}

	if manifests == nil || len(*manifests) == 0 {
		return fmt.Errorf("no manifests found, please come back later")
	}

	state := MenuPage(*manifests)
	if state == nil {
		return fmt.Errorf("failed to create menu page")
	}

	if menu := tea.NewProgram(state); menu != nil {
		if _, err := menu.Run(); err != nil {
			return err
		}
	}

	if len(state.Selected) == 0 {
		return nil
	}

	if state.Selected[state.Cursor] == "" {
		return fmt.Errorf("no manifest selected")
	}

	if download := DownloadPage(state.Selected[state.Cursor]); download != nil {
		if _, err := tea.NewProgram(download).Run(); err != nil {
			return err
		}
	}

	return nil
}