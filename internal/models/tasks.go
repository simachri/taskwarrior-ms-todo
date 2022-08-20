package models

import (
	"errors"
	"fmt"
)

type TaskStatus int

const (
	// Taskwarrior User Defined Attribute (UDA): Microsoft To-Do Task ID as received from the
	// API
	UDANameTodoTaskID string = "ms_todo_taskid"
	// Taskwarrior User Defined Attribute (UDA): Microsoft To-Do List ID as received from the
	// API
	UDANameTodoListID string = "ms_todo_listid"

	TODO_TASKSTATUS_NOTSTARTED string = "notStarted"

	TW_TASKSTATUS_PENDING TaskStatus = iota
	TW_TASKSTATUS_COMPLETED
	TW_TASKSTATUS_DELETED
)

type Task struct {
	ToDoListID *string
	ToDoTaskID *string
	Title      *string
	// Format is yyyy-MM-DDThh:mm:ss, example: 2022-08-02T00:00:00.0000000
	CompletedAt *string
	Status      TaskStatus
}

type TaskwarriorTask struct {
	Task
	TaskWarriorUUID *string
}

// IsUpToDate compares the data fileds of two tasks ignoring the MS To-Do IDs and
// Taskwarrior UUID.
func (this *Task) IsUpToDate(that *Task) bool {
	if that == nil {
		return false
	}

	if this == that {
		return true
	}

	if *this.Title == *that.Title &&
		*this.CompletedAt == *that.CompletedAt &&
		this.Status == that.Status {
		return true
	}

	return false
}

func ConvStatusFromToDo(todoStatus *string) (TaskStatus, error) {
	if todoStatus == nil || *todoStatus == "" {
		return -1, errors.New("[ConvStatusFromToDo] Failed to convert status. " +
			"Status is 'nil' or empty.")
	}

	switch *todoStatus {
	case TODO_TASKSTATUS_NOTSTARTED:
		return TW_TASKSTATUS_PENDING, nil
	}

	return -1, errors.New(fmt.Sprintf("[ConvStatusFromToDo] Failed to convert status. "+
		"Status '%s' is unknown.", *todoStatus))
}

func ConvStatusFromTW(twStatus *string) (TaskStatus, error) {
	if twStatus == nil || *twStatus == "" {
		return -1, errors.New("[ConvStatusFromTW] Failed to convert status. " +
			"Status is 'nil' or empty.")
	}

	switch *twStatus {
	case "pending":
		return TW_TASKSTATUS_PENDING, nil
	case "completed":
		return TW_TASKSTATUS_COMPLETED, nil
	case "deleted":
		return TW_TASKSTATUS_DELETED, nil
	}

	return -1, errors.New(fmt.Sprintf("[ConvStatusFromTW] Failed to convert status. "+
		"Status '%s' is unknown.", *twStatus))
}
