package netopt

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/isolateminds/go-conduit-cli/internal/docker/endpointopt"
)

// SetCreateNetworkOptions is a function type for configuring options when creating a Docker network.
type SetCreateNetworkOptions func(options *types.NetworkCreate)

// CheckDuplicate sets whether duplicate networks should be checked during network creation.
// Use this function to indicate if Docker should check for existing networks with the same name.
func CheckDuplicate() SetCreateNetworkOptions {
	return func(options *types.NetworkCreate) {
		options.CheckDuplicate = true
	}
}

// Driver sets the network driver to be used when creating the Docker network.
// Use this function to specify the network driver that will manage the network's communication.
func Driver(name string) SetCreateNetworkOptions {
	return func(options *types.NetworkCreate) {
		options.Driver = name
	}
}

// Scope sets the scope of the Docker network.
// Use this function to define the network's scope, such as "local" or "global".
func Scope(scope string) SetCreateNetworkOptions {
	return func(options *types.NetworkCreate) {
		options.Scope = scope
	}
}

// EnableIPV6 sets whether IPv6 support should be enabled for the Docker network.
// Use this function to indicate if IPv6 support should be enabled for network communication.
func EnableIPV6() SetCreateNetworkOptions {
	return func(options *types.NetworkCreate) {
		options.EnableIPv6 = true
	}
}

// Internal sets whether the Docker network is intended to be internal.
// Use this function to define whether the network should only be accessible within the host environment.
// If set to true, the network is restricted to communication within containers on the same host.
// If set to false (the default), the network can allow communication between containers across different hosts.
func Internal() SetCreateNetworkOptions {
	return func(options *types.NetworkCreate) {
		options.Internal = true
	}
}

// Attachable sets whether the Docker network is attachable.
// Use this function to specify if other containers can attach to this network.
func Attachable() SetCreateNetworkOptions {
	return func(options *types.NetworkCreate) {
		options.Attachable = true
	}
}

// Ingress sets whether the Docker network is an ingress network.
// Use this function to indicate if the network is the ingress network used for routing externally.
func Ingress() SetCreateNetworkOptions {
	return func(options *types.NetworkCreate) {
		options.Ingress = true
	}
}

// ConfigOnly sets whether the Docker network is a config-only network.
// Use this function to indicate if the network is only used for storing service configuration details.
func ConfigOnly() SetCreateNetworkOptions {
	return func(options *types.NetworkCreate) {
		options.ConfigOnly = true
	}
}

// ConfigFrom specifies the source which provides a network's configuration
func ConfigFrom(net fmt.Stringer) SetCreateNetworkOptions {
	return func(options *types.NetworkCreate) {
		options.ConfigFrom = &network.ConfigReference{
			Network: net.String(),
		}
	}
}

// Options sets custom options for the Docker network during creation.
// Use this function to provide additional key-value pairs for network configuration.
// These options allow you to customize specific behaviors and settings of the network.
func Options(key, value string) SetCreateNetworkOptions {
	return func(options *types.NetworkCreate) {
		if options.Options == nil {
			options.Options = map[string]string{}
		}
		options.Options[key] = value
	}
}

// Label sets labels for the Docker network during creation.
// Use this function to assign custom labels to the network for better organization and identification.
// Labels are key-value pairs that can provide metadata and context to the network.
func Label(key, value string) SetCreateNetworkOptions {
	return func(options *types.NetworkCreate) {
		if options.Labels == nil {
			options.Labels = map[string]string{}
		}
		options.Labels[key] = value
	}
}

// FOR ENDPOINTS ON CONTAINER CREATION
type SetContainerNetworkOptFn func(options *network.NetworkingConfig)

/*
Adds a networking endpoint option for the networking configuration.
*/
func Endpoint(name string, endpoint *endpointopt.Endpoint) SetContainerNetworkOptFn {
	return func(net *network.NetworkingConfig) {
		if net.EndpointsConfig == nil {
			net.EndpointsConfig = make(map[string]*network.EndpointSettings)
		}
		net.EndpointsConfig[name] = endpoint.Settings
	}
}
