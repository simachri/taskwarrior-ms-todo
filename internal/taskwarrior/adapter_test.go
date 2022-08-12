package taskwarrior

import (
	"fmt"
	"math/rand"
	"os/exec"
	"testing"

	"github.com/simachri/taskwarrior-ms-todo/internal/models"
	testUtils "github.com/simachri/taskwarrior-ms-todo/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestUDAExists_isTrue(t *testing.T) {
	testUtils.NewTaskwarriorEnv(t)

	udaName := "foo"
	udaLabel := "bar"

	err := CreateUDA(udaName, udaLabel)
	assert.NoError(t, err, "CreateUDA returned an error.")

	udaExists, err := UDAExists(udaName)
	assert.NoError(t, err)
	assert.True(t, udaExists)
}

func TestUDAExists_isFalse(t *testing.T) {
	testUtils.CreateTempTaskRC(t)

	udaName := "foo"
	udaLabel := "bar"

	err := CreateUDA(udaName, udaLabel)
	assert.NoError(t, err, "CreateUDA returned an error.")

	anotherUDAName := "baz"
	udaExists, err := UDAExists(anotherUDAName)
	assert.NoError(t, err)
	assert.False(t, udaExists)
}

func TestCreateUDAs_noUDAs_existAfterwards(t *testing.T) {
	taskrc := testUtils.CreateTempTaskRC(t)
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

func TestTaskDelete_taskExists_isFalse(t *testing.T) {
	testUtils.NewTaskwarriorEnv(t)
	err := CreateIntegrationUDAs()
	assert.NoError(t, err)

	taskTitle := "foo"
	toDoListID := generateRandomString(10)
	toDoTaskID := generateRandomString(10)

	uuid, err := createTask(&taskTitle, &toDoListID, &toDoTaskID)
	assert.NoError(t, err)

	cmdStr := fmt.Sprintf(
		// Deleting a task shows a prompt to confirm the deletion.
		"'yes' | task %s delete",
		uuid,
	)
	cmd := exec.Command(
		"bash",
		"-c",
		cmdStr,
	)
	err = cmd.Run()
	assert.NoError(t, err)

	exists, err := taskExists(&toDoListID, &toDoTaskID)
	assert.NoError(
		t,
		err,
		"TaskExists must not fail.")
	assert.False(
		t,
		exists,
		fmt.Sprintf(
			"Task with To-Do List ID '%s' and Task ID '%s' must not exist.",
			toDoListID,
			toDoTaskID,
		),
	)
}

func TestTaskExists_notExists_isFalse(t *testing.T) {
	testUtils.NewTaskwarriorEnv(t)
	err := CreateIntegrationUDAs()
	assert.NoError(t, err)

	toDoListID := generateRandomString(10)
	toDoTaskID := generateRandomString(10)

	exists, err := taskExists(&toDoListID, &toDoTaskID)

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
	testUtils.NewTaskwarriorEnv(t)
	err := CreateIntegrationUDAs()
	assert.NoError(t, err)

	taskTitle := "foo"
	toDoListID := generateRandomString(10)
	toDoTaskID := generateRandomString(10)
	taskUUID, err := createTask(&taskTitle, &toDoListID, &toDoTaskID)
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
	testUtils.NewTaskwarriorEnv(t)
	err := CreateIntegrationUDAs()
	assert.NoError(t, err)

	taskTitle := "foo"
	toDoListID := generateRandomString(10)
	toDoTaskID := generateRandomString(10)
	createTask(&taskTitle, &toDoListID, &toDoTaskID)

	exists, err := taskExists(&toDoListID, &toDoTaskID)

	assert.NoError(
		t,
		err,
		"TaskExists returned an error")
	assert.True(t, exists)
}

func TestCreateTask_taskHasUDAs(t *testing.T) {
	testUtils.NewTaskwarriorEnv(t)
	err := CreateIntegrationUDAs()
	assert.NoError(t, err)

	taskTitle := "foo"
	toDoListID := generateRandomString(10)
	toDoTaskID := generateRandomString(10)
	taskUUID, _ := createTask(&taskTitle, &toDoListID, &toDoTaskID)

	cmd := fmt.Sprintf("task _get %s.%s", taskUUID, models.UDANameTodoListID)
	out, err := exec.Command("bash", "-c", cmd).
		CombinedOutput()
	assert.NoError(t, err)
	outListID := string(out[:len(out)-1])
	assert.Equal(t, toDoListID, outListID)

	cmd = fmt.Sprintf("task _get %s.%s", taskUUID, models.UDANameTodoTaskID)
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

func TestReadTasksAll_isOK(t *testing.T) {
	testUtils.NewTaskwarriorEnv(t)
	err := CreateIntegrationUDAs()
	assert.NoError(t, err)

	taskTitleA := "foo"
	taskTitleB := "bar"
	toDoListID := generateRandomString(10)
	toDoTaskIDA := generateRandomString(10)
	toDoTaskIDB := generateRandomString(10)
	createTask(&taskTitleA, &toDoListID, &toDoTaskIDA)
	createTask(&taskTitleB, &toDoListID, &toDoTaskIDB)

	tasks, err := ReadTasksAll()

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
