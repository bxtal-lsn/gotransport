# GoTransport

A modern TLS certificate and DNS management tool written in Go. GoTransport helps you generate and manage TLS certificates for your infrastructure and provides local DNS services for your development environments with a clean CLI interface.

## Features

- Create CA (Certificate Authority) certificates
- Create server and client certificates signed by your CA
- Generate RSA keys
- Support for multiple domains (SAN certificates)
- Mutual TLS (mTLS) support
- Lightweight DNS server for local development
- DNS record management (A and CNAME records)
- Wildcard domain support
- Easy configuration via YAML

## Installation

### From Source

```bash
git clone https://github.com/bxtal-lsn/gotransport.git
cd gotransport
go install
```

### Using Go Install

```bash
go install github.com/bxtal-lsn/gotransport@latest
```

## Configuration

GoTransport uses a YAML configuration file to define certificate properties. Create a file named `tls.yaml` with the following structure:

```yaml
caCert:
  serial: 1
  validForYears: 10
  subject:
    country: US
    organization: Your Organization
    commonName: Your CA

certs:
  server:
    serial: 1
    validForYears: 1
    dnsNames: ["localhost", "example.com"]
    subject:
      country: US
      organization: Your Organization
      commonName: example.com
  
  client:
    serial: 2
    validForYears: 1
    subject:
      country: US
      organization: Your Organization
      commonName: client
```

## Usage

### Create a CA Certificate

```bash
gotransport ca --key-out ca.key --cert-out ca.crt
```

### Create a Server Certificate

```bash
gotransport cert --name server --ca-key ca.key --ca-cert ca.crt --key-out server.key --cert-out server.crt
```

### Create a Client Certificate

```bash
gotransport cert --name client --ca-key ca.key --ca-cert ca.crt --key-out client.key --cert-out client.crt
```

### Generate an RSA Key

```bash
gotransport key --key-out key.pem --key-length 4096
```

### Start DNS Server

```bash
# Start on port 53 (requires root/admin privileges)
sudo gotransport dns serve

# Start on non-privileged port 5353
gotransport dns serve --insecure

# Specify storage location
gotransport dns serve --storage /path/to/dns.json
```

### Add DNS Record

```bash
# Add A record
gotransport dns add --domain harbor.local --type A --value 192.168.1.100

# Add CNAME record
gotransport dns add --domain www.harbor.local --type CNAME --value harbor.local

# Add wildcard record
gotransport dns add --domain "*.harbor.local" --type A --value 192.168.1.100
```

### List DNS Records

```bash
gotransport dns list
```

### Remove DNS Record

```bash
gotransport dns remove --domain harbor.local
```

## Setting up Harbor with HTTPS and DNS

You can use GoTransport to generate certificates and set up DNS for your Harbor registry:

1. Create a CA:
```bash
gotransport ca --key-out harbor-ca.key --cert-out harbor-ca.crt
```

2. Create a server certificate for Harbor:
```bash
gotransport cert --name harbor --ca-key harbor-ca.key --ca-cert harbor-ca.crt --key-out harbor.key --cert-out harbor.crt
```

3. Configure Harbor to use the certificates in your `harbor.yml`:
```yaml
hostname: harbor.local
https:
  port: 443
  certificate: /path/to/harbor.crt
  private_key: /path/to/harbor.key
```

4. Set up a local DNS server for harbor.local:
```bash
# Add DNS record for Harbor
gotransport dns add --domain harbor.local --type A --value 192.168.1.100

# Add wildcard record for subdomains
gotransport dns add --domain "*.harbor.local" --type A --value 192.168.1.100

# Start DNS server (requires root/admin privileges)
sudo gotransport dns serve
```

5. Configure your system to use the local DNS server:
   - Linux: Add `nameserver 127.0.0.1` to `/etc/resolv.conf`
   - macOS: Add `127.0.0.1` as DNS server in Network settings
   - Windows: Add `127.0.0.1` as DNS server in network adapter settings

Now you can access Harbor using `https://harbor.local` and pull/push images using `harbor.local/project/image:tag`.

## Examples

See the `examples/` directory for:
- mTLS Server Example
- mTLS Client Example

## License

MIT
