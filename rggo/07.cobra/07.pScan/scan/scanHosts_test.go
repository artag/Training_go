package scan_test

import (
	"fmt"
	"net"
	"rggo/cobra/pScan/scan"
	"strconv"
	"testing"
)

func TestStateString(t *testing.T) {
	ps := scan.PortState{}
	if ps.Open.String() != "closed" {
		t.Errorf("Expected %q, got %q instead\n", "closed", ps.Open.String())
	}

	ps.Open = true
	if ps.Open.String() != "open" {
		t.Errorf("Expected %q, got %q instead\n", "open", ps.Open.String())
	}
}

func TestRunHostFound(t *testing.T) {
	testCases := []struct {
		name        string
		expectState string
	}{
		{"OpenPort", "open"},
		{"ClosedPort", "closed"},
	}

	host := "localhost"
	hl := &scan.HostsList{}
	hl.Add(host)

	ports := []int{}
	// Init ports, 1 open, 1 closed
	for _, tc := range testCases {
		// 0 - port automatically chosen
		ln, err := net.Listen("tcp", net.JoinHostPort(host, "0"))
		if err != nil {
			t.Fatal(err)
		}
		defer ln.Close()

		hostport := ln.Addr().String() // 127.0.0.1:12345
		fmt.Printf("hostport: %q\n", hostport)
		_, portStr, err := net.SplitHostPort(hostport)
		if err != nil {
			t.Fatal(err)
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Printf("Choosen port: %d\n", port) // 12345
		ports = append(ports, port)

		if tc.name == "ClosedPort" {
			ln.Close()
		}
	}

	scanParams := scan.ScanParameters{
		Ports: ports,
		Settings: scan.ScanSettings{
			UDPScanMode: false,
			TimeoutMs:   1000,
		},
	}

	// Act
	res := scan.Run(hl, scanParams)

	// Assert. Verify results for HostFound test
	if len(res) != 1 {
		t.Fatalf("Expected 1 results, got %d instead\n", len(res))
	}

	if res[0].Host != host {
		t.Errorf("Expected host %q, got %q instead\n", host, res[0].Host)
	}

	if res[0].NotFound {
		t.Errorf("Expected host %q to be found\n", host)
	}

	if len(res[0].PortStates) != 2 {
		t.Fatalf("Expected 2 port states, got %d instead\n", len(res[0].PortStates))
	}

	for i, tc := range testCases {
		if res[0].PortStates[i].Port != ports[i] {
			t.Errorf("Expected port %d, got %d instead\n", ports[0], res[0].PortStates[i].Port)
		}

		if res[0].PortStates[i].Open.String() != tc.expectState {
			t.Errorf("Expected port %d to be %s\n", ports[i], tc.expectState)
		}
	}
}

func TestRunHostNotFound(t *testing.T) {
	// Arrange
	host := "389.389.389.389"
	hl := &scan.HostsList{}
	hl.Add(host)

	scanParams := scan.ScanParameters{
		Ports: []int{},
		Settings: scan.ScanSettings{
			UDPScanMode: false,
			TimeoutMs:   1000,
		},
	}

	// Act
	res := scan.Run(hl, scanParams)

	// Assert. Verify results for HostNotFound test
	if len(res) != 1 {
		t.Fatalf("Expected 1 results, got %d instead\n", len(res))
	}
	if res[0].Host != host {
		t.Errorf("Expected host %q, got %q instead\n", host, res[0].Host)
	}
	if !res[0].NotFound {
		t.Errorf("Expected host %q NOT be found\n", host)
	}
	if len(res[0].PortStates) != 0 {
		t.Fatalf("Expected 0 port states, got %d instead\n", len(res[0].PortStates))
	}
}
