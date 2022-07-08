package main

import (
	"context"
	"fmt"
	"os"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/joho/godotenv"
	a "github.com/microsoft/kiota-authentication-azure-go"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

func main() {
	godotenv.Load()
	client, err := authenticate(os.Getenv("TENANT_ID"), os.Getenv("CLIENT_ID"))

	todoLists, err := client.Me().Todo().Lists().Get()
	if err != nil {
		fmt.Printf("Error getting the to-do lists: %v\n", err)
	}
	for _, todoItm := range todoLists.GetValue() {
		fmt.Printf("Found to-do list: %v\n", *todoItm.GetDisplayName())
	}
}

func authenticate(tenantID string, clientID string) (*msgraphsdk.GraphServiceClient, error) {
	cred, err := azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
		TenantID: tenantID,
		ClientID: clientID,
		UserPrompt: func(ctx context.Context, message azidentity.DeviceCodeMessage) error {
			fmt.Println(message.Message)
			return nil
		},
	})
	if err != nil {
		fmt.Printf("Error creating credentials: %v\n", err)
	}

	auth, err := a.NewAzureIdentityAuthenticationProviderWithScopes(cred, []string{"Tasks.Read"})
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
