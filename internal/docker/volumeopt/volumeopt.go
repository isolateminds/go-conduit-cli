package volumeopt

import (
	"github.com/docker/docker/api/types/volume"
)

// SetVolumeOptFn is a function type that configures options for creating a Docker volume.
type SetVolumeOptFn func(options *volume.CreateOptions)

// Driver sets the driver to use for creating the Docker volume.
func Driver(driver string) SetVolumeOptFn {
	return func(options *volume.CreateOptions) {
		options.Driver = driver
	}
}

// DriverOptions sets driver-specific options for the Docker volume.
// Use this function to provide additional parameters that are specific to the chosen driver.
func DriverOptions(key, value string) SetVolumeOptFn {
	return func(options *volume.CreateOptions) {
		if options.DriverOpts == nil {
			options.DriverOpts = map[string]string{}
		}
		options.DriverOpts[key] = value
	}
}

// Name sets the name of the Docker volume.
// Use this function to assign a custom name to the volume during creation.
func Name(name string) SetVolumeOptFn {
	return func(options *volume.CreateOptions) {
		options.Name = name
	}
}

// Label adds a label to the Docker volume with the specified key-value pair.
// Labels provide a way to attach metadata to volumes for categorization or organization.
func Label(key, value string) SetVolumeOptFn {
	return func(options *volume.CreateOptions) {
		if options.Labels == nil {
			options.Labels = map[string]string{}
		}
		options.Labels[key] = value
	}
}
