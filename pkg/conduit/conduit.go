package conduit

import (
	"context"
	"errors"
	"fmt"

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
	DatabaseName string
	Composer     *compose.Composer
}

func (c *Conduit) Up(ctx context.Context) error {
	return c.Composer.Up(ctx)
}

// For bootsrapping conduit projects and enabling profiles
func NewConduitBootstrapper(ctx context.Context, name string, detatched bool, profiles []string) (*Conduit, error) {
	//Automatically checks if connected to daemon
	client, err := docker.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewConduitError: %s", err)
	}
	db, err := getDatabaseName(profiles)
	if err != nil {
		return nil, fmt.Errorf("NewConduitError: %s", err)
	}
	composer, err := compose.NewComposer(
		name,
		composeopt.Client(client),
		composeopt.YamlFetchUrl(yamlURL),
		composeopt.Profiles(profiles...),
		applyDetachedFlag(ctx, detatched),
		applyEnvFromDatabaseProfile(ctx, db),
	)
	if err != nil {
		return nil, fmt.Errorf("NewConduitError: %s", err)
	}

	return &Conduit{Composer: composer, DatabaseName: db}, nil
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
