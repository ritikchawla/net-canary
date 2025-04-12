package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/ritikchawla/net-canary/internal/services"
)

// Logger implements the EventHandler interface for logging security events
type Logger struct {
	file     *os.File
	mu       sync.Mutex
	filepath string
}

type logEntry struct {
	Timestamp  string            `json:"timestamp"`
	Service    string            `json:"service"`
	RemoteAddr string            `json:"remote_addr"`
	EventType  string            `json:"event_type"`
	Details    map[string]string `json:"details"`
}

// NewLogger creates a new file-based logger
func NewLogger(filepath string) (*Logger, error) {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	return &Logger{
		file:     file,
		filepath: filepath,
	}, nil
}

// HandleEvent implements the EventHandler interface
func (l *Logger) HandleEvent(event services.Event) {
	entry := logEntry{
		Timestamp:  event.Timestamp,
		Service:    event.ServiceName,
		RemoteAddr: event.RemoteAddr,
		EventType:  event.EventType,
		Details:    event.Details,
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	data, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling log entry: %v\n", err)
		return
	}

	if _, err := l.file.Write(append(data, '\n')); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to log file: %v\n", err)
		// Attempt to reopen the file
		if err := l.reopenFile(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reopening log file: %v\n", err)
		}
	}
}

// Close closes the log file
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.file.Close()
}

// reopenFile attempts to reopen the log file
func (l *Logger) reopenFile() error {
	if err := l.file.Close(); err != nil {
		return err
	}

	file, err := os.OpenFile(l.filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	l.file = file
	return nil
}
