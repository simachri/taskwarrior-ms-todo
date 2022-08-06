package models

import (
	"errors"
	"fmt"
)

type TaskStatus int

const (
	PENDING TaskStatus = iota
	COMPLETED
	DELETED
)

type Task struct {
	ToDoListID *string
	ToDoTaskID *string
	Title      *string
	// Format is yyyy-MM-DDThh:mm:ss, example: 2022-08-02T00:00:00.0000000
	CompletedAt *string
	Status      TaskStatus
}

func ConvStatusFromTW(twStatus *string) (TaskStatus, error) {
	if twStatus == nil || *twStatus == "" {
		return -1, errors.New("[ConvStatusFromTW] Failed to convert status. " +
			"Status is 'nil' or empty.")
	}

	switch *twStatus {
	case "pending":
		return PENDING, nil
	case "completed":
		return COMPLETED, nil
	case "deleted":
		return DELETED, nil
	}

	return -1, errors.New(fmt.Sprintf("[ConvStatusFromTW] Failed to convert status. "+
		"Status '%s' is unknown.", *twStatus))
}
