package util

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

const (
	// DefaultPort is the default port to use if none is specified
	DefaultPort = "8080"
	// MaxPortAttempts is the maximum number of ports to try
	MaxPortAttempts = 50
	// MinValidPort is the minimum valid port number
	MinValidPort = 1024
	// MaxValidPort is the maximum valid port number
	MaxValidPort = 65535
)

// findAvailablePort tries to find an available port starting from the given port
// It will try up to maxAttempts ports (incrementing by 1 each time)
func findAvailablePort(startPort int, maxAttempts int) (int, error) {
	for port := startPort; port < startPort+maxAttempts; port++ {
		// Don't go beyond the valid port range
		if port > MaxValidPort {
			return 0, fmt.Errorf("exceeded maximum valid port %d", MaxValidPort)
		}

		if isPortAvailable(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("could not find an available port after %d attempts", maxAttempts)
}

// isPortAvailable checks if a port is available by trying to listen on it
func isPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

// GetPort returns an available port, starting from the provided defaultPort
// If defaultPort is unavailable, it will find the next available port
func GetPort(defaultPort string) string {
	if defaultPort == "" {
		defaultPort = DefaultPort
		log.Printf("No port specified, using default port: %s", DefaultPort)
	}

	// Try to parse the default port
	port, err := strconv.Atoi(defaultPort)
	if err != nil {
		log.Printf("Invalid port number '%s', using %s instead", defaultPort, DefaultPort)
		port, _ = strconv.Atoi(DefaultPort)
	}

	// Validate port range
	if port < MinValidPort || port > MaxValidPort {
		log.Printf("Port %d is outside valid range (%d-%d), using %s instead",
			port, MinValidPort, MaxValidPort, DefaultPort)
		port, _ = strconv.Atoi(DefaultPort)
	}

	// First, check if the default port is available
	if isPortAvailable(port) {
		log.Printf("Using port: %d", port)
		return strconv.Itoa(port)
	}

	log.Printf("Port %d is already in use. Searching for an available port...", port)

	// Find an available port by incrementing
	availablePort, err := findAvailablePort(port+1, MaxPortAttempts)
	if err != nil {
		log.Printf("Warning: %v. Trying OS-assigned port.", err)
		// Try with port 0 which lets the OS assign a random available port
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			log.Fatalf("Failed to find available port: %v", err)
		}

		// Get the port that was assigned
		addr := listener.Addr().(*net.TCPAddr)
		availablePort = addr.Port
		listener.Close()

		log.Printf("Using OS-assigned port: %d", availablePort)
	} else {
		log.Printf("Found available port: %d", availablePort)
	}

	return strconv.Itoa(availablePort)
}
