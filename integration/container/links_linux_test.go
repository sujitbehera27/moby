package container // import "github.com/docker/docker/integration/container"

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/integration/internal/container"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
	"gotest.tools/skip"
)

func TestLinksEtcHostsContentMatch(t *testing.T) {
	skip.If(t, testEnv.IsRemoteDaemon)

	hosts, err := ioutil.ReadFile("/etc/hosts")
	skip.If(t, os.IsNotExist(err))

	defer setupTest(t)()
	client := testEnv.APIClient()
	ctx := context.Background()

	cID := container.Run(t, ctx, client, container.WithNetworkMode("host"))
	res, err := container.Exec(ctx, client, cID, []string{"cat", "/etc/hosts"})
	assert.NilError(t, err)
	assert.Assert(t, is.Len(res.Stderr(), 0))
	assert.Equal(t, 0, res.ExitCode)

	assert.Check(t, is.Equal(string(hosts), res.Stdout()))
}

func TestLinksContainerNames(t *testing.T) {
	skip.If(t, testEnv.DaemonInfo.OSType == "windows")

	defer setupTest(t)()
	client := testEnv.APIClient()
	ctx := context.Background()

	containerA := "first_" + t.Name()
	containerB := "second_" + t.Name()
	container.Run(t, ctx, client, container.WithName(containerA))
	container.Run(t, ctx, client, container.WithName(containerB), container.WithLinks(containerA+":"+containerA))

	f := filters.NewArgs(filters.Arg("name", containerA))

	containers, err := client.ContainerList(ctx, types.ContainerListOptions{
		Filters: f,
	})
	assert.NilError(t, err)
	assert.Check(t, is.Equal(1, len(containers)))
	assert.Check(t, is.DeepEqual([]string{"/" + containerA, "/" + containerB + "/" + containerA}, containers[0].Names))
}
