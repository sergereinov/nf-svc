package config

type Logs struct {
	KeepDays      int
	MaxFileSizeMB int
	Path          string
}

func (c *Logs) GetKeepDays() int {
	return c.KeepDays
}

func (c *Logs) GetMaxFileSizeMB() int {
	return c.MaxFileSizeMB
}

func (c *Logs) GetPath() string {
	return c.Path
}
