package cli

import (
	"fmt"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/simachri/taskwarrior-ms-todo/internal/mstodo"
	tw "github.com/simachri/taskwarrior-ms-todo/internal/taskwarrior"
)

type command interface {
	exec(client *msgraphsdk.GraphServiceClient) error
}

type tasksPull struct {
	listID string
}

func (t tasksPull) exec(client *msgraphsdk.GraphServiceClient) error {
	fmt.Printf(
		"[taskPull] Fetching tasks from MS To-Do list '%s'...\n",
		t.listID,
	)
	tasks, err := mstodo.ReadOpenTasks(client, t.listID)
	if err != nil {
		return err
	}

	tList := tasks.GetValue()
	fmt.Printf(
		"[taskPull] '%v' tasks fetched.\n",
		len(tList),
	)

	for _, task := range tList {
		todoTaskID := task.GetId()
        taskExists, err := tw.TaskExists(*todoTaskID)
        if (err != nil) {
            fmt.Println(err)
            continue
        }
		if  (taskExists){
			fmt.Printf(
				"[taskPull] SKIP - task already exists in Taskwarrior: '%s'\n",
				*task.GetTitle(),
			)
			continue
		}

		fmt.Printf(
			"[taskPull] NEW - Create new Taskwarrior task: '%s'\n",
			*task.GetTitle(),
		)
		tUUID, err := tw.CreateTask(*task.GetTitle(), *todoTaskID)
		if err != nil {
			fmt.Printf("[taskPull] Failed to create Taskwarrior task: %v\n", err)
			continue
		}
		fmt.Printf(
			"[taskPull] NEW - New Taskwarrior task created with UUID: %s\n",
			tUUID,
		)
	}
	return nil
}
