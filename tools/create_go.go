package tools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func CreateGoProject(projectName string, useCobraCLI bool, addTest bool) error {
	fmt.Println("Creating Go project...")

	hub := "github.com"
	owner := "rudifa"

	// Check if directory already exists
	if _, err := os.Stat(projectName); !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exists", projectName)
	}

	// Create project directory
	err := os.MkdirAll(projectName, 0755)
	if err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Change to project directory
	err = os.Chdir(projectName)
	if err != nil {
		return fmt.Errorf("failed to change to project directory: %w", err)
	}

	// Initialize go module
	modulePath := fmt.Sprintf("%s/%s/%s", hub, owner, projectName)
	cmd := exec.Command("go", "mod", "init", modulePath)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to initialize go module: %w", err)
	}

	// Create main.go
	mainContent := `/*
Copyright © 2024 Rudolf Farkas @rudifa rudi.farkas@gmail.com
*/

package main

import "fmt"

func main() {
	fmt.Println("Here we go")
}
`
	err = os.WriteFile("main.go", []byte(mainContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to create main.go: %w", err)
	}

	// Optionally create main_test.go
	if addTest {
		testContent := `package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	assert.Equal(t, 2, 1+1)
}
`
		err = os.WriteFile("main_test.go", []byte(testContent), 0644)
		if err != nil {
			return fmt.Errorf("failed to create main_test.go: %w", err)
		}
	}

	// Create README.md
	err = os.WriteFile("README.md", []byte("# "+projectName), 0644)
	if err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	// Initialize cobra-cli if specified
	if useCobraCLI {
		cmd = exec.Command("cobra-cli", "init")
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to initialize cobra-cli: %w", err)
		}
	}

	fmt.Printf("\rGo project created successfully in %s\n", filepath.Join(projectName))
	return nil
}
