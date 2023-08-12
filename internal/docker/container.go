package docker

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/isolateminds/go-conduit-cli/internal/docker/containeropt"
	"github.com/isolateminds/go-conduit-cli/internal/docker/hostopt"
	"github.com/isolateminds/go-conduit-cli/internal/docker/netopt"
	"github.com/isolateminds/go-conduit-cli/internal/docker/platformopt"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// Container represents a Docker container along with its configuration options.
type Container struct {
	Id                string
	Name              string
	options           *container.Config
	hostOptions       *container.HostConfig
	networkingOptions *network.NetworkingConfig
	platformOptions   *v1.Platform
}

// String returns the name of the Docker container.
func (c *Container) String() string {
	return c.Name
}

// SetHostOptions configures host-related options for the Docker container.
// Use this method to set various host options using functions from the hostopt package.
func (c *Container) SetHostOptions(setHOFns ...hostopt.SetHostOptFn) {
	for _, set := range setHOFns {
		set(c.hostOptions)
	}
}

// SetNetworkOptions configures network-related options for the Docker container.
// Use this method to set various network options using functions from the netopt package.
func (c *Container) SetNetworkOptions(setNwOptFns ...netopt.SetContainerNetworkOptFn) {
	for _, set := range setNwOptFns {
		set(c.networkingOptions)
	}
}

// SetOptions configures options for the Docker container.
// Use this method to set various container options using functions from the containeropt package.
func (c *Container) SetOptions(setOFns ...containeropt.SetOptionsFns) {
	for _, set := range setOFns {
		set(c.options)
	}
}
func (c *Container) SetPlatformOptions(setPOFns ...platformopt.SetPlatformOptions) {
	for _, set := range setPOFns {
		set(c.platformOptions)
	}
}

// NewContainer creates a new Container instance with the specified name.
// The Container instance contains configuration options for creating a Docker container.
func (*Client) NewContainer(name string) *Container {
	container := &Container{
		Id:                "",
		Name:              name,
		options:           &container.Config{},
		hostOptions:       &container.HostConfig{},
		networkingOptions: &network.NetworkingConfig{},
		platformOptions:   &v1.Platform{},
	}

	return container
}
