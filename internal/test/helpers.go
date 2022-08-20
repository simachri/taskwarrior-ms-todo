package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// CreateTempTaskRC creates a .taskrc that is used for a single test only.
func CreateTempTaskRC(t *testing.T) (taskrcPath string) {
	taskrcPath = filepath.Join(t.TempDir(), "taskrc")
	file, err := os.Create(taskrcPath)
	assert.NoError(t, err, "Failed to create .taskrc for testing.")
	err = file.Close()
	assert.NoError(t, err, "Failed to close file handle to .taskrc.")
	os.Setenv("TASKRC", taskrcPath)
	return taskrcPath
}

// NewTaskwarriorEnv creates a new Taskwarrior environment with a custom .taskrc, a
// custom TASKDATA directory that is purged after each test and the
// User-Defined-Attributes (UDAs).
func NewTaskwarriorEnv(t *testing.T) {
	CreateTempTaskRC(t)

	taskdataPath := filepath.Join(t.TempDir(), "task")
	err := os.Mkdir(taskdataPath, 0755)
	assert.NoError(t, err, "Failed to create TASKDATA directory for testing.")
	os.Setenv("TASKDATA", taskdataPath)
}
