package steps

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestMarkerTag(t *testing.T) {
	_, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{})

	markerTag, err := stepCleanupMarkerTag(ctx)

	require.NoError(t, err)

	require.Equal(t, "ecsdeployer/project", markerTag.key)
	require.Equal(t, "dummy", markerTag.value)

}
