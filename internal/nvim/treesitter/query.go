package treesitter

type Query struct {
	Language string
	Query    string
	Range    *LineRange
}

func (q Query) MapQuery(f func(query string) string) Query {
	return Query{
		Language: q.Language,
		Query:    f(q.Query),
		Range:    q.Range,
	}
}

func (q Query) WithRange(range_ LineRange) Query {
	return Query{
		Language: q.Language,
		Query:    q.Query,
		Range:    &range_,
	}
}
