package main

type GroupRow interface{}

type Data[R GroupRow] interface {
	Partition() string
	GroupKey() string
	Aggregate(acc R) R
}

type Summary[R GroupRow] struct {
	groups map[string]map[string]R
}

func (s *Summary[R]) Add(data Data[R]) {
	if s == nil {
		return
	}

	if s.groups == nil {
		s.groups = make(map[string]map[string]R)
	}

	partitionKey := data.Partition()
	partition, ok := s.groups[partitionKey]
	if !ok {
		partition = make(map[string]R)
		s.groups[partitionKey] = partition
	}

	groupKey := data.GroupKey()
	groupRow := partition[groupKey]
	groupRow = data.Aggregate(groupRow)
	partition[groupKey] = groupRow
}

func (s *Summary[R]) Dump() map[string]map[string]R {
	if s == nil {
		return map[string]map[string]R{}
	}
	return s.groups
}
