package tester

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupDir_ValidFiles(t *testing.T) {
	tempDir := t.TempDir()
	sourceFile := filepath.Join(tempDir, "source.txt")
	err := os.WriteFile(sourceFile, []byte("test data"), 0644)
	assert.NoError(t, err)

	pathsToCopy := map[string]string{sourceFile: ""}
	destDir, cleanup := SetupDir(t, pathsToCopy)
	defer cleanup()

	assert.NoError(t, err)
	destFile := filepath.Join(destDir, filepath.Base(sourceFile))
	assert.FileExists(t, destFile)

	content, err := os.ReadFile(destFile)
	assert.NoError(t, err)
	assert.Equal(t, "test data", string(content))
}

func TestSetupDir_DirectoryCopy(t *testing.T) {
	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "sourceDir")
	assert.NoError(t, os.Mkdir(sourceDir, 0755))
	sourceFile := filepath.Join(sourceDir, "file.txt")
	assert.NoError(t, os.WriteFile(sourceFile, []byte("test data"), 0644))

	pathsToCopy := map[string]string{sourceDir: ""}
	destDir, cleanup := SetupDir(t, pathsToCopy)
	defer cleanup()

	destPath := filepath.Join(destDir, filepath.Base(sourceDir))
	assert.DirExists(t, destPath)

	destFile := filepath.Join(destPath, "file.txt")
	assert.FileExists(t, destFile)

	content, err := os.ReadFile(destFile)
	assert.NoError(t, err)
	assert.Equal(t, "test data", string(content))
}

func TestSetupDir_DestinationPath(t *testing.T) {
	tempDir := t.TempDir()
	sourceFile := filepath.Join(tempDir, "source.txt")
	err := os.WriteFile(sourceFile, []byte("test data"), 0644)
	assert.NoError(t, err)

	pathsToCopy := map[string]string{sourceFile: "custom/destination.txt"}
	destDir, cleanup := SetupDir(t, pathsToCopy)
	defer cleanup()

	destFile := filepath.Join(destDir, "custom/destination.txt")
	assert.FileExists(t, destFile)

	content, err := os.ReadFile(destFile)
	assert.NoError(t, err)
	assert.Equal(t, "test data", string(content))
}

func TestCopyDirectory_Ok(t *testing.T) {
	source := t.TempDir()
	destination := t.TempDir()

	err := os.WriteFile(filepath.Join(source, "file1.txt"), []byte("content1"), 0644)
	assert.NoError(t, err)
	err = os.Mkdir(filepath.Join(source, "subdir"), 0755)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(source, "subdir", "file2.txt"), []byte("content2"), 0644)
	assert.NoError(t, err)

	err = copyDirectory(source, destination)
	assert.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(destination, "file1.txt"))
	assert.NoError(t, err)
	assert.Equal(t, "content1", string(content))
	content, err = os.ReadFile(filepath.Join(destination, "subdir", "file2.txt"))
	assert.NoError(t, err)
	assert.Equal(t, "content2", string(content))
}

func TestCopyDirectory_Error(t *testing.T) {
	source := t.TempDir()
	destination := t.TempDir()

	err := copyDirectory(filepath.Join(source, "nonexistent"), destination)
	assert.ErrorContains(t, err, "no such file or directory")
}

func TestCopyFile(t *testing.T) {
	source := t.TempDir() + "/source.txt"
	destination := t.TempDir() + "/destination.txt"

	content := []byte("This is a test content")
	err := os.WriteFile(source, content, os.ModePerm)
	assert.NoError(t, err)

	err = copyFile(source, destination)
	assert.NoError(t, err)

	copiedContent, err := os.ReadFile(destination)
	assert.NoError(t, err)
	assert.Equal(t, content, copiedContent)
}

func TestSetupDir_NonExistentSource(t *testing.T) {
	assert.Panics(t, func() {
		pathsToCopy := map[string]string{"nonexistent.txt": ""}
		_, cleanup := SetupDir(t, pathsToCopy)
		defer cleanup()
	})
}
