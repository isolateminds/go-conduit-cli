package conduit

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
)

/*
This is what tells the  go-conduit-cli that the current directory
(if it exists) is a conduit project think of it as a conduit only package.json
*/
type ConduitJson struct {
	ProjectName string   `json:"projectName"`
	Version     string   `json:"version"`
	Database    string   `json:"database"`
	Profiles    []string `json:"profiles"`
}

func WriteConduitJSON(projectName, version, database string, profiles []string) error {
	cj := &ConduitJson{projectName, version, database, profiles}
	b, err := json.Marshal(cj)
	if err != nil {
		return fmt.Errorf("WriteConduitJSONError: %s", err)
	}
	err = os.WriteFile("conduit.json", b, fs.ModePerm)
	if err != nil {
		return fmt.Errorf("WriteConduitJSONError: %s", err)
	}
	return nil
}
