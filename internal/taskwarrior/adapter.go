package taskwarrior

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"

	models "github.com/simachri/taskwarrior-ms-todo/internal/models"
)

// Taskwarrior User Defined Attribute (UDA): Microsoft To-Do Task ID as received from the
// API
const UDANameTodoTaskID = "ms_todo_taskid"

// Taskwarrior User Defined Attribute (UDA): Microsoft To-Do List ID as received from the
// API
const UDANameTodoListID = "ms_todo_listid"

// taskExists returns 'true' if a Taskwarrior task for the given Microsoft To-Do List and
// Task ID exists in the given task list, otherwise 'false'.
func taskExists(toDoListID *string, toDoTaskID *string) (bool, error) {
	cmd := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf(
			"task %s:%s %s:%s",
			UDANameTodoListID,
			*toDoListID,
			UDANameTodoTaskID,
			*toDoTaskID,
		),
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok == true && exitErr.ExitCode() == 1 {
			// Exit code is 1 if task does not exist.
			return false, nil
		}

		fmt.Printf("[taskExists] Command:\n%s\n", cmd.Args)
		fmt.Printf("[taskExists] Output:\n%s\n", out)
		return false, fmt.Errorf(
			"[taskExists] Failed to check task existence:\nTo-Do List ID: %v\n"+
				"To-Do Task ID: %v\n"+
				"Error: %w\n",
			toDoListID,
			toDoTaskID,
			err,
		)
	}

	return true, nil
}

// createTask creates a Taskwarrior task using the 'task' CLI.
// The Microsoft To-Do task and list IDs are stored as user-defined attribute (UDA) in
// the Taskwarrior task.
func createTask(
	title *string,
	todoListID *string,
	todoTaskID *string,
) (taskUUID string, err error) {
	// 'task add ...' returns a message 'Created task 42.'
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
			"[createTask] Failed to create task and extract the UUID: %w\n",
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

func GetAllToDoTasks() (*[]models.Task, error) {
	// Get JSON representation of all tasks with an MS To-Do Task ID.
	cmdExport := fmt.Sprintf("task %s.any: export", UDANameTodoTaskID)
	// If a TASKRC or TASKDATA override is active for Taskwarrior, for example when
	// running unit tests, additional lines are printed to stderr to show the overrides
	// used for the export. Thus, only use Output() instead of CombinedOutput().
	tasksJSONExport, err := exec.Command("bash", "-c", cmdExport).Output()
	if err != nil {
		return nil, fmt.Errorf(
			"[GetAllToDoTasks] Failed to get JSON representation of tasks: %w\n"+
				"Output of command: %s\n",
			err,
			string(tasksJSONExport),
		)
	}

	var tasksJSON []map[string]interface{}
	err = json.Unmarshal(tasksJSONExport, &tasksJSON)
	if err != nil {
		return nil, fmt.Errorf(
			"[GetAllToDoTasks] Failed to unmarshall JSON representation of tasks: %w\n"+
				"Run '%s' to get the JSON.\n",
			err,
			cmdExport,
		)
	}

	return parseTasksFromJSON(&tasksJSON)
}

func parseTaskStringAttrFromJSON(
	attrName string,
	taskJSON *map[string]interface{},
) (string, error) {
	attr, ok := (*taskJSON)[attrName].(string)
	if ok == false {
		return "", errors.New(fmt.Sprintf(
			"[parseTaskStringAttrFromJSON] Failed to parse '%s' of task as string.\n"+
				"Task JSON: \n%v\n",
			attrName,
			taskJSON,
		),
		)
	}

	return attr, nil
}

func parseTasksFromJSON(tasksJSON *[]map[string]interface{}) (*[]models.Task, error) {
	var tasks []models.Task
	for _, taskJSON := range *tasksJSON {
		toDoListID, err := parseTaskStringAttrFromJSON(UDANameTodoListID, &taskJSON)
		if err != nil {
			return nil, err
		}

		toDoTaskID, err := parseTaskStringAttrFromJSON(UDANameTodoTaskID, &taskJSON)
		if err != nil {
			return nil, err
		}

		taskDescr, err := parseTaskStringAttrFromJSON("description", &taskJSON)
		if err != nil {
			return nil, err
		}

		taskStatusStr, err := parseTaskStringAttrFromJSON("status", &taskJSON)
		if err != nil {
			return nil, err
		}
		taskStatus, err := models.ConvStatusFromTW(&taskStatusStr)
		if err != nil {
			return nil, fmt.Errorf(
				"[GetAllToDoTasks] Failed to parse 'status' of task.\n"+
					"Task JSON: \n%v\n"+
					"Error: %w",
				taskJSON,
				err,
			)
		}

		tasks = append(tasks, models.Task{
			ToDoListID:  &toDoListID,
			ToDoTaskID:  &toDoTaskID,
			Title:       &taskDescr,
			CompletedAt: &taskStatusStr,
			Status:      taskStatus,
		})
	}
	return &tasks, nil
}
func UDAExists(udaName *string) (bool, error) {
	if udaName == nil || *udaName == "" {
		return false, errors.New("Cannot check UDA existence. Provided UDA is empty.")
	}

	err := exec.Command("bash", "-c", fmt.Sprintf("task udas | grep %s", *udaName)).Run()
	if err != nil {
		if err.(*exec.ExitError).ExitCode() == 1 {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
