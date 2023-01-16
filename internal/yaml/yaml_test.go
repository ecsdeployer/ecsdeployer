package yaml_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalStrict(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		val := []byte("something: yar")
		res := somethingYaml{}
		err := yaml.UnmarshalStrict(val, &res)
		require.NoError(t, err)
		require.Equal(t, "yar", *res.SomeStr)
	})

	t.Run("bad", func(t *testing.T) {
		val := []byte("bloop: yar")
		res := somethingYaml{}
		err := yaml.UnmarshalStrict(val, &res)
		require.Error(t, err)
	})
}

func TestUnmarshal(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		val := []byte("something: yar\nblah: foo")
		res := somethingYaml{}
		err := yaml.Unmarshal(val, &res)
		require.NoError(t, err)
		require.Equal(t, "yar", *res.SomeStr)
	})

	t.Run("bad", func(t *testing.T) {
		val := []byte("bloop: yar")
		res := somethingYaml{}
		err := yaml.Unmarshal(val, &res)
		require.NoError(t, err)
	})
}

func TestMarshal(t *testing.T) {
	someStr := "blah"
	val := somethingYaml{
		SomeStr: &someStr,
	}

	res, err := yaml.Marshal(val)
	require.NoError(t, err)
	require.Equal(t, "something: blah\n", string(res))

}
