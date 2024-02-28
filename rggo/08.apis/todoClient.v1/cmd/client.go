package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// To print a shorter version of the time data. Include month, day, hour, minute.
const timeFormat = "Jan/02 @15:04"

var (
	//lint:ignore ST1005 Ignore warning
	ErrConnection = errors.New("Connection error")
	//lint:ignore ST1005 Ignore warning
	ErrNotFound = errors.New("Not found")
	//lint:ignore ST1005 Ignore warning
	ErrInvalidResponse = errors.New("Invalid server response")
	//lint:ignore ST1005 Ignore warning
	ErrInvalid = errors.New("Invalid data")
	//lint:ignore ST1005 Ignore warning
	ErrNotNumber = errors.New("Not a number")
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type response struct {
	Results      []item `json:"results"`
	Date         int64  `json:"date"`
	TotalResults int    `json:"total_results"`
}

func newClient() *http.Client {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	return c
}

func getItems(url string) ([]item, error) {
	r, err := newClient().Get(url)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrConnection, err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		msg, err := io.ReadAll(r.Body)
		if err != nil {
			//lint:ignore ST1005 Ignore warning
			return nil, fmt.Errorf("Cannot read body: %w", err)
		}
		err = ErrInvalidResponse
		if r.StatusCode == http.StatusNotFound {
			err = ErrNotFound
		}
		return nil, fmt.Errorf("%w: %s", err, msg)
	}

	var resp response
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return nil, err
	}

	if resp.TotalResults == 0 {
		return nil, fmt.Errorf("%w: No results found", ErrNotFound)
	}

	return resp.Results, nil
}

func getAll(apiRoot string) ([]item, error) {
	u := fmt.Sprintf("%s/todo", apiRoot)
	return getItems(u)
}

func getOne(apiRoot string, id int) (item, error) {
	u := fmt.Sprintf("%s/todo/%d", apiRoot, id)

	items, err := getItems(u)
	if err != nil {
		return item{}, err
	}

	if len(items) != 1 {
		return item{}, fmt.Errorf("%w: Invalid results", ErrInvalid)
	}

	return items[0], nil
}
