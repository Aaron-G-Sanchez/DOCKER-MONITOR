package main

import (
	"context"
	"fmt"

	"github.com/moby/moby/client"
)

func main() {

	// Sample code from the Docker SDK docs
	ctx := context.Background()
	apiClient, err := client.New(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	images, err := apiClient.ImageList(ctx, client.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Images: ")

	for _, image := range images.Items {
		fmt.Println(image.RepoTags)
	}

}
