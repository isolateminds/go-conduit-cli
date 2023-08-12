package platformopt

import (
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// SetPlatformOptions is a type representing a function to set platform options for an image.
type SetPlatformOptions func(option *v1.Platform)

// Arch sets the CPU architecture for the image platform.
// Use this function to specify the architecture of the target platform, such as 'amd64' or 'ppc64'.
func Arch(arch string) SetPlatformOptions {
	return func(option *v1.Platform) {
		option.Architecture = arch
	}
}

// OS sets the operating system for the image platform.
// Use this function to define the operating system of the target platform, such as 'linux' or 'windows'.
func OS(os string) SetPlatformOptions {
	return func(option *v1.Platform) {
		option.OS = os
	}
}

// OSVersion sets the version of the operating system for the image platform.
// Use this function to specify the version of the operating system, for example, '10.0.14393.1066' on Windows.
func OSVersion(osV string) SetPlatformOptions {
	return func(option *v1.Platform) {
		option.OSVersion = osV
	}
}

// OSFeatures adds required operating system features to the image platform.
// Use this function to provide an array of strings representing required OS features,
// such as 'win32k' on Windows.
func OSFeatures(features ...string) SetPlatformOptions {
	return func(option *v1.Platform) {
		if option.OSFeatures == nil {
			option.OSFeatures = []string{}
		}
		option.OSFeatures = append(option.OSFeatures, features...)
	}
}

// Variant sets the CPU variant for the image platform.
// Use this function to specify a variant of the CPU architecture, such as 'v7' for ARMv7 when architecture is 'arm'.
func Variant(v string) SetPlatformOptions {
	return func(option *v1.Platform) {
		option.Variant = v
	}
}
