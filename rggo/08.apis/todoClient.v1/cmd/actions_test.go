package cmd

import (
	"bytes"
	"errors"
	"fmt"
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
				fmt.Printf("The url of mock server: %q\n", url)
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
			fmt.Printf("The url of mock server: %q\n", url)

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
