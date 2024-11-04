/*
Copyright © 2024 Rudolf Farkas @rudifa rudi.farkas@gmail.com
*/
package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rudifa/newpro/tools"
)

type model struct {
	choices       []string
	cursor        int
	selected      string
	err           error
	projectName   string
	useCobraCLI   bool
	addTest       bool
	currentOption string
}

func initialModel() model {
	return model{
		choices: []string{
			"Create new Go project",
			"Create new Astro project",
		},
		cursor:        0,
		projectName:   "",
		useCobraCLI:   false,
		addTest:       false,
		currentOption: "projectType",
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

		case "enter":
			switch m.currentOption {
			case "projectType":
				m.selected = m.choices[m.cursor]
				m.currentOption = "projectName"
			case "projectName":
				if m.projectName != "" {
					if m.selected == "Create new Go project" {
						m.currentOption = "useCobraCLI"
					} else {
						return m.handleSelection()
					}
					m.cursor = 0
				}
			case "useCobraCLI":
				m.useCobraCLI = m.cursor == 0
				m.currentOption = "addTest"
			case "addTest":
				m.addTest = m.cursor == 0
				return m.handleSelection()
			}

		default:
			if m.currentOption == "projectName" {
				switch msg.String() {
				case "backspace":
					if len(m.projectName) > 0 {
						m.projectName = m.projectName[:len(m.projectName)-1]
					}
				case "space":
					m.projectName += " "
				default:
					if len(msg.String()) == 1 {
						m.projectName += msg.String()
					}
				}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	var s strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		MarginBottom(1)

	menuStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		PaddingLeft(4)

	selectedStyle := menuStyle.Copy().
		Foreground(lipgloss.Color("170")).
		Bold(true)

	s.WriteString(titleStyle.Render("Project Creator"))
	s.WriteString("\n\n")

	switch m.currentOption {
	case "projectType":
		s.WriteString("Select project type:\n\n")
		for i, choice := range m.choices {
			if m.cursor == i {
				s.WriteString(selectedStyle.Render("> " + choice + "\n"))
			} else {
				s.WriteString(menuStyle.Render("  " + choice + "\n"))
			}
			s.WriteString("\n") // Add a newline after each option
		}
	case "projectName":
		s.WriteString(fmt.Sprintf("Enter project name for %s:\n", m.selected))
		s.WriteString(m.projectName)
		s.WriteString("_") // Add a cursor
		s.WriteString("\n\n(Press Enter when done)\n")
	case "useCobraCLI", "addTest":
		prompt := "Use Cobra CLI?"
		if m.currentOption == "addTest" {
			prompt = "Add test file?"
		}
		s.WriteString(prompt + "\n\n")
		for i, choice := range []string{"Yes", "No"} {
			if m.cursor == i {
				s.WriteString(selectedStyle.Render("> " + choice + "\n"))
			} else {
				s.WriteString(menuStyle.Render("  " + choice + "\n"))
			}
			s.WriteString("\n") // Add a newline after each option
		}
	}

	s.WriteString("\n(Press q to quit)")
	return s.String()
}

func (m model) handleSelection() (tea.Model, tea.Cmd) {
	var err error
	projectName := m.projectName

	if projectName == "" || strings.ContainsAny(projectName, "/\\:*?\"<>|") {
		m.err = fmt.Errorf("invalid project name: %s", projectName)
		fmt.Println(m.err)
		return m, tea.Quit
	}

	switch m.selected {
	case "Create new Go project":
		err = tools.CreateGoProject(projectName, m.useCobraCLI, m.addTest)
	case "Create new Astro project":
		err = tools.CreateAstroProject(projectName)
	}

	if err != nil {
		m.err = err
		fmt.Println(err) // Print the error immediately
		return m, tea.Quit
	}

	m.err = nil
	fmt.Printf("\rProject '%s' created successfully in the current directory\n", projectName)
	return m, tea.Quit
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
