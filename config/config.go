package config

type Config struct {
	Port             int
	SummaryIntervals []int
	SummaryTopCount  int
	TrackingClients  []string
	Logs             Logs
}

func (c *Config) GetSummaryIntervals() []int {
	return c.SummaryIntervals
}

func (c *Config) GetSummaryTopCount() int {
	return c.SummaryTopCount
}

func (c *Config) GetTrackingClients() []string {
	return c.TrackingClients
}
