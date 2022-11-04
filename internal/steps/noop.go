package steps

// Step that does no work at all. It's used to skip other steps (returned in their constructor)
func NoopStep() *Step {
	return NewStep(&Step{
		Label: "Noop",
	})
}
