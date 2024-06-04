# GoLog

GoLog is a logger for GoLang, designed to provide flexible and easy-to-use logging capabilities for Go projects. It aims to simplify logging tasks while offering customization options to suit various project requirements.

## Installation 🛠

You can install GoLog as a project dependency using `go get`:

```bash
go get -u github.com/blackmamoth/GoLog
```

## How It Works ⚙

GoLog simplifies logging tasks by providing convenient methods for logging messages at various severity levels, including `TRACE`, `DEBUG`, `INFO`, `WARN`, `ERROR`, and `CRITICAL`. Additionally, it offers the following customization options:

- `New()`: Creates a new GoLog instance for logging.

- `Set_Log_Level(log_level LOG_LEVEL)`: Sets the minimum log level.

- `Set_Log_Format(log_format string)`: Updates the log format displayed. The default format is `[%(asctime)] %(levelname) - [%(filename).%(lineno)]: %(message)`.

- `Set_Log_Stream(log_stream LOG_STREAM)`: Sets the log stream to console, log file, or both.

- `Set_Datetime_Format(format string)`: Sets the format string for time representation in logs. The default format is `Mon, 02 Jan, 2006 15:04:05`.

- `Set_File_Name(file_name string)`: Sets the log file name. The default is `access.log`.

- `Set_Max_File_Size(max_size int64)`: Sets the maximum size for a log file. If exceeded, the current log file is archived, and a new one is created. The default size is `50 * 1024 * 1024` bytes.

- `Set_Max_Days(max_days int)`: Sets the maximum number of days a log file should store logs. If exceeded, the current file is archived, and a new one is created.

- `Zip_Archive(zip_archive bool)`: Specifies whether a log file should be archived while it is rotated. If `false`, the current log file will be removed, and a new one will be generated.

- `With_Emoji(with_emoji bool)`: Specifies whether log statements should contain emoji.

## Example Usage 📝

```go
package main

import (
	golog "github.com/blackmamoth/GoLog"
)

func main() {
	// Initialize the logger
	logger := golog.New()

	// Set logging level
	logger.Set_Log_Level(golog.LOG_LEVEL_DEBUG)

	// Set log format
	logger.Set_Log_Format("[%(asctime)] %(levelname) - %(message)")

	// Set log stream
	logger.Set_Log_Stream(golog.LOG_STREAM_MULTIPLE)

	// Set max size
	logger.Set_Max_File_Size(10 * 1024 * 1024)

	// Make use of emojis in log statement
	logger.With_Emoji(true)

	// Log messages
	logger.DEBUG("Debug message")
	logger.INFO("Info message")
	logger.WARN("Warning message")
	logger.ERROR("Error message")
	logger.CRITICAL("Critical message")
}
```

## Planned/Pending Features 📅

- ~~DRY principle~~

- Thread safe logging

- Write tests
