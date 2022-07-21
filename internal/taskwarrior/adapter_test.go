package taskwarrior

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTempTaskRC(t *testing.T) (taskrcPath string) {
	taskrcPath = filepath.Join(t.TempDir(), "taskrc")
	file, err := os.Create(taskrcPath)
	assert.NoError(t, err, "Failed to create .taskrc for testing.")
	err = file.Close()
	assert.NoError(t, err, "Failed to close file handle to .taskrc.")
	os.Setenv("TASKRC", taskrcPath)
	return taskrcPath
}

func TestCreateUDAs_noUDAs_existAfterwards(t *testing.T) {
	// Use a custom .taskrc for testing.
	taskrc := createTempTaskRC(t)
	udaName := "foo"
	udaLabel := "bar"

	err := CreateUDA(udaName, udaLabel)
	assert.NoError(t, err, "CreateUDA returned an error.")

	out, err := exec.Command("bash", "-c", fmt.Sprintf("cat %s", taskrc)).
		CombinedOutput()
	t.Log(string(out))
	err = exec.Command("bash", "-c", fmt.Sprintf("cat %s | grep uda.%s.label=%s", taskrc, udaName, udaLabel)).
		Run()
	assert.NoError(
		t,
		err,
		fmt.Sprintf(
			"'%s' has no config entry 'uda.%s.label=%s'.",
			taskrc,
			udaName,
			udaLabel,
		),
	)
	err = exec.Command("bash", "-c", fmt.Sprintf("cat %s | grep uda.%s.type=string", taskrc, udaName)).
		Run()
	assert.NoError(
		t,
		err,
		fmt.Sprintf("'%s' has no config entry 'uda.%s.type=string'.", taskrc, udaName),
	)
}

func setup(t *testing.T) {
	// Use a custom .taskrc for testing.
	createTempTaskRC(t)

	// Use a custom TASKDATA directory for testing.
	taskdataPath := filepath.Join(t.TempDir(), "task")
	err := os.Mkdir(taskdataPath, 0755)
	assert.NoError(t, err, "Failed to create TASKDATA directory for testing.")
	os.Setenv("TASKDATA", taskdataPath)

	// Create UDA for the To-Do ID.
	err = CreateUDA(UDANameTodoID, "foo")
	assert.NoError(t, err, "CreateUDA returned an error.")
}

func TestTaskExists_notExists_returnsFalse(t *testing.T) {
	setup(t)

	toDoID := generateRandomString(10)
	exists, err := TaskExists(toDoID)
	assert.NoError(
		t,
		err,
		"TaskExists returned an error")
	assert.False(
		t,
		exists,
		fmt.Sprintf("Task with To-Do ID '%s' must not exist.", toDoID),
	)
}

func TestCreateTask_isOK(t *testing.T) {
	setup(t)

	toDoID := generateRandomString(10)
	taskUUID, err := CreateTask("foo", toDoID)
	assert.NoError(
		t,
		err,
		"CreateTask returned an error")
	assert.NotEmpty(
		t,
		taskUUID,
		fmt.Sprint("CreateTask returned an empty taskUUID."),
	)
}

func TestTaskExists_exists_returnsTrue(t *testing.T) {
	setup(t)

	toDoID := generateRandomString(10)
	CreateTask("foo", toDoID)
    exists, err := TaskExists(toDoID)
	assert.NoError(
		t,
		err,
		"TaskExists returned an error")
    assert.True(t, exists)
}

func generateRandomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, length)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
