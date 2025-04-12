package services

import (
	"context"
	"fmt"
	"net"
	"time"
)

// SSHService implements a fake SSH server
type SSHService struct {
	BaseService
	port   int
	host   string
	banner string
}

func NewSSHService(port int, host string, banner string, handler EventHandler) *SSHService {
	return &SSHService{
		BaseService: NewBaseService("ssh", handler),
		port:        port,
		host:        host,
		banner:      banner,
	}
}

func (s *SSHService) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to start SSH service: %v", err)
	}

	s.listener = listener
	fmt.Printf("SSH honeypot listening on %s\n", addr)

	go s.acceptConnections(ctx)
	return nil
}

func (s *SSHService) acceptConnections(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				continue
			}
			go s.handleConnection(conn)
		}
	}
}

func (s *SSHService) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Log the connection attempt
	s.handler.HandleEvent(Event{
		ServiceName: s.Name(),
		RemoteAddr:  conn.RemoteAddr().String(),
		Timestamp:   time.Now().Format(time.RFC3339),
		EventType:   "connection_attempt",
		Details: map[string]string{
			"local_addr": conn.LocalAddr().String(),
		},
	})

	// Send SSH banner
	conn.Write([]byte(fmt.Sprintf("SSH-2.0-%s\r\n", s.banner)))

	// Read initial client message (we don't process it, just log it)
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return
	}

	// Log the client identification
	s.handler.HandleEvent(Event{
		ServiceName: s.Name(),
		RemoteAddr:  conn.RemoteAddr().String(),
		Timestamp:   time.Now().Format(time.RFC3339),
		EventType:   "client_identification",
		Details: map[string]string{
			"client_version": string(buffer[:n]),
		},
	})

	// Sleep briefly to simulate server processing
	time.Sleep(500 * time.Millisecond)

	// Always reject the connection after logging
	conn.Write([]byte("Protocol mismatch.\r\n"))
}
