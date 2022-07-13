package taskwarrior

import (
	"fmt"
	"os/exec"
)

// Taskwarrior User Defined Attribute (UDA): Microsoft To-Do Task ID as received from the
// API
const udaNameTodoID = "ms_todo_id"

// TaskExists returns 'true' if a Taskwarrior task for the given Microsoft To-Do task ID
// exists, otherwise 'false'.
func TaskExists(todoID string) (bool, error) {
	cmd := exec.Command("bash", "-c",
		fmt.Sprintf("task %s:%s", udaNameTodoID, todoID))
	err := cmd.Run()
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok == true && exitErr.ExitCode() == 1 {
			// Exit code is 1 if task does not exist.
			return false, nil
		}

	    fmt.Println(cmd.Args)
		return false, fmt.Errorf(
			"[TaskExists] Failed to check task existence for todoID '%v': %v\n",
			todoID,
			err,
		)
	}

	return true, nil
}

// CreateTask creates a Taskwarrior task using the 'task' CLI.
// The last 13 characters of the Microsoft To-Do task ID are stored as user-defined 
// attribute (UDA) in the Taskwarrior task.
func CreateTask(title string, todoID string) (taskUUID string, err error) {
	// 'task add "TITLE" returns a message 'Created task 42.'
	cmd := exec.Command("bash", "-c",
		fmt.Sprintf("task add '%s' %s:'%s'", title, udaNameTodoID, todoID)+
			// Extract the task ID
			" | grep -oP '[0-9]+'"+
			// Extract the task UUID.
			" | xargs -I '{id}' task _get {id}.uuid")

	uuid, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf(
			"[CreateTask] Failed to create task and extract the UUID: %w\n",
			err,
		)
	}

	return fmt.Sprintf("%s", uuid), nil
}

// CreateUDAs creates the following User Defined Attributes (UDAs):
// ms_todo_id: Microsoft To-Do Task ID as received from the API
func CreateUDAs() (err error) {
	// 'echo "yes"' is required to answer the prompt 'Are you sure?'.
	err = exec.Command("echo", "'yes' |",
		fmt.Sprintf("task config uda.%s.type string", udaNameTodoID),
		"&&",
		fmt.Sprintf("task config uda.%s.label ToDo-ID", udaNameTodoID)).Run()
	if err != nil {
		return fmt.Errorf(
			"[CreateUDAs] Failed to create UDA '%s': %w\n",
			udaNameTodoID,
			err,
		)
	}
	return nil
}
