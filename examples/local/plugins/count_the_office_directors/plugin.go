package tf

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

type APIRoot struct {
	Data []Episode `json:"data"`
}

type Episode struct {
	Writer   Employee `json:"writer"`
	Director Employee `json:"director"`
}

type Employee struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type EmployeeCountResult struct {
	Name             string `json:"name"`
	WrittenEpisodes  int    `json:"directed"`
	DirectedEpisodes int    `json:"written"`
}

// ProcessorPluginV1 will load the results of https://officeapi.dev/api/episodes/ and
// count the number of episodes directed and/or written by person.
func ProcessorPluginV1(ctx context.Context, inputData string, vars map[string]string) (string, error) {
	// Load options.
	countDirected, _ := strconv.ParseBool(vars["directed"])
	countWritten, _ := strconv.ParseBool(vars["written"])
	if !countDirected && !countWritten {
		return "", fmt.Errorf("directed and written counts disabled, at least one must be used")
	}

	// Load input data.
	root := APIRoot{}
	err := json.Unmarshal([]byte(inputData), &root)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal JSON input: %w", err)
	}

	// Count info from episodes.
	res := map[string]*EmployeeCountResult{}
	for _, episode := range root.Data {
		employeeDirector, ok := res[episode.Director.ID]
		if !ok {
			employeeDirector = &EmployeeCountResult{Name: episode.Director.Name}
			res[episode.Director.ID] = employeeDirector
		}

		employeeWriter, ok := res[episode.Writer.ID]
		if !ok {
			employeeWriter = &EmployeeCountResult{Name: episode.Writer.Name}
			res[episode.Writer.ID] = employeeWriter
		}

		if countDirected {
			employeeDirector.DirectedEpisodes++
		}

		if countWritten {
			employeeWriter.WrittenEpisodes++
		}
	}

	result, err := json.Marshal(res)
	if err != nil {
		return "", fmt.Errorf("could not marshall JSON result: %w", err)
	}

	return string(result), nil
}
