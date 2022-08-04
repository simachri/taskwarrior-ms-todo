package taskwarrior

import (
	"fmt"
	"os/exec"
)

// Taskwarrior User Defined Attribute (UDA): Microsoft To-Do Task ID as received from the
// API
const UDANameTodoTaskID = "ms_todo_taskid"

// Taskwarrior User Defined Attribute (UDA): Microsoft To-Do List ID as received from the
// API
const UDANameTodoListID = "ms_todo_listid"

// TaskExists returns 'true' if a Taskwarrior task for the given Microsoft To-Do task ID
// exists, otherwise 'false'.
func TaskExists(todoTaskID string) (bool, error) {
	cmd := exec.Command("bash", "-c",
		fmt.Sprintf("task %s:%s", UDANameTodoTaskID, todoTaskID))
	out, err := cmd.CombinedOutput()
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok == true && exitErr.ExitCode() == 1 {
			// Exit code is 1 if task does not exist.
			return false, nil
		}

		fmt.Printf("[TaskExists] Command:\n%s\n", cmd.Args)
		fmt.Printf("[TaskExists] Output:\n%s\n", out)
		return false, fmt.Errorf(
			"[TaskExists] Failed to check task existence for todoID '%v': %v\n",
			todoTaskID,
			err,
		)
	}

	return true, nil
}

// CreateTask creates a Taskwarrior task using the 'task' CLI.
// The Microsoft To-Do task and list IDs are stored as user-defined attribute (UDA) in
// the Taskwarrior task.
func CreateTask(
	title *string,
	todoListID *string,
	todoTaskID *string,
) (taskUUID string, err error) {
	// 'task add "TITLE" returns a message 'Created task 42.'
	cmd := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf(
			"task add '%s' %s:'%s' %s:'%s'",
			*title,
			UDANameTodoListID,
			*todoListID,
			UDANameTodoTaskID,
			*todoTaskID,
		)+
			// Extract the task ID
			" | grep -oP '[0-9]+'"+
			// Extract the task UUID.
			" | xargs -I '{id}' task _get {id}.uuid",
	)

	uuid, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf(
			"[CreateTask] Failed to create task and extract the UUID: %w\n",
			err,
		)
	}

	// The UUID has a trailing newline character that needs to be stripped.
	return string(uuid[:len(uuid)-1]), nil
}

// CreateUDA creates a User Defined Attribute (UDA) in Taskwarrior.
func CreateUDA(name string, label string) (err error) {
	// 'echo "yes"' is required to answer the prompt 'Are you sure?'.
	out, err := exec.Command("bash", "-c",
		fmt.Sprintf("'yes' | task config uda.%s.type string", name)).CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"[CreateUDAs] Failed to create UDA '%s': %w\nOutput of command: %s\n",
			name,
			err,
			out,
		)
	}
	out, err = exec.Command("bash", "-c",
		fmt.Sprintf("'yes' | task config uda.%s.label %s", name, label)).CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"[CreateUDAs] Failed to create UDA '%s': %w\nOutput of command: %s\n",
			name,
			err,
			out,
		)
	}

	return nil
}
