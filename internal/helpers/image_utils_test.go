package helpers

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestResolveImageUri(t *testing.T) {

	testutil.MockSimpleStsProxy(t)

	project := &config.Project{
		ProjectName: "dummy",
	}
	ctx := config.New(project)

	alreadyResolved := config.NewImageUri("fakefake")
	alreadyResolved.SetResolved(alreadyResolved.Value())

	tables := []struct {
		obj      config.ImageUri
		expected string
	}{
		{config.NewImageUri("fakefake"), "fakefake"},
		{config.ImageUri{Ecr: new("test/thing"), Tag: new("blah")}, "555555555555.dkr.ecr.us-east-1.amazonaws.com/test/thing:blah"},
		{config.ImageUri{Ecr: new("test/thing"), Digest: new("sha256:deadbeef")}, "555555555555.dkr.ecr.us-east-1.amazonaws.com/test/thing@sha256:deadbeef"},
		{config.ImageUri{Docker: new("user/reponame"), Tag: new("sometag")}, "user/reponame:sometag"},
		{config.ImageUri{Docker: new("user/reponame"), Digest: new("sha256:deadbeef")}, "user/reponame@sha256:deadbeef"},
		{alreadyResolved, "fakefake"},
	}

	for _, table := range tables {

		if err := table.obj.Validate(); err != nil {
			// none of the test cases should ever be invalid
			require.NoErrorf(t, err, "This should not occur unless you wrote a bad test case")
		}

		uri, err := ResolveImageUri(ctx, &table.obj)
		require.NoError(t, err)
		require.Equal(t, table.expected, uri)
	}
}
