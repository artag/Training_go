package main

import (
	"fmt"
	"io"                // To read the response body
	"net/http"          // To deal with HTTP requests
	"net/http/httptest" // Provides HTTP testing utilities (test HTTP server)
	"strings"           // To compare strings
	"testing"           // Provides testing utilities
)

const (
	ContentType = "Content-Type"
)

func setupAPI(t *testing.T) (string, func()) {
	t.Helper()

	ts := httptest.NewServer(newMux(""))
	fmt.Printf("test server url: %q\n", ts.URL)
	return ts.URL, func() {
		ts.Close()
	}
}

func TestGet(t *testing.T) {
	testCases := []struct {
		// Test name
		name string
		// Server URL path to test
		path string
		// Expected return code from the server
		expCode int
		// Expected number of items returned when querying the to-do API
		expItems int
		// Expected body content of the response
		expContent string
	}{
		{
			name:       "GetRoot",
			path:       "/",
			expCode:    http.StatusOK,
			expContent: "There's an API here",
		},
		{
			name:    "NotFound",
			path:    "/todo/500",
			expCode: http.StatusNotFound,
		},
	}

	url, cleanup := setupAPI(t)
	defer cleanup()

	for _, tc := range testCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				var (
					body []byte
					err  error
				)

				r, err := http.Get(url + tc.path)
				if err != nil {
					t.Error(err)
				}
				defer r.Body.Close()

				if r.StatusCode != tc.expCode {
					t.Fatalf(
						"Expected %q, got %q",
						http.StatusText(tc.expCode),
						http.StatusText(r.StatusCode))
				}

				switch {
				case strings.Contains(r.Header.Get(ContentType), "text/plain"):
					if body, err = io.ReadAll(r.Body); err != nil {
						t.Error(err)
					}
					content := string(body)
					if !strings.Contains(content, tc.expContent) {
						t.Errorf("Expected %q, got %q", tc.expContent, content)
					}
				default:
					t.Fatalf("Unsupported Content-Type: %q", r.Header.Get(ContentType))
				}
			})
	}
}
