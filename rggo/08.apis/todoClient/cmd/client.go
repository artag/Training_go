package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	//lint:ignore ST1005 Doesn't matter
	ErrConnection = errors.New("Connection error")
	//lint:ignore ST1005 Doesn't matter
	ErrNotFound = errors.New("Not found")
	//lint:ignore ST1005 Doesn't matter
	ErrInvalidResponse = errors.New("Invalid server response")
	//lint:ignore ST1005 Doesn't matter
	ErrInvalid = errors.New("Invalid data")
	//lint:ignore ST1005 Doesn't matter
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
			//lint:ignore ST1005 Doesn't matter
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
