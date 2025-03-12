package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

// CreateNetworkHandler creates a Docker network.
func CreateNetworkHandler(ctx context.Context, parameters map[string]interface{}) (interface{}, error) {
	name, _ := parameters["name"].(string)
	if name == "" {
		return nil, fmt.Errorf("missing network name")
	}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	_, err = cli.NetworkCreate(ctx, name, network.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("error creating network: %v", err)
	}
	return fmt.Sprintf("Network '%s' created successfully", name), nil
}

// PullImageHandler pulls a Docker image.
func PullImageHandler(ctx context.Context, parameters map[string]interface{}) (interface{}, error) {
	imgName, ok := parameters["image"].(string)
	if !ok || imgName == "" {
		name, nameOk := parameters["name"].(string)
		tag, tagOk := parameters["tag"].(string)
		if !nameOk || name == "" {
			return nil, fmt.Errorf("missing image name for pull_image")
		}
		if !tagOk || tag == "" {
			tag = "latest"
		}
		imgName = fmt.Sprintf("%s:%s", name, tag)
	}
	return pullImage(ctx, imgName)
}

func pullImage(ctx context.Context, imgName string) (interface{}, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	pullCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	out, err := cli.ImagePull(pullCtx, imgName, image.PullOptions{})
	if err != nil {
		return nil, err
	}
	defer out.Close()
	_, err = io.Copy(io.Discard, out)
	if err != nil {
		return nil, fmt.Errorf("error pulling image: %v", err)
	}
	return fmt.Sprintf("Image '%s' pulled successfully", imgName), nil
}

// Handler is the exported symbol that Yaegi will look for.
func Handler(ctx context.Context, parameters map[string]interface{}) (interface{}, error) {
	action, ok := parameters["action"].(string)
	if !ok {
		return nil, fmt.Errorf("missing action parameter")
	}
	switch action {
	case "create_network":
		return CreateNetworkHandler(ctx, parameters)
	case "pull_image":
		return PullImageHandler(ctx, parameters)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}
