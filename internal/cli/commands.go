package cli

import (
	"fmt"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/simachri/taskwarrior-ms-todo/internal/mstodo"
)

type command interface {
	exec(client *msgraphsdk.GraphServiceClient) error
}

type tasksPull struct {
	listID string
}

func (t tasksPull) exec(client *msgraphsdk.GraphServiceClient) error {
	tasks, err := mstodo.ReadTasks(client, t.listID)
	if err != nil {
        return err
	}
	for _, task := range tasks.GetValue() {
		fmt.Printf("[taskPull] Found task in list '%s': %v\n", t.listID, *task.GetTitle())
	}
	return nil
}
