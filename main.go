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
func CreateNetworkHandler(ctx context.Context, parameters map[string]interface{}) error {
	name, _ := parameters["name"].(string)
	if name == "" {
		return fmt.Errorf("missing network name")
	}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	_, err = cli.NetworkCreate(ctx, name, network.CreateOptions{})
	if err != nil {
		fmt.Printf("error creating Network :%v", err)
	}
	return nil
}

// PullImageHandler pulls a Docker image.
func PullImageHandler(ctx context.Context, parameters map[string]interface{}) error {
	imgName, ok := parameters["image"].(string)
	if !ok || imgName == "" {
		name, nameOk := parameters["name"].(string)
		tag, tagOk := parameters["tag"].(string)
		if !nameOk || name == "" {
			return fmt.Errorf("missing image name for pull_image")
		}
		if !tagOk || tag == "" {
			tag = "latest"
		}
		imgName = fmt.Sprintf("%s:%s", name, tag)
	}
	return pullImage(ctx, imgName)
}

func pullImage(ctx context.Context, imgName string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	pullCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	out, err := cli.ImagePull(pullCtx, imgName, image.PullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(io.Discard, out)
	if err != nil {
		fmt.Printf("error pulling image: %v", err)
	}
	return nil
}

// Handler is the exported symbol that Yaegi will look for.
// It selects the proper operation based on the provided "action" parameter.
func Handler(ctx context.Context, parameters map[string]interface{}) error {
	action, ok := parameters["action"].(string)
	if !ok {
		return fmt.Errorf("missing action parameter")
	}
	switch action {
	case "create_network":
		return CreateNetworkHandler(ctx, parameters)
	case "pull_image":
		return PullImageHandler(ctx, parameters)
	default:
		return fmt.Errorf("unknown action: %s", action)
	}
}
