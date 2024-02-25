package main

import (
	"encoding/json"         // To customize JSON output (marshal input data to bytes slice)
	"rggo/interacting/todo" // To-do application (from previous chapters)
	"time"                  // To worj with time functions (get current time)
)

type todoResponse struct {
	Results todo.List `json:"results"`
}

func (r *todoResponse) MarshalJSON() ([]byte, error) {
	resp := struct {
		Results      todo.List `json:"results"`
		Date         int64     `json:"date"`
		TotalResults int       `json:"total_results"`
	}{
		Results:      r.Results,
		Date:         time.Now().Unix(),
		TotalResults: len(r.Results),
	}
	return json.Marshal(resp)
}
