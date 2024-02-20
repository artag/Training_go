package scan

import (
	"fmt"
	"net"
	"time"
)

// The state of a single TCP port
type PortState struct {
	Port int
	Open state
}

// Port state
type state bool

// Results represents the scan results for a single port
type Results struct {
	Host       string
	NotFound   bool
	PortStates []PortState
}

type ScanParameters struct {
	Ports    []int
	Settings ScanSettings
}

type ScanSettings struct {
	UDPScanMode bool
	TimeoutMs   int
}

// Converts the boolean value of state to a human readable string
func (s state) String() string {
	if s {
		return "open"
	}
	return "closed"
}

// Performs a port scan on a single TCP or UDP port
func scanPort(host string, port int, settings ScanSettings) PortState {
	p := PortState{
		Port: port,
	}

	portStr := fmt.Sprintf("%d", port)
	address := net.JoinHostPort(host, portStr)

	scanTimeout := time.Duration(settings.TimeoutMs) * time.Millisecond

	var scanConn net.Conn
	var err error

	if settings.UDPScanMode {
		scanConn, err = net.DialTimeout("udp", address, scanTimeout)
	} else {
		scanConn, err = net.DialTimeout("tcp", address, scanTimeout)
	}

	if err != nil {
		return p // If error - port closed
	}

	scanConn.Close()
	p.Open = true // Port open
	return p
}

func Run(hl *HostsList, params ScanParameters) []Results {
	res := make([]Results, 0, len(hl.Hosts))
	for _, h := range hl.Hosts {
		r := Results{
			Host: h,
		}

		// Resolve the host name into a valid IP address
		if _, err := net.LookupHost(h); err != nil {
			r.NotFound = true // IP address not resolved
			res = append(res, r)
			continue
		}

		// Scan ports
		for _, port := range params.Ports {
			portState := scanPort(h, port, params.Settings)
			r.PortStates = append(r.PortStates, portState)
		}

		res = append(res, r)
	}

	return res
}
