package cmd

import (
	"fmt"
	"io" // Use io.Writer interface
	"os" // Use os.Stdout
	"slices"
	"strconv"
	"strings"

	"rggo/cobra/pScan/scan"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Run a port scan on the hosts",
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile := viper.GetString("hosts-file")
		portsStr, err := cmd.Flags().GetString("ports")
		if err != nil {
			return err
		}

		ports, err := parsePorts(portsStr)
		if err != nil {
			return err
		}

		return scanAction(os.Stdout, hostsFile, ports)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	// 22, 80, 443 - port values by default
	scanCmd.Flags().StringP("ports", "p", "22, 80, 443", "ports to scan")
}

func scanAction(out io.Writer, hostsFile string, ports []int) error {
	hl := &scan.HostsList{}
	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	results := scan.Run(hl, ports)
	return printResults(out, results)
}

func printResults(out io.Writer, results []scan.Results) error {
	message := ""
	for _, r := range results {
		message += fmt.Sprintf("%s:", r.Host)
		if r.NotFound {
			message += " Host not found\n\n"
			continue
		}

		message += fmt.Sprintln()

		for _, p := range r.PortStates {
			message += fmt.Sprintf("\t%d: %s\n", p.Port, p.Open)
		}

		message += fmt.Sprintln()
	}

	_, err := fmt.Fprint(out, message)
	return err
}

func parsePorts(str string) ([]int, error) {
	if str == "" {
		return nil, fmt.Errorf("No ports provided.\n")
	}

	split := strings.Split(str, ",")
	parsed := make([]int, 0, len(split))
	for i := 0; i < len(split); i++ {
		numStr := strings.TrimSpace(split[i])
		num, err := strconv.Atoi(numStr)
		if err == nil {
			if err := validateInt(num); err != nil {
				return nil, err
			}
			parsed = append(parsed, num)
		} else {
			split2 := strings.Split(numStr, "-")
			if len(split2) != 2 {
				return nil, fmt.Errorf("Invalid port value %q.\n", numStr)
			}
			begin, err := strconv.Atoi(split2[0])
			if err != nil {
				return nil, fmt.Errorf("Invalid port value %q.\n", split2[0])
			}
			if err := validateInt(begin); err != nil {
				return nil, err
			}
			end, err := strconv.Atoi(split2[1])
			if err != nil {
				return nil, fmt.Errorf("Invalid port value %q.\n", split2[1])
			}
			if err := validateInt(end); err != nil {
				return nil, err
			}
			if begin >= end {
				return nil, fmt.Errorf("Invalid ports range %d-%d.\n", begin, end)
			}
			for port := begin; port <= end; port++ {
				parsed = append(parsed, port)
			}
		}
	}

	slices.Sort(parsed)
	allPorts := make([]int, 0, 0)
	prev := -1
	for _, port := range parsed {
		if port == prev {
			continue
		}

		allPorts = append(allPorts, port)
		prev = port
	}

	return allPorts, nil
}

func validateInt(num int) error {
	if 0 < num && num <= 65535 {
		return nil
	}

	return fmt.Errorf("Invalid port value %d. Port value must be between 1 and 65535.\n", num)
}
