package cli

import (
	"fmt"
	"net/rpc"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/simachri/taskwarrior-ms-todo/internal/server"
	"github.com/spf13/cobra"
)

type command interface {
	exec(client *msgraphsdk.GraphServiceClient) error
}

type tasksPullCmd struct {
	listID *string
	cmd    *cobra.Command
}

func (cmd *tasksPullCmd) exec() error {
	rpcClient, err := rpc.Dial("tcp", "127.0.0.1:41001")
	if err != nil {
		return err
	}
	defer rpcClient.Close()

	resp := new(server.Response)
	err = rpcClient.Call(server.TasksPullCmd, &server.Request{
		ListID: *cmd.listID,
	}, resp)
	if err != nil {
		return err
	}
	fmt.Println(resp.Message)

	return nil
}

func addPullCmd(parentCmd *cobra.Command) {
	pullCmd := &tasksPullCmd{}

	c := &cobra.Command{
		Use:   "pull",
		Short: "Pull tasks",
		Long:  `Pulls the tasks from a MS To-Do list and creates them as tasks in Taskwarrior`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return pullCmd.exec()
		},
	}
    listIDFlagName := "list"
	pullCmd.listID = c.PersistentFlags().
		StringP(listIDFlagName, "l", "", "MS To-Do Tasklist ID (required)")
    c.MarkPersistentFlagRequired(listIDFlagName)
	pullCmd.cmd = c

	parentCmd.AddCommand(c)
}
