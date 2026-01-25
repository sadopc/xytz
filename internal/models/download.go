package models

import (
	"strings"
	"time"
	"xytz/internal/styles"
	"xytz/internal/types"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type DownloadModel struct {
	Progress     progress.Model
	CurrentSpeed string
	CurrentETA   string
	Completed    bool
}

func NewDownloadModel() DownloadModel {
	pr := progress.New(progress.WithSolidFill(string(styles.InfoColor)))

	return DownloadModel{Progress: pr}
}

func (m DownloadModel) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return progress.FrameMsg{}
	})
}

func (m DownloadModel) Update(msg tea.Msg) (DownloadModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case types.ProgressMsg:
		cmd = m.Progress.SetPercent(msg.Percent / 100.0)
		m.CurrentSpeed = msg.Speed
		m.CurrentETA = msg.Eta
	case tea.KeyMsg:
		if m.Completed && msg.Type == tea.KeyEnter {
			cmd = func() tea.Msg {
				return types.GoBackMsg{}
			}
		}
	}

	newModel, downloadCmd := m.Progress.Update(msg)
	if newModel, ok := newModel.(progress.Model); ok {
		m.Progress = newModel
	}

	return m, tea.Batch(cmd, downloadCmd)
}

func (m DownloadModel) HandleResize(w, h int) DownloadModel {
	if w > 100 {
		m.Progress.Width = (w / 2) - 10
	} else {
		m.Progress.Width = w - 10
	}
	return m
}

func (m DownloadModel) View() string {
	var s strings.Builder

	statusText := "⇣ Downloading"
	if m.Completed {
		statusText = "Download Complete"
	}

	s.WriteString(styles.SectionHeaderStyle.Foreground(styles.InfoColor).Render(statusText))
	s.WriteRune('\n')

	if m.Completed {
		s.WriteString(styles.CompletionMessageStyle.Render("Video saved to current directory."))
		s.WriteRune('\n')
	} else {
		percent := m.Progress.Percent()
		if percent > 0 || m.CurrentSpeed != "" || m.CurrentETA != "" {
			bar := styles.ProgressContainer.Render(m.Progress.View())
			s.WriteString(bar)
			s.WriteRune('\n')

			if m.CurrentSpeed != "" {
				s.WriteString("Speed: " + styles.SpeedStyle.Render(m.CurrentSpeed))
				s.WriteRune('\n')
			}

			if m.CurrentETA != "" {
				s.WriteString("Time remaining: " + styles.TimeRemainingStyle.Render(m.CurrentETA))
				s.WriteRune('\n')
			}

			dest := "./"
			s.WriteString("Destination: " + styles.DestinationStyle.Render(dest))
		} else {
			s.WriteString("Starting download...")
			s.WriteRune('\n')
		}
	}

	return s.String()
}
