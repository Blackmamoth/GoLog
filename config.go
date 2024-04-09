package golog

type LOG_LEVEL int
type LOG_STREAM int

const (
	LOG_LEVEL_TRACE LOG_LEVEL = iota + 1
	LOG_LEVEL_DEBUG
	LOG_LEVEL_INFO
	LOG_LEVEL_WARN
	LOG_LEVEL_ERROR
	LOG_LEVEL_CRITICAL
)

const (
	LOG_STREAM_CONSOLE LOG_STREAM = iota + 1
	LOG_STREAM_FILE
	LOG_STREAM_MULTIPLE
)

type log_rotation_config struct {
	rotate_file       bool
	file_name         string
	zip_archive       bool
	max_file_size     uint32
	max_rotation_days uint
}

type logger_config struct {
	log_rotation_config log_rotation_config
	log_level           LOG_LEVEL
	log_stream          LOG_STREAM
	log_format          string
	datetime_format     string
}

type logger struct {
	logger_config logger_config
}
