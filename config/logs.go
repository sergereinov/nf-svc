package config

type Logs struct {
	keepDays      int
	maxFileSizeMB int
	dir           string
}

func (c *Logs) KeepDays() int {
	return c.keepDays
}

func (c *Logs) MaxFileSizeMB() int {
	return c.maxFileSizeMB
}

func (c *Logs) Dir() string {
	return c.dir
}
