package main

import "github.com/lasfh/eletrize/vscode"

func loadVSCodeLaunch(workspaceDir string) (*Eletrize, error) {
	launch, err := vscode.LoadLaunch(workspaceDir)
	if err != nil {
		return nil, err
	}

	eletrize := Eletrize{
		launch: true,
	}

	for _, config := range launch.Configurations {
		schema, ok := config.Schema(workspaceDir)
		if ok {
			eletrize.Schema = append(eletrize.Schema, schema)
		}
	}

	if len(eletrize.Schema) == 0 {
		return nil, vscode.ErrNoLaunchDetected
	}

	return &eletrize, nil
}
