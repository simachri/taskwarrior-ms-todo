package mstodo

import (
	"context"
	"errors"
	"fmt"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	a "github.com/microsoft/kiota-authentication-azure-go"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	graphconfig "github.com/microsoftgraph/msgraph-sdk-go/me/todo/lists/item/tasks"

	models "github.com/simachri/taskwarrior-ms-todo/internal/models"
)

var authenticatedGraphClient *GraphClient

type ClientFactory struct {
	// Using functions is required as Viper parses the config not before a command's
	// Execute() function is called.
	GetTenantID func() string
	GetClientID func() string
}

type ClientFacade interface {
	ReadOpenTasks(listID *string) (*[]models.Task, error)
	ReadTaskByID(listID *string, taskID *string) (*models.Task, error)
}

type GraphClient struct {
	authenticatedClient *msgraphsdk.GraphServiceClient
}

// Get returns a singleton instance of a Microsoft Graph client using the Device Code
// Authentication Provider.
func (fact *ClientFactory) GetGraphClient() (*GraphClient,
	error,
) {
	if authenticatedGraphClient != nil {
		return authenticatedGraphClient, nil
	}

	tenantID := fact.GetTenantID()
	clientID := fact.GetClientID()

	if tenantID == "" || clientID == "" {
		return nil, errors.New(
			"[AzureAuth] Empty Azure Tenant ID and/or Client ID. Check the credentials.yaml file.",
		)
	}

	client, err := authenticate(tenantID, clientID)
	if err != nil {
		return nil, err
	}

	me, err := client.Me().Get()
	if err != nil {
		return nil, fmt.Errorf(
			"[AzureAuth] Failed to retrieve data about authenticated user: %v",
			err,
		)
	}

	fmt.Printf("[AzureAuth] Authenticated as %s\n", *me.GetDisplayName())

	authenticatedClient := &GraphClient{authenticatedClient: client}
	return authenticatedClient, nil
}

// ReadTaskByID retrieves task data for a given task by a task ID from a list, given by a
// list ID.
func (graph GraphClient) ReadTaskByID(
	listID *string,
	taskID *string,
) (*models.Task, error) {
	taskData, err := graph.authenticatedClient.Me().
		Todo().
		ListsById(*listID).
		TasksById(*taskID).
		Get()
	if err != nil {
		return nil, fmt.Errorf(
			"[ReadTaskByID] Failed to fetch the task with ID '%s' from To-Do list "+
				"'%s':\n%w\n",
			*taskID,
			*listID,
			err,
		)
	}
	if taskData == nil {
		return nil, fmt.Errorf(
			"[ReadTaskByID] Task with ID '%s' does not exist in To-Do list '%s':\n%w\n",
			*taskID,
			*listID,
			err,
		)
	}

	fmt.Printf(
		"[ReadTaskByID] Data of task read: '%s'\n",
		*taskData.GetTitle(),
	)

	completedAt := ""
	if taskData.GetCompletedDateTime() != nil {
		completedAt = *taskData.GetCompletedDateTime().GetDateTime()
	}

	todoTaskStatus := taskData.GetStatus().String()
	taskStatus, err := models.ConvStatusFromToDo(&todoTaskStatus)
	if err != nil {
		return nil, fmt.Errorf(
			"[ReadTaskByID] Task with ID '%s': Failed to parse task status '%s':\n%w\n",
			*taskID,
			todoTaskStatus,
			err,
		)
	}

	return &models.Task{
		ToDoTaskID:  taskData.GetId(),
		ToDoListID:  listID,
		Title:       taskData.GetTitle(),
		CompletedAt: &completedAt,
		Status:      taskStatus,
	}, nil
}

// ReadOpenTasks uses the Microsoft Graph API to fetch the To-Do tasks with status
// 'notStarted'.
func (graph GraphClient) ReadOpenTasks(
	listID *string,
) (*[]models.Task, error) {
	openTasksFilter := fmt.Sprintf("status eq '%s'", models.TODO_TASKSTATUS_NOTSTARTED)
	reqParams := &graphconfig.TasksRequestBuilderGetQueryParameters{
		Filter: &openTasksFilter,
	}
	reqConf := &graphconfig.TasksRequestBuilderGetRequestConfiguration{
		QueryParameters: reqParams,
	}

	tasksResponse, err := graph.authenticatedClient.Me().
		Todo().
		ListsById(*listID).
		Tasks().
		GetWithRequestConfigurationAndResponseHandler(reqConf, nil)
	if err != nil {
		return nil, fmt.Errorf(
			"[ReadOpenTasks] Failed to fetch the tasks of To-Do list '%s': %w\n",
			*listID,
			err,
		)
	}

	tasksRespVal := tasksResponse.GetValue()
	fmt.Printf(
		"[ReadOpenTasks] %v tasks fetched.\n",
		len(tasksRespVal),
	)

	var tasks []models.Task
	for _, task := range tasksRespVal {
		tasks = append(tasks, models.Task{
			ToDoListID: listID,
			ToDoTaskID: task.GetId(),
			Title:      task.GetTitle(),
		})
	}
	return &tasks, nil
}

func authenticate(
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
		fmt.Printf("[AzureAuth] Error creating credentials: %v\n", err)
	}

	auth, err := a.NewAzureIdentityAuthenticationProviderWithScopes(
		cred,
		[]string{"Tasks.Read"},
	)
	if err != nil {
		fmt.Printf("[AzureAuth] Error authentication provider: %v\n", err)
		return nil, err
	}

	adapter, err := msgraphsdk.NewGraphRequestAdapter(auth)
	if err != nil {
		fmt.Printf("[AzureAuth] Error creating adapter: %v\n", err)
		return nil, err
	}

	client := msgraphsdk.NewGraphServiceClient(adapter)

	return client, nil
}
