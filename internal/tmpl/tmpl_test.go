package tmpl

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestTemplate_Functions(t *testing.T) {

	testutil.MockSimpleStsProxy(t)

	ctx := config.New(&config.Project{})
	ctx.ImageTag = "IMAGETAG"
	ctx.Stage = ""

	tpl := New(ctx).WithExtraFields(Fields{
		"FalseThing":   false,
		"TrueThing":    true,
		"EmptyThing":   "",
		"PresentThing": "yep yep yep",
		"NullThing":    nil,
	})

	tables := []struct {
		input    string
		expected string
	}{
		{"{{ prefix .ImageTag 4 }}", "IMAG"},
		{`{{ join "-" "thing" "yar" "blah" 5 }}`, "thing-yar-blah-5"},
		{"{{ AwsRegion }}", ctx.AwsRegion()},
		{"{{ AwsAccountId }}", ctx.AwsAccountId()},
	}

	for _, table := range tables {
		actual, err := tpl.Apply(table.input)
		require.NoError(t, err)

		if actual != table.expected {
			t.Errorf("Mismatch. Expected=%s, received=%s", table.expected, actual)
			break
		}

	}
}

func TestTemplateConditional(t *testing.T) {

	ctx := config.New(&config.Project{})
	ctx.ImageTag = "XimageXtagX"
	ctx.Stage = ""

	tpl := New(ctx).WithExtraFields(Fields{
		"FalseThing":   false,
		"TrueThing":    true,
		"EmptyThing":   "",
		"PresentThing": "yep yep yep",
		"NullThing":    nil,
	})

	tables := []struct {
		input    string
		expected string
	}{
		{"start-{{ if .Stage }}{{ .Stage }}-{{ end }}end", "start-end"},
		{"start-{{ if .ImageTag }}{{ .ImageTag }}-{{ end }}end", "start-XimageXtagX-end"},
		{"start-{{ if .TrueThing }}{{ .ImageTag }}-{{ end }}end", "start-XimageXtagX-end"},
		{"start-{{ if .FalseThing }}{{ .ImageTag }}-{{ end }}end", "start-end"},
		{"start-{{ if .NullThing }}{{ .ImageTag }}-{{ end }}end", "start-end"},
	}

	for _, table := range tables {
		actual, err := tpl.Apply(table.input)

		require.NoError(t, err)

		if actual != table.expected {
			t.Errorf("Mismatch. Expected=%s, received=%s", table.expected, actual)
			break
		}

	}
}

func TestTplFuncJoin(t *testing.T) {

	tables := []struct {
		params   []interface{}
		expected string
	}{
		{[]interface{}{"test", util.Ptr("thing"), true, false, int(5), int32(6), int16(50), int64(123123), float32(1.33), float64(2.67), nil}, "test/thing/true/false/5/6/50/123123/1.33/2.67/"},
		{[]interface{}{int(1), int8(2), int16(3), int32(4), int64(5)}, "1/2/3/4/5"},
		{[]interface{}{uint(1), uint8(2), uint16(3), uint32(4), uint64(5)}, "1/2/3/4/5"},
	}

	for _, table := range tables {
		result := tplFuncJoin("/", table.params...)
		require.Equal(t, table.expected, result)
	}
}

func TestTplFuncPrefix_Valid(t *testing.T) {

	tables := []struct {
		str      string
		length   int
		expected string
	}{
		{"something", 4, "some"},
		{"test", 4, "test"},
	}

	for _, table := range tables {
		result, err := tplFuncPrefix(table.str, table.length)
		require.NoError(t, err)
		require.Equal(t, table.expected, result)
	}
}

func TestTplFuncPrefix_Invalid(t *testing.T) {

	tables := []struct {
		str    string
		length int
	}{
		{"something", 0},
		{"something", -1},
	}

	for _, table := range tables {
		_, err := tplFuncPrefix(table.str, table.length)
		require.Error(t, err)
	}
}
