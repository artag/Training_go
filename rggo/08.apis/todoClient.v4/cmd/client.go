package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// To print a shorter version of the time data. Include month, day, hour, minute.
	timeFormat             = "Jan/02 @15:04"
	ContentType            = "Content-Type"
	ContentTextPlain       = "text/plain"
	ContentApplicationJson = "application/json"
)

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

func sendRequest(
	url, method, contentType string, expStatus int, body io.Reader) error {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	if contentType != "" {
		request.Header.Set(ContentType, contentType)
	}

	response, err := newClient().Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != expStatus {
		msg, err := io.ReadAll(response.Body)
		if err != nil {
			//lint:ignore ST1005 Ignore warning
			return fmt.Errorf("Cannot read body: %w", err)
		}
		err = ErrInvalidResponse
		if response.StatusCode == http.StatusNotFound {
			err = ErrNotFound
		}
		return fmt.Errorf("%w: %s", err, msg)
	}

	return nil
}

func getAll(apiRoot string) ([]item, error) {
	u := fmt.Sprintf("%s/todo", apiRoot)
	return getItems(u)
}

func getOne(apiRoot string, id int) (item, error) {
	url := fmt.Sprintf("%s/todo/%d", apiRoot, id)

	items, err := getItems(url)
	if err != nil {
		return item{}, err
	}

	if len(items) != 1 {
		return item{}, fmt.Errorf("%w: Invalid results", ErrInvalid)
	}

	return items[0], nil
}

func addItem(apiRoot, task string) error {
	url := fmt.Sprintf("%s/todo", apiRoot)
	item := struct {
		Task string `json:"task"`
	}{
		Task: task,
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(item); err != nil {
		return err
	}

	return sendRequest(url, http.MethodPost, ContentApplicationJson, http.StatusCreated, &body)
}

func completeItem(apiRoot string, id int) error {
	url := fmt.Sprintf("%s/todo/%d?complete", apiRoot, id)
	return sendRequest(url, http.MethodPatch, "", http.StatusNoContent, nil)
}

func deleteItem(apiRoot string, id int) error {
	url := fmt.Sprintf("%s/todo/%d", apiRoot, id)
	return sendRequest(url, http.MethodDelete, "", http.StatusNoContent, nil)
}
