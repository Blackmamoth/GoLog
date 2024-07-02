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
	max_file_size     int64
	max_rotation_days int
}

type logger_config struct {
	log_rotation_config log_rotation_config
	log_level           LOG_LEVEL
	log_stream          LOG_STREAM
	log_format          string
	datetime_format     string
	with_emoji          bool
	exit_on_critical    bool
}

type logger struct {
	logger_config logger_config
}

type Logger interface {
	Set_Log_Level(log_level LOG_LEVEL)
	Set_Log_Format(log_format string)
	Set_Log_Stream(log_stream LOG_STREAM)
	Set_Datetime_Format(format string)
	Set_File_Name(file_name string)
	Set_Max_File_Size(max_size int64)
	Set_Max_Days(max_days int)
	Zip_Archive(zip_archive bool)
	With_Emoji(with_emoji bool)
	EXIT_ON_CRITICAL(should_exit bool)
	TRACE(msg string)
	DEBUG(msg string)
	INFO(msg string)
	WARN(msg string)
	ERROR(msg string)
	CRITICAL(msg string)
}
