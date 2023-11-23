package main

import (
	"bytes"  // To create buffers to capture the output
	"errors" // To validate errors
	"fmt"
	"io" // To use the op.Reader interface
	"math"
	"testing"        // To execute tests
	"testing/iotest" // To assist in executing tests that fail to read data
)

func TestOperations(t *testing.T) {
	data := [][]float64{
		{10, 20, 15, 30, 45, 50, 100, 30},
		{5.5, 8, 2.2, 9.75, 8.45, 3, 2.5, 10.25, 4.75, 6.1, 7.67, 12.287, 5.47},
		{-10, -20},
		{102, 37, 44, 57, 67, 129},
	}

	// Test cases for Operations Test
	testCases := []struct {
		name string
		op   statsFunc
		exp  []float64
	}{
		{"Sum", sum, []float64{300, 85.927, -30, 436}},
		{"Avg", avg, []float64{37.5, 6.609769230769231, -15, 72.66666666666667}},
		{"Min", min, []float64{10, 2.2, -20, 37}},
		{"Max", max, []float64{100, 12.287, -10, 129}},
	}

	// Operations Tests execution
	for _, tc := range testCases {
		for k, exp := range tc.exp {
			name := fmt.Sprintf("%sData%d", tc.name, k)
			t.Run(
				name,
				func(t *testing.T) {
					act := tc.op(data[k])
					if math.Abs(act-exp) > math.SmallestNonzeroFloat64 {
						t.Errorf("\nExpected: %g\nActual: %g\n", exp, act)
					}
				})
		}
	}
}

func TestCV2Float(t *testing.T) {
	csvData := `IP Address,Requests,Response Time
192.168.0.199,2056,236
192.168.0.88,899,220
192.168.0.199,3054,226
192.168.0.100,4133,218
192.168.0.199,950,238
`
	// Test cases for Operations Test
	testCases := []struct {
		name   string
		col    int
		exp    []float64
		expErr error
		reader io.Reader
	}{
		{
			name:   "Column2",
			col:    2,
			exp:    []float64{2056, 899, 3054, 4133, 950},
			expErr: nil,
			reader: bytes.NewBufferString(csvData)},
		{
			name:   "Column3",
			col:    3,
			exp:    []float64{236, 220, 226, 218, 238},
			expErr: nil,
			reader: bytes.NewBufferString(csvData)},
		{
			name:   "FailRead",
			col:    1,
			exp:    nil,
			expErr: iotest.ErrTimeout,
			reader: iotest.TimeoutReader(bytes.NewReader([]byte{0}))},
		{
			name:   "FailedNotNumber",
			col:    1,
			exp:    nil,
			expErr: ErrNotNumber,
			reader: bytes.NewBufferString(csvData)},
		{
			name:   "FailedInvalidColumn",
			col:    4,
			exp:    nil,
			expErr: ErrInvalidColumn,
			reader: bytes.NewBufferString(csvData)},
	}

	// CSV2Float Tests execution
	for _, tc := range testCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				act, err := csv2float(tc.reader, tc.col)
				// Check for errors if expErr is not null
				if tc.expErr != nil {
					if err == nil {
						t.Errorf("Expected error. Got nil instead")
					}

					if !errors.Is(err, tc.expErr) {
						t.Errorf("\nExpected error: %q\nActual error: %q\n", tc.expErr, err)
					}

					return
				}
				// Check results if errors are not expected
				if err != nil {
					t.Errorf("Unexpected error: %q", err)
				}
				for i, exp := range tc.exp {
					if act[i] != exp {
						t.Errorf("\nExpected: %g\nActual: %g\n", exp, act[i])
					}
				}
			})
	}
}
