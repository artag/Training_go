package cmd

import (
	"fmt"
	"io"             // To use the io.Writer interface
	"os"             // To use the os.Stdout for output
	"text/tabwriter" // To print formatted tabulated data

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:          "list",
	Short:        "List todo items",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiRoot := viper.GetString("api-root")
		listActiveOnly, err := cmd.Flags().GetBool("active")
		if err != nil {
			return err
		}

		if listActiveOnly {
			return listOnlyActiveAction(os.Stdout, apiRoot)
		} else {
			return listAction(os.Stdout, apiRoot)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("active", "a", false, "List only active tasks")
}

func listOnlyActiveAction(out io.Writer, apiRoot string) error {
	items, err := getAll(apiRoot)
	if err != nil {
		return err
	}

	return printAll(out, items, true)
}

func listAction(out io.Writer, apiRoot string) error {
	items, err := getAll(apiRoot)
	if err != nil {
		return err
	}

	return printAll(out, items, false)
}

func printAll(out io.Writer, items []item, listActiveOnly bool) error {
	// minimum column width - 3 characters
	// tabwidth - 2 characters
	// padding - 0 characters
	// pad character - whitespace (' ')
	// disable additional flags - 0
	w := tabwriter.NewWriter(out, 3, 2, 0, ' ', 0)
	for index, item := range items {
		done := "-"
		if item.Done {
			done = "X"
		}
		if listActiveOnly && done == "X" {
			continue
		}
		fmt.Fprintf(w, "%s\t%d\t%s\t\n", done, index+1, item.Task)
	}

	// Flush the output to the io.Writer
	return w.Flush()
}
