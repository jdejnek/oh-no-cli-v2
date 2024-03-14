package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Menu struct {
	MainMenuItems  []string
	Submenus       map[string][]string
	CurrentSubmenu int
	SelectedIndex  int
	SubmenuIndex   int
}

var redText = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#de1b62"))

var header = lipgloss.NewStyle().Bold(true)

const (
	MainMenu int = iota
	Submenu
)

func initialModel() Menu {
	return Menu{
		MainMenuItems: []string{"Sims", "Connectors", "Webhooks", "Live Monitor"},
		Submenus: map[string][]string{
			"Sims":         {"View", "Create", "Delete"},
			"Connectors":   {"View", "Create", "Delete"},
			"Webhooks":     {"View", "Create", "Delete"},
			"Live Monitor": {"Traffic Monitor", "Network Logs", "Webhooks feed"},
		},
		CurrentSubmenu: -1,
		SelectedIndex:  0,
		SubmenuIndex:   0,
	}
}

func (m Menu) Init() tea.Cmd {
	return nil
}

func (m Menu) MoveSelection(direction int) Menu {
	newIndex := m.SelectedIndex + direction
	if newIndex < 0 {
		newIndex = len(m.MainMenuItems) - 1
	} else if newIndex >= len(m.MainMenuItems) {
		newIndex = 0
	}
	m.SelectedIndex = newIndex
	return m
}

func (m Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.CurrentSubmenu == -1 {
				return m.MoveSelection(-1), nil
			} else {
				return m.MoveSubmenuSelection(-1), nil
			}
		case "down", "j":
			if m.CurrentSubmenu == -1 {
				return m.MoveSelection(1), nil
			} else {
				return m.MoveSubmenuSelection(1), nil
			}
		case "enter":
			if m.CurrentSubmenu == -1 {
				m.CurrentSubmenu = m.SelectedIndex
			} else {
				m.CurrentSubmenu = -1
				m.SubmenuIndex = 0
			}
			return m, nil
		case "esc":
			if m.CurrentSubmenu != -1 {
				m.CurrentSubmenu = -1
				m.SubmenuIndex = 0
				return m, nil
			}
		}
	}
	return m, nil
}

func (m Menu) View() string {
	var output string
	logo := fmt.Sprintln(redText.Render(`
    ____     __  __           _   __   ____           ______    __     ____
   / __ \   / / / /          / | / /  / __ \         / ____/   / /    /  _/
  / / / /  / /_/ /  ______  /  |/ /  / / / / ______ / /       / /     / /
 / /_/ /  / __  /  /_____/ / /|  /  / /_/ / /_____// /___    / /___ _/ /
 \____/  /_/ /_/          /_/ |_/   \____/         \____/   /_____//___/


`))
	output += logo

	if m.CurrentSubmenu == -1 {
		// Main menu view

		for i, item := range m.MainMenuItems {
			if i == m.SelectedIndex {
				selectedItem := fmt.Sprintf("  %s", item)
				output += fmt.Sprintln(redText.Render(selectedItem))
			} else {
				output += fmt.Sprintf(" %s\n", item)
			}
		}
		output += "\n\nPress 'enter' to select an option, 'q' to quit"
	} else {
		// Submenu view
		submenuKey := m.MainMenuItems[m.CurrentSubmenu]
		submenuItems := m.Submenus[submenuKey]

		submenuTitle := fmt.Sprintf(" %s", submenuKey)
		output += fmt.Sprintln(header.Render(submenuTitle))
		for i, item := range submenuItems {
			if i == m.SubmenuIndex {
				selectedItem := fmt.Sprintf("  %s", item)
				output += fmt.Sprintln(redText.Render(selectedItem))
			} else {
				output += fmt.Sprintf(" %s\n", item)
			}
		}
		output += "\n\nPress 'esc' to go back to main menu, 'q' to quit"
	}
	return output
}

func (m Menu) MoveSubmenuSelection(direction int) Menu {
	submenuKey := m.MainMenuItems[m.CurrentSubmenu]
	submenuItems := m.Submenus[submenuKey]
	newIndex := m.SubmenuIndex + direction
	if newIndex < 0 {
		newIndex = len(submenuItems) - 1
	} else if newIndex >= len(submenuItems) {
		newIndex = 0
	}
	m.SubmenuIndex = newIndex
	return m
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Whoops. %v", err)
		os.Exit(1)
	}
}
