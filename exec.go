package main

import (
	"bufio"
	"context"
	"os"
	"strings"

	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/puppetlabs/transparent-containers/cli/logging"
)

func main() {
	containerID := os.Args[1]
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	cmd := strings.Split(string("ls -lisa"), " ")

	execOpts := types.ExecConfig{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: false,
	}
	ctx := context.Background()
	execInstance, _ := cli.ContainerExecCreate(ctx, containerID, execOpts)
	att, _ := cli.ContainerExecAttach(ctx, execInstance.ID, execOpts)
	defer att.Close()
	execStartOpts := types.ExecStartCheck{}
	cli.ContainerExecStart(ctx, execInstance.ID, execStartOpts)
	p := make([]byte, 8)
	att.Reader.Read(p)
	scanner := bufio.NewScanner(att.Reader)
	var data []string
	logging.Stderr("Debug")
	for scanner.Scan() {
		txt := strings.Replace(scanner.Text(), "'", "", 2)
		if txt != "" {
			data = append(data, txt)
		}
	}
	for i, line := range data {
		fmt.Printf("%d: %s\n", i, line)
	}

}
