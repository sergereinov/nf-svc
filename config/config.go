package config

type Config struct {
	Port             int
	SummaryIntervals []int
	TrackingClients  []string
	Logs             Logs
}

func (c *Config) GetSummaryIntervals() []int {
	return c.SummaryIntervals
}

func (c *Config) GetTrackingClients() []string {
	return c.TrackingClients
}

func (c *Config) GetLogs() *Logs {
	return &c.Logs
}
