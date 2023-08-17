package conduit

import (
	"encoding/json"
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

// Writes the conduit.json file to current path
func (cj *ConduitJson) WriteFile() error {
	b, err := json.MarshalIndent(cj, "", "	") // The last argument is the indentation prefix
	if err != nil {
		return err
	}
	err = os.WriteFile("conduit.json", b, fs.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
