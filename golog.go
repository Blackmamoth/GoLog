package golog

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

func New() (*logger, error) {
	pc, _, _, _ := runtime.Caller(1)

	caller_file := runtime.FuncForPC(pc).Name()

	caller_working_dir, err := filepath.Abs(caller_file)

	if err != nil {
		color.Red("An error occured while fetching current directory details.\n")
		return nil, err
	}

	logger := logger{
		logger_config: logger_config{
			datetime_format: "Mon, 02 Jan, 2006 15:04:05",
			log_format:      "[%(asctime)] %(levelname) - [%(filename).%(lineno)]: %(message)",
			log_level:       LOG_LEVEL_INFO,
			log_stream:      LOG_STREAM_MULTIPLE,
			log_rotation_config: log_rotation_config{
				file_name:         path.Join(caller_working_dir, "access.log"),
				max_file_size:     50 * 1024 * 1024,
				max_rotation_days: 7,
				rotate_file:       true,
				zip_archive:       true,
			},
		},
	}
	return &logger, nil
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

func (l *logger) compress_file() error {
	file_name := l.logger_config.log_rotation_config.file_name
	log_file, err := os.Open(file_name)
	if err != nil {
		color.Red("Can't access the log file for rotation.\n")
		return err
	}
	defer log_file.Close()
	var dir_slice []string
	if runtime.GOOS == "windows" {
		dir_slice = strings.Split(file_name, "\\")
	} else {
		dir_slice = strings.Split(file_name, "/")
	}
	log_file_name := dir_slice[len(dir_slice)-1]
	log_dir := filepath.Dir(file_name)
	log_zip, err := os.Create(path.Join(log_dir, "access.log.zip"))
	if err != nil {
		color.Red("An error occured while creating zip file: %s\n", err)
		return err
	}
	defer log_zip.Close()
	zip_writer := zip.NewWriter(log_zip)
	file_info, err := log_file.Stat()
	if err != nil {
		color.Red("An error occured while retrieving file info: %s\n", err)
		return err
	}
	header, err := zip.FileInfoHeader(file_info)
	if err != nil {
		color.Red("An error occured while creating file header: %s\n", err)
		return err
	}
	header.Name = log_file_name
	writer, err := zip_writer.CreateHeader(header)
	if err != nil {
		color.Red("An error occured while creating file in zip archive: %s\n", err)
		return err
	}
	_, err = io.Copy(writer, log_file)
	if err != nil {
		color.Red("An error occured while copying file to zip archive: %s\n", err)
		return err
	}
	err = zip_writer.Close()
	if err != nil {
		color.Red("An error occured while closing zip file: %s\n", err)
		return err
	}
	color.Magenta("Log file [%s] has successfully been archived.\n", log_file_name)
	return nil
}

func (l *logger) write_file(log string) error {
	file_name := l.logger_config.log_rotation_config.file_name
	if file_name == "" {
		color.Red("File name not specified for logging.\n")
		return errors.New("file name not specified for logging")
	}
	file, err := os.OpenFile(file_name, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		color.Red("An error occured while accessing the log file.\n")
		return err
	}
	defer file.Close()
	file_stat, err := file.Stat()
	if err != nil {
		file.Close()
		color.Red("An error occured while retrieving file stats.\n")
		return err
	}
	size := file_stat.Size()
	if size >= l.logger_config.log_rotation_config.max_file_size {
		color.Magenta("Rotating old log file, max size reached.\n")
		l.compress_file()
		file.Close()
		err := os.Remove(l.logger_config.log_rotation_config.file_name)
		if err != nil {
			color.Red("Could not remove old log file.\n")
			return err
		}
		file, err := os.OpenFile(file_name, os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			color.Red("An error occured while creating new log file.\n")
			return err
		}
		defer file.Close()
	}
	mod_time := file_stat.ModTime()
	max_days := l.logger_config.log_rotation_config.max_rotation_days
	if time.Since(mod_time) > time.Duration(max_days)*24*time.Hour {
		color.Magenta("Rotating old log file, max days reached.\n")
		l.compress_file()
		file.Close()
		err := os.Remove(l.logger_config.log_rotation_config.file_name)
		if err != nil {
			color.Red("Could not remove old log file.\n")
			return err
		}
		file, err := os.OpenFile(file_name, os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			color.Red("An error occured while creating new log file")
			return err
		}
		defer file.Close()
	}
	file.WriteString(fmt.Sprintf("%s\n", log))
	return nil
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

func (l *logger) Set_Max_File_Size(max_size int64) {
	l.logger_config.log_rotation_config.max_file_size = max_size
}

func (l *logger) Set_Max_Days(max_days int) {
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
			color.Magenta("%s\n", log)
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
