package cli

import (
	"os"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

// Sync interprets the command and executes it.
func Sync(client *msgraphsdk.GraphServiceClient) {
    parseCommand().exec(client)
}

func parseCommand() command {
	switch os.Args[1] {
	case "pull":
		taskListID := os.Args[2]
		return tasksPull{
			listID: taskListID,
		}
	}

	return nil
}
