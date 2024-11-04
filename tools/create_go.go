package tools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func CreateGoProject(projectName string, useCobraCLI bool, addTest bool) error {
	fmt.Printf("\rAttempting to create Go project %s...\n", projectName)

	hub := "github.com"
	owner := "rudifa"

	// Check if directory already exists
	if _, err := os.Stat(projectName); err == nil {
		return fmt.Errorf("\rerror: directory '%s' already exists", projectName)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("\rerror checking if directory exists: %w", err)
	}

	// Create project directory
	fmt.Printf("\rCreating project directory...")
	err := os.Mkdir(projectName, 0o755)
	if err != nil {
		return fmt.Errorf("\rfailed to create project directory: %w", err)
	}

	// Change to project directory
	fmt.Printf("\rChanging to project directory...")
	err = os.Chdir(projectName)
	if err != nil {
		return fmt.Errorf("\rfailed to change to project directory: %w", err)
	}

	// Initialize go module
	fmt.Printf("\rInitializing go module...")
	modulePath := fmt.Sprintf("%s/%s/%s", hub, owner, projectName)
	cmd := exec.Command("go", "mod", "init", modulePath)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("\rfailed to initialize go module: %w", err)
	}

	// Create main.go
	fmt.Printf("\rCreating main.go...")
	mainContent := `/*
Copyright © 2024 Rudolf Farkas @rudifa rudi.farkas@gmail.com
*/

package main

import "fmt"

func main() {
    fmt.Println("Here we go")
}
`
	err = os.WriteFile("main.go", []byte(mainContent), 0o644)
	if err != nil {
		return fmt.Errorf("\rfailed to create main.go: %w", err)
	}

	// Optionally create main_test.go
	if addTest {
		fmt.Printf("\rCreating main_test.go...")
		testContent := `package main

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
    assert.Equal(t, 2, 1+1)
}
`
		err = os.WriteFile("main_test.go", []byte(testContent), 0o644)
		if err != nil {
			return fmt.Errorf("\rfailed to create main_test.go: %w", err)
		}
	}

	// Create README.md
	fmt.Printf("\rCreating README.md...")
	err = os.WriteFile("README.md", []byte("# "+projectName), 0o644)
	if err != nil {
		return fmt.Errorf("\rfailed to create README.md: %w", err)
	}

	// Initialize cobra-cli if specified
	if useCobraCLI {
		fmt.Printf("\rInitializing cobra-cli...")
		cmd = exec.Command("cobra-cli", "init")
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("\rfailed to initialize cobra-cli: %w", err)
		}
	}

	fmt.Printf("\rGo project created successfully in %s\n", filepath.Join(projectName))
	return nil
}
