package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateCppProject(projectName string) error {
	fmt.Printf("\rAttempting to create a C++ project %s...\n", projectName)

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

	// Create main.cpp
	fmt.Printf("\rCreating main.cpp...")
	mainContent := `/*
Copyright © 2024 Rudolf Farkas @rudifa rudi.farkas@gmail.com
*/

#include <iostream>

int main(int, const char *[])
{
    std::cout << "Hello, World!\n"
              << std::endl;
	return 0;
}
`
	err = os.WriteFile("main.cpp", []byte(mainContent), 0o644)
	if err != nil {
		return fmt.Errorf("\rfailed to create main.cpp: %w", err)
	}

	// Create Makefile
	fmt.Printf("\rCreating Makefile...")
	mainContent = `
CXX = /usr/bin/clang++
CXXFLAGS = -std=c++23 -Wall -Wextra -O2
TARGET = main
SRC = main.cpp
BUILD_DIR = build
TARGET_PATH = $(BUILD_DIR)/$(TARGET)

.PHONY: all run clean

all: $(TARGET_PATH)

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

$(TARGET_PATH): $(SRC) | $(BUILD_DIR)
	$(CXX) $(CXXFLAGS) -o $(TARGET_PATH) $(SRC)

run: $(TARGET_PATH)
	./$(TARGET_PATH)

clean:
	rm -rf $(BUILD_DIR)
`
	err = os.WriteFile("Makefile", []byte(mainContent), 0o644)
	if err != nil {
		return fmt.Errorf("\rfailed to create Makefile: %w", err)
	}

	// Create README.md
	fmt.Printf("\rCreating README.md...")
	err = os.WriteFile("README.md", []byte("# "+projectName), 0o644)
	if err != nil {
		return fmt.Errorf("\rfailed to create README.md: %w", err)
	}

	fmt.Printf("\rC++ project created successfully in %s\n", filepath.Join(projectName))
	fmt.Printf("\rTo test the new project:\n")
	fmt.Printf("\rcd %s\n", filepath.Join(projectName))
	fmt.Printf("\rmake run\n")
	return nil
}
