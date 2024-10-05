package tester

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// mustFindProjectRoot finds the project root by looking for "go.mod" file.
func FindProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("failed to os.Getwd: %v", err))
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break // Reached the root directory
		}
		dir = parentDir
	}
	panic("project root with go.mod file not found")
}

// SetupDir sets up a temporary test environment by copying necessary files and directories.
func SetupDir(t *testing.T, pathsToCopy map[string]string) (string, func()) {
	// Store the current working directory
	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("SetupDir: failed to os.Getwd: %v", err))
	}

	// Create a temporary directory
	tempDir := t.TempDir()

	// Defer cleanup to reset the working directory after the test
	cleanup := func() {
		_ = os.Chdir(wd)
	}
	defer func() {
		if err := os.Chdir(tempDir); err != nil {
			panic(fmt.Errorf("SetupDir: failed to os.Chdir: %v", err))
		}
	}()

	projectRoot := FindProjectRoot()

	// Copy specified paths to the temporary directory
	for source, destination := range pathsToCopy {
		sourcePath := strings.ReplaceAll(source, "@", projectRoot)
		if destination == "" {
			destination = filepath.Join(tempDir, filepath.Base(sourcePath))
		} else {
			destination = filepath.Join(tempDir, destination)
		}

		if err := copyPath(sourcePath, destination); err != nil {
			panic(fmt.Errorf("SetupDir: failed to copyPath: %v", err))
		}
	}

	return tempDir, cleanup
}

// copyPath copies a file or directory from sourcePath to destinationPath.
func copyPath(sourcePath, destinationPath string) error {
	info, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return copyDirectory(sourcePath, destinationPath)
	}
	return copyFile(sourcePath, destinationPath)
}

// copyDirectory copies all files and subdirectories from sourceDir to destinationDir.
func copyDirectory(sourceDir, destinationDir string) error {
	if err := os.MkdirAll(destinationDir, os.ModePerm); err != nil {
		return err
	}

	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourceEntryPath := filepath.Join(sourceDir, entry.Name())
		destinationEntryPath := filepath.Join(destinationDir, entry.Name())

		if entry.IsDir() {
			if err := copyDirectory(sourceEntryPath, destinationEntryPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(sourceEntryPath, destinationEntryPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// copyFile copies a single file from sourceFile to destinationFile.
func copyFile(sourceFile, destinationFile string) error {
	if err := os.MkdirAll(filepath.Dir(destinationFile), os.ModePerm); err != nil {
		return err
	}

	source, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(destinationFile)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}
