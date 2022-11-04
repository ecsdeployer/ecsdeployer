package rshell

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDockerLabel_Basics(t *testing.T) {
	require.Equal(t, "cloud87.rshell", LabelName)

	label1 := &DockerLabel{
		Cluster:          "fake",
		Path:             "fake",
		Port:             1234,
		AssignPublicIp:   true,
		SubnetIds:        []string{"subnet-11111"},
		SecurityGroupIds: []string{"sg-11111"},
	}
	var _ string = label1.ToJSON()

}
