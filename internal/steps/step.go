package steps

import (
	"errors"
	"fmt"
	"sync"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/caarlos0/log"
	"golang.org/x/sync/errgroup"
)

// type ParallelDeps []Step

var (
	ErrFuncNotImplemented    = errors.New("Function not implemented")
	ErrStepFailed            = errors.New("step failed")
	ErrStepDependencyFailure = errors.New("step dependency failure")
	ErrStepInvalid           = errors.New("step invalid")
)

const useLogIndentation = false

type OutputFields map[string]interface{}

// Passed to the CRUD functions
type StepMetadata struct {
	// NOT USED YET. PLACEHOLDER FOR FUTURE
	// ParentStep *Step

	// TmpCtx cp
}

type Step struct {
	Label string
	ID    string

	Parent *Step

	Create CreateFunc
	Read   ReadFunc
	Update UpdateFunc

	// NOT IMPLEMENTED YET
	Delete DeleteFunc

	// This runs immediately before Apply and can be used to add dependencies
	PreApply PreApplyFunc

	// For Building a graph
	// GraphBuilder GraphBuilderFunc

	ParallelDeps bool

	// this step does no "work", it only has descendants
	IsNoOp bool

	Attributes map[string]interface{}

	Outputs map[interface{}]OutputFields

	OutputFields OutputFields

	Dependencies []*Step

	Resource interface{}

	isExisting       *bool
	ExistingResource interface{}

	inited  bool
	applied bool

	mu sync.RWMutex

	Logger *log.Entry
}

type CreateFunc func(*config.Context, *Step, *StepMetadata) (OutputFields, error)
type UpdateFunc func(*config.Context, *Step, *StepMetadata) (OutputFields, error)
type DeleteFunc func(*config.Context, *Step, *StepMetadata) (OutputFields, error)
type ReadFunc func(*config.Context, *Step, *StepMetadata) (any, error)
type DifferFunc func(*config.Context, *Step, *StepMetadata) (any, error)
type PreApplyFunc func(*config.Context, *Step, *StepMetadata) error

func (s *Step) init() {
	if s.inited {
		return
	}

	// if s.mu == nil {
	// 	s.mu = &sync.Mutex{}
	// }

	if s.Logger == nil {
		s.Logger = s.setupLogger()
	}

	if s.OutputFields == nil {
		s.OutputFields = make(OutputFields)
	}

	if s.Attributes == nil {
		s.Attributes = make(map[string]interface{})
	}

	if s.Outputs == nil {
		s.Outputs = make(map[interface{}]OutputFields)
	}

	if s.Dependencies != nil && len(s.Dependencies) > 0 {
		for i := range s.Dependencies {
			s.Dependencies[i].Parent = s
		}
	}

	s.inited = true
}

func (s *Step) SetAttr(key string, value interface{}) {
	s.Attributes[key] = value
}

func (s *Step) GetAttr(key string) (interface{}, bool) {
	val, ok := s.Attributes[key]
	return val, ok
}

func (s *Step) GetAttrMust(key string) interface{} {
	return s.Attributes[key]
}

// looks in this step and any children
func (s *Step) SearchAttr(key string) (interface{}, bool) {
	val, ok := s.Attributes[key]
	if ok {
		return val, true
	}

	for _, dep := range s.Dependencies {
		val, ok := dep.SearchAttr(key)
		if ok {
			return val, true
		}
	}

	return nil, false
}

func (s *Step) LookupOutput(key string) (interface{}, bool) {
	val, ok := s.OutputFields[key]
	if ok {
		return val, true
	}

	for _, dep := range s.Dependencies {
		val, ok := dep.LookupOutput(key)
		if ok {
			return val, true
		}
	}

	return nil, false
}

func (s *Step) Applied() bool {
	return s.applied
}

func (s *Step) MarkApplied() {
	s.applied = true
}

func (s *Step) Indent() {
	if !useLogIndentation {
		return
	}
	s.Logger.IncreasePadding()
}

func (s *Step) Outdent() {
	if !useLogIndentation {
		return
	}
	s.Logger.DecreasePadding()
}

func (s *Step) Apply(ctx *config.Context) error {

	if s.Applied() {
		// s.Logger.Debug("Already Applied! Skipping")
		return nil
	}

	if s.PreApply != nil {
		s.Logger.Debugf("Execute:PreApply")
		if err := s.PreApply(ctx, s, s.buildMetadata()); err != nil {
			// s.Logger.WithError(err).WithField("phase", "PreApply").Error("Failure")
			s.Logger.WithField("phase", "PreApply").Error("Failure")
			return err
		}
	}

	// s.Logger.Debug("Applying step")

	if len(s.Dependencies) > 0 {
		s.Logger.Debugf("Applying Dependencies")
		s.Indent()

		if s.ParallelDeps {

			eg := new(errgroup.Group)
			eg.SetLimit(5)

			for _, dep := range s.Dependencies {
				dep := dep
				eg.Go(func() error {
					return dep.Apply(ctx)
				})
			}

			if err := eg.Wait(); err != nil {
				s.Outdent()
				// s.Logger.WithError(err).WithField("phase", "Apply").Error("Failure")
				s.Logger.WithField("phase", "Apply").Error("Failure")
				return err
			}
		} else {
			// sequential
			for _, dep := range s.Dependencies {
				if err := dep.Apply(ctx); err != nil {
					s.Outdent()
					// s.Logger.WithError(err).WithField("phase", "Apply").Error("Failure")
					s.Logger.WithField("phase", "Apply").Error("Failure")
					return err
				}
			}
		}

		s.Outdent()
	}

	isExist, err := s.exists(ctx)
	if err != nil {
		return err
	}

	if isExist {
		if err := s.update(ctx); err != nil {
			return err
		}
		return nil
	}

	if err := s.create(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Step) updateOutputFields(fields OutputFields) {
	if fields == nil {
		return
	}

	for k, v := range fields {
		s.OutputFields[k] = v
	}
}

func (s *Step) exists(ctx *config.Context) (bool, error) {

	if s.isExisting != nil {
		return *s.isExisting, nil
	}

	if err := s.read(ctx); err != nil {
		return false, err
	}

	return *s.isExisting, nil
}

func (s *Step) read(ctx *config.Context) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	s.Logger.Debug("Execute:Read")

	existing, err := s.Read(ctx, s, s.buildMetadata())
	if err != nil {
		// s.Logger.WithError(err).WithField("phase", "Read").Error("Failure")
		s.Logger.WithField("phase", "Read").Error("Failure")
		return err
	}

	if existing == nil {
		s.isExisting = aws.Bool(false)
		s.ExistingResource = nil
		return nil
	}

	s.isExisting = aws.Bool(true)
	s.ExistingResource = existing
	return nil
}

func (s *Step) create(ctx *config.Context) error {
	if s.Create == nil {
		// CANT CREATE
		return nil
	}
	s.Logger.Debug("Execute:Create")

	outputs, err := s.Create(ctx, s, s.buildMetadata())
	if err != nil {
		s.Logger.WithField("phase", "Create").Error("Failure")
		// s.Logger.WithError(err).WithField("phase", "Create").Error("Failure")
		return err
	}

	s.updateOutputFields(outputs)

	// s.Logger.WithField("outputs", outputs).Debug("Resource created")

	return nil
}

func (s *Step) update(ctx *config.Context) error {
	if s.Update == nil {
		// CANT UPDATE
		return nil
	}
	s.Logger.Debug("Execute:Update")

	outputs, err := s.Update(ctx, s, s.buildMetadata())
	if err != nil {
		s.Logger.WithError(err).WithField("phase", "Update").Error("Failure")
		return err
	}

	s.updateOutputFields(outputs)

	// s.Logger.WithField("outputs", outputs).Debug("Resource updated")

	return nil
}

func (s *Step) Name() string {
	return s.Label
}

func (s *Step) Identifier() string {
	if s.ID != "" {
		return fmt.Sprintf("%s(%s)", s.Label, s.ID)
	}
	return s.Label
}

func (s *Step) Log() *log.Entry {
	return s.setupLogger()
}

func (s *Step) setupLogger() *log.Entry {
	entry := log.WithField("step", s.Label)

	if s.ID != "" {
		entry = entry.WithField("id", s.ID)
	}

	return entry
}

func (s *Step) buildMetadata() *StepMetadata {
	return &StepMetadata{}
}

func NewStep(step *Step) *Step {
	step.init()

	if step.Read == nil && step.Create == nil && step.Update == nil && step.Delete == nil {
		step.IsNoOp = true
	}

	// they didnt give a Read, so assume it never exists
	if step.Read == nil {
		step.Read = func(_ *config.Context, _ *Step, _ *StepMetadata) (any, error) {
			return nil, nil
		}
	}

	if err := step.Validate(); err != nil {
		// programmer error
		panic(err)
	}

	return step
}

func (s *Step) Validate() error {

	if s.Label == "" {
		return fmt.Errorf("%w: Please include a Label", ErrStepInvalid)
	}

	if !s.IsNoOp {
		if s.Create == nil && s.Update == nil {
			return fmt.Errorf("%w: You must implement Update or Create", ErrStepInvalid)
		}

		if s.Update != nil && s.Read == nil {
			return fmt.Errorf("%w: You can't have an Update func without a Read func", ErrStepInvalid)
		}

		if s.Delete != nil && s.Read == nil {
			return fmt.Errorf("%w: You can't have a Delete func without a Read func", ErrStepInvalid)
		}

	}

	return nil
}

// find all children that have a specific label. Or all children if the label is nil
func (s *Step) FindAllChildren(label *string) []*Step {

	if len(s.Dependencies) == 0 {
		// return []*Step{}
		return nil
	}

	children := []*Step{}
	for i := range s.Dependencies {

		if label != nil && s.Dependencies[i].Label == *label {
			children = append(children, s.Dependencies[i])
		}

		depChilds := s.Dependencies[i].FindAllChildren(label)
		if len(depChilds) > 0 {
			children = append(children, depChilds...)
		}
	}

	return children
}
