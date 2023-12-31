package conduit

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/isolateminds/go-conduit-cli/internal/compose"
	"github.com/isolateminds/go-conduit-cli/internal/compose/composeopt"
	"github.com/isolateminds/go-conduit-cli/internal/docker"
	"github.com/isolateminds/go-conduit-cli/internal/utils"
	"github.com/isolateminds/go-conduit-cli/pkg/conduit/errordefs"
)

const (
	//Hard coded for now I guess

	mongoEnvTemplateURL    = "https://raw.githubusercontent.com/isolateminds/go-conduit-cli/main/templates/mongo.env"
	postgresEnvTemplateURL = "https://raw.githubusercontent.com/isolateminds/go-conduit-cli/main/templates/postgres.env"

	dockerComposeTemplateURL = "https://raw.githubusercontent.com/isolateminds/go-conduit-cli/main/templates/docker-compose.yml"
)

type Conduit struct {
	composer *compose.Composer
	json     *ConduitJson
}

func (c *Conduit) Remove(ctx context.Context, services []string) error {
	return c.composer.Remove(ctx, services)
}
func (c *Conduit) Stop(ctx context.Context, services []string) error {
	return c.composer.Stop(ctx, services)
}

func (c *Conduit) Up(ctx context.Context) error {
	return c.composer.Up(ctx)
}
func (c *Conduit) Start(ctx context.Context, services []string) error {
	return c.composer.Start(ctx, services)
}
func (c *Conduit) Create(ctx context.Context, services []string) error {
	return c.composer.Create(ctx, services)
}

// For already bootstrapped projects must be in project root dir when you call this
func NewConduitFromProject(ctx context.Context, detached bool, profiles []string) (*Conduit, error) {
	//Automatically checks if connected to daemon
	client, err := docker.NewClient(ctx)
	if err != nil {
		return nil, errordefs.NewConduitFromProjectError(err)
	}
	//Load persisted data
	data := &ConduitJson{}
	b, err := ioutil.ReadFile("conduit.json")
	if err != nil {
		return nil, errordefs.NewConduitFromProjectError(err)
	}
	err = json.Unmarshal(b, data)
	if err != nil {
		return nil, errordefs.NewConduitFromProjectError(err)
	}

	//block databases from being added because they where added already during bootstrapping
	//block the profiles specified inside conduit.json
	blockList := append(data.Profiles, "postgres", "mongodb")
	//new profiles are ones who haven't been added before
	newProfiles := blockProfiles(profiles, blockList...)
	//the updated profiles are a mix of new and saved profiles (conduit.json)
	updatedProfiles := append(data.Profiles, newProfiles...)

	composer, err := compose.NewComposer(
		data.ProjectName,
		withDetachedFlag(ctx, detached),
		composeopt.WithClient(client),
		composeopt.WithEnvFromFile(".env"),
		composeopt.WithYamlFromFile("docker-compose.yaml"),
		//the profiles added here will be used with dcoker compose
		composeopt.WithProfiles(updatedProfiles...),
	)
	if err != nil {
		return nil, errordefs.NewConduitFromProjectError(err)
	}

	return &Conduit{
		composer: composer,
		json: &ConduitJson{
			ProjectName: data.ProjectName,
			Version:     data.Version,
			Database:    data.Database,
			//filter the profiles here to save the actual profiles defined in the schema
			Profiles: composer.FilterYamlProfiles(updatedProfiles),
		},
	}, nil
}

type BootstrapperOptions struct {
	ProjectName   string
	Profiles      []string
	Detached      bool
	ImageTag      string
	UIImageTag    string
	MountDatabase bool
}

// For bootsrapping conduit projects and enabling profiles
func NewConduitBootstrapper(ctx context.Context, options *BootstrapperOptions) (*Conduit, error) {
	//Automatically checks if connected to daemon
	client, err := docker.NewClient(ctx)
	if err != nil {
		return nil, errordefs.NewConduitBootstrapperError(err)
	}
	db, err := ensureProperDatabase(options.Profiles)
	if err != nil {
		return nil, errordefs.NewConduitBootstrapperError(err)
	}
	composer, err := compose.NewComposer(
		options.ProjectName,
		composeopt.WithClient(client),
		withYamlBasedOnDatabaseBind(ctx, db, options),
		composeopt.WithProfiles(options.Profiles...),
		withDetachedFlag(ctx, options.Detached),
		withEnvBasedOnDatabaseProfile(ctx, db, options),
	)
	if err != nil {
		return nil, errordefs.NewConduitBootstrapperError(err)
	}
	return &Conduit{
		composer: composer,
		json: &ConduitJson{
			ProjectName: options.ProjectName,
			Database:    db,
			//filter the profiles here to save the actual profiles defined in the schema
			Profiles: composer.FilterYamlProfiles(options.Profiles),
		},
	}, nil
}

// Writes the docker-compose.yaml to current path
func (c *Conduit) WriteComposeFile() error {
	return os.WriteFile("docker-compose.yaml", c.composer.Options.Yaml.Bytes, fs.ModePerm)
}

// Writes the .env to current path
func (c *Conduit) WriteEnvFile() error {
	return os.WriteFile(".env", c.composer.Options.Environment.Bytes, fs.ModePerm)
}

// Writes the json file to current path
func (c *Conduit) WriteConduitJsonFile() error {
	return c.json.WriteFile()
}

// modifies the docker compose file if MountDatabase set
func withYamlBasedOnDatabaseBind(ctx context.Context, db string, options *BootstrapperOptions) composeopt.SetComposerOptions {
	if options.MountDatabase {
		return composeopt.WithYamlFromUrlFormatter(dockerComposeTemplateURL, newComposeBindDbFormatter(db))
	}
	return composeopt.WithYamlFromUrl(dockerComposeTemplateURL)
}

// If detatched no logging will be done, same as --detach or -d flag in docker compose
func withDetachedFlag(ctx context.Context, detached bool) composeopt.SetComposerOptions {
	if detached {
		return composeopt.WithCustomLogConsumer(nil)
	}
	return composeopt.WithDefaultComposeLogConsumer(ctx)
}

/*
Fetches either the mongodb .env  template or the postgres one depending on profiles
and formats the env template
*/
func withEnvBasedOnDatabaseProfile(ctx context.Context, db string, options *BootstrapperOptions) composeopt.SetComposerOptions {
	dbPass := utils.GenerateRandomString(32)
	masterKey := utils.GenerateRandomString(64)
	vMap := variableMap{
		"MasterKey":   masterKey,
		"ImageTag":    options.ImageTag,
		"UIImageTag":  options.UIImageTag,
		"ProjectName": options.ProjectName,
	}
	switch db {
	case "mongodb":
		vMap["MongoPassword"] = dbPass
		return composeopt.WithEnvFromUrlFormatter(mongoEnvTemplateURL, newEnvFormatter(vMap))
	case "postgres":
		vMap["PostgresPassword"] = dbPass
		return composeopt.WithEnvFromUrlFormatter(postgresEnvTemplateURL, newEnvFormatter(vMap))
	default:
		return composeopt.WithError("a database profile has not been given use")
	}
}

// Helper function 'blocks' profiles from being included in the result.
// It takes a list of profiles and a variable number of profiles to be blocked.
// Blocked profiles are excluded from the result list.
func blockProfiles(profiles []string, blocked ...string) []string {
	blockedSet := make(map[string]struct{})
	for _, b := range blocked {
		blockedSet[b] = struct{}{}
	}

	result := []string{}
	for _, p := range profiles {
		if _, ok := blockedSet[p]; !ok {
			result = append(result, p)
		}
	}
	return result
}

// Gets the database name either mongodb or postgres if it finds two there will be an error
func ensureProperDatabase(profiles []string) (database string, err error) {
	var db string
	var dbs []string
	for _, profile := range profiles {
		if profile == "mongodb" {
			db = profile
			dbs = append(dbs, db)
		}
		if profile == "postgres" {
			db = profile
			dbs = append(dbs, db)
		}
	}
	if len(dbs) > 1 {
		return "", errors.New("cannot use multiple database profiles")
	}
	return db, nil
}
