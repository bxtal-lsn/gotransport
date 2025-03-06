package dnsrecords

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// RecordType represents DNS record types
type RecordType string

const (
	// A record type
	A RecordType = "A"
	// CNAME record type
	CNAME RecordType = "CNAME"
)

// Record represents a DNS record
type Record struct {
	Domain string     `json:"domain"`
	Type   RecordType `json:"type"`
	Value  string     `json:"value"`
	TTL    uint32     `json:"ttl"`
}

// Storage manages DNS records
type Storage struct {
	records map[string]Record
	mu      sync.RWMutex
	file    string
}

// NewStorage creates a new DNS record storage
func NewStorage(storagePath string) (*Storage, error) {
	// Create storage directory if it doesn't exist
	dir := filepath.Dir(storagePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("failed to create storage directory: %w", err)
		}
	}

	storage := &Storage{
		records: make(map[string]Record),
		file:    storagePath,
	}

	// Load existing records if file exists
	if _, err := os.Stat(storagePath); err == nil {
		if err := storage.load(); err != nil {
			return nil, err
		}
	}

	return storage, nil
}

// Add adds a new DNS record
func (s *Storage) Add(domain string, recordType RecordType, value string, ttl uint32) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Normalize domain name
	domain = normalizeDomain(domain)

	// Validate record type
	if recordType != A && recordType != CNAME {
		return fmt.Errorf("unsupported record type: %s", recordType)
	}

	// Validate value based on record type
	if recordType == A {
		if net.ParseIP(value) == nil {
			return fmt.Errorf("invalid IP address: %s", value)
		}
	}

	// Add the record
	s.records[domain] = Record{
		Domain: domain,
		Type:   recordType,
		Value:  value,
		TTL:    ttl,
	}

	// Save changes
	return s.save()
}

// Remove removes a DNS record
func (s *Storage) Remove(domain string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Normalize domain name
	domain = normalizeDomain(domain)

	// Check if record exists
	if _, exists := s.records[domain]; !exists {
		return fmt.Errorf("record not found: %s", domain)
	}

	// Remove the record
	delete(s.records, domain)

	// Save changes
	return s.save()
}

// Get retrieves a DNS record
func (s *Storage) Get(domain string) (Record, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Normalize domain name
	domain = normalizeDomain(domain)

	record, exists := s.records[domain]
	return record, exists
}

// List returns all DNS records
func (s *Storage) List() []Record {
	s.mu.RLock()
	defer s.mu.RUnlock()

	records := make([]Record, 0, len(s.records))
	for _, record := range s.records {
		records = append(records, record)
	}
	return records
}

// save persists DNS records to disk
func (s *Storage) save() error {
	data, err := json.MarshalIndent(s.records, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal records: %w", err)
	}

	if err := os.WriteFile(s.file, data, 0o644); err != nil {
		return fmt.Errorf("failed to save records: %w", err)
	}

	return nil
}

// load reads DNS records from disk
func (s *Storage) load() error {
	data, err := os.ReadFile(s.file)
	if err != nil {
		return fmt.Errorf("failed to read records: %w", err)
	}

	if err := json.Unmarshal(data, &s.records); err != nil {
		return fmt.Errorf("failed to unmarshal records: %w", err)
	}

	return nil
}

// normalizeDomain ensures a domain ends with a period
func normalizeDomain(domain string) string {
	domain = strings.ToLower(domain)
	if !strings.HasSuffix(domain, ".") {
		domain = domain + "."
	}
	return domain
}
