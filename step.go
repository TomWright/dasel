package dasel

// Step is a single step in the query.
// Each function call has its own step.
// Each value in the output is simply a pointer to the actual data point in the context data.
type Step struct {
	context  *Context
	selector Selector
	index    int
	output   Values
}

func (s *Step) Selector() Selector {
	return s.selector
}

func (s *Step) Index() int {
	return s.index
}

func (s *Step) Output() Values {
	return s.output
}

func (s *Step) execute() error {
	f, err := s.context.functions.Get(s.selector.funcName)
	if err != nil {
		return err
	}
	output, err := f(s.context, s, s.selector.funcArgs)
	s.output = output
	return err
}

func (s *Step) inputs() Values {
	prevStep := s.context.Step(s.index - 1)
	if prevStep == nil {
		return Values{}
	}
	return prevStep.output
}
