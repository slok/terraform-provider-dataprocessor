package tf

import (
	"context"
	"encoding/json"
	"fmt"
)

type WorkspaceIn struct {
	Name     string `json:"name"`
	RepoPath string `json:"repo_path"`
	Owner    string `json:"owner"`
}

type WorkspaceOut struct {
	Name     string `json:"name"`
	RepoPath string `json:"repo_path"`
}

type TeamOut struct {
	Name       string         `json:"name"`
	Workspaces []WorkspaceOut `json:"workspaces"`
}

// ProcessorPluginV1 will transform the received data by aggregating the workspaces by owner.
func ProcessorPluginV1(ctx context.Context, inputData string, vars map[string]string) (string, error) {
	workspaces := []WorkspaceIn{}
	err := json.Unmarshal([]byte(inputData), &workspaces)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal JSON input: %w", err)
	}

	// Transform data aggregating workspaces by owner(team).
	teams := map[string]*TeamOut{}
	for _, wks := range workspaces {
		team, ok := teams[wks.Owner]
		if !ok {
			team = &TeamOut{Name: wks.Owner}
			teams[wks.Owner] = team
		}
		team.Workspaces = append(team.Workspaces, WorkspaceOut{
			Name:     wks.Name,
			RepoPath: wks.RepoPath,
		})
	}

	result, err := json.Marshal(teams)
	if err != nil {
		return "", fmt.Errorf("could not marshall JSON result: %w", err)
	}

	return string(result), nil
}
