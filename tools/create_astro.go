package tools

import (
    "fmt"
    "os"
    "os/exec"
    "time"
		"path/filepath"

)

func CreateAstroProject(projectName string) error {
    if _, err := os.Stat(projectName); !os.IsNotExist(err) {
        return fmt.Errorf("directory '%s' already exists", projectName)
    }

    // fmt.Println("Creating Astro project...")

    // Create a channel to signal when the command is done
    done := make(chan bool)

    // Start the command in a goroutine
    go func() {
        cmd := exec.Command("npm", "create", "astro@latest", projectName, "--", "--template", "minimal", "--install", "--git", "--typescript", "strict", "-y")
        err := cmd.Run()
        if err != nil {
            fmt.Printf("\nFailed to create Astro project: %v\n", err)
        }
        done <- true
    }()

    // Display a spinner while waiting for the command to complete
    spinner := []string{"|", "/", "-", "\\"}
    i := 0
    for {
        select {
        case <-done:
			fmt.Printf("\rAstro project created successfully in %s\n", filepath.Join(projectName))
            return nil
        default:
            fmt.Printf("\rCreating Astro project... %s", spinner[i])
            time.Sleep(100 * time.Millisecond)
            i = (i + 1) % len(spinner)
        }
    }
}
