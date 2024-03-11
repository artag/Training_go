package cmd

import (
	"fmt"
	"io"             // To use io.Writer interface
	"os"             // To use os.Stdout for output
	"strconv"        // To convert string to integer
	"text/tabwriter" // To print tabulated data

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var viewCmd = &cobra.Command{
	Use:          "view <id>",
	Short:        "View details about a single item",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1), // required one argument
	RunE: func(cmd *cobra.Command, args []string) error {
		apiRoot := viper.GetString("api-root")
		return viewAction(os.Stdout, apiRoot, args[0])
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
}

func viewAction(out io.Writer, apiRoot, arg string) error {
	id, err := strconv.Atoi(arg)
	if err != nil {
		return fmt.Errorf("%w: Item id must be a number", ErrNotNumber)
	}

	item, err := getOne(apiRoot, id)
	if err != nil {
		return err
	}

	return printOne(out, item)
}

func printOne(out io.Writer, i item) error {
	// 14 - min width
	// 2 - tab width
	// 0 - padding
	// ' ' - pad char
	// 0 - flags
	w := tabwriter.NewWriter(out, 14, 2, 0, ' ', 0)
	fmt.Fprintf(w, "Task:\t%s\n", i.Task)
	fmt.Fprintf(w, "Created at:\t%s\n", i.CreatedAt.Format(timeFormat))
	if i.Done {
		fmt.Fprintf(w, "Completed:\t%s\n", "Yes")
		fmt.Fprintf(w, "Completed at:\t%s\n", i.CompletedAt.Format(timeFormat))
		return w.Flush()
	}
	fmt.Fprintf(w, "Completed:\t%s\n", "No")
	return w.Flush()
}
