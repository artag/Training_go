package main

import (
	"bytes"
	"context" // To define command contexts
	"errors"
	"fmt"       // To print formatted output
	"io/ioutil" // To create the temporary directory
	"os"        // To interact with the operating system
	"os/exec"   // To execute external programs
	"os/signal"
	"path/filepath" // To handle path operations
	"syscall"
	"testing"
	"time" // To simulate a timeout
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string // Test name
		proj     string // Target project
		out      string // Actual output
		expErr   error  // Expected error
		setupGit bool   // true - use mock Git
		// Function used to mock a command if required
		mockCmd func(ctx context.Context, name string, arg ...string) *exec.Cmd
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
			mockCmd:  nil,
		},
		{
			name: "successMock",
			proj: "./testdata/tool",
			out: "Go Build: SUCCESS\n" +
				"Go Test: SUCCESS\n" +
				"Gofmt: SUCCESS\n" +
				"Git Push: SUCCESS\n",
			expErr:   nil,
			setupGit: false,
			mockCmd:  mockCmdContext,
		},
		{
			name:     "fail",
			proj:     "./testdata/toolErr",
			out:      "",
			expErr:   &stepErr{step: "go build"},
			setupGit: false,
			mockCmd:  nil,
		},
		{
			name:     "failFormat",
			proj:     "./testdata/toolFmtErr",
			out:      "",
			expErr:   &stepErr{step: "go fmt"},
			setupGit: false,
			mockCmd:  nil,
		},
		{
			name:     "failTimeout",
			proj:     "./testdata/tool",
			out:      "",
			expErr:   context.DeadlineExceeded,
			setupGit: false,
			mockCmd:  mockCmdTimeout,
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				// Setup real git
				if tc.setupGit {
					_, err := exec.LookPath("git")
					if err != nil {
						// Skip the test if git isn't available.
						t.Skip("Git not installed. Skipping test.")
					}

					cleanup := setupGit(t, tc.proj)
					defer cleanup()
				}

				// Setup mock function to mock a command
				if tc.mockCmd != nil {
					command = tc.mockCmd
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

// Test the signal handling
func TestRunKill(t *testing.T) {
	testCases := []struct {
		name   string
		proj   string
		sig    syscall.Signal
		expErr error
	}{
		{
			name:   "Handle signal SIGINT",
			proj:   "./testdata/tool",
			sig:    syscall.SIGINT,
			expErr: ErrSignal,
		},
		{
			name:   "Handle signal SIGTERM",
			proj:   "./testdata/tool",
			sig:    syscall.SIGTERM,
			expErr: ErrSignal,
		},
		{
			name:   "Not handle a different signal SIGQUIT",
			proj:   "./testdata/tool",
			sig:    syscall.SIGQUIT,
			expErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				command = mockCmdTimeout

				// To get error
				errCh := make(chan error)

				// To trap expected signal
				ignSigCh := make(chan os.Signal, 1)
				signal.Notify(ignSigCh, syscall.SIGQUIT)
				defer signal.Stop(ignSigCh)

				// To trap ignored signal
				expSigCh := make(chan os.Signal, 1)
				signal.Notify(expSigCh, tc.sig)
				defer signal.Stop(expSigCh)

				// Sends the error to the error channel.
				go func() {
					errCh <- run(tc.proj, ioutil.Discard)
				}()

				// Sends the desired signal to the test executable.
				go func() {
					time.Sleep(2 * time.Second)
					pid := syscall.Getpid()   // To obtain the process ID
					syscall.Kill(pid, tc.sig) // To send the signal
				}()

				// select error
				select {
				case err := <-errCh:
					if err == nil {
						t.Errorf("Expected error. Got 'nil' instead.")
						return
					}
					if !errors.Is(err, tc.expErr) {
						t.Errorf("Expected error: %q. Got %q", tc.expErr, err)
					}

					// select signal
					select {
					case rec := <-expSigCh:
						if rec != tc.sig {
							t.Errorf("Expected signal %q, got %q", tc.sig, rec)
						}
					default:
						t.Errorf("Signal not received")
					}

				case <-ignSigCh:
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

// Simulates the command.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	// Simulate long-running process.
	if os.Getenv("GO_HELPER_TIMEOUT") == "1" {
		time.Sleep(15 * time.Second)
	}

	if os.Args[2] == "git" {
		fmt.Fprintln(os.Stdout, "Everything up-to-date")
		// The command completed successfully.
		os.Exit(0)
	}

	// Indicating error.
	os.Exit(1)
}

func mockCmdContext(ctx context.Context, exe string, args ...string) *exec.Cmd {
	// The argument list that will be passed to the command.
	cs := []string{"-test.run=TestHelperProcess"}
	// Append the command and arguments that would be passed to the real command.
	cs = append(cs, exe)
	cs = append(cs, args...)
	fmt.Printf("Arguments to mock function: %q\n", cs)

	cmd := exec.CommandContext(ctx, os.Args[0], cs...)
	// Add environment variable to run 'TestHelperProcess'.
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// To simulate a command that times out.
func mockCmdTimeout(ctx context.Context, exe string, args ...string) *exec.Cmd {
	cmd := mockCmdContext(ctx, exe, args...)
	cmd.Env = append(cmd.Env, "GO_HELPER_TIMEOUT=1")
	return cmd
}
