package user

type ListUserInput struct {
	Sort      *string
	Ascending *bool
	Offset    *int64
	Limit     *int64
}

func (s *ListUserInput) SetSort(sort string) *ListUserInput {
	s.Sort = &sort
	return s
}

func (s *ListUserInput) SetOffset(offset int64) *ListUserInput {
	s.Offset = &offset
	return s
}

func (s *ListUserInput) SetLimit(limit int64) *ListUserInput {
	s.Limit = &limit
	return s
}

func (s *ListUserInput) SetAscending(ascending bool) *ListUserInput {
	s.Ascending = &ascending
	return s
}
