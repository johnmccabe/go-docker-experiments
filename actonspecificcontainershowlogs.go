package main

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	options := types.ContainerLogsOptions{ShowStdout: true, Since: "1481500800", ShowStderr: true}
	out, err := cli.ContainerLogs(ctx, "1b2c8c9cc13e", options)
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, out)
}
