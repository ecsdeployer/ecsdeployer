package yaml_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"github.com/stretchr/testify/require"
)

type somethingYaml struct {
	SomeStr  *string `yaml:"something,omitempty"`
	SomeInt  *int    `yaml:"someint,omitempty"`
	SomeBool *bool   `yaml:"somebool,omitempty"`
}

func TestParseYAMLFile(t *testing.T) {

	t.Run("normal", func(t *testing.T) {
		obj, err := yaml.ParseYAMLFile[somethingYaml]("testdata/something.yml")

		require.NoError(t, err)
		require.IsType(t, &somethingYaml{}, obj)

		require.NotNil(t, obj.SomeStr)
		require.Equal(t, "blah", *obj.SomeStr)

		require.NotNil(t, obj.SomeInt)
		require.Equal(t, 1234, *obj.SomeInt)

		require.Nil(t, obj.SomeBool)

	})

	t.Run("error cases", func(t *testing.T) {
		_, err := yaml.ParseYAMLFile[somethingYaml]("testdata/_badfile.yml")
		require.Error(t, err)
	})
}

func TestParseYAMLString(t *testing.T) {

	t.Run("normal", func(t *testing.T) {
		obj, err := yaml.ParseYAMLString[somethingYaml]("something: foo\nsomeint: 4567")

		require.NoError(t, err)
		require.IsType(t, &somethingYaml{}, obj)

		require.NotNil(t, obj.SomeStr)
		require.Equal(t, "foo", *obj.SomeStr)

		require.NotNil(t, obj.SomeInt)
		require.Equal(t, 4567, *obj.SomeInt)

		require.Nil(t, obj.SomeBool)

	})

	t.Run("error cases", func(t *testing.T) {
		_, err := yaml.ParseYAMLString[somethingYaml]("1234")
		require.Error(t, err)
	})
}
