package proto

import (
	"context"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Manager holds multiple gRPC client connections
type Manager struct {
	mu      sync.Mutex
	clients map[string]*grpc.ClientConn
}

// NewManager initializes the manager
func NewManager() *Manager {
	return &Manager{
		clients: make(map[string]*grpc.ClientConn),
	}
}

// Connect dials a new gRPC server and stores it with a key
func (m *Manager) Connect(name, addr string) (*grpc.ClientConn, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if conn, ok := m.clients[name]; ok {
		return conn, nil // already connected
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // for local/dev; replace with TLS in prod
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	m.clients[name] = conn
	return conn, nil
}

// Get retrieves an existing client connection
func (m *Manager) Get(name string) (*grpc.ClientConn, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	conn, ok := m.clients[name]
	return conn, ok
}

// CloseAll closes all connections
func (m *Manager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, conn := range m.clients {
		if err := conn.Close(); err != nil {
			log.Printf("failed to close gRPC client %s: %v", name, err)
		}
	}
}
