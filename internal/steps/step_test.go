package steps

import (
	"errors"
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestStep_FindAllChildren(t *testing.T) {
	project, _ := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{})
	step := DeploymentStep(project)

	t.Run("list all children", func(t *testing.T) {
		children := step.FindAllChildren(nil)
		for _, child := range children {
			require.NotNil(t, child)
			require.IsType(t, &Step{}, child)
		}
	})

	t.Run("list specific type", func(t *testing.T) {
		children := step.FindAllChildren(aws.String("TaskDefinition"))

		for _, child := range children {
			require.NotNil(t, child)
			require.IsType(t, &Step{}, child)
			require.Equal(t, "TaskDefinition", child.Label)
		}

	})
}

func TestStep(t *testing.T) {
	t.Run("names", func(t *testing.T) {
		require.Equal(t, "TestStep", (&Step{Label: "TestStep"}).Name())
		require.Equal(t, "TestStep", (&Step{Label: "TestStep"}).Identifier())
		require.Equal(t, "TestStep(Thinger)", (&Step{Label: "TestStep", ID: "Thinger"}).Identifier())
	})

	t.Run("error handling", func(t *testing.T) {

		errCreateErr := errors.New("TestStepError_Create")
		errReadErr := errors.New("TestStepError_Read")
		errUpdateErr := errors.New("TestStepError_Update")
		// errDeleteErr := errors.New("TestStepError_Delete")
		errPreApplyErr := errors.New("TestStepError_PreApply")

		_, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{})

		t.Run("in PreApply", func(t *testing.T) {
			step := NewStep(&Step{
				Label:    "Test",
				PreApply: func(ctx *config.Context, s *Step, sm *StepMetadata) error { return errPreApplyErr },
				Create:   func(ctx *config.Context, s *Step, sm *StepMetadata) (OutputFields, error) { return nil, nil },
			})

			err := step.Apply(ctx)
			require.Error(t, err)
			require.ErrorIs(t, err, errPreApplyErr)
		})

		t.Run("in Read", func(t *testing.T) {
			step := NewStep(&Step{
				Label:  "Test",
				Read:   func(ctx *config.Context, s *Step, sm *StepMetadata) (any, error) { return nil, errReadErr },
				Create: func(ctx *config.Context, s *Step, sm *StepMetadata) (OutputFields, error) { return nil, nil },
			})

			err := step.Apply(ctx)
			require.Error(t, err)
			require.ErrorIs(t, err, errReadErr)
		})

		t.Run("in Update", func(t *testing.T) {
			step := NewStep(&Step{
				Label:  "Test",
				Read:   func(ctx *config.Context, s *Step, sm *StepMetadata) (any, error) { return true, nil },
				Update: func(ctx *config.Context, s *Step, sm *StepMetadata) (OutputFields, error) { return nil, errUpdateErr },
			})

			err := step.Apply(ctx)
			require.Error(t, err)
			require.ErrorIs(t, err, errUpdateErr)
		})

		t.Run("in Create", func(t *testing.T) {
			step := NewStep(&Step{
				Label:  "Test",
				Read:   func(ctx *config.Context, s *Step, sm *StepMetadata) (any, error) { return nil, nil },
				Create: func(ctx *config.Context, s *Step, sm *StepMetadata) (OutputFields, error) { return nil, errCreateErr },
			})

			err := step.Apply(ctx)
			require.Error(t, err)
			require.ErrorIs(t, err, errCreateErr)
		})

		t.Run("in dependent step", func(t *testing.T) {

			failRead := func(ctx *config.Context, s *Step, sm *StepMetadata) (any, error) { return nil, errReadErr }
			noOpRead := func(ctx *config.Context, s *Step, sm *StepMetadata) (any, error) { return nil, nil }
			noOpCreate := func(ctx *config.Context, s *Step, sm *StepMetadata) (OutputFields, error) { return nil, nil }

			t.Run("parallel deps", func(t *testing.T) {
				step := NewStep(&Step{
					Label:        "TestParallel",
					Read:         noOpRead,
					Create:       noOpCreate,
					ParallelDeps: true,
					Dependencies: []*Step{
						NewStep(&Step{Label: "Dep1", Create: noOpCreate}),
						NewStep(&Step{Label: "Dep2", Read: failRead, Create: noOpCreate}),
						NewStep(&Step{Label: "Dep3", Create: noOpCreate}),
					},
				})

				err := step.Apply(ctx)
				require.Error(t, err)
				require.ErrorIs(t, err, errReadErr)
			})

			t.Run("sequential deps", func(t *testing.T) {
				step := NewStep(&Step{
					Label:        "TestSeq",
					Read:         noOpRead,
					Create:       noOpCreate,
					ParallelDeps: false,
					Dependencies: []*Step{
						NewStep(&Step{Label: "Dep1", Create: noOpCreate}),
						NewStep(&Step{Label: "Dep2", Read: failRead, Create: noOpCreate}),
						NewStep(&Step{Label: "Dep3", Create: noOpCreate}),
					},
				})

				err := step.Apply(ctx)
				require.Error(t, err)
				require.ErrorIs(t, err, errReadErr)
			})
		})

	})

	t.Run("Validate", func(t *testing.T) {
		tables := []struct {
			step   *Step
			valid  bool
			errMsg string
			panics bool
		}{
			{
				step:  NoopStep(),
				valid: true,
			},
			{
				step: &Step{
					Label: "",
				},
				errMsg: "include a Label",
				panics: true,
			},
			{
				step: &Step{
					Label: "NoCRUDFuncs",
				},
				errMsg: "must implement",
			},
			{
				step: &Step{
					Label:  "UpdatedWithoutRead",
					Update: func(ctx *config.Context, s *Step, sm *StepMetadata) (OutputFields, error) { return nil, nil },
				},
				errMsg: "Update func without a Read",
			},
			{
				step: &Step{
					Label:  "DeleteWithoutRead",
					Create: func(ctx *config.Context, s *Step, sm *StepMetadata) (OutputFields, error) { return nil, nil },
					Delete: func(ctx *config.Context, s *Step, sm *StepMetadata) (OutputFields, error) { return nil, nil },
				},
				errMsg: "Delete func without a Read",
			},
		}

		for i, table := range tables {
			t.Run(fmt.Sprintf("test_%02d", i+1), func(t *testing.T) {
				err := table.step.Validate()
				if table.valid {
					require.NoError(t, err)
					return
				}

				require.Error(t, err)
				require.ErrorIs(t, err, ErrStepInvalid)

				if table.errMsg != "" {
					require.ErrorContains(t, err, table.errMsg)
				}

				if table.panics {
					require.Panics(t, func() { NewStep(table.step) })
				}
			})
		}
	})
}
