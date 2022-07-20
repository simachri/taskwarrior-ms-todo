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

func createTempTaskRC (t *testing.T) (taskrcName string, taskrcFile *os.File) {
	taskrc := filepath.Join(t.TempDir(), "taskrc")
	file, err := os.Create(taskrc)
	assert.NoError(t, err, "Failed to create .taskrc for testing.")
	os.Setenv("TASKRC", taskrc)
    return taskrc, file
}

func TestCreateUDAs_noUDAs_existAfterwards(t *testing.T) {
	// Use a custom .taskrc for testing.
    taskrc, file := createTempTaskRC(t)
	defer file.Close()
	udaName := "foo"
	udaLabel := "bar"

    err := CreateUDA(udaName, udaLabel)
	assert.NoError(t, err, "CreateUDA returned an error.")

	out, err := exec.Command("bash", "-c", fmt.Sprintf("cat %s", taskrc)).
		CombinedOutput()
	fmt.Println(string(out))
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

func TestTaskExists_notExistsOK(t *testing.T) {
	// Use a custom .taskrc for testing.
    _, taskrcFile := createTempTaskRC(t)
	defer taskrcFile.Close()

	// Use a custom TASKDATA directory for testing.
	taskdataPath := filepath.Join(t.TempDir(), "task")
	err := os.Mkdir(taskdataPath, 0755)
	assert.NoError(t, err, "Failed to create TASKDATA directory for testing.")
	os.Setenv("TASKDATA", taskdataPath)

    // Create UDA for the To-Do ID.
    err = CreateUDA(UDANameTodoID, "foo")
	assert.NoError(t, err, "CreateUDA returned an error.")

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

func generateRandomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, length)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
