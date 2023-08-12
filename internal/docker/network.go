package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/isolateminds/go-conduit-cli/internal/docker/netopt"
)

type Network struct {
	Id      string
	Name    string
	options *types.NetworkCreate
}

func (n *Network) String() string {
	return n.Name
}

func (n *Network) SetOptions(setNOFns ...netopt.SetCreateNetworkOptions) {
	for _, set := range setNOFns {
		set(n.options)
	}
}
func (c *Client) NewNetwork(name string) *Network {
	return &Network{
		Name:    name,
		options: &types.NetworkCreate{},
	}
}
