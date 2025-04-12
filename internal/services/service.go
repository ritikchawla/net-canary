package services

import (
	"context"
	"net"
)

// Event represents a security event detected by a service
type Event struct {
	ServiceName string
	RemoteAddr  string
	Timestamp   string
	EventType   string
	Details     map[string]string
}

// EventHandler processes security events
type EventHandler interface {
	HandleEvent(event Event)
}

// Service defines the interface that all honeypot services must implement
type Service interface {
	// Start begins listening for connections
	Start(ctx context.Context) error

	// Stop gracefully shuts down the service
	Stop() error

	// Name returns the service identifier
	Name() string
}

// BaseService provides common functionality for all services
type BaseService struct {
	name     string
	listener net.Listener
	handler  EventHandler
}

func NewBaseService(name string, handler EventHandler) BaseService {
	return BaseService{
		name:    name,
		handler: handler,
	}
}

func (b *BaseService) Name() string {
	return b.name
}

func (b *BaseService) Stop() error {
	if b.listener != nil {
		return b.listener.Close()
	}
	return nil
}
