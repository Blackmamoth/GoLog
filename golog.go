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

func New() *logger {
	_, file, _, _ := runtime.Caller(1)

	// caller_file := runtime.FuncForPC(pc)

	caller_working_dir := filepath.Dir(file)
	logger := logger{
		logger_config: logger_config{
			datetime_format: "Mon, 02 Jan, 2006 15:04:05",
			log_format:      "[%(asctime)] %(levelname) - [%(filename).%(lineno)]: %(message)",
			log_level:       LOG_LEVEL_INFO,
			log_stream:      LOG_STREAM_FILE,
			log_rotation_config: log_rotation_config{
				file_name:         filepath.Join(caller_working_dir, "access.log"),
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

func (l *logger) compress_file() error {
	file_path := l.logger_config.log_rotation_config.file_name
	log_file, err := os.Open(file_path)

	if err != nil {
		color.Red("Can't access the log file for rotation.\n")
		return err
	}

	defer log_file.Close()

	var dir_slice []string

	if runtime.GOOS == "windows" {
		dir_slice = strings.Split(file_path, "\\")
	} else {
		dir_slice = strings.Split(file_path, "/")
	}

	log_file_name := dir_slice[len(dir_slice)-1]
	log_dir := filepath.Dir(file_path)
	file_name := filepath.Base(file_path)
	log_zip, err := os.Create(path.Join(log_dir, fmt.Sprintf("%s.zip", file_name)))

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

func (l *logger) rotate_file(file *os.File, file_name string, roatate_msg string) error {
	color.Magenta("Rotating old log file, %s.\n", roatate_msg)
	l.compress_file()
	file.Close()
	err := os.Remove(l.logger_config.log_rotation_config.file_name)

	if err != nil {
		color.Red("Could not remove old log file.\n")
		return err
	}
	file, err = os.OpenFile(file_name, os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		color.Red("An error occured while creating new log file.\n")
		return err
	}

	defer file.Close()
	return nil
}

func (l *logger) write_file(log string) error {
	file_path := l.logger_config.log_rotation_config.file_name

	if file_path == "" {
		color.Red("File name not specified for logging.\n")
		return errors.New("file name not specified for logging")
	}

	file, err := os.OpenFile(file_path, os.O_APPEND|os.O_CREATE, 0644)

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
		err = l.rotate_file(file, file_path, "max size reached")
		defer file.Close()
		if err != nil {
			return err
		}
	}

	mod_time := file_stat.ModTime()
	max_days := l.logger_config.log_rotation_config.max_rotation_days

	if time.Since(mod_time) > time.Duration(max_days)*24*time.Hour {
		err = l.rotate_file(file, file_path, "max days reached")
		defer file.Close()
		if err != nil {
			return err
		}
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

func (l *logger) log(level string, msg string, levelColor func(string, ...interface{})) {
	_, file, line, _ := runtime.Caller(2)

	file = filepath.Base(file)
	log := l.generate_log(level, msg, file, line)

	if l.logger_config.log_level <= LOG_LEVEL_TRACE {

		if l.logger_config.log_stream == LOG_STREAM_CONSOLE || l.logger_config.log_stream == LOG_STREAM_MULTIPLE {
			levelColor("%s\n", log)
		}

		if l.logger_config.log_stream == LOG_STREAM_FILE || l.logger_config.log_stream == LOG_STREAM_MULTIPLE {
			l.write_file(log)
		}
	}
}

func (l *logger) TRACE(msg string) {
	l.log("TRACE", msg, color.White)
}

func (l *logger) DEBUG(msg string) {
	l.log("DEBUG", msg, color.Blue)
}

func (l *logger) INFO(msg string) {
	l.log("INFO", msg, color.Magenta)
}

func (l *logger) WARN(msg string) {
	l.log("WARN", msg, color.Yellow)
}

func (l *logger) ERROR(msg string) {
	l.log("ERROR", msg, color.Red)
}

func (l *logger) CRITICAL(msg string) {
	critical := color.New(color.FgRed, color.Bold, color.Underline)
	l.log("CRITICAL", msg, critical.PrintfFunc())
}
