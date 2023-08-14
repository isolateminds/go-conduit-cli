package conduit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/isolateminds/go-conduit-cli/internal/compose"
	"github.com/isolateminds/go-conduit-cli/internal/compose/composeopt"
	"github.com/isolateminds/go-conduit-cli/internal/docker"
	"github.com/isolateminds/go-conduit-cli/internal/utils"
)

const (
	//Hard coded for now I guess

	mongoENVTemplateURL    = "https://raw.githubusercontent.com/isolateminds/go-conduit-cli/main/content/env-mongo-temlate.env"
	postgresENVTemplateURL = "https://raw.githubusercontent.com/isolateminds/go-conduit-cli/main/content/env-postgres-template.env"

	yamlURL = "https://raw.githubusercontent.com/isolateminds/go-conduit-cli/main/content/docker-compose.yml"
)

type Conduit struct {
	Composer *compose.Composer
	Json     *ConduitJson
}

func (c *Conduit) Stop(ctx context.Context, services []string) error {
	return c.Composer.Stop(ctx, services)
}

func (c *Conduit) Up(ctx context.Context) error {
	return c.Composer.Up(ctx)
}

// For already bootstrapped projects must be in project root dir when you call this
func NewConduitFromProject(ctx context.Context, detached bool, profiles []string) (*Conduit, error) {
	//Automatically checks if connected to daemon
	client, err := docker.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewConduitFromProject: %s", err)
	}
	//Load persisted data
	data := &ConduitJson{}
	b, err := ioutil.ReadFile("conduit.json")
	if err != nil {
		return nil, fmt.Errorf("NewConduitFromProject: %s", err)
	}
	err = json.Unmarshal(b, data)
	if err != nil {
		return nil, fmt.Errorf("NewConduitFromProject: %s", err)
	}

	//block databases from being added
	blockList := append(data.Profiles, "postgres", "mongodb")
	//new profiles are ones who haven't been added before and any databases
	newProfiles := BlockProfiles(profiles, blockList...)
	//the updated profiles are a mix of new and saved profiles (conduit.json)
	updatedProfiles := append(data.Profiles, newProfiles...)

	composer, err := compose.NewComposer(
		data.ProjectName,
		composeopt.Client(client),
		composeopt.EnvFromFile(".env"),
		composeopt.YamlFromFile("docker-compose.yaml"),
		applyDetachedFlag(ctx, detached),
		composeopt.Profiles(updatedProfiles...),
	)
	if err != nil {
		return nil, fmt.Errorf("NewConduitFromProject: %s", err)
	}

	return &Conduit{
		Composer: composer,
		Json: &ConduitJson{
			ProjectName: data.ProjectName,
			Version:     data.Version,
			Database:    data.Database,
			//filter the profiles here to save the actual profiles defined in the schema
			Profiles: composer.FilterYamlProfiles(updatedProfiles),
		},
	}, nil
}

// For bootsrapping conduit projects and enabling profiles
func NewConduitBootstrapper(ctx context.Context, name string, detached bool, profiles []string) (*Conduit, error) {
	//Automatically checks if connected to daemon
	client, err := docker.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewConduitBootstrapperError: %s", err)
	}
	db, err := getDatabaseName(profiles)
	if err != nil {
		return nil, fmt.Errorf("NewConduitBootstrapperError: %s", err)
	}
	composer, err := compose.NewComposer(
		name,
		composeopt.Client(client),
		composeopt.YamlFetchUrl(yamlURL),
		composeopt.Profiles(profiles...),
		applyDetachedFlag(ctx, detached),
		applyEnvFromDatabaseProfile(ctx, db),
	)
	if err != nil {
		return nil, fmt.Errorf("NewConduitBootstrapperError: %s", err)
	}

	return &Conduit{
		Composer: composer,
		Json: &ConduitJson{
			ProjectName: name,
			Database:    db,
			//filter the profiles here to save the actual profiles defined in the schema
			Profiles: composer.FilterYamlProfiles(profiles),
		},
	}, nil
}

// If detatched no logging will be done same as --detach or -d flag in docker compose
func applyDetachedFlag(ctx context.Context, detached bool) composeopt.SetComposerOptions {
	if detached {
		return composeopt.CustomLogConsumer(nil)
	}
	return composeopt.DefaultComposeLogConsumer(ctx)
}

// Fetches either the mongodb .env  template or the postgres depending on profiles
func applyEnvFromDatabaseProfile(ctx context.Context, db string) composeopt.SetComposerOptions {

	dbPass := utils.GenerateRandomString(32)
	switch db {
	case "mongodb":
		formatter := newTemplateFormatter(variableMap{
			"MasterKey":     utils.GenerateRandomString(64),
			"MongoPassword": dbPass,
		})
		return composeopt.TemplateEnvFetchUrl(mongoENVTemplateURL, formatter)
	case "postgres":
		formatter := newTemplateFormatter(variableMap{
			"MasterKey":        utils.GenerateRandomString(64),
			"PostgresPassword": dbPass,
		})
		return composeopt.TemplateEnvFetchUrl(postgresENVTemplateURL, formatter)
	default:
		return composeopt.Error("DatabaseEnvFromProfileError: a database profile has not been given")
	}
}

func BlockProfiles(profiles []string, blocked ...string) []string {
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

func getDatabaseName(profiles []string) (string, error) {
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
		return "", errors.New("GetDatabaseNameError: cannot use multiple database profiles")
	}
	return db, nil
}
