package endpointopt

import (
	"github.com/docker/docker/api/types/network"
)

type SetEndpointSettingsFn func(settings *network.EndpointSettings)
type Endpoint struct {
	Settings *network.EndpointSettings
}

func (ew *Endpoint) SetEndpointSetting(setEpSFns ...SetEndpointSettingsFn) {
	for _, set := range setEpSFns {
		set(ew.Settings)
	}
}

func NewEndpoint() *Endpoint {
	return &Endpoint{Settings: &network.EndpointSettings{}}
}

// Adds DriverOpts to the endpoint settings
func DriverOpts(key, value string) SetEndpointSettingsFn {
	return func(es *network.EndpointSettings) {
		if es.DriverOpts == nil {
			es.DriverOpts = make(map[string]string)
		}
		es.DriverOpts[key] = value
	}
}

// Adds a mac address to the endpoint settings
func Mac(mac string) SetEndpointSettingsFn {
	return func(es *network.EndpointSettings) {
		es.MacAddress = mac
	}
}

// Adds GlobalIPv6PrefixLen to endpoint settings
func GlobalIPv6PrefixLen(prefixLen int) SetEndpointSettingsFn {
	return func(es *network.EndpointSettings) {
		es.GlobalIPv6PrefixLen = prefixLen
	}
}

// Adds GlobalIPv6Address to endpoint settings
func GlobalIPv6Address(addr string) SetEndpointSettingsFn {
	return func(es *network.EndpointSettings) {
		es.GlobalIPv6Address = addr
	}
}

// Adds a ipv6 gateway to the endpoint settings
func IPv6Gateway(gateway string) SetEndpointSettingsFn {
	return func(es *network.EndpointSettings) {
		es.IPv6Gateway = gateway
	}
}

// Adds IPPrefixLen to the endpoint settings
func IPPrefixLen(prefixLen int) SetEndpointSettingsFn {
	return func(es *network.EndpointSettings) {
		es.IPPrefixLen = prefixLen
	}
}

// Adds a IP address to the endpoint settings
func IPAddress(addr string) SetEndpointSettingsFn {
	return func(es *network.EndpointSettings) {
		es.IPAddress = addr
	}
}

// Adds a gateway to the endpoint settings
func Gateway(gateway string) SetEndpointSettingsFn {
	return func(es *network.EndpointSettings) {
		es.Gateway = gateway
	}
}

// Adds EnpointID to endpoint settings
func EndpointID(id string) SetEndpointSettingsFn {
	return func(es *network.EndpointSettings) {
		es.EndpointID = id
	}
}

// Adds links to the endpoint settings
func Links(links ...string) SetEndpointSettingsFn {
	return func(es *network.EndpointSettings) {
		if es.Links == nil {
			es.Links = make([]string, 0)
		}
		es.Links = append(es.Links, links...)
	}
}

// Adds aliases to the endpoint settings
func Aliases(aliases ...string) SetEndpointSettingsFn {
	return func(es *network.EndpointSettings) {
		if es.Aliases == nil {
			es.Aliases = make([]string, 0)
		}
		es.Aliases = append(es.Aliases, aliases...)
	}
}

// EndpointIPAMConfig represents IPAM Configurations for the endpoint
func IPAM(IPv4, IPv6 string, linkLocalIPs []string) SetEndpointSettingsFn {

	return func(es *network.EndpointSettings) {
		es.IPAMConfig = &network.EndpointIPAMConfig{
			IPv4Address:  IPv4,
			IPv6Address:  IPv6,
			LinkLocalIPs: linkLocalIPs,
		}
	}
}
