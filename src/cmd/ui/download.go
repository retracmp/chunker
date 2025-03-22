package ui

import (
	"acid/chunker/src/compile"
	"acid/chunker/src/downloader"
	"fmt"
	"path"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type download_stage int
const (
	enter_path download_stage = iota
	download_build
)

type download struct {
	manifest string
	path textinput.Model
	stage download_stage
}

func DownloadPage(selected string) *download {
	input := textinput.New()
	input.Placeholder = "Enter path to download"
	input.Focus()
	input.Width = 50

	return &download{
		manifest: selected,
		path: input,
		stage: enter_path,
	}
}

func (d *download) Init() tea.Cmd {
	return tea.Batch(
		tea.ClearScreen,
		tea.SetWindowTitle(fmt.Sprintf("Downloading %s", d.manifest)),
		textinput.Blink,
	)
}

func (d *download) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return d, tea.Quit
		case "enter":
			if d.stage == enter_path {
				d.stage = download_build
				return d, d.DownloadManifest()
			}
		}
	case DownloadResult:
		return d, tea.Quit
	}

	cmd := tea.Cmd(nil)
	d.path, cmd = d.path.Update(msg)
	return d, cmd
}

func (d *download) View() string {
	switch d.stage {
	case enter_path:
		return fmt.Sprintf(
			"Downlading Version %s\n\n%s\n\n%s", 
			d.manifest, 
			d.path.View(),
			gloss.NewStyle().Foreground(gloss.Color("#626262")).Render("q/esc/ctrl+c to quit"),
		)
	case download_build:
		return fmt.Sprintf("Downloading %s", d.manifest)
	default:
		return fmt.Sprintf("Downloading %s", d.manifest)
	}
}


type DownloadResult struct {}

func (d *download) DownloadManifest() tea.Cmd {
	options := downloader.NewDownloadOptions(
		fmt.Sprintf("https://cdn.retrac.site/manifest/%s.acidmanifest", d.manifest),
		d.path.Value(),	
	)

	return func() tea.Msg {
		downloader := downloader.NewDownloader(options.BaseURL, path.Join(options.TempDownloadDir, "chunks"), options.BuildDir)
		if _, err := downloader.FetchManifest(options.Manifest, options.TempDownloadDir); err != nil {
			return err
		}
		
		compiler := compile.NewCompiler(path.Join(options.TempDownloadDir, options.Manifest), path.Join(options.TempDownloadDir, "chunks"), options.BuildDir)
		if err := compiler.Check(); err == nil {
			if err := compiler.Cleanup(options.TempDownloadDir); err != nil {
				return err
			}
			return err
		}

		if err := downloader.DownloadThreaded(4); err != nil {
			return err
		}
		if err := compiler.Compile(); err != nil {
			return err
		}
		if err := compiler.Cleanup(options.TempDownloadDir); err != nil {
			return err
		}

		return DownloadResult{}
	}
}