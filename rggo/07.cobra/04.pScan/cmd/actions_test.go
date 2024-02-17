package cmd

import (
	"bytes" // bytes.Buffer - to testing STDOUT
	"fmt"
	"io" // io.Writer
	"net"
	"os" // os.CreateTemp - create temp file
	"strconv"
	"strings"
	"testing"

	"rggo/cobra/pScan/scan"
)

func setup(t *testing.T, hosts []string, initList bool) (string, func()) {
	// Create temp file
	tf, err := os.CreateTemp("", "pScan")
	if err != nil {
		t.Fatal(err)
	}
	tf.Close()

	// Initialize list if needed
	if initList {
		hl := &scan.HostsList{}
		for _, h := range hosts {
			hl.Add(h)
		}
		if err := hl.Save(tf.Name()); err != nil {
			t.Fatal(err)
		}
	}

	// Return temp file name and cleanup function
	return tf.Name(), func() {
		os.Remove(tf.Name())
	}
}

func TestHostActions(t *testing.T) {
	hosts := []string{
		"host1",
		"host2",
		"host3",
	}

	testCases := []struct {
		name           string
		args           []string
		expectedOut    string
		initList       bool
		actionFunction func(io.Writer, string, []string) error
	}{
		{
			name:           "AddAction",
			args:           hosts,
			expectedOut:    "Added host: host1\nAdded host: host2\nAdded host: host3\n",
			initList:       false,
			actionFunction: addAction,
		},
		{
			name:           "ListAction",
			expectedOut:    "host1\nhost2\nhost3\n",
			initList:       true,
			actionFunction: listAction,
		},
		{
			name:           "DeleteAction",
			args:           []string{"host1", "host2"},
			expectedOut:    "Deleted host: host1\nDeleted host: host2\n",
			initList:       true,
			actionFunction: deleteAction,
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				// Setup Action List
				tf, cleanup := setup(t, hosts, tc.initList)
				defer cleanup()

				// Define var to capture Action output
				var out bytes.Buffer

				// Execute Action and capture output
				if err := tc.actionFunction(&out, tf, tc.args); err != nil {
					t.Fatalf("Expected no error, got %q\n", err)
				}

				// Test Actions output
				if out.String() != tc.expectedOut {
					t.Errorf("Expected output %q, got %q\n", tc.expectedOut, out.String())
				}
			})
	}
}

func TestIntegration(t *testing.T) {
	// Define hosts for integration test
	hosts := []string{
		"host1",
		"host2",
		"host3",
	}

	// Setup integration test
	tf, cleanup := setup(t, hosts, false)
	defer cleanup()

	// Host to delete, and expected list of hosts at the end of the test
	hostToDelete := "host2"
	hostsAtTheEnd := []string{
		"host1",
		"host3",
	}

	// Define var to capture output
	var out bytes.Buffer

	// Define expected output for all actions
	expectedOut := ""
	for _, v := range hosts {
		expectedOut += fmt.Sprintf("Added host: %s\n", v)
	}
	expectedOut += strings.Join(hosts, "\n")
	expectedOut += fmt.Sprintln()
	expectedOut += fmt.Sprintf("Deleted host: %s\n", hostToDelete)
	expectedOut += strings.Join(hostsAtTheEnd, "\n")
	expectedOut += fmt.Sprintln()
	for _, v := range hostsAtTheEnd {
		expectedOut += fmt.Sprintf("%s: Host not found\n", v)
		expectedOut += fmt.Sprintln()
	}

	// Execute all operations in the sequence add->list->delete->list->scan

	// Add hosts to the list
	if err := addAction(&out, tf, hosts); err != nil {
		t.Fatalf("Expected no error, got %q\n", err)
	}

	// List hosts
	if err := listAction(&out, tf, nil); err != nil {
		t.Fatalf("Expected no error, got %q\n", err)
	}

	// Delete host2
	if err := deleteAction(&out, tf, []string{hostToDelete}); err != nil {
		t.Fatalf("Expected no error, got %q\n", err)
	}

	// List hosts after delete
	if err := listAction(&out, tf, nil); err != nil {
		t.Fatalf("Expected no error, got %q\n", err)
	}

	// Scan hosts
	if err := scanAction(&out, tf, nil); err != nil {
		t.Fatalf("Expected no error, got %q\n", err)
	}

	// Assert. Test integration output
	if out.String() != expectedOut {
		t.Errorf("Expected output %q, got %q\n", expectedOut, out.String())
	}
}

func TestScanAction(t *testing.T) {
	hosts := []string{
		"localhost",
		"unknownhost",
	}

	tf, cleanup := setup(t, hosts, true)
	defer cleanup()

	// Init ports, 1 open, 1 closed
	ports := []int{}
	for i := 0; i < 2; i++ {
		ln, err := net.Listen("tcp", net.JoinHostPort("localhost", "0"))
		if err != nil {
			t.Fatal(err)
		}
		defer ln.Close()

		_, portStr, err := net.SplitHostPort(ln.Addr().String())
		if err != nil {
			t.Fatal(err)
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatal(err)
		}

		ports = append(ports, port)
		if i == 1 {
			ln.Close()
		}
	}

	// Define expected output for scan action
	expectedOut := fmt.Sprintln("localhost:")
	expectedOut += fmt.Sprintf("\t%d: open\n", ports[0])
	expectedOut += fmt.Sprintf("\t%d: closed\n", ports[1])
	expectedOut += fmt.Sprintln()
	expectedOut += fmt.Sprintln("unknownhost: Host not found")
	expectedOut += fmt.Sprintln()

	// Define var to capture scan output
	var out bytes.Buffer

	// Execute scan and capture output
	if err := scanAction(&out, tf, ports); err != nil {
		t.Fatalf("Expected no error, got %q\n", err)
	}
	// Test scan output
	if out.String() != expectedOut {
		t.Errorf("Expected output %q, got %q\n", expectedOut, out.String())
	}
}
