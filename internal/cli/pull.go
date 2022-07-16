package cli

import (
	"fmt"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/simachri/taskwarrior-ms-todo/internal/mstodo"
	tw "github.com/simachri/taskwarrior-ms-todo/internal/taskwarrior"
	"github.com/spf13/cobra"
)

type command interface {
	exec(client *msgraphsdk.GraphServiceClient) error
}

type tasksPullCmd struct {
	listID *string
	cmd    *cobra.Command
}

func (cmd tasksPullCmd) exec(client *msgraphsdk.GraphServiceClient) error {
	fmt.Printf(
		"[taskPull] Fetching tasks from MS To-Do list '%s'...\n",
		*cmd.listID,
	)
	tasks, err := mstodo.ReadOpenTasks(client, cmd.listID)
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
		if err != nil {
			fmt.Println(err)
			continue
		}
		if taskExists {
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

func addPullCmd(parentCmd *cobra.Command, client *msgraphsdk.GraphServiceClient) {
	pullCmd := &tasksPullCmd{}

    c := &cobra.Command{
		Use:   "pull [MS To-Do Tasklist ID]",
		Short: "Pull tasks",
		Long:  `Pulls the tasks from a MS To-Do list and creates them as tasks in  Taskwarrior`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return pullCmd.exec(client)
		},
	}
	pullCmd.listID = c.PersistentFlags().
		StringP("list", "l", "", "MS To-Do Tasklist ID")
    pullCmd.cmd = c
   
	parentCmd.AddCommand(c)
}
