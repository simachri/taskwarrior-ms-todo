package taskwarrior

import models "github.com/simachri/taskwarrior-ms-todo/internal/models"

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
