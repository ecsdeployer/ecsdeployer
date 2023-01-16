package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestImageUri_ParsingValid(t *testing.T) {
	// only test strings not structured imagetags
	tables := []struct {
		str string
	}{
		{"blah"},
	}

	for _, table := range tables {
		imgUri := config.NewImageUri(table.str)

		err := imgUri.Validate()
		require.NoError(t, err)

	}
}

func TestImageUri_ParsingFailures(t *testing.T) {

}

func TestImageUri_Validate(t *testing.T) {

	tables := []struct {
		obj   config.ImageUri
		valid bool
	}{
		{config.NewImageUri("fakefake"), true},
		{config.ImageUri{Ecr: util.Ptr("test/thing"), Tag: util.Ptr("blah")}, true},
		{config.ImageUri{Ecr: util.Ptr("test/thing"), Digest: util.Ptr("sha256:deadbeef")}, true},
		{config.ImageUri{Docker: util.Ptr("user/reponame"), Tag: util.Ptr("sometag")}, true},
		{config.ImageUri{Docker: util.Ptr("user/reponame"), Digest: util.Ptr("sha256:deadbeef")}, true},
		{config.ImageUri{Docker: util.Ptr("user/reponame")}, false},
		{config.ImageUri{Ecr: util.Ptr("user/reponame")}, false},
		{config.ImageUri{Docker: util.Ptr("user/reponame"), Ecr: util.Ptr("user/reponame")}, false},
		{config.ImageUri{Tag: util.Ptr("xxx")}, false},
	}

	for i, table := range tables {
		err := table.obj.Validate()
		require.Equalf(t, table.valid, (err == nil), "Row %d: %s", i, err)
	}
}

func TestImageUri_Value(t *testing.T) {

	tables := []struct {
		obj      config.ImageUri
		expected string
	}{
		{config.NewImageUri("fakefake"), "fakefake"},
		{config.ImageUri{Ecr: util.Ptr("test/thing"), Tag: util.Ptr("blah")}, "{{ AwsAccountId }}.dkr.ecr.{{ AwsRegion }}.amazonaws.com/test/thing:blah"},
		{config.ImageUri{Ecr: util.Ptr("test/thing"), Digest: util.Ptr("sha256:deadbeef")}, "{{ AwsAccountId }}.dkr.ecr.{{ AwsRegion }}.amazonaws.com/test/thing@sha256:deadbeef"},
		{config.ImageUri{Docker: util.Ptr("user/reponame"), Tag: util.Ptr("sometag")}, "user/reponame:sometag"},
		{config.ImageUri{Docker: util.Ptr("user/reponame"), Digest: util.Ptr("sha256:deadbeef")}, "user/reponame@sha256:deadbeef"},
	}

	for i, table := range tables {

		require.NoErrorf(t, table.obj.Validate(), "BAD TEST CASE. ALL TEST CASES SHOULD BE VALID IMAGEURI OBJECTS")

		imgValue := table.obj.Value()
		require.Equalf(t, table.expected, imgValue, "Entry %d", i)

	}
}
