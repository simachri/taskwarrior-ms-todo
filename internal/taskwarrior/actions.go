package taskwarrior

import (
	"errors"
	"fmt"

	models "github.com/simachri/taskwarrior-ms-todo/internal/models"
)

type ImportResult int32

const (
	TASK_CREATED ImportResult = iota
	TASK_EXISTS_AND_SKIPPED
)

func Import(task *models.Task) (ImportResult, error) {
	toDoListID := task.ToDoListID
	toDoTaskID := task.ToDoTaskID

	taskExists, err := taskExists(toDoListID, toDoTaskID)
	if err != nil {
		return -1, err
	}
	if taskExists {
		return TASK_EXISTS_AND_SKIPPED, nil
	}

	_, err = createTask(task.Title, toDoListID, toDoTaskID)
	if err != nil {
		return -1, err
	}

	return TASK_CREATED, nil
}

func Update(task *models.TaskwarriorTask) error {
	taskExists, err := taskExists(task.ToDoListID, task.ToDoTaskID)
	if err != nil {
		return err
	}
	if !taskExists {
		return errors.New(
			fmt.Sprintf("[Update] Failed - no Taskwarrior task exists for\n"+
				"MS To-Do List ID: %s\n"+
				"MS To-Do Task ID: %s\n", *task.ToDoListID, *task.ToDoTaskID),
		)
	}

	return update(task)
}
