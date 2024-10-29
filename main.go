/*
Copyright © 2024 Rudolf Farkas @rudifa rudi.farkas@gmail.com
*/
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices  []string
	cursor   int
	selected string
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
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "What kind of project would you like to create?\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\nPress q to quit.\n"

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	// Handle the selected option
	if m.(model).selected != "" {
		switch m.(model).selected {
		case "Create new Go project":
			if err := createGoProject(); err != nil {
				fmt.Printf("Error creating Go project: %v\n", err)
				os.Exit(1)
			}
		case "Create new Astro project":
			if err := createAstroProject(); err != nil {
				fmt.Printf("Error creating Astro project: %v\n", err)
				os.Exit(1)
			}
		}
	}
}


func createGoProject() error {
    // Your existing Go project creation logic
	print("Creating Go project...")
	fmt.Println("not yet.")
    return nil
}

func createAstroProject() error {
    // Your Astro project creation logic
	print("Creating Astro project...")
	fmt.Println("not yet.")
    return nil
}
