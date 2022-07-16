package cli

import (
	"fmt"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFileName         string
	credentialsFileName string

	rootCmd = &cobra.Command{
		Use:   "twtodo",
		Short: "Sync Microsoft To-do with Taskwarrior",
		Args:  cobra.MinimumNArgs(1),
	}
)

func Execute(client *msgraphsdk.GraphServiceClient) error {
	addPullCmd(rootCmd, client)

	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().
		StringVar(&cfgFileName, "config", "", "config filename - default is $XDG_CONFIG_HOME/twtodo/config.yaml")
	rootCmd.PersistentFlags().
		StringVar(&credentialsFileName, "credentials", "", "credentials filename - default is $XDG_CONFIG_HOME/twtodo/credentials.yaml")
}

func initConfig() {
	if cfgFileName != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFileName)
	} else {
		viper.AddConfigPath("$XDG_CONFIG_HOME")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
