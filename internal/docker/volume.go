package docker

import (
	"github.com/docker/docker/api/types/volume"
	"github.com/isolateminds/go-conduit-cli/internal/docker/volumeopt"
)

// Volume represents a Docker volume along with its creation options.
type Volume struct {
	options *volume.CreateOptions
}

// returns the volume name
func (v *Volume) String() string {
	return v.options.Name
}

// SetOptions configures options for the Docker volume.
// Use this method to set various volume options using functions from the volumeopt package.
func (v *Volume) SetOptions(setVOFns ...volumeopt.SetVolumeOptFn) {
	for _, set := range setVOFns {
		set(v.options)
	}
}

// NewVolume creates a new Volume instance with the specified name.
// The Volume instance contains configuration options for creating a Docker volume.
func (c *Client) NewVolume(name string) *Volume {
	return &Volume{options: &volume.CreateOptions{Name: name}}
}
