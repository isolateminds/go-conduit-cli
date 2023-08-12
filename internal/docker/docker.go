package docker

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	. "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Client struct {
	client         *client.Client
	imageResWriter io.Writer
	statsResWriter io.Writer
}

func (c *Client) String() string {
	return c.client.DaemonHost()
}

func (c *Client) CreateNetwork(ctx context.Context, network *Network) error {
	res, err := c.client.NetworkCreate(ctx, network.Name, *network.options)
	if err != nil {
		return err
	}
	network.Id = res.ID
	return nil
}

func (c *Client) CreateVolume(ctx context.Context, volume *Volume) error {
	_, err := c.client.VolumeCreate(ctx, *volume.options)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetContainerStats(ctx context.Context, container *Container) error {

	res, err := c.client.ContainerStats(ctx, container.Name, true)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if _, err := io.Copy(c.statsResWriter, res.Body); err != nil {
		return err
	}
	return nil
}
func (c *Client) RemoveContainer(ctx context.Context, container *Container, force bool) error {
	return c.client.ContainerRemove(ctx, container.Name, types.ContainerRemoveOptions{
		RemoveLinks:   force,
		RemoveVolumes: force,
		Force:         force,
	})
}
func (c *Client) UnpauseContainer(ctx context.Context, container *Container) error {
	return c.client.ContainerUnpause(ctx, container.Name)
}
func (c *Client) PauseContainer(ctx context.Context, container *Container) error {
	return c.client.ContainerPause(ctx, container.Name)
}
func (c *Client) RestartContainer(ctx context.Context, container *Container) error {
	return c.client.ContainerRestart(ctx, container.Name, StopOptions{})
}

func (c *Client) StopContainer(ctx context.Context, container *Container) error {
	return c.client.ContainerStop(ctx, container.Name, StopOptions{})
}
func (c *Client) StartContainer(ctx context.Context, container *Container) error {
	return c.client.ContainerStart(ctx, container.Name, types.ContainerStartOptions{})
}
func (c *Client) CreateContainer(ctx context.Context, container *Container) error {
	res, err := c.client.ContainerCreate(
		ctx,
		container.options,
		container.hostOptions,
		container.networkingOptions,
		container.platformOptions,
		container.Name,
	)
	if err != nil {
		return err
	}

	container.Id = res.ID

	return nil
}

func (c *Client) BuildImage(ctx context.Context, image *Image) error {
	if image.buildOptions.Context == nil {
		return errors.New("BuildImageError: no context was supplied use image.NewImageFromSrc(dir) or supply the context manually.")
	}
	res, err := c.client.ImageBuild(ctx, image.buildOptions.Context, *image.buildOptions)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if _, err = io.Copy(c.imageResWriter, res.Body); err != nil {
		return err
	}
	return nil
}
func (c *Client) PullImage(ctx context.Context, image *Image) error {
	rc, err := c.client.ImagePull(ctx, image.ref, *image.pullOptions)
	if err != nil {
		return err
	}
	defer rc.Close()
	if _, err = io.Copy(c.imageResWriter, rc); err != nil {
		return err
	}
	return nil
}

// This sets the image response writer for Docker's API.
// If this is not set, the client wrapper will default to stdout.
func (c *Client) SetImageResponeWriter(dst io.Writer) {
	c.imageResWriter = dst
}

// This sets the stats response writer for Docker's API.
// If this is not set, the client wrapper will default to stdout.
func (c *Client) SetStatsResponeWriter(dst io.Writer) {
	c.statsResWriter = dst
}

func NewClient(ctx context.Context) (*Client, error) {
	client, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}
	ok, err := isDaemonRunning(ctx, client)
	if ok {
		return &Client{
			client:         client,
			imageResWriter: os.Stdout,
			statsResWriter: os.Stdout,
		}, nil
	} else {
		return nil, err
	}
}

// checks if the docker daemon is running by pinging it
func isDaemonRunning(ctx context.Context, client *client.Client) (bool, error) {
	if _, err := client.Ping(ctx); err != nil {
		return false, err
	}
	return true, nil
}
