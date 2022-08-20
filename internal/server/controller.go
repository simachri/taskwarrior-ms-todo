package server

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"

	"github.com/simachri/taskwarrior-ms-todo/internal/models"
	"github.com/simachri/taskwarrior-ms-todo/internal/mstodo"
	"github.com/simachri/taskwarrior-ms-todo/internal/taskwarrior"
)

var TasksPullCmd = "Handler.OnTasksPull"

type Handler struct {
	client mstodo.ClientFacade
}

type updateStatistics struct {
	taskCountTotal    int
	taskCountUpdated  int32
	taskCountUpToDate int32
	taskCountError    int32
}

type importStatistics struct {
	taskCountFetched int
	taskCountCreated int32
	taskCountExisted int32
	taskCountError   int32
}

func updateTaskwarriorTasks(
	client mstodo.ClientFacade,
	toDoListID *string,
) (stat *updateStatistics, err error) {
	stat = &updateStatistics{
		taskCountTotal:    0,
		taskCountUpdated:  0,
		taskCountUpToDate: 0,
		taskCountError:    0,
	}

	fmt.Println("[updateTaskWarriorTasks] Reading all imported Taskwarrior tasks.")
	tasks, err := taskwarrior.ReadTasksAll()
	if err != nil {
		return stat, err
	}

	stat.taskCountTotal = len(*tasks)

	for _, task := range *tasks {
		taskFromMSToDo, err := client.ReadTaskByID(task.ToDoListID, task.ToDoTaskID)
		if err != nil {
			fmt.Printf(
				"[updateTaskWarriorTasks] Failed to read task from MS To-Do by ID: %v\n",
				err,
			)
			stat.taskCountError = stat.taskCountError + 1.
			continue
		}

		if taskFromMSToDo.IsUpToDate(&task.Task) {
			fmt.Printf("[updateTaskWarriorTasks] Task is up to date: %s\n", *task.Title)
			stat.taskCountUpToDate = stat.taskCountUpToDate + 1
			continue
		}

		err = taskwarrior.Update(&models.TaskwarriorTask{
			Task:            *taskFromMSToDo,
			TaskWarriorUUID: task.TaskWarriorUUID,
		})
		if err != nil {
			fmt.Printf("[updateTaskWarriorTasks] Failed to update task: %v\n", err)
			stat.taskCountError = stat.taskCountError + 1.
			continue
		}
		fmt.Printf("[updateTaskWarriorTasks] Task updated: %s\n", *task.Title)
	}

	return stat, nil
}

func importOpenTasks(
	client mstodo.ClientFacade,
	toDoListID *string,
) (stat *importStatistics, err error) {
	stat = &importStatistics{
		taskCountFetched: 0,
		taskCountCreated: 0,
		taskCountExisted: 0,
		taskCountError:   0,
	}

	fmt.Printf(
		"[importOpenTasks] Fetching tasks from MS To-Do list '%s'...\n",
		*toDoListID,
	)
	tasks, err := client.ReadOpenTasks(toDoListID)
	if err != nil {
		return stat, err
	}

	stat.taskCountFetched = len(*tasks)

	for _, task := range *tasks {
		fmt.Println("[importOpenTasks] Start import of task...")

		result, err := taskwarrior.Import(&task)
		if err != nil {
			fmt.Printf("[importOpenTasks] Error: %v", err)
			stat.taskCountError = stat.taskCountError + 1
			continue
		}

		switch result {
		case taskwarrior.TASK_CREATED:
			fmt.Printf(
				"[importOpenTasks] NEW - Taskwarrior task created: '%s'\n",
				*task.Title,
			)
			stat.taskCountCreated = stat.taskCountCreated + 1
			continue

		case taskwarrior.TASK_EXISTS_AND_SKIPPED:
			fmt.Printf(
				"[importOpenTasks] SKIP - task already exists in Taskwarrior: '%s'\n",
				*task.Title,
			)
			stat.taskCountExisted = stat.taskCountExisted + 1
			continue
		}
	}

	return stat, nil
}

func (h *Handler) OnTasksPull(req Request, res *Response) error {
	fmt.Println("[OnTasksPull] Handling 'pull' command...")

	updateStat, err := updateTaskwarriorTasks(h.client, &req.ListID)
	if err != nil {
		return err
	}

	importStat, err := importOpenTasks(h.client, &req.ListID)
	if err != nil {
		return err
	}

	res.Message = fmt.Sprintf(
		"[OnTasksPull] Pull succesful:\n"+
			"    [Update] MS To-Do tasks existing in Taskwarrior: %v\n"+
			"    [Update] Taskwarrior tasks up-to-date: %v\n"+
			"    [Update] Taskwarrior tasks updated: %v\n"+
			"    [Update] Errors: %v\n"+
			"    [Import] Open Tasks fetched from MS To-Do: %v\n"+
			"    [Import] New Tasks created in Taskwarrior: %v\n"+
			"    [Import] Tasks already existed in Taskwarrior: %v\n"+
			"    [Import] Errors: %v",
		updateStat.taskCountTotal,
		updateStat.taskCountUpToDate,
		updateStat.taskCountUpToDate,
		updateStat.taskCountError,
		importStat.taskCountFetched,
		importStat.taskCountCreated,
		importStat.taskCountExisted,
		importStat.taskCountError,
	)

	fmt.Println("[OnTasksPull] 'pull' command finished.")
	return nil
}

// Start starts the server to handle commands from the CLI.
func Start(client mstodo.ClientFacade, port *int32) error {
	rpc.Register(&Handler{client: client})

	fmt.Println("[Server] Starting...")

	fmt.Println("[Server] Performing health checks...")
	err := checkHealth()
	if err != nil {
		return err
	}
	fmt.Println("[Server] All health check passed.")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", *port))
	if err != nil {
		return err
	}
	defer listener.Close()
	fmt.Println(fmt.Sprintf("[Server] Listening on port '%v'.", *port))

	rpc.Accept(listener)

	return nil
}

func checkHealth() error {
	fmt.Println("[healthCheck] Checking existence of User-Defined-Attributes (UDAs).")
	if !udasExist() {
		return errors.New(fmt.Sprintf("[healthCheck] The following Taskwarrior "+
			"User-Defined-Attributes have to exist. Create them by running the command "+
			"'twtodo setup'.\n"+
			"              %s\n"+
			"              %s\n",
			models.UDANameTodoListID,
			models.UDANameTodoTaskID))
	}

	return nil
}

func udasExist() bool {
	udas := [2]string{models.UDANameTodoListID, models.UDANameTodoTaskID}
	for _, udaName := range udas {
		if exists, _ := taskwarrior.UDAExists(udaName); !exists {
			return false
		}
	}
	return true
}
