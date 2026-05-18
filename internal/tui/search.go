package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/Abhiram86/echotune/internal/models"
)

type searchModel struct {
	results  []models.SearchResult
	cursor   int
	selected int
	quitting bool
}

func (m searchModel) Init() tea.Cmd {
	return nil
}

func (m searchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			m.selected = -1
			return m, tea.Quit
		case "up", "n":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.results) - 1
			}
		case "down", "b":
			m.cursor++
			if m.cursor >= len(m.results) {
				m.cursor = 0
			}
		case "enter":
			m.quitting = true
			m.selected = m.cursor
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m searchModel) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}

	var b strings.Builder
	b.WriteString("\n  Select a song to play:\n\n")

	for i, res := range m.results {
		cursor := " "
		if m.cursor == i {
			cursor = "▶"
		}
		
		title := res.Title
		// Truncate title if it's too long
		if len(title) > 60 {
			title = title[:57] + "..."
		}

		if m.cursor == i {
			fmt.Fprintf(&b, "  %s %s\n      %s\n", cursor, title, res.Channel)
		} else {
			fmt.Fprintf(&b, "  %s %s\n      %s\n", cursor, title, res.Channel)
		}
	}
	
	b.WriteString("\n  [n/↑] up  [b/↓] down  [enter] select  [q/esc] quit\n")

	return tea.NewView(b.String())
}

func SearchSelection(results []models.SearchResult) (int, error) {
	m := searchModel{
		results:  results,
		selected: -1,
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return -1, err
	}

	sm := finalModel.(searchModel)
	if sm.selected == -1 {
		return -1, fmt.Errorf("selection cancelled")
	}

	return sm.selected, nil
}
