package crawler

type statistics struct {
	missingLocation      map[string]struct{}
	missingOrigin        map[string]struct{}
	missingDocumentation map[string]struct{}
}

type statisticsReporter func(path string)

func (s *statistics) Report(path string, reporters ...statisticsReporter) {
	for _, reporter := range reporters {
		reporter(path)
	}
}

func (s *statistics) MissingLocation() statisticsReporter {
	return func(path string) {
		s.missingLocation[path] = struct{}{}
	}
}

func (s *statistics) MissingOrigin() statisticsReporter {
	return func(path string) {
		s.missingOrigin[path] = struct{}{}
	}
}

func (s *statistics) MissingDocumentation() statisticsReporter {
	return func(path string) {
		s.missingDocumentation[path] = struct{}{}
	}
}

func NewStatistics() *statistics {
	return &statistics{
		missingLocation:      map[string]struct{}{},
		missingOrigin:        map[string]struct{}{},
		missingDocumentation: map[string]struct{}{},
	}
}
