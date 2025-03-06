package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/bxtal-lsn/gotransport/pkg/tls"
)

func main() {
	// Create TLS configuration for the client
	tlsConfig, err := tls.NewClientTLSConfig(tls.Config{
		CAPath:     "ca.crt",
		CertPath:   "client.crt",
		KeyPath:    "client.key",
		ServerName: "localhost", // Must match the server certificate's Common Name
	})
	if err != nil {
		log.Fatalf("Failed to create TLS config: %v", err)
	}

	// Create HTTP client with our custom TLS configuration
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// Make request to the server
	resp, err := client.Get("https://localhost:8443/client")
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read and display the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Response: %s\n", body)
}
