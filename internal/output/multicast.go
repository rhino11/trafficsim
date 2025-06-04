package output

import (
	"fmt"
	"log"
	"net"
	"time"
)

// MulticastPublisher handles publishing CoT messages to multicast groups
type MulticastPublisher struct {
	conn          *net.UDPConn
	multicastAddr *net.UDPAddr
	generator     *CoTGenerator
	interval      time.Duration
}

// NewMulticastPublisher creates a new multicast publisher for CoT messages
func NewMulticastPublisher(multicastIP string, port int) (*MulticastPublisher, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", multicastIP, port))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve multicast address: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create UDP connection: %w", err)
	}

	return &MulticastPublisher{
		conn:          conn,
		multicastAddr: addr,
		generator:     NewCoTGenerator(),
		interval:      5 * time.Second, // Default 5 second interval
	}, nil
}

// SetPublishInterval sets the interval for publishing CoT messages
func (p *MulticastPublisher) SetPublishInterval(interval time.Duration) {
	p.interval = interval
}

// PublishPlatformState publishes a single platform state as CoT message
func (p *MulticastPublisher) PublishPlatformState(state PlatformState) error {
	cotMessage, err := p.generator.GenerateCoTMessage(state)
	if err != nil {
		return fmt.Errorf("failed to generate CoT message: %w", err)
	}

	_, err = p.conn.Write(cotMessage)
	if err != nil {
		return fmt.Errorf("failed to send CoT message: %w", err)
	}

	return nil
}

// StartPublishing starts a goroutine that continuously publishes platform states
func (p *MulticastPublisher) StartPublishing(platformStates chan PlatformState, stopChan chan bool) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	var latestStates = make(map[string]PlatformState)

	go func() {
		for {
			select {
			case state := <-platformStates:
				latestStates[state.ID] = state
			case <-ticker.C:
				// Publish all latest platform states
				for _, state := range latestStates {
					if err := p.PublishPlatformState(state); err != nil {
						log.Printf("Error publishing CoT message for platform %s: %v", state.ID, err)
					}
				}
			case <-stopChan:
				return
			}
		}
	}()
}

// Close closes the multicast connection
func (p *MulticastPublisher) Close() error {
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

// GetMulticastAddress returns the configured multicast address
func (p *MulticastPublisher) GetMulticastAddress() string {
	return p.multicastAddr.String()
}
