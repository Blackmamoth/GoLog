package golog

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

func New() *logger {
	logger := logger{
		logger_config: logger_config{
			datetime_format: "Mon, 02 Jan, 2006 15:04:05",
			log_format:      "[%(asctime)] %(levelname) - [%(filename).%(lineno)]: %(message)",
			log_level:       LOG_LEVEL_INFO,
			log_stream:      LOG_STREAM_MULTIPLE,
			log_rotation_config: log_rotation_config{
				file_name:         "access.log",
				max_file_size:     50 * 1024 * 1024,
				max_rotation_days: 7,
				rotate_file:       true,
				zip_archive:       true,
			},
		},
	}
	return &logger
}

func (l *logger) generate_log(log_level string, msg string, caller_file string, caller_line int) string {
	var log string
	log = l.logger_config.log_format
	now := time.Now().Format(l.logger_config.datetime_format)
	log = strings.ReplaceAll(log, "%(asctime)", now)
	log = strings.ReplaceAll(log, "%(levelname)", log_level)
	log = strings.ReplaceAll(log, "%(filename)", caller_file)
	log = strings.ReplaceAll(log, "%(lineno)", strconv.Itoa(caller_line))
	log = strings.ReplaceAll(log, "%(message)", msg)
	return log
}

func (l *logger) write_file(log string) {
	file_name := l.logger_config.log_rotation_config.file_name
	if file_name == "" {
		color.Red("File name not specified for logging.\n")
		os.Exit(1)
	}
	file, err := os.OpenFile(file_name, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		color.Red("An error occured while creating file [%s].\n", file_name)
		os.Exit(1)
	}
	defer file.Close()
	file.WriteString(fmt.Sprintf("%s\n", log))
}

func (l *logger) Set_Log_Level(log_level LOG_LEVEL) {
	l.logger_config.log_level = log_level
}

func (l *logger) Set_Log_Format(log_format string) {
	l.logger_config.log_format = log_format
}

func (l *logger) Set_Log_Stream(log_stream LOG_STREAM) {
	l.logger_config.log_stream = log_stream
}

func (l *logger) Set_Datetime_Format(format string) {
	l.logger_config.datetime_format = format
}

func (l *logger) Set_File_Name(file_name string) {
	l.logger_config.log_rotation_config.file_name = file_name
}

func (l *logger) Set_Max_File_Size(max_size uint32) {
	l.logger_config.log_rotation_config.max_file_size = max_size
}

func (l *logger) Set_Max_Days(max_days uint) {
	l.logger_config.log_rotation_config.max_rotation_days = max_days
}

func (l *logger) Zip_Archive(zip_archive bool) {
	l.logger_config.log_rotation_config.zip_archive = zip_archive
}

func (l *logger) TRACE(msg string) {
	_, file, line, _ := runtime.Caller(1)
	log := l.generate_log("TRACE", msg, file, line)
	if l.logger_config.log_level <= LOG_LEVEL_TRACE {
		if l.logger_config.log_stream == LOG_STREAM_CONSOLE || l.logger_config.log_stream == LOG_STREAM_MULTIPLE {
			color.White("%s\n", log)
		}
		if l.logger_config.log_stream == LOG_STREAM_FILE || l.logger_config.log_stream == LOG_STREAM_MULTIPLE {
			l.write_file(log)
		}
	}
}

func (l *logger) DEBUG(msg string) {
	_, file, line, _ := runtime.Caller(1)
	log := l.generate_log("DEBUG", msg, file, line)
	if l.logger_config.log_level <= LOG_LEVEL_DEBUG {
		if l.logger_config.log_stream == LOG_STREAM_CONSOLE || l.logger_config.log_stream == LOG_STREAM_MULTIPLE {
			color.Blue("%s\n", log)
		}
		if l.logger_config.log_stream == LOG_STREAM_FILE || l.logger_config.log_stream == LOG_STREAM_MULTIPLE {
			l.write_file(log)
		}
	}
}

func (l *logger) INFO(msg string) {
	_, file, line, _ := runtime.Caller(1)
	log := l.generate_log("INFO", msg, file, line)
	if l.logger_config.log_level <= LOG_LEVEL_INFO {
		if l.logger_config.log_stream == LOG_STREAM_CONSOLE || l.logger_config.log_stream == LOG_STREAM_MULTIPLE {
			color.Green("%s\n", log)
		}
		if l.logger_config.log_stream == LOG_STREAM_FILE || l.logger_config.log_stream == LOG_STREAM_MULTIPLE {
			l.write_file(log)
		}
	}
}

func (l *logger) WARN(msg string) {
	_, file, line, _ := runtime.Caller(1)
	log := l.generate_log("WARN", msg, file, line)
	if l.logger_config.log_level <= LOG_LEVEL_WARN {
		if l.logger_config.log_stream == LOG_STREAM_CONSOLE || l.logger_config.log_stream == LOG_STREAM_MULTIPLE {
			color.Yellow("%s\n", log)
		}
		if l.logger_config.log_stream == LOG_STREAM_FILE || l.logger_config.log_stream == LOG_STREAM_MULTIPLE {
			l.write_file(log)
		}
	}
}

func (l *logger) ERROR(msg string) {
	_, file, line, _ := runtime.Caller(1)
	log := l.generate_log("ERROR", msg, file, line)
	if l.logger_config.log_level <= LOG_LEVEL_ERROR {
		if l.logger_config.log_stream == LOG_STREAM_CONSOLE || l.logger_config.log_stream == LOG_STREAM_MULTIPLE {
			color.Red("%s\n", log)
		}
		if l.logger_config.log_stream == LOG_STREAM_FILE || l.logger_config.log_stream == LOG_STREAM_MULTIPLE {
			l.write_file(log)
		}
	}
}

func (l *logger) CRITICAL(msg string) {
	_, file, line, _ := runtime.Caller(1)
	log := l.generate_log("CRITICAL", msg, file, line)
	if l.logger_config.log_level <= LOG_LEVEL_CRITICAL {
		if l.logger_config.log_stream == LOG_STREAM_CONSOLE || l.logger_config.log_stream == LOG_STREAM_MULTIPLE {
			critical := color.New(color.FgRed, color.Bold)
			critical.Printf("%s\n", log)
		}
		if l.logger_config.log_stream == LOG_STREAM_FILE || l.logger_config.log_stream == LOG_STREAM_MULTIPLE {
			l.write_file(log)
		}
	}
}
