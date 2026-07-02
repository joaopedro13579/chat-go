package log

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Logger struct {
	file *os.File
	mu   sync.Mutex
}

var logger *Logger

func Init(filename string) error {
	file, err := os.OpenFile(
		filename,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return err
	}

	logger = &Logger{
		file: file,
	}

	return nil
}

func WriteLog(message string) error {
	if logger == nil {
		return fmt.Errorf("logger not initialized")
	}

	logger.mu.Lock()
	defer logger.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	_, err := logger.file.WriteString(
		fmt.Sprintf("[%s] %s\n", timestamp, message),
	)

	return err
}

func Close() error {
	if logger == nil {
		return nil
	}
	return logger.file.Close()
}
