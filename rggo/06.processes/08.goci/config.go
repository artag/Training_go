package main

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

const (
	ConfigFilename = "config.json"
)

type stepType int

const (
	Step stepType = iota
	ExceptionStep
	TimeoutStep
)

type dataJson struct {
	Steps []stepJson `json:"steps"`
}

type stepJson struct {
	Type    stepType `json:"type"`
	Name    string   `json:"name"`
	Exe     string   `json:"exe"`
	Args    []string `json:"args"`
	Message string   `json:"message"`
	Timeout int64    `json:"timeout"`
}

func LoadPipeline(proj string) ([]executer, error) {
	data, err := load()
	if err != nil {
		return make([]executer, 0), err
	}

	return parseData(data, proj), nil
}

func load() (*dataJson, error) {
	file, err := ioutil.ReadFile(ConfigFilename)
	if err != nil {
		return nil, err
	}

	if len(file) == 0 {
		return nil, ErrEmptyConfigFile
	}

	var js dataJson
	if err = json.Unmarshal(file, &js); err != nil {
		return nil, err
	}

	return &js, nil
}

func parseData(data *dataJson, proj string) []executer {
	result := make([]executer, 0)
	for _, step := range data.Steps {
		switch step.Type {
		case Step:
			s := newStep(step.Name, step.Exe, step.Message, proj, step.Args)
			result = append(result, s)

		case ExceptionStep:
			s := newExceptionStep(step.Name, step.Exe, step.Message, proj, step.Args)
			result = append(result, s)

		case TimeoutStep:
			s := newTimeoutStep(step.Name, step.Exe, step.Message, proj, step.Args, time.Duration(step.Timeout))
			result = append(result, s)
		}
	}
	return result
}
