package sshclient

import (
	"bytes"
	"fmt"
	"os"

	"ssht/internal/config"

	"golang.org/x/crypto/ssh"
)

type CommandResult struct {
	Output string
	Error  error
}

func ExecuteCommand(host config.HostConfig, command string) CommandResult {
	// Configure SSH client
	clientConfig := &ssh.ClientConfig{
		User:            host.Username,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Set authentication method
	switch host.AuthMethod {
	case "password":
		clientConfig.Auth = append(clientConfig.Auth, ssh.Password(host.Password))
	case "private_key":
		key, err := os.ReadFile(os.ExpandEnv(host.PrivateKeyPath))
		if err != nil {
			return CommandResult{Error: fmt.Errorf("failed to read private key: %v", err)}
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return CommandResult{Error: fmt.Errorf("failed to parse private key: %v", err)}
		}
		clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(signer))
	default:
		return CommandResult{Error: fmt.Errorf("unsupported authentication method: %s", host.AuthMethod)}
	}

	// Connect to SSH server
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.IP, host.Port), clientConfig)
	if err != nil {
		return CommandResult{Error: fmt.Errorf("connection failed: %v", err)}
	}
	defer client.Close()

	// Create session and execute command
	session, err := client.NewSession()
	if err != nil {
		return CommandResult{Error: fmt.Errorf("failed to create session: %v", err)}
	}
	defer session.Close()

	// Capture command output
	var stdoutBuf, stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	if err := session.Run(command); err != nil {
		return CommandResult{
			Error:  fmt.Errorf("command execution failed: %v", err),
			Output: stderrBuf.String(),
		}
	}

	// Return command output
	return CommandResult{
		Output: stdoutBuf.String(),
	}
}
