package tmpl

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

func TestTemplate_Functions(t *testing.T) {

	closeMock := testutil.MockSimpleStsProxy(t)
	defer closeMock()

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

		if err != nil {
			t.Errorf("unexpected error: %s", err)
			break
		}

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

		if err != nil {
			t.Errorf("unexpected error: %s", err)
			break
		}

		if actual != table.expected {
			t.Errorf("Mismatch. Expected=%s, received=%s", table.expected, actual)
			break
		}

	}
}
