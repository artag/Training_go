package main

import (
	"bytes"
	"errors"
	"fmt"           // To print formatted output
	"io/ioutil"     // To create the temporary directory
	"os"            // To interact with the operating system
	"os/exec"       // To execute external programs
	"path/filepath" // To handle path operations
	"testing"
)

func TestRun(t *testing.T) {
	_, err := exec.LookPath("git")
	if err != nil {
		// Skip the test if git isn't available.
		t.Skip("Git not installed. Skipping test.")
	}

	testCases := []struct {
		name     string // Test name
		proj     string // Target project
		out      string // Actual output
		expErr   error  // Expected error
		setupGit bool   // true - use mock Git
	}{
		{
			name: "success",
			proj: "./testdata/tool",
			out: "Go Build: SUCCESS\n" +
				"Go Test: SUCCESS\n" +
				"Gofmt: SUCCESS\n" +
				"Git Push: SUCCESS\n",
			expErr:   nil,
			setupGit: true,
		},
		{
			name:     "fail",
			proj:     "./testdata/toolErr",
			out:      "",
			expErr:   &stepErr{step: "go build"},
			setupGit: false,
		},
		{
			name:     "failFormat",
			proj:     "./testdata/toolFmtErr",
			out:      "",
			expErr:   &stepErr{step: "go fmt"},
			setupGit: false,
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				if tc.setupGit {
					cleanup := setupGit(t, tc.proj)
					defer cleanup()
				}

				var out bytes.Buffer
				err := run(tc.proj, &out)
				if tc.expErr != nil {
					if err == nil {
						t.Errorf("Exected error: %q. Got 'nil'.", tc.expErr)
						return
					}

					if !errors.Is(err, tc.expErr) {
						t.Errorf("Expected error: %q. Got: %q.", tc.expErr, err)
					}

					return
				}

				if err != nil {
					t.Errorf("Unexpected error: %q", err)
					return
				}

				if out.String() != tc.out {
					t.Errorf("Expected output: %q. Got %q", tc.out, out.String())
				}
			})
	}
}

func setupGit(t *testing.T, proj string) func() {
	t.Helper()

	// Check git availability.
	gitExec, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("git: %q\n", gitExec)

	// Create a temporary directory with prefix gocitest (for mock git repository).
	tempDir, err := ioutil.TempDir("", "gocitest")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Temporary dir: %q\n", tempDir)

	// Get full path of the target project directory.
	projPath, err := filepath.Abs(proj)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Target project dir: %q\n", projPath)

	// Get URI to local git repository
	remoteUri := fmt.Sprintf("file://%s", tempDir)
	fmt.Printf("URI to local git repository: %q\n", remoteUri)

	// Git commands to set up the test environment.
	var gitCmdList = []struct {
		// Arguments for the git command.
		args []string
		// The directory on which to execute the command.
		dir string
		// Environent variables.
		env []string
	}{
		{
			[]string{"init", "--bare"},
			tempDir,
			nil,
		},
		{
			[]string{"init"},
			projPath,
			nil,
		},
		{
			[]string{"remote", "add", "origin", remoteUri},
			projPath,
			nil,
		},
		{
			[]string{"add", "."},
			projPath,
			nil,
		},
		{
			[]string{"commit", "-m", "test"},
			projPath,
			[]string{
				"GIT_COMMITTER_NAME=test",
				"GIT_COMMITTER_EMAIL=test@example.com",
				"GIT_AUTHOR_NAME=test",
				"GIT_AUTHOR_EMAIL=test@example.com",
			},
		},
	}

	for _, g := range gitCmdList {
		gitCmd := exec.Command(gitExec, g.args...)
		gitCmd.Dir = g.dir

		if g.env != nil {
			gitCmd.Env = append(os.Environ(), g.env...)
		}

		if err := gitCmd.Run(); err != nil {
			t.Fatal(err)
		}
	}

	return func() {
		fmt.Printf("Removing temporary git repository %q\n", tempDir)
		os.RemoveAll(tempDir)

		dotGit := filepath.Join(projPath, ".git")
		fmt.Printf("Removing %q\n", dotGit)
		os.RemoveAll(dotGit)
	}
}
