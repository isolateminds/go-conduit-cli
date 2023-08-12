package imageopt

import (
	"io"

	"github.com/docker/docker/api/types"
)

type SetBuildOptFn func(options *types.ImageBuildOptions)

/*
Specify to remove intermediate containers after a successful build

	img := client.NewImage("node:latest")
	img.BuildOptions(
		Remove(),
	)
*/
func Remove() SetBuildOptFn {
	return func(options *types.ImageBuildOptions) {
		options.Remove = true
	}
}

/*
Specify not to remove intermediate containers after a successful build

	img := client.NewImage("node:latest")
	img.BuildOptions(
		DontRemove(),
	)
*/
func DontRemove() SetBuildOptFn {
	return func(options *types.ImageBuildOptions) {
		options.Remove = false
	}
}

/*
Adds a tag to the image reference

	img := client.NewImage("node:latest")
	img.BuildOptions(
		Tag("latest"),
	)
*/
func Tag(tag string) SetBuildOptFn {
	return func(options *types.ImageBuildOptions) {
		if options.Tags == nil {
			options.Tags = make([]string, 0)
		}
		options.Tags = append(options.Tags, tag)
	}
}

/*
Sets the Dockerfile location for the image

	img := client.NewImage("node:latest")
	img.BuildOptions(
		DisableCache()
	)
*/
func DisableCache() SetBuildOptFn {
	return func(options *types.ImageBuildOptions) {
		options.NoCache = true
	}
}

/*
Sets the Dockerfile location for the image

	img := client.NewImage("node:latest")
	img.BuildOptions(
		Dockerfile("/app/Dockerfile")
	)
*/
func Dockerfile(path string) SetBuildOptFn {
	return func(options *types.ImageBuildOptions) {
		options.Dockerfile = path
	}
}

/*
Sets the build context for the image

	img := client.NewImage("node:latest")
	img.BuildOptions(
		BuildContext(bytes.NewReader(tarFile)),
	)
*/
func BuildContext(context io.Reader) SetBuildOptFn {
	return func(options *types.ImageBuildOptions) {
		options.Context = context
	}
}
