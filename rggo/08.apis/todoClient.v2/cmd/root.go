package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "todoClient",
	Short: "A Todo API client",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&cfgFile, "config", "", "config file (default is $HOME/.todoClient.yaml)")

	rootCmd.PersistentFlags().String(
		"api-root", "http://localhost:8080", "Todo API URL")

	// Environment variable TODO_API_ROOT
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("TODO")

	// Bind an environment variable TODO_API_ROOT
	viper.BindPFlag("api-root", rootCmd.PersistentFlags().Lookup("api-root"))
}
