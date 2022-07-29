package cli

import (
	"fmt"

	"github.com/simachri/taskwarrior-ms-todo/internal/mstodo"
	"github.com/simachri/taskwarrior-ms-todo/internal/server"
	"github.com/spf13/cobra"
)

type UpCmdConfig struct {
	Port int32
}

type upCmd struct {
	cmd *cobra.Command
	// Using a function is required as Viper parses the config not before a command's
	// Execute() function is called.
	GetConfig func() (*UpCmdConfig, error)
}

func (upCmd *upCmd) exec(client mstodo.ClientFacade) error {
	config, err := upCmd.GetConfig()
	if err != nil {
		return fmt.Errorf("[upCmd] Error: %v", err)
	}
	return server.Start(client, &config.Port)
}

func addUpCmd(
	parentCmd *cobra.Command,
	clientFactory *mstodo.ClientFactory,
	getConfig func() (*UpCmdConfig, error),
) {
	upCmd := &upCmd{GetConfig: getConfig}

	c := &cobra.Command{
		Use:   "up",
		Short: "Start the sync server",
		Long:  `Starts the sync server and authenticates to MS Azure.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			authenticatedClient, err := clientFactory.GetGraphClient()
			if err != nil {
				return err
			}
			return upCmd.exec(authenticatedClient)
		},
	}
	upCmd.cmd = c

	parentCmd.AddCommand(c)
}
