<div align="center">
    <br>
    <a href="https://getconduit.dev" target="_blank"><img src="https://getconduit.dev/conduitLogo.svg" height="80px" alt="logo"/></a>
    <br/>
    <h3>The only Backend you'll ever need.</h3>
</div>

# go-conduit-cli
A Golang port of [ConduitPlatform/CLI](https://github.com/ConduitPlatform/CLI) with some additional features

GoConduitCli is a multipurpose tool that's meant to facilitate your development experience and speed up your work
regardless of whether you're deploying a Conduit instance for your project, developing custom modules or even
contributing to the upstream project in your spare time.


# Commands
<!-- commands -->
<!-- * [`conduit cli update`](#conduit-cli-update) -->
* [`goconduit deploy setup`](#goconduit-deploy-setup)
* [`goconduit deploy start`](#goconduit-deploy-start)
* [`goconduit deploy stop`](#goconduit-deploy-stop)
* [`goconduit deploy rm`](#goconduit-deploy-rm)
* [`goconduit deploy recreate`](#goconduit-deploy-recreate)

<!-- * [`conduit deploy update`](#conduit-deploy-update) -->
<!-- * [`conduit generateClient graphql`](#conduit-generateclient-graphql) -->
<!-- * [`conduit generateClient rest`](#conduit-generateclient-rest) -->
<!-- * [`conduit generateSchema [PATH]`](#conduit-generateschema-path) -->
<!-- * [`conduit help [COMMAND]`](#conduit-help-command) -->
<!-- * [`conduit init`](#conduit-init) -->

<!-- ## `conduit cli update`

Update your CLI

```
USAGE
  $ conduit cli update

DESCRIPTION
  Update your CLI
``` -->

## `goconduit deploy setup`

Bootstrap a local Conduit deployment

```
USAGE
  $ goconduit deploy setup --profiles <value>,<value> [--project-name <value>] [--ui-image-tag <value>] [--image-tag <value>] [--detach] [--mount-database]

FLAGS
  --profiles        profiles to enable (one database profile is required either mongodb or postgres)

  --project-name    set the project name (defaults to conduit)

  --ui-image-tag    set conduit-ui image tag (defaults to latest)

  --image-tag       set all other conduit image tag (defaults to latest)

  --detach          set detach mode to disable console log output (defaults to false)

  --mount-database  enable this to bind mount postgres or mongodb container to project directory (defaults to false). if this is not set it will use persistent volumes


```

## `goconduit deploy start`

Bring up your local Conduit deployment

```
USAGE
  $ goconduit deploy start [--profiles <value>,<value>] [--detach]

DESCRIPTION
  Bring up your local Conduit deployment

FLAGS
  --profiles    profiles to enable

  --detach      set detach mode to disable console log output (defaults to false)

```

## `goconduit deploy stop`

Bring down your local Conduit deployment

```
USAGE
  $ goconduit deploy stop [--services <value>,<value>]
FLAGS
  --services    services to stop
```

## `goconduit deploy rm`

Remove your local Conduit deployment

```
USAGE
  $ goconduit deploy rm [--services <value>,<value>]
FLAGS
  --services    services to remove
```

## `goconduit deploy recreate`

recreate your local Conduit deployment containers

```
USAGE
  $ goconduit deploy recreate
```