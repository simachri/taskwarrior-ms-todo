package taskwarrior

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/simachri/taskwarrior-ms-todo/internal/models"
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

func TestCreateUDAs_emptyString_raisesErr(t *testing.T) {
	// Use a custom .taskrc for testing.
	createTempTaskRC(t)
	udaName := ""
	udaLabel := ""

	err := CreateUDA(udaName, udaLabel)
	assert.Error(t, err)
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
	err = exec.Command("bash", "-c",
		fmt.Sprintf("cat %s | grep uda.%s.label=%s", taskrc, udaName, udaLabel)).
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
	err = exec.Command("bash", "-c",
		fmt.Sprintf("cat %s | grep uda.%s.type=string", taskrc, udaName)).
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

	// Create UDA for the To-Do List ID.
	err = CreateUDA(UDANameTodoListID, "todo_list_id")
	assert.NoError(t, err, "CreateUDA returned an error.")
	// Create UDA for the To-Do Task ID.
	err = CreateUDA(UDANameTodoTaskID, "todo_task_id")
	assert.NoError(t, err, "CreateUDA returned an error.")
}

func TestTaskExists_notExists_returnsFalse(t *testing.T) {
	setup(t)

	toDoListID := generateRandomString(10)
	toDoTaskID := generateRandomString(10)

	exists, err := TaskExists(&toDoListID, &toDoTaskID)

	assert.NoError(
		t,
		err,
		"TaskExists returned an error")
	assert.False(
		t,
		exists,
		fmt.Sprintf("Task with To-Do ID '%s' must not exist.", toDoTaskID),
	)
}

func TestCreateTask_isOK(t *testing.T) {
	setup(t)

	taskTitle := "foo"
	toDoListID := generateRandomString(10)
	toDoTaskID := generateRandomString(10)
	taskUUID, err := CreateTask(&taskTitle, &toDoListID, &toDoTaskID)
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

	taskTitle := "foo"
	toDoListID := generateRandomString(10)
	toDoTaskID := generateRandomString(10)
	CreateTask(&taskTitle, &toDoListID, &toDoTaskID)

	exists, err := TaskExists(&toDoListID, &toDoTaskID)

	assert.NoError(
		t,
		err,
		"TaskExists returned an error")
	assert.True(t, exists)
}

func TestCreateTask_taskHasUDAs(t *testing.T) {
	setup(t)

	taskTitle := "foo"
	toDoListID := generateRandomString(10)
	toDoTaskID := generateRandomString(10)
	taskUUID, _ := CreateTask(&taskTitle, &toDoListID, &toDoTaskID)

	cmd := fmt.Sprintf("task _get %s.%s", taskUUID, UDANameTodoListID)
	out, err := exec.Command("bash", "-c", cmd).
		CombinedOutput()
	assert.NoError(t, err)
	outListID := string(out[:len(out)-1])
	assert.Equal(t, toDoListID, outListID)

	cmd = fmt.Sprintf("task _get %s.%s", taskUUID, UDANameTodoTaskID)
	out, err = exec.Command("bash", "-c", cmd).
		CombinedOutput()
	assert.NoError(t, err)
	outTaskID := string(out[:len(out)-1])
	assert.Equal(t, toDoTaskID, outTaskID)
}

func generateRandomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, length)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func TestGetTasks_isOK(t *testing.T) {
	setup(t)

	taskTitleA := "foo"
	taskTitleB := "bar"
	toDoListID := generateRandomString(10)
	toDoTaskIDA := generateRandomString(10)
	toDoTaskIDB := generateRandomString(10)
	CreateTask(&taskTitleA, &toDoListID, &toDoTaskIDA)
	CreateTask(&taskTitleB, &toDoListID, &toDoTaskIDB)

	tasks, err := GetAllToDoTasks()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(*tasks))

	assert.Equal(t, taskTitleA, *(*tasks)[0].Title)
	assert.Equal(t, toDoListID, *(*tasks)[0].ToDoListID)
	assert.Equal(t, toDoTaskIDA, *(*tasks)[0].ToDoTaskID)
	assert.Equal(t, models.PENDING, (*tasks)[0].Status)

	assert.Equal(t, taskTitleB, *(*tasks)[1].Title)
	assert.Equal(t, toDoListID, *(*tasks)[1].ToDoListID)
	assert.Equal(t, toDoTaskIDB, *(*tasks)[1].ToDoTaskID)
	assert.Equal(t, models.PENDING, (*tasks)[1].Status)
}
