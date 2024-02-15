package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: "0.1",
	Use:     "pScan",
	Short:   "Fast TCP port scanner",
	Long: `pScan - short for Port Scanner - executes TCP port scan on a list of hosts.

pScan allows you to add, list, and delete hosts from the list.

pScan executes a port scan on specified TCP ports. You can customize the
target ports using a command line flag.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.
		PersistentFlags().
		StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pScan.yaml)")

	rootCmd.
		PersistentFlags().
		StringP("hosts-file", "f", "pScan.hosts", "pScan hosts file")

	// Template for version
	versionTemplate := `{{printf "%s: %s - version %s\n" .Name .Short .Version}}`
	rootCmd.SetVersionTemplate(versionTemplate)
}

func initConfig() {
	// Uses the package "viper" to include configuration management for application.
}
