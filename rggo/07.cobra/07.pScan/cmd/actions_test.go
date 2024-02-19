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

func TestPortParse(t *testing.T) {
	testCases := []struct {
		name     string
		arg      string
		expPorts []int
		expErr   string
	}{
		{
			name:     "Empty_String",
			arg:      "",
			expPorts: nil,
			expErr:   "No ports provided.\n",
		},
		{
			name:     "One_integer",
			arg:      "80",
			expPorts: []int{80},
			expErr:   "",
		},
		{
			name:     "One_invalid_integer",
			arg:      "80a",
			expPorts: nil,
			expErr:   "Invalid port value \"80a\".\n",
		},
		{
			name:     "One_integer_range",
			arg:      "33-35",
			expPorts: []int{33, 34, 35},
			expErr:   "",
		},
		{
			name:     "One_range_two_integers",
			arg:      "1, 42, 33-35, 34",
			expPorts: []int{1, 33, 34, 35, 42},
			expErr:   "",
		},
		{
			name:     "Two_ranges_three_integers",
			arg:      "34-35, 443, 80, 33-35, 42",
			expPorts: []int{33, 34, 35, 42, 80, 443},
			expErr:   "",
		},
		{
			name:     "Invalid_integer_begin_range",
			arg:      "33a-35b",
			expPorts: nil,
			expErr:   "Invalid port value \"33a\".\n",
		},
		{
			name:     "Invalid_integer_end_range",
			arg:      "33-35b",
			expPorts: nil,
			expErr:   "Invalid port value \"35b\".\n",
		},
		{
			name:     "Invalid_ports_range_1",
			arg:      "31,33-33",
			expPorts: nil,
			expErr:   "Invalid ports range 33-33.\n",
		},
		{
			name:     "Invalid_ports_range_2",
			arg:      " 2, 35-33, 1",
			expPorts: nil,
			expErr:   "Invalid ports range 35-33.\n",
		},
		{
			name:     "Invalid_ports_range_3",
			arg:      " 2, -1-33, 1",
			expPorts: nil,
			expErr:   "Invalid port value \"-1-33\".\n",
		},
		{
			name:     "Invalid_ports_range_4",
			arg:      " 2, 1-65536, 1",
			expPorts: nil,
			expErr:   "Invalid port value 65536. Port value must be between 1 and 65535.\n",
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				ports, err := parsePorts(tc.arg)
				if tc.expErr != "" && err.Error() != tc.expErr {
					t.Fatalf("Expected %q, actual %q error.\n", tc.expErr, err.Error())
				}
				if tc.expErr == "" && err != nil {
					t.Fatalf("Got error %q.\n", err)
				}
				if tc.expPorts == nil && ports != nil {
					t.Fatal("Expected nil ports.")
				}
				if len(tc.expPorts) != len(ports) {
					t.Fatalf("Expected %d ports, actual %d ports.\n", len(tc.expPorts), len(ports))
				}
				for i, act := range ports {
					exp := tc.expPorts[i]
					if act != exp {
						t.Fatalf("Expected %d port[%d], actual %d port[%d]\n", exp, i, act, i)
					}
				}
			})
	}
}
