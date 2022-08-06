package server

import (
	"fmt"
	"net"
	"net/rpc"

	"github.com/simachri/taskwarrior-ms-todo/internal/mstodo"
	"github.com/simachri/taskwarrior-ms-todo/internal/taskwarrior"
)

var TasksPullCmd = "Handler.OnTasksPull"

type Handler struct {
	client mstodo.ClientFacade
}

func (h *Handler) OnTasksPull(req Request, res *Response) error {
	var (
		taskCountCreated int32
		taskCountExisted int32
		taskCountError   int32
	)
	todoListID := req.ListID

	fmt.Printf(
		"[OnTasksPull] Fetching tasks from MS To-Do list '%s'...\n",
		todoListID,
	)
	tasks, err := h.client.ReadOpenTasks(&todoListID)
	if err != nil {
		return err
	}

	for _, task := range *tasks {
		toDoListID := task.ToDoListID
		toDoTaskID := task.ToDoTaskID
		taskExists, err := taskwarrior.TaskExists(toDoListID, toDoTaskID)
		if err != nil {
			fmt.Println(err)
			taskCountError = taskCountError + 1
			continue
		}
		if taskExists {
			fmt.Printf(
				"[OnTasksPull] SKIP - task already exists in Taskwarrior: '%s'\n",
				*task.Title,
			)
			taskCountExisted = taskCountExisted + 1
			continue
		}

		fmt.Printf(
			"[OnTasksPull] NEW - Create new Taskwarrior task: '%s'\n",
			*task.Title,
		)
		tUUID, err := taskwarrior.CreateTask(task.Title, &todoListID, toDoTaskID)
		if err != nil {
			fmt.Printf("[OnTasksPull] Failed to create Taskwarrior task: %v\n", err)
			taskCountError = taskCountError + 1
			continue
		}
		fmt.Printf(
			"[OnTasksPull] NEW - New Taskwarrior task created with UUID: %s\n",
			tUUID,
		)
		taskCountCreated = taskCountCreated + 1
	}

	res.Message = fmt.Sprintf(
		"Pull succesful:\nTasks fetched: %v\nTasks created: %v\nTasks that had already existed: %v\nErrors: %v",
		len(*tasks),
		taskCountCreated,
		taskCountExisted,
		taskCountError,
	)

	return nil
}

// Start starts the server to handle commands from the CLI.
func Start(client mstodo.ClientFacade, port *int32) error {
	rpc.Register(&Handler{client: client})

	fmt.Println("[Server] Starting...")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", *port))
	if err != nil {
		return err
	}
	defer listener.Close()
	fmt.Println(fmt.Sprintf("[Server] Listening on port '%v'.", *port))

	rpc.Accept(listener)

	return nil
}
