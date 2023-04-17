package loggers

const (
	_SUMMARY_LOG = "-summary.log"
)

func NewSummaryWriter(cfg LoggersConfig) chan<- string {
	logger := newLogger(cfg, baseExecutableName()+_SUMMARY_LOG)
	input := make(chan string)
	NewLoggerWriter(input, logger)
	return input
}
