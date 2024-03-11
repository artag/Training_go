//go:build !integration
// +build !integration

package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestListAction(t *testing.T) {
	testCases := []struct {
		// Test case name
		name string
		// Expected error
		expError error
		// Expected output
		expOut string
		// A mock server response
		resp struct {
			Status int
			Body   string
		}
		// To close server immediately to test error conditions
		closeServer bool
	}{
		{
			name:     "Results",
			expError: nil,
			expOut:   "-  1  Task 1\n-  2  Task 2\n",
			resp:     testResp["resultsMany"],
		},
		{
			name:     "NoResults",
			expError: ErrNotFound,
			resp:     testResp["noResults"],
		},
		{
			name:        "InvalidURL",
			expError:    ErrConnection,
			resp:        testResp["noResults"],
			closeServer: true,
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				url, cleanup := mockServer(
					func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(tc.resp.Status)
						fmt.Fprintln(w, tc.resp.Body)
					})
				fmt.Printf("The url of mock server for list test: %q\n", url)
				defer cleanup()
				if tc.closeServer {
					cleanup()
				}

				var out bytes.Buffer
				err := listAction(&out, url)

				if tc.expError != nil {
					if err == nil {
						t.Fatalf("Expected error %q, got no error.", tc.expError)
					}
					if !errors.Is(err, tc.expError) {
						t.Errorf("Expected error %q, got %q.", tc.expError, err)
					}
					return
				}

				if err != nil {
					t.Fatalf("Expected no error, got %q.", err)
				}
				if tc.expOut != out.String() {
					t.Errorf("Expected output %q, got %q.", tc.expOut, out.String())
				}
			})
	}
}

func TestViewAction(t *testing.T) {
	testCases := []struct {
		// Test case name
		name string
		// Expected error
		expError error
		// Expected output
		expOut string
		// A mock server response
		resp struct {
			Status int
			Body   string
		}
		// Identifier of the task
		id string
	}{
		{
			name:     "ResultsOne",
			expError: nil,
			expOut: `Task:         Task 1
Created at:   Oct/28 @08:23
Completed:    No
`,
			resp: testResp["resultsOne"],
			id:   "1",
		},
		{
			name:     "NotFound",
			expError: ErrNotFound,
			resp:     testResp["notFound"],
			id:       "1",
		},
		{
			name:     "InvalidID",
			expError: ErrNotNumber,
			resp:     testResp["noResults"],
			id:       "a",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, cleanup := mockServer(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tc.resp.Status)
					fmt.Fprintf(w, tc.resp.Body)
				})
			defer cleanup()
			fmt.Printf("The url of mock server for view test: %q\n", url)

			var out bytes.Buffer
			err := viewAction(&out, url, tc.id)

			if tc.expError != nil {
				if err == nil {
					t.Fatalf("Expected error %q, got no error.", tc.expError)
				}
				if !errors.Is(err, tc.expError) {
					t.Errorf("Expected error %q, got %q.", tc.expError, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no error, got %q.", err)
			}
			if tc.expOut != out.String() {
				t.Errorf("Expected output %q, got %q.", tc.expOut, out.String())
			}
		})
	}
}

func TestAddAction(t *testing.T) {
	expUrlPath := "/todo"
	expMethod := http.MethodPost
	expBody := "{\"task\":\"Task 1\"}\n"
	expContentType := ContentApplicationJson
	expOut := "Added task \"Task 1\" to the list.\n"
	args := []string{"Task", "1"}

	// Instatiate a test server for Add test
	url, cleanup := mockServer(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != expUrlPath {
				t.Errorf("Expected path %q, got %q", expUrlPath, r.URL.Path)
			}
			if r.Method != expMethod {
				t.Errorf("Expected method %q, got %q", expMethod, r.Method)
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}
			r.Body.Close()

			if string(body) != expBody {
				t.Errorf("Expected body %q, got %q", expBody, string(body))
			}

			contentType := r.Header.Get(ContentType)
			if contentType != expContentType {
				t.Errorf("Expected Content-Type %q, got %q", expContentType, contentType)
			}

			w.WriteHeader(testResp["created"].Status)
			fmt.Fprintln(w, testResp["created"].Body)
		})
	defer cleanup()
	fmt.Printf("The url of mock server for add test: %q\n", url)

	// Execute Add
	var out bytes.Buffer
	if err := addAction(&out, url, args); err != nil {
		t.Fatalf("Expected no error, got %q.", err)
	}

	if expOut != out.String() {
		t.Errorf("Expected output %q, got %q", expOut, out.String())
	}
}

func TestCompleteAction(t *testing.T) {
	expURLPath := "/todo/1"
	expMethod := http.MethodPatch
	expQuery := "complete"
	expOut := "Item number 1 marked as completed.\n"
	arg := "1"

	// Instatiate a test server for Complete test
	url, cleanup := mockServer(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != expURLPath {
				t.Errorf("Expected path %q, got %q", expURLPath, r.URL.Path)
			}
			if r.Method != expMethod {
				t.Errorf("Expected query %q, got %q", expMethod, r.Method)
			}
			if _, ok := r.URL.Query()[expQuery]; !ok {
				t.Errorf("Expected query %q not found in URL", expQuery)
			}

			w.WriteHeader(testResp["noContent"].Status)
			fmt.Fprintln(w, testResp["noContent"].Body)
		})
	defer cleanup()
	fmt.Printf("The url of mock server for complete test: %q\n", url)

	// Execute complete test
	var out bytes.Buffer
	if err := completeAction(&out, url, arg); err != nil {
		t.Fatalf("Expected no error, got %q.", err)
	}
	if expOut != out.String() {
		t.Errorf("Expected output %q, got %q", expOut, out.String())
	}
}

func TestDelAction(t *testing.T) {
	expURLPath := "/todo/1"
	expMethod := http.MethodDelete
	expOut := "Item number 1 deleted.\n"
	arg := "1"

	// Instatiate a test server for Del test
	url, cleanup := mockServer(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != expURLPath {
				t.Errorf("Expected path %q, got %q", expURLPath, r.URL.Path)
			}
			if r.Method != expMethod {
				t.Errorf("Expected method %q, got %q", expMethod, r.Method)
			}

			w.WriteHeader(testResp["noContent"].Status)
			fmt.Fprintf(w, testResp["noContent"].Body)
		})
	defer cleanup()
	fmt.Printf("The url of mock server for del test: %q\n", url)

	// Execute Del test
	var out bytes.Buffer
	if err := delAction(&out, url, arg); err != nil {
		t.Fatalf("Expected no error, got %q.", err)
	}
	if expOut != out.String() {
		t.Errorf("Expected output %q, got %q.", expOut, out.String())
	}
}
