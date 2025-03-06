package dns

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
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
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.Address, s.Port)

	// Create a new DNS server
	dnsServer := &dns.Server{
		Addr: addr,
		Net:  "udp",
	}

	// Set request handler
	dns.HandleFunc(".", s.handleRequest)

	// Store server reference
	s.server = dnsServer

	// Start server
	fmt.Printf("Starting DNS server on %s\n", addr)
	return dnsServer.ListenAndServe()
}

// handleRequest handles DNS requests
func (s *Server) handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	// Create a new response message
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	// Process each question
	for _, question := range r.Question {
		fmt.Printf("Query: %s, type: %d\n", question.Name, question.Qtype)

		switch question.Qtype {
		case dns.TypeA:
			s.handleARecord(question, m)
		case dns.TypeAAAA:
			// We don't support IPv6 yet, but include for future
			m.Rcode = dns.RcodeNameError
		default:
			// For other types, return NXDOMAIN
			m.Rcode = dns.RcodeNameError
		}
	}

	// Send response
	w.WriteMsg(m)
}

// handleARecord handles A record requests
func (s *Server) handleARecord(q dns.Question, m *dns.Msg) {
	// Look for an exact match first
	record, exists := s.Storage.Get(q.Name)

	// If not found, try to match wildcards
	if !exists {
		// Split the domain into parts
		parts := strings.Split(strings.TrimSuffix(q.Name, "."), ".")

		// Try matching with wildcards (e.g., *.example.com)
		for i := 0; i < len(parts)-1; i++ {
			wildcardDomain := "*." + strings.Join(parts[i+1:], ".") + "."
			record, exists = s.Storage.Get(wildcardDomain)
			if exists {
				break
			}
		}
	}

	if exists {
		if record.Type == dnsrecords.A {
			rr, err := dns.NewRR(fmt.Sprintf("%s IN A %s", q.Name, record.Value))
			if err == nil {
				m.Answer = append(m.Answer, rr)
			}
		} else if record.Type == dnsrecords.CNAME {
			// Handle CNAME records by adding both the CNAME and resolving it
			cnameRR, err := dns.NewRR(fmt.Sprintf("%s IN CNAME %s", q.Name, record.Value))
			if err == nil {
				m.Answer = append(m.Answer, cnameRR)

				// Try to resolve the CNAME target
				if targetRecord, targetExists := s.Storage.Get(record.Value); targetExists && targetRecord.Type == dnsrecords.A {
					targetRR, err := dns.NewRR(fmt.Sprintf("%s IN A %s", record.Value, targetRecord.Value))
					if err == nil {
						m.Answer = append(m.Answer, targetRR)
					}
				}
			}
		}
	} else {
		// Record not found
		m.Rcode = dns.RcodeNameError
	}
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
