package conduit

import (
	"context"
	"fmt"

	"github.com/isolateminds/go-conduit-cli/internal/compose"
	"github.com/isolateminds/go-conduit-cli/internal/compose/composeopt"
	"github.com/isolateminds/go-conduit-cli/internal/docker"
	"github.com/isolateminds/go-conduit-cli/internal/utils"
)

const (
	//Hard coded for now I guess

	mongoENVTemplateURL    = "https://raw.githubusercontent.com/isolateminds/go-conduit-cli/main/content/env-mongo-temlate.env"
	postgresENVTemplateURL = "https://github.com/isolateminds/go-conduit-cli/blob/main/content/env-postgres-template.env"

	yamlURL = "https://raw.githubusercontent.com/isolateminds/go-conduit-cli/main/content/docker-compose.yml"
)

type Conduit struct {
	compose *compose.Composer
}

func (c *Conduit) Setup(ctx context.Context) error {
	return c.compose.Up(ctx)
}

func NewConduit(ctx context.Context, name string, detatched bool, profiles []string) (*Conduit, error) {
	//Automatically checks if connected to daemon
	client, err := docker.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	compose, err := compose.NewComposer(
		name,
		composeopt.Client(client),
		composeopt.YamlFetchUrl(yamlURL),
		composeopt.Profiles(profiles...),
		applyDetachedFlag(ctx, detatched),
		applyDatabaseEnvFromProfile(ctx, profiles),
	)
	if err != nil {
		return nil, fmt.Errorf("NewConduitError: %s", err)
	}

	return &Conduit{compose: compose}, nil
}

// If detatched no logging will be done same as --detach or -d flag in docker compose
func applyDetachedFlag(ctx context.Context, detached bool) composeopt.SetComposerOptions {
	if detached {
		return composeopt.CustomLogConsumer(nil)
	}
	return composeopt.DefaultComposeLogConsumer(ctx)
}

// Fetches either the mongodb .env  template or the postgres depending on profiles
func applyDatabaseEnvFromProfile(ctx context.Context, profiles []string) composeopt.SetComposerOptions {
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
		return composeopt.Error("DatabaseEnvFromProfileError: cannot use mulpie database profiles")
	}
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
		return composeopt.Error("DatabaseEnvFromProfileError: a database profile has not been selected")
	}
}
