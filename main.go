package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices  []string         // which item is selected
	cursor   int              // which item cursor is pointing at
	selected map[int]struct{} // which items are selected
}

func initialModel() model {
	return model{
		choices:  []string{"Sims", "Usage", "Connectors", "Webhooks", "Live Monitor"},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ", "l":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	s := `
    ____     __  __           _   __   ____           ______    __     ____
   / __ \   / / / /          / | / /  / __ \         / ____/   / /    /  _/
  / / / /  / /_/ /  ______  /  |/ /  / / / / ______ / /       / /     / /
 / /_/ /  / __  /  /_____/ / /|  /  / /_/ / /_____// /___    / /___ _/ /
 \____/  /_/ /_/          /_/ |_/   \____/         \____/   /_____//___/


`
	for i, choice := range m.choices {

		// is the cursor pointing at this chocice?
		cursor := " "
		if m.cursor == i {
			cursor = "  "
		}
		// is this choice selected?
		if _, ok := m.selected[i]; ok {
			s += "Selected. "
		}

		// render the row
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	s += "\n\nPress q to quit."
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Whoops. %v", err)
		os.Exit(1)
	}
}
