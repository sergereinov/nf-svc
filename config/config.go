package config

type Config struct {
	Port             int
	SummaryIntervals []int
	TrackingClients  []string
	Logs             Logs
}
