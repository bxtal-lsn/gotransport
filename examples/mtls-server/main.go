package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/bxtal-lsn/gotransport/pkg/tls"
)

func index(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Welcome to the mTLS secured server!")
}

func showClientCert(w http.ResponseWriter, req *http.Request) {
	// Check if client provided a valid certificate
	var commonName string
	if req.TLS != nil && len(req.TLS.VerifiedChains) > 0 && len(req.TLS.VerifiedChains[0]) > 0 {
		commonName = req.TLS.VerifiedChains[0][0].Subject.CommonName
		fmt.Fprintf(w, "Hello, %s! Your client certificate was verified successfully.", commonName)
	} else {
		fmt.Fprintf(w, "No valid client certificate provided.")
	}
}

func main() {
	// Define HTTP handlers
	http.HandleFunc("/", index)
	http.HandleFunc("/client", showClientCert)

	// Create TLS configuration for the server
	tlsConfig, err := tls.NewServerTLSConfig(tls.Config{
		CAPath:     "ca.crt",
		CertPath:   "server.crt",
		KeyPath:    "server.key",
		ClientAuth: tls.RequireAndVerifyClientCert,
		MinVersion: tls.VersionTLS12,
	})
	if err != nil {
		log.Fatalf("Failed to create TLS config: %v", err)
	}

	// Create HTTPS server
	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}

	// Start the server
	log.Printf("Starting mTLS server on https://localhost:8443")
	log.Printf("Use the client example to connect with a client certificate")
	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
