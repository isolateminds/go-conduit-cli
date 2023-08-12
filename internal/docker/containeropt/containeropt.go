package containeropt

import (
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
)

type SetOptionsFns func(options *container.Config)

/*
Adds a health check to the container configuration that exec arguments directly

	container := client.NewContainer("my_container")
	container.SetOptions(
		containeropt.HealthCheckExec(
			time.Second*5,  //start
			time.Second*20, //timeout
			time.Second*5, //interval
			10, //timeout
			"CMD-SHELL", "curl", "-f", "http://localhost", "||", "exit", "1",
		)
	)
*/
func HealthCheckExec(start, timeout, interval time.Duration, retries int, args ...string) SetOptionsFns {
	return func(Config *container.Config) {
		Config.Healthcheck = &container.HealthConfig{
			StartPeriod: start,
			Timeout:     timeout,
			Interval:    interval,
			Test:        args,
			Retries:     retries,
		}
	}
}

/*
Disables the health check

	container := client.NewContainer("my_container")
	container.SetOptions(
		containeropt.DisableHealthCheck()
	)
*/
func DisableHealthCheck() SetOptionsFns {
	return func(Config *container.Config) {
		Config.Healthcheck = &container.HealthConfig{
			Test: []string{"None"},
		}
	}
}

/*
Exposes a port

	container := client.NewContainer("my_container")
	container.SetOptions(
		containeropt.Expose("8000")
	)
*/
func Expose(containerPort string) SetOptionsFns {
	return func(Config *container.Config) {
		if Config.ExposedPorts == nil {
			Config.ExposedPorts = make(nat.PortSet)
		}
		Config.ExposedPorts[nat.Port(containerPort)] = struct{}{}
	}
}

/*
Adds a hostname to the container configuration.

	container := client.NewContainer("my_container")
	container.SetOptions(
		containeropt.Hostname("com.example.localhost"),
	)
*/
func Hostname(name string) SetOptionsFns {
	return func(Config *container.Config) {
		Config.Hostname = name
	}
}

/*
Adds a domain name to the container configuration.

	container := client.NewContainer("my_container")
	container.SetOptions(
		containeropt.DomainName("com.example.localhost"),
	)
*/
func DomainName(name string) SetOptionsFns {
	return func(Config *container.Config) {
		Config.Domainname = name
	}
}

/*
Adds a image to the container configuration.

	image := client.NewImage("alpine")
	container := client.NewContainer("my_container")
	container.SetOptions(
		containeropt.Image(image),
	)
*/
func Image(image fmt.Stringer) SetOptionsFns {
	return func(Config *container.Config) {
		Config.Image = image.String()
	}
}

/*
Adds a command to the container configuration.

	container := client.NewContainer("my_container")
	container.SetOptions(
		CMD("/bin/sh", "-c", "echo hello"),
	)
*/
func CMD(cmd ...string) SetOptionsFns {
	return func(Config *container.Config) {
		if Config.Cmd == nil {
			Config.Cmd = make(strslice.StrSlice, 0)
		}
		Config.Cmd = append(Config.Cmd, cmd...)
	}
}

/*
Sets a user for the container configuration

	container := client.NewContainer("my_container")
	container.SetOptions(
		containeropt.User("my_user"),
	)
*/
func User(user string) SetOptionsFns {
	return func(Config *container.Config) {
		Config.User = user
	}
}

/*
Sets attatch stdin to true in the container configuration

	container := client.NewContainer("my_container")
	container.SetOptions(
		containeropt.AttachStdin(),
	)
*/
func AttachStdin() SetOptionsFns {
	return func(Config *container.Config) {
		Config.AttachStdin = true
	}
}

/*
Sets attach stdout to true in the container configuration

	container := client.NewContainer("my_container")
	container.SetOptions(
		containeropt.AttachStdout(),
	)
*/
func AttachStdout() SetOptionsFns {
	return func(Config *container.Config) {
		Config.AttachStdout = true
	}
}

/*
Sets attach stderr to true in the container configuration

	container := client.NewContainer("my_container")
	container.SetOptions(
		AttachStderr(),
	)
*/
func AttachStderr() SetOptionsFns {
	return func(Config *container.Config) {
		Config.AttachStderr = true
	}
}

/*
Sets TTY to true in the container configuration

	container := client.NewContainer("my_container")
	container.SetOptions(
		containeropt.TTY(),
	)
*/
func TTY() SetOptionsFns {
	return func(Config *container.Config) {
		Config.Tty = true
	}
}

/*
Sets OpenStdin to true in the container configuration

	container := client.NewContainer("my_container")
	container.SetOptions(
		containeropt.OpenStdin(),
	)
*/
func OpenStdin() SetOptionsFns {
	return func(Config *container.Config) {
		Config.OpenStdin = true
	}
}

/*
Sets StdinOnce to true in the container configuration
that closes stdin after the 1 attached client disconnects.

	container := client.NewContainer("my_container")
	container.SetOptions(
		containeropt.StdinOnce(),
	)
*/
func StdinOnce() SetOptionsFns {
	return func(Config *container.Config) {
		Config.StdinOnce = true
	}
}

/*
Sets ArsExcaped to true in the container configuration.
Use if command is already escaped (meaning treat as a command line) (Windows specific).

	container := client.NewContainer("my_container")
	container.SetOptions(
		containeropt.EscapeArgs(),
	)
*/
func EscapeArgs() SetOptionsFns {
	return func(Config *container.Config) {
		Config.ArgsEscaped = true
	}
}

/*
Adds a volume to the container configuration.

	volume := clinet.NewVolume("my_volume")
	container := client.NewContainer("my_container")
	container.SetOptions(
		containeropt.Volume(volume),
	)
*/
func Volume(volume fmt.Stringer) SetOptionsFns {
	return func(Config *container.Config) {
		if Config.Volumes == nil {
			Config.Volumes = make(map[string]struct{})
		}
		Config.Volumes[volume.String()] = struct{}{}
	}
}

/*
Sets the working directory for the container configuration.

	container := client.NewContainer("my_container")
	container.SetOptions(
		containeropt.WorkingDir("/my/working/directory"),
	)
*/
func WorkingDir(dir string) SetOptionsFns {
	return func(Config *container.Config) {
		Config.WorkingDir = dir
	}
}

/*
Sets The network to diabled in the container configuration.

	container := client.NewContainer("my_container")
	container.SetOptions(
		containeropt.DisableNetwork(),
	)
*/
func DisableNetwork() SetOptionsFns {
	return func(Config *container.Config) {
		Config.NetworkDisabled = true
	}
}

/*
Sets the ONBUILD metadata that were defined on the image Dockerfile

	container := client.NewContainer("my_container")
	container.SetOptions(
		containeropt.OnBuild("ADD", "."),
	)
*/
func OnBuild(args ...string) SetOptionsFns {
	return func(Config *container.Config) {
		if Config.OnBuild == nil {
			Config.OnBuild = []string{}
		}
		Config.OnBuild = append(Config.OnBuild, args...)
	}
}

/*
Adds a label to the container configuration.

	container := client.NewContainer("my_container")
	container.SetOptions(
		containeropt.Label("my_label", "my_value"),
	)
*/
func Label(label, value string) SetOptionsFns {
	return func(Config *container.Config) {
		if Config.Labels == nil {
			Config.Labels = make(map[string]string)
		}
		Config.Labels[label] = value
	}
}

/*
Adds a StopSignal to the container configuration.

	container := client.NewContainer("my_container")
	container.SetOptions(
		containeropt.StopSignal("SIGTERM"),
	)
*/
func StopSignal(signal string) SetOptionsFns {
	return func(Config *container.Config) {
		Config.StopSignal = signal
	}
}
