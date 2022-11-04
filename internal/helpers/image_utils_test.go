package helpers

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

func TestResolveImageUri(t *testing.T) {

	closeMock := testutil.MockSimpleStsProxy(t)
	defer closeMock()

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
		{config.ImageUri{Ecr: util.Ptr("test/thing"), Tag: util.Ptr("blah")}, "555555555555.dkr.ecr.us-east-1.amazonaws.com/test/thing:blah"},
		{config.ImageUri{Ecr: util.Ptr("test/thing"), Digest: util.Ptr("sha256:deadbeef")}, "555555555555.dkr.ecr.us-east-1.amazonaws.com/test/thing@sha256:deadbeef"},
		{config.ImageUri{Docker: util.Ptr("user/reponame"), Tag: util.Ptr("sometag")}, "user/reponame:sometag"},
		{config.ImageUri{Docker: util.Ptr("user/reponame"), Digest: util.Ptr("sha256:deadbeef")}, "user/reponame@sha256:deadbeef"},
		{alreadyResolved, "fakefake"},
	}

	for _, table := range tables {

		if err := table.obj.Validate(); err != nil {
			// none of the test cases should ever be invalid
			t.Fatal("YOU BLEW IT!")
		}

		uri, err := ResolveImageUri(ctx, &table.obj)
		if err != nil {
			t.Errorf("Error when resolving img: <%s> got err: %v", table.obj.Value(), err)
			continue
		}

		if uri != table.expected {
			t.Errorf("Got incorrect image uri: expected=<%v> but got=<%v>", table.expected, uri)
		}
	}
}
