package imageopt

import (
	"encoding/base64"
	"runtime"

	"github.com/docker/docker/api/types"
)

type SetPullOptFn func(options *types.ImagePullOptions)

/*
Download all tagged images in the repository.

Short hand equivalent:
"--all-tags , -a"

	img := client.NewImage("alpine")
	img.PullOptions(
		pulloptions.DownloadAllTaggedImages()
	)
*/
func DownloadAllTaggedImages() SetPullOptFn {
	return func(options *types.ImagePullOptions) {
		options.All = true
	}
}

/*
Sets platform if server is multi-platform capable

Short hand equivalent:
"--platform"

	img := client.NewImage("alpine")
	img.PullOptions(
		pulloptions.Platform("linux/arm64")
	)
*/
func Platform(platform string) SetPullOptFn {
	return func(options *types.ImagePullOptions) {
		options.Platform = platform
	}
}

// Matches the platform to the runtime's arch
func MatchPlatform() SetPullOptFn {
	return func(options *types.ImagePullOptions) {
		options.Platform = runtime.GOARCH
	}
}

/*
Sets the privilege function sed to
authenticate or authorize the pull
operation for images that have restricted access.

	img := client.NewImage("alpine")
	img.PullOptions(
		pulloptions.Privilege(func() (string, error) {
			// Perform your authentication logic here
			// For example, return an authentication token
			return "Bearer <your-auth-token>", nil
		})
	)
*/
func Privilege(pFn func() (string, error)) SetPullOptFn {
	return func(options *types.ImagePullOptions) {
		options.PrivilegeFunc = pFn
	}
}

/*
Sets the base64 encoded credentials for the registry

	img := client.NewImage("alpine")
	img.PullOptions(
		pulloptions.RegistryAuth("user", "pass")
	)
*/
func RegistryAuth(username, password string) SetPullOptFn {
	encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return func(options *types.ImagePullOptions) {
		options.RegistryAuth = encoded
	}
}
