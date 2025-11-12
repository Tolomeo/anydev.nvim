package crawler

type statistics struct {
	missingLocation      map[string]struct{}
	missingOrigin        map[string]struct{}
	missingDocumentation map[string]struct{}
}

func (s *statistics) LocationNotFound(path string) {
	s.missingLocation[path] = struct{}{}
}

func (s *statistics) OriginNotFound(path string) {
	s.missingOrigin[path] = struct{}{}
}

func (s *statistics) DocumentationNotFound(path string) {
	s.missingDocumentation[path] = struct{}{}
}
