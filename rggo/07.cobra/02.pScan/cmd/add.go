package cmd

import (
	"fmt"
	"io"
	"os"
	"rggo/cobra/pScan/scan"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:          "add <host 1>...<host n>",
	Aliases:      []string{"a"},
	Short:        "Add new host(s) to list",
	SilenceUsage: true,
	Args:         cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile, err := cmd.Flags().GetString("hosts-file")
		if err != nil {
			return err
		}

		return addAction(os.Stdout, hostsFile, args)
	},
}

func init() {
	hostsCmd.AddCommand(addCmd)
}

func addAction(out io.Writer, hostsFile string, args []string) error {
	hl := &scan.HostsList{}

	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	for _, h := range args {
		if err := hl.Add(h); err != nil {
			return err
		}

		fmt.Fprintln(out, "Added host:", h)
	}

	return hl.Save(hostsFile)
}
