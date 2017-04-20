package main

import (
	"bytes"
	"context"
	"os"

	"fmt"

	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	containerID := os.Args[1]
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	cmd := []string{"rpm", "-qa", "--queryformat", "%{NAME},%{VERSION}-%{RELEASE}-%{ARCH}\n"}

	execOpts := types.ExecConfig{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}
	ctx := context.Background()
	execInstance, _ := cli.ContainerExecCreate(ctx, containerID, execOpts)
	att, _ := cli.ContainerExecAttach(ctx, execInstance.ID, execOpts)
	defer att.Close()
	execStartOpts := types.ExecStartCheck{}
	cli.ContainerExecStart(ctx, execInstance.ID, execStartOpts)

	var data []string

	allbytes, _ := ioutil.ReadAll(att.Reader)
	fmt.Printf("All bytes: %v", allbytes)
	var bytedata [][]byte
	bytedata = bytes.Split(allbytes, []byte("\n"))

	// Looks like the stream can sometimes return two lines in a single payload
	// [1 0 0 0 0 0 0 33 108 105 98 100 98 45 117 116 105 108 115 44 53 46 51 46 50 49 45 49 57 46 101 108 55 45 120 56 54 95 54 52]
	// [1 0 0 0 0 0 0 61 110 99 117 114 115 101 115 44 53 46 57 45 49 51 46 50 48 49 51 48 53 49 49 46 101 108 55 45 120 56 54 95 54 52]
	// [108 105 98 116 97 115 110 49 44 51 46 56 45 51 46 101 108 55 45 120 56 54 95 54 52]
	// [1 0 0 0 0 0 0 39 99 97 45 99 101 114 116 105 102 105 99 97 116 101 115 44 50 48 49 53 46 50 46 54 45 55 51 46 101 108 55 45 110 111 97 114 99 104]

	for _, b := range bytedata {
		fmt.Printf("%v\n", b)
	}

}
