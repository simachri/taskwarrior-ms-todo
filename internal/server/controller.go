package server

import (
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

type fetchStatistics struct {
	taskCountFetched int
	taskCountCreated int32
	taskCountExisted int32
	taskCountError   int32
}

func importOpenTasks(
	client mstodo.ClientFacade,
	toDoListID *string,
) (statistics *fetchStatistics, err error) {
	statistics = &fetchStatistics{
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
		return statistics, err
	}

	statistics.taskCountFetched = len(*tasks)

	for _, task := range *tasks {
		result, err := taskwarrior.Import(&task)
		if err != nil {
			fmt.Printf("[importOpenTasks] Error: %v", err)
			statistics.taskCountError = statistics.taskCountError + 1
			continue
		}

		switch result {
		case taskwarrior.TASK_CREATED:
			fmt.Printf(
				"[importOpenTasks] SKIP - task already exists in Taskwarrior: '%s'\n",
				*task.Title,
			)
			statistics.taskCountCreated = statistics.taskCountCreated + 1
			continue

		case taskwarrior.TASK_EXISTS_AND_SKIPPED:
			fmt.Printf(
				"[importOpenTasks] NEW - Taskwarrior task created: '%s'\n",
				*task.Title,
			)
			statistics.taskCountExisted = statistics.taskCountExisted + 1
			continue
		}
	}

	return statistics, nil
}

func (h *Handler) OnTasksPull(req Request, res *Response) error {
	statistics, err := importOpenTasks(h.client, &req.ListID)
	if err != nil {
		return err
	}

	res.Message = fmt.Sprintf(
		"Pull succesful:\n"+
			"Tasks fetched: %v\n"+
			"Tasks created: %v\n"+
			"Tasks existing: %v\n"+
			"Errors: %v",
		statistics.taskCountFetched,
		statistics.taskCountCreated,
		statistics.taskCountExisted,
		statistics.taskCountError,
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
