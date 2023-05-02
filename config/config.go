package config

type Config struct {
	port    int
	Summary Summary
	Logs    Logs
}

func (c *Config) Port() int {
	return c.port
}
