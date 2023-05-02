package config

type Logs struct {
	keepDays         int
	maxFileSizeMB    int
	dir              string
	enableSummaryLog bool
	enableNetFlowLog bool
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

func (c *Logs) EnableSummaryLog() bool {
	return c.enableSummaryLog
}

func (c *Logs) EnableNetFlowLog() bool {
	return c.enableNetFlowLog
}
