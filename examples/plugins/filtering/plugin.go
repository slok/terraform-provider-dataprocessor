package tf

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
)

type User struct {
	Username string
	Age      int
}

// ProcessorPluginV1 Will take a list of users as input and will filter
// them by a regex against its username, it will return the list again
// without the ones that matched.
func ProcessorPluginV1(ctx context.Context, inputData string, vars map[string]string) (string, error) {
	// Get filter regex.
	regexOpt := vars["username_filter"]
	if regexOpt == "" {
		regexOpt = ".*"
	}
	regex, err := regexp.Compile(regexOpt)
	if err != nil {
		return "", fmt.Errorf("regex %q could not be compiled: %w", regexOpt, err)
	}

	// Load input data.
	users := []User{}
	err = json.Unmarshal([]byte(inputData), &users)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal input into JSON: %w", err)
	}

	// Filter users if no match and sort result.
	resultUsers := []User{}
	for _, u := range users {
		if !regex.MatchString(u.Username) {
			resultUsers = append(resultUsers, u)
		}
	}
	sort.SliceStable(resultUsers, func(i, j int) bool { return resultUsers[i].Age < resultUsers[j].Age })

	result, err := json.Marshal(resultUsers)
	if err != nil {
		return "", fmt.Errorf("could not marshal result into JSON: %w", err)
	}

	return string(result), nil
}
