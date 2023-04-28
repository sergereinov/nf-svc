package config

type Config struct {
	port             int
	summaryIntervals []int
	summaryTopCount  int
	Logs             Logs
}

func (c *Config) Port() int {
	return c.port
}

func (c *Config) SummaryIntervals() []int {
	return c.summaryIntervals
}

func (c *Config) SummaryTopCount() int {
	return c.summaryTopCount
}
