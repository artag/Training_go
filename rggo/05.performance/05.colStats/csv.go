package main

import (
	"encoding/csv" // To read data as string from CSV files
	"fmt"          // To print formatted results out
	"io"           // To provide the io.Reader interface
	"math"
	"strconv" // To convert string data into numeric data
)

// statsFunc defines a generic statistical function
type statsFunc func(data []float64) float64

func sum(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum
}

func avg(data []float64) float64 {
	len := len(data)
	if len < 1 {
		return 0.0
	}
	return sum(data) / float64(len)
}

func min(data []float64) float64 {
	min := math.MaxFloat64
	for _, v := range data {
		if min > v {
			min = v
		}
	}
	return min
}

func max(data []float64) float64 {
	max := -math.MaxFloat64
	for _, v := range data {
		if max < v {
			max = v
		}
	}
	return max
}

func csv2float(csvSource io.Reader, column int) ([]float64, error) {
	// Create the CSV Reader used to read in data from CSV files
	cr := csv.NewReader(csvSource)
	cr.ReuseRecord = true // Reuse slice for performance

	// Adjusting for 0 based index
	column--

	// Looping through all records
	var data []float64
	for i := 0; ; i++ {
		// Read one row from csv file
		row, err := cr.Read()
		if err == io.EOF {
			break // End of csv file
		}
		if err != nil {
			return nil, fmt.Errorf("Cannot read data from file: %w", err)
		}

		if i == 0 {
			continue // discard the first (title) line
		}

		// Checking number of columns in CSV file
		if len(row) <= column {
			return nil, fmt.Errorf("%w: File has only %d columns", ErrInvalidColumn, len(row))
		}
		// Try to convert data read into a float number
		v, err := strconv.ParseFloat(row[column], 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrNotNumber, err)
		}

		data = append(data, v)
	}

	// Return the slice of float64 and nil error
	return data, nil
}
