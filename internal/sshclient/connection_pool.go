package sshclient

import (
	"fmt"
	"os"
	"sync"
	"time"

	"ssht/internal/config"

	"golang.org/x/crypto/ssh"
)

type ConnectionPool struct {
	mu          sync.Mutex
	connections map[string]*ssh.Client
	configs     map[string]*ssh.ClientConfig
}

func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{
		connections: make(map[string]*ssh.Client),
		configs:     make(map[string]*ssh.ClientConfig),
	}
}

func (p *ConnectionPool) GetConnection(host config.HostConfig) (*ssh.Client, error) {
	key := p.getConnectionKey(host)

	p.mu.Lock()
	defer p.mu.Unlock()

	// Return existing connection if available
	if conn, exists := p.connections[key]; exists {
		// Check if connection is still alive
		if _, _, err := conn.SendRequest("keepalive@openssh.com", true, nil); err == nil {
			return conn, nil
		}
		// Connection is dead, remove it
		conn.Close()
		delete(p.connections, key)
	}

	// Create new connection if not in pool
	clientConfig, err := p.getClientConfig(host)
	if err != nil {
		return nil, err
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.IP, host.Port), clientConfig)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}

	p.connections[key] = client
	return client, nil
}

func (p *ConnectionPool) getClientConfig(host config.HostConfig) (*ssh.ClientConfig, error) {
	key := p.getConnectionKey(host)

	if config, exists := p.configs[key]; exists {
		return config, nil
	}

	clientConfig := &ssh.ClientConfig{
		User:            host.Username,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	switch host.AuthMethod {
	case "password":
		clientConfig.Auth = append(clientConfig.Auth, ssh.Password(host.Password))
	case "private_key":
		key, err := p.readPrivateKey(host.PrivateKeyPath)
		if err != nil {
			return nil, err
		}
		clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(key))
	default:
		return nil, fmt.Errorf("unsupported authentication method: %s", host.AuthMethod)
	}

	p.configs[key] = clientConfig
	return clientConfig, nil
}

func (p *ConnectionPool) readPrivateKey(path string) (ssh.Signer, error) {
	keyBytes, err := os.ReadFile(os.ExpandEnv(path))
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return signer, nil
}

func (p *ConnectionPool) getConnectionKey(host config.HostConfig) string {
	return fmt.Sprintf("%s@%s:%d", host.Username, host.IP, host.Port)
}

func (p *ConnectionPool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var firstErr error
	for key, conn := range p.connections {
		if err := conn.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		delete(p.connections, key)
	}

	return firstErr
}