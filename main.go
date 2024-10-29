/*
Copyright © 2024 Rudolf Farkas @rudifa rudi.farkas@gmail.com
*/
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rudifa/newproject/tools"
	"strings"
)

type model struct {
	choices  []string
	cursor   int
	selected string
	err      error
}

func initialModel() model {
	return model{
		choices: []string{
			"Create new Go project",
			"Create new Astro project",
		},
		cursor: 0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// If there's an error and user presses 'q', quit immediately
		if m.err != nil && (msg.String() == "q" || msg.String() == "ctrl+c") {
			return m, tea.Quit
		}
		// Clear any existing error when a key is pressed
		m.err = nil
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

		case "enter":
			m.selected = m.choices[m.cursor]
			return m.handleSelection()
		}
	}

	return m, nil
}

func (m model) handleSelection() (tea.Model, tea.Cmd) {
    var err error
    var projectName string
    switch m.selected {
    case "Create new Go project":
        projectName = ".tmp/my-go-project"
        err = tools.CreateGoProject(projectName, false, true) // No Cobra, with test
    case "Create new Astro project":
        projectName = ".tmp/my-astro-project"
        err = tools.CreateAstroProject(projectName)
    }

    if err != nil {
        m.err = err
        return m, tea.Quit
    }

    m.err = nil
    return m, tea.Quit
}

func (m model) View() string {
    // Define styles
    titleStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("205")).
        Bold(true).
        MarginBottom(1)

    selectedStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("170")).
        Bold(true)

    normalStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("244"))

    quitStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("240")).
        Italic(true).
        MarginTop(1)

    errorStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("196")).
        Bold(true).
        MarginTop(1).
        MarginBottom(1)

    var s strings.Builder

    // Always start with a newline
    s.WriteString("\n")

    if m.err != nil {
        s.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
        s.WriteString("\n")
        s.WriteString(quitStyle.Render("Press q to quit."))
        return s.String()
    }

    s.WriteString(titleStyle.Render("What kind of project would you like to create?"))
    s.WriteString("\n\n")

    for i, choice := range m.choices {
        cursor := "  "
        if m.cursor == i {
            cursor = "❯ "
            s.WriteString(selectedStyle.Render(cursor + choice))
        } else {
            s.WriteString(normalStyle.Render(cursor + choice))
        }
        s.WriteString("\n")
    }

    s.WriteString("\n")
    s.WriteString(quitStyle.Render("Press q to quit."))

    return s.String()
}
func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
