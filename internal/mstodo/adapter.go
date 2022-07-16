package mstodo

import (
	"context"
	"fmt"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	a "github.com/microsoft/kiota-authentication-azure-go"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	graphconfig "github.com/microsoftgraph/msgraph-sdk-go/me/todo/lists/item/tasks"
)

// ReadOpenTasks uses the Microsoft Graph API to fetch the To-Do tasks with status
// 'notStarted'.
func ReadOpenTasks(
	client *msgraphsdk.GraphServiceClient,
	listID *string,
) (models.TodoTaskCollectionResponseable, error) {
    openTasksFilter := "status eq 'notStarted'"
	reqParams := &graphconfig.TasksRequestBuilderGetQueryParameters{
		Filter: &openTasksFilter,
	}
	reqConf := &graphconfig.TasksRequestBuilderGetRequestConfiguration{
		QueryParameters: reqParams,
	}

	tasks, err := client.Me().
		Todo().
		ListsById(*listID).
		Tasks().
		GetWithRequestConfigurationAndResponseHandler(reqConf, nil)

	if err != nil {
		return nil, fmt.Errorf(
			"[ReadTasks] Failed to fetch the tasks of To-Do list '%s': %w\n",
			*listID,
			err,
		)
	}
	return tasks, nil
}

// Authenticate creates a Microsoft Graph client using the Device Code Authentication
// Provider.
func Authenticate(
	tenantID string,
	clientID string,
) (*msgraphsdk.GraphServiceClient, error) {
	cred, err := azidentity.NewDeviceCodeCredential(
		&azidentity.DeviceCodeCredentialOptions{
			TenantID: tenantID,
			ClientID: clientID,
			UserPrompt: func(ctx context.Context, message azidentity.DeviceCodeMessage) error {
				fmt.Println(message.Message)
				return nil
			},
		},
	)
	if err != nil {
		fmt.Printf("Error creating credentials: %v\n", err)
	}

	auth, err := a.NewAzureIdentityAuthenticationProviderWithScopes(
		cred,
		[]string{"Tasks.Read"},
	)
	if err != nil {
		fmt.Printf("Error authentication provider: %v\n", err)
		return nil, err
	}

	adapter, err := msgraphsdk.NewGraphRequestAdapter(auth)
	if err != nil {
		fmt.Printf("Error creating adapter: %v\n", err)
		return nil, err
	}

	client := msgraphsdk.NewGraphServiceClient(adapter)

	return client, nil
}
