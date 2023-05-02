package config

type Summary struct {
	intervals []int
	topCount  int
}

func (s *Summary) Intervals() []int {
	return s.intervals
}

func (s *Summary) TopCount() int {
	return s.topCount
}
