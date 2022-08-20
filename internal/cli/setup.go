package cli

import (
	"fmt"

	"github.com/simachri/taskwarrior-ms-todo/internal/taskwarrior"
	"github.com/spf13/cobra"
)

func addSetupCmd(parentCmd *cobra.Command) {
	c := &cobra.Command{
		Use:   "setup",
		Short: "Setup the integration",
		Long: `Creates the User-Defined-Attributes (UDAs) in Taskwarrior required for ` +
			`the integration`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := taskwarrior.CreateIntegrationUDAs()
			if err != nil {
				return err
			}

			fmt.Println("[Setup] Finished - run 'twtodo up' to start the server.")
			return nil
		},
	}

	parentCmd.AddCommand(c)
}
