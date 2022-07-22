package cli

import (
	"fmt"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/joho/godotenv"
	"github.com/simachri/taskwarrior-ms-todo/internal/mstodo"
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

func Execute() error {
	addPullCmd(rootCmd, &mstodo.Client{})

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
		viper.SetConfigFile(cfgFileName)
	} else {
		viper.AddConfigPath(filepath.Join(xdg.ConfigHome, twtodoCfg))
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if credentialsFileName != "" {
		godotenv.Load(credentialsFileName)
	} else {
		godotenv.Load(filepath.Join(xdg.ConfigHome, twtodoCfg, "credentials.env"))
	}
}
