package crawler

import (
	"github.com/Tolomeo/anydev.nvim/internal/utils/set"
)

type statistics struct {
	NoLocation      set.StringSet `json:"missingLocation"`
	NoOrigin        set.StringSet `json:"missingOrigin"`
	NoDocumentation set.StringSet `json:"missingDocumentation"`
}

type statisticsReporter func(path string)

func (s *statistics) Report(path string, reporters ...statisticsReporter) {
	for _, reporter := range reporters {
		reporter(path)
	}
}

func (s *statistics) MissingLocation() statisticsReporter {
	return func(path string) {
		s.NoLocation[path] = struct{}{}
	}
}

func (s *statistics) MissingOrigin() statisticsReporter {
	return func(path string) {
		s.NoOrigin[path] = struct{}{}
	}
}

func (s *statistics) MissingDocumentation() statisticsReporter {
	return func(path string) {
		s.NoDocumentation[path] = struct{}{}
	}
}

func NewStatistics() *statistics {
	return &statistics{
		NoLocation:      set.StringSet{},
		NoOrigin:        set.StringSet{},
		NoDocumentation: set.StringSet{},
	}
}
