package cli

import (
	"fmt"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/simachri/taskwarrior-ms-todo/internal/mstodo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFileName          string
	credentialsFileName  string
	cfgFileViper         = viper.New()
	credentialsFileViper = viper.New()

	rootCmd = &cobra.Command{
		Use:   "twtodo",
		Short: "Sync Microsoft To-do with Taskwarrior",
		Args:  cobra.MinimumNArgs(1),
	}
)

func Execute() error {
	graphClientFactory := &mstodo.ClientFactory{
		// Passing this as function is required as Viper parses the config not before
		// a command's Execute() function is called.
		GetTenantID: func() string { return credentialsFileViper.GetString("tenant_id") },
		GetClientID: func() string { return credentialsFileViper.GetString("client_id") },
	}

    addSetupCmd(rootCmd)

	getUpCmdConfig := func() (*UpCmdConfig, error) {
		configKey := "server"
		var config UpCmdConfig
		err := cfgFileViper.UnmarshalKey(configKey, &config)
		if err != nil {
			return nil, fmt.Errorf(
				"[Config] Failed to read key '%s' from config.yaml.",
				configKey,
			)
		}
		return &config, nil
	}
	addUpCmd(rootCmd, graphClientFactory, getUpCmdConfig)

	addPullCmd(rootCmd, cfgFileViper)

	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().
		StringVar(&cfgFileName, "config", "", "config filename - default is $XDG_CONFIG_HOME/twtodo/config.yaml")
	rootCmd.PersistentFlags().
		StringVar(&credentialsFileName, "credentials", "", "credentials filename - default is $XDG_CONFIG_HOME/twtodo/credentials.env")
}

// initConfig is run when each command's Execute function is called.
func initConfig() {
	const twtodoCfg = "twtodo"

	if cfgFileName != "" {
		// Use config file from the flag.
		cfgFileViper.SetConfigFile(cfgFileName)
	} else {
		cfgFileViper.AddConfigPath(filepath.Join(xdg.ConfigHome, twtodoCfg))
		cfgFileViper.SetConfigType("yaml")
		cfgFileViper.SetConfigName("config")
	}
	cfgFileViper.AutomaticEnv()
	if err := cfgFileViper.ReadInConfig(); err == nil {
		fmt.Println("[Config] Using config file:", cfgFileViper.ConfigFileUsed())
	}

	if credentialsFileName != "" {
		credentialsFileViper.SetConfigFile(credentialsFileName)
	} else {
		credentialsFileViper.AddConfigPath(filepath.Join(xdg.ConfigHome, twtodoCfg))
		credentialsFileViper.SetConfigType("yaml")
		credentialsFileViper.SetConfigName("credentials")
	}
	credentialsFileViper.AutomaticEnv()
	if err := credentialsFileViper.ReadInConfig(); err == nil {
		fmt.Println(
			"[Config] Using credentials file:",
			credentialsFileViper.ConfigFileUsed(),
		)
	}
}
