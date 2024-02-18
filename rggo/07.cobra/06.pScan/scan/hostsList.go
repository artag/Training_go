package scan

import (
	"bufio"  // To read data from files
	"errors" // To define error values
	"fmt"    // To print formatted output
	"os"     // For operating system-related functions
	"sort"   // To sort the hosts list of hosts
)

var (
	ErrExists    = errors.New("Host already in the list")
	ErrNotExists = errors.New("Host not in the list")
)

// A list og hosts torun port scan
type HostsList struct {
	Hosts []string
}

// Searches for hosts in the list
func (hl *HostsList) search(host string) (bool, int) {
	sort.Strings(hl.Hosts)
	i := sort.SearchStrings(hl.Hosts, host)
	if i < len(hl.Hosts) && hl.Hosts[i] == host {
		return true, i
	}

	return false, -1
}

// Adds a host to the list
func (hl *HostsList) Add(host string) error {
	if found, _ := hl.search(host); found {
		return fmt.Errorf("%w: %s", ErrExists, host)
	}

	hl.Hosts = append(hl.Hosts, host)
	return nil
}

// Deletes a host from the list
func (hl *HostsList) Remove(host string) error {
	if found, i := hl.search(host); found {
		hl.Hosts = append(hl.Hosts[:i], hl.Hosts[i+1:]...)
		return nil
	}

	return fmt.Errorf("%w: %s", ErrNotExists, host)
}

// Obtains hosts from a hosts file
func (hl *HostsList) Load(hostsFile string) error {
	f, err := os.Open(hostsFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		hl.Hosts = append(hl.Hosts, scanner.Text())
	}

	return nil
}

// Saves hosts to a hosts file
func (hl *HostsList) Save(hostsFile string) error {
	output := ""
	for _, h := range hl.Hosts {
		output += fmt.Sprintln(h)
	}

	// 644 means you can read and write the file or directory and
	// other users can only read it. Suitable for public text files.
	return os.WriteFile(hostsFile, []byte(output), 0644)
}
