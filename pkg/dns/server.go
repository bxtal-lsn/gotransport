package dns

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/bxtal-lsn/gotransport/internal/dnsrecords"
	"github.com/miekg/dns"
)

// Server represents a DNS server
type Server struct {
	Address string
	Port    int
	Storage *dnsrecords.Storage
	server  *dns.Server
}

// NewServer creates a new DNS server
func NewServer(address string, port int, storagePath string) (*Server, error) {
	storage, err := dnsrecords.NewStorage(storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	server := &Server{
		Address: address,
		Port:    port,
		Storage: storage,
	}

	return server, nil
}

// Start starts the DNS server
// In pkg/dns/server.go, modify the Start() method
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.Address, s.Port)

	// Try to bind to the port to check if it's available
	listener, err := net.ListenPacket("udp", addr)
	if err != nil {
		// If the original port is unavailable, try to find a free port
		if s.Port > 1024 { // Only for non-privileged ports
			for testPort := s.Port + 1; testPort < s.Port+100; testPort++ {
				testAddr := fmt.Sprintf("%s:%d", s.Address, testPort)
				listener, err = net.ListenPacket("udp", testAddr)
				if err == nil {
					s.Port = testPort
					addr = testAddr
					fmt.Printf("Original port was busy, using port %d instead\n", testPort)
					break
				}
			}
		}

		if err != nil {
			return fmt.Errorf("failed to bind to any port: %w", err)
		}
	}

	// Create a DNS server using the established listener
	dnsServer := &dns.Server{
		PacketConn: listener,
		Net:        "udp",
	}

	// Set request handler
	dns.HandleFunc(".", s.handleRequest)

	// Store server reference
	s.server = dnsServer

	// Start server
	fmt.Printf("Starting DNS server on %s\n", addr)
	return dnsServer.ActivateAndServe()
}

func (s *Server) handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	// Create a new response message
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false
	m.Authoritative = true

	// Copy recursion desired flag from request (but don't actually do recursion)
	m.RecursionDesired = r.RecursionDesired
	m.RecursionAvailable = false

	// Initialize response code as success
	m.Rcode = dns.RcodeSuccess

	// Flag to track if we found any valid records
	recordFound := false

	// Process each question
	for _, question := range r.Question {
		fmt.Printf("Query: %s, type: %d\n", question.Name, question.Qtype)

		switch question.Qtype {
		case dns.TypeA:
			foundRecord := s.handleARecord(question, m)
			recordFound = recordFound || foundRecord
		case dns.TypeAAAA:
			// We don't support IPv6 yet
			// Don't set NXDOMAIN here, just don't add any records
		default:
			// For other types, just don't add any records
		}
	}

	// Only set NXDOMAIN if we didn't find any records for any question
	if !recordFound && len(m.Answer) == 0 {
		m.Rcode = dns.RcodeNameError
	}

	// Send response
	w.WriteMsg(m)
}

// Modified handleARecord to return whether it found a record
func (s *Server) handleARecord(q dns.Question, m *dns.Msg) bool {
	// Look for an exact match first
	record, exists := s.Storage.Get(q.Name)

	// If not found, try to match wildcards
	if !exists {
		// Wildcard matching code...
	}

	if exists {
		if record.Type == dnsrecords.A {
			rr, err := dns.NewRR(fmt.Sprintf("%s %d IN A %s", q.Name, record.TTL, record.Value))
			if err == nil {
				m.Answer = append(m.Answer, rr)
				return true
			}
		} else if record.Type == dnsrecords.CNAME {
			// CNAME handling...
			return true
		}
	}

	return false
}

// Stop stops the DNS server
func (s *Server) Stop() error {
	if s.server != nil {
		return s.server.Shutdown()
	}
	return nil
}

// StartWithSignalHandling starts the DNS server and handles termination signals
func (s *Server) StartWithSignalHandling() error {
	// Create a channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	errChan := make(chan error)
	go func() {
		if err := s.Start(); err != nil {
			errChan <- err
		}
	}()

	// Wait for either an error or a signal
	select {
	case err := <-errChan:
		return err
	case sig := <-sigChan:
		fmt.Printf("Received signal: %v\n", sig)
		fmt.Println("Shutting down DNS server...")
		return s.Stop()
	}
}
