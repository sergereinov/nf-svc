package config

type Logs struct {
	KeepDays      int
	MaxFileSizeMB int
	Dir           string
}

func (c *Logs) GetKeepDays() int {
	return c.KeepDays
}

func (c *Logs) GetMaxFileSizeMB() int {
	return c.MaxFileSizeMB
}

func (c *Logs) GetDir() string {
	return c.Dir
}
