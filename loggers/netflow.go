package loggers

const (
	_NETFLOW_LOG = "-netflow.log"
)

func NewNetflowWriter(cfg LoggersConfig) chan<- string {
	logger := newLogger(cfg, baseExecutableName()+_NETFLOW_LOG)
	input := make(chan string)
	NewLoggerWriter(input, logger)
	return input
}
