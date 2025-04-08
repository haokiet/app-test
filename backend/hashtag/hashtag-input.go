package hashtag

type ListHashtagInput struct {
	Sort      *string
	Ascending *bool
	Offset    *int64
	Limit     *int64
}

func (s *ListHashtagInput) SetSort(sort string) *ListHashtagInput {
	s.Sort = &sort
	return s
}

func (s *ListHashtagInput) SetOffset(offset int64) *ListHashtagInput {
	s.Offset = &offset
	return s
}

func (s *ListHashtagInput) SetLimit(limit int64) *ListHashtagInput {
	s.Limit = &limit
	return s
}

func (s *ListHashtagInput) SetAscending(ascending bool) *ListHashtagInput {
	s.Ascending = &ascending
	return s
}
