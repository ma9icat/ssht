package app

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"ssht/internal/config"
	"ssht/internal/sshclient"

	"github.com/sirupsen/logrus"
)

func Run() error {
	// Initialize configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger
	log := logrus.New()

	// Configure log format
	switch cfg.LogFormat {
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "time",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	default: // text
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
			DisableQuote:    true,
			PadLevelText:    true,
		})
	}

	// Determine target nodes to execute
	targetNodes := cfg.Groups.Default
	if len(cfg.Nodes) > 0 {
		targetNodes = cfg.Nodes
	}

	// Set log level
	if cfg.Debug {
		log.SetLevel(logrus.DebugLevel)
		log.Debug("Debug mode enabled")
		log.Debugf("Loaded configuration: %+v", cfg)
		log.Debugf("Target nodes: %v", targetNodes)
		log.Debugf("Available hosts: %v", cfg.Hosts)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}

	// Configure log output
	if cfg.LogFile != "" {
		file, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		defer file.Close()
		log.SetOutput(file)
	}

	// Log task start information
	log.WithFields(logrus.Fields{
		"command": cfg.Command,
		"nodes":   len(targetNodes),
	}).Info("Starting SSH task execution")

	// Execute SSH commands concurrently
	var wg sync.WaitGroup
	taskCh := make(chan struct{}, 10) // Limit concurrency to 10

	for _, host := range cfg.Hosts {
		// Check if host is in target nodes list
		shouldExecute := false
		for _, node := range targetNodes {
			if host.Name == node {
				shouldExecute = true
				break
			}
		}

		if shouldExecute {
			wg.Add(1)
			taskCh <- struct{}{} // Acquire a concurrency slot

			host := host // Create local copy for goroutine
			go func(h config.HostConfig) {
				defer wg.Done()
				defer func() { <-taskCh }()

				startTime := time.Now()
				log.Debugf("[%s] Executing command: %q", h.Name, cfg.Command)
				log.Debugf("[%s] Connection details: %s@%s:%d", h.Name, h.Username, h.IP, h.Port)

				result := sshclient.ExecuteCommand(h, cfg.Command)
				duration := time.Since(startTime).Round(time.Millisecond)

				log.Debugf("[%s] Execution completed in %s", h.Name, duration)
				if result.Output != "" {
					log.Debugf("[%s] Command output:\n%s", h.Name, result.Output)
				}
				output := strings.TrimSpace(result.Output)
				if result.Error != nil {
					log.Errorf("[%s] Command failed", h.Name)
				} else if output != "" {
					log.Infof("[%s] %s", h.Name, output)
				} else {
					log.Infof("[%s] Command executed", h.Name)
				}
				log.Infof("[%s] Duration: %s", h.Name, duration)
			}(host)
		}
	}

	wg.Wait() // Wait for all goroutines to complete

	// Log task completion
	log.Info("SSH task completed")

	return nil
}

// Removed unused nodeLogHook
