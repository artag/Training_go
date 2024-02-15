package cmd

import (
	"fmt" // To use format fmt.Fprintln
	"io"  // To use io.Writer interface
	"os"  // To use os.Stdout for outputs
	"rggo/cobra/pScan/scan"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:          "delete <host1>...<hostn>",
	Aliases:      []string{"d"},
	Short:        "Delete host(s) from list",
	SilenceUsage: true,
	Args:         cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile, err := cmd.Flags().GetString("hosts-file")
		if err != nil {
			return err
		}

		return deleteAction(os.Stdout, hostsFile, args)
	},
}

func init() {
	hostsCmd.AddCommand(deleteCmd)
}

func deleteAction(out io.Writer, hostsFile string, args []string) error {
	hl := &scan.HostsList{}

	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	for _, h := range args {
		if err := hl.Remove(h); err != nil {
			return err
		}

		fmt.Fprintln(out, "Deleted host:", h)
	}

	return hl.Save(hostsFile)
}
