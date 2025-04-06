# Muraena Deployment Guide

This document provides detailed instructions for deploying the Muraena project in both Docker and standalone environments.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Docker Deployment](#docker-deployment)
- [Manual Deployment](#manual-deployment)
- [Configuration](#configuration)
- [SSL Certificates](#ssl-certificates)
- [Troubleshooting](#troubleshooting)
- [Local Development Setup](#local-development-setup)

## Prerequisites

### Docker Deployment
- Docker Engine 20.10.x or later
- Docker Compose v2.x or later
- Git
- Valid SSL certificates
- Redis (included in Docker Compose)

### Manual Deployment
- Go 1.21 or later
- Redis 7.x
- Git
- Make
- Valid SSL certificates

## Docker Deployment

1. Clone the repository:
   ```bash
   git clone https://github.com/your-org/muraena.git
   cd muraena
   ```

2. Create necessary directories:
   ```bash
   mkdir -p certificates
   ```

3. Configure your environment:
   - Copy the example configuration:
     ```bash
     cp config/config.toml.example config/config.toml
     ```
   - Edit `config/config.toml` according to your needs
   - See [Configuration Guide](CONFIG.md) for detailed configuration options

4. Place your SSL certificates in the `certificates` directory:
   - `certificates/cert.pem` (fullchain certificate)
   - `certificates/key.pem` (private key)
   - `certificates/chain.pem` (certificate chain)

5. Update the TLS configuration in `config.toml`:
   ```toml
   [tls]
   enable = true
   certificate = "/app/certificates/cert.pem"
   key = "/app/certificates/key.pem"
   root = "/app/certificates/chain.pem"
   ```
   See [TLS Configuration](CONFIG.md#tls-settings) for more options.

6. Configure proxy settings in `config.toml`:
   ```toml
   [proxy]
   phishing = "your-phishing-domain.com"
   destination = "target-site.com"
   ```
   See [Proxy Configuration](CONFIG.md#proxy-configuration) for more options.

7. Start the containers:
   ```bash
   docker-compose up -d
   ```

8. Verify deployment:
   ```bash
   docker-compose ps
   ```

## Manual Deployment

1. Clone the repository:
   ```bash
   git clone https://github.com/your-org/muraena.git
   cd muraena
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the application:
   ```bash
   make build
   ```

4. Install and configure Redis:
   ```bash
   # Ubuntu/Debian
   sudo apt install redis-server
   
   # CentOS/RHEL
   sudo yum install redis
   ```
   See [Redis Configuration](CONFIG.md#redis-configuration) for setup options.

5. Configure the application:
   ```bash
   cp config/config.toml.example config/config.toml
   # Edit config.toml according to your needs
   ```
   Refer to the [Configuration Guide](CONFIG.md) for detailed settings.

6. Run the application:
   ```bash
   ./build/muraena -config config/config.toml
   ```

## Configuration

The main configuration file is `config.toml`. This file controls all aspects of Muraena's operation. See our comprehensive [Configuration Guide](CONFIG.md) for detailed information about:

- [Proxy Settings](CONFIG.md#proxy-configuration)
- [Origins Configuration](CONFIG.md#origins-configuration)
- [Transform Rules](CONFIG.md#transform-rules)
- [TLS Settings](CONFIG.md#tls-settings)
- [Module Configurations](CONFIG.md#modules-configuration)

Key configuration sections include:

1. **Proxy Configuration**
   - Set your phishing and target domains
   - Configure ports and protocols
   - Set up HTTP to HTTPS redirection

2. **TLS Configuration**
   - Certificate paths
   - TLS versions and security settings
   - Certificate validation options

3. **Module Configuration**
   - Tracking module for credential capture
   - Static server for serving content
   - Optional modules (NecroBrowser, Watchdog, Telegram)

4. **Transform Rules**
   - Request/response transformation
   - Header modifications
   - Content replacements

For detailed configuration options and examples, see the [Configuration Guide](CONFIG.md).

## SSL Certificates

SSL certificates are required for HTTPS support. You can:

1. Use Let's Encrypt:
   ```bash
   certbot certonly --standalone -d your-domain.com
   ```

2. Use self-signed certificates (for testing only):
   ```bash
   openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
     -keyout certificates/key.pem \
     -out certificates/cert.pem
   ```

After obtaining certificates, update the [TLS configuration](CONFIG.md#tls-settings) in your `config.toml`.

## Troubleshooting

### Common Issues

1. Port conflicts:
   - Check if ports 80/443 are already in use
   - Modify the port mapping in docker-compose.yml if needed
   - See [Proxy Configuration](CONFIG.md#proxy-configuration) for port settings

2. Redis connection issues:
   - Verify Redis is running: `redis-cli ping`
   - Check Redis connection settings in [Redis Configuration](CONFIG.md#redis-configuration)

3. Certificate errors:
   - Ensure certificates exist in the correct location
   - Verify certificate permissions
   - Check certificate validity: `openssl x509 -in cert.pem -text -noout`
   - Review [TLS Settings](CONFIG.md#tls-settings)

### Logs

- Docker logs:
  ```bash
  docker-compose logs -f muraena
  ```

- Standalone logs:
  - Check the `log` directory
  - Use `journalctl` if running as a service
  - See [Logging Configuration](CONFIG.md#logging)

### Health Check

Monitor the application health:
```bash
curl -k https://your-domain.com/health
```

## Security Considerations

1. Always run as non-root user (handled in Docker setup)
2. Keep Redis secured and not exposed to public
3. Regularly update dependencies
4. Monitor logs for suspicious activities
5. Use strong SSL certificates
6. Follow security best practices for your deployment environment
7. Review [Security Best Practices](CONFIG.md#best-practices)

## Updates and Maintenance

1. Update the application:
   ```bash
   git pull
   docker-compose build
   docker-compose up -d
   ```

2. Backup important data:
   - Configuration files
   - SSL certificates
   - Redis data (if persistent storage is needed)

3. Regular maintenance tasks:
   - Update GeoIP database
   - Review and update watchdog rules
   - Monitor Redis performance
   - See [Maintenance Best Practices](CONFIG.md#best-practices)

## Local Development Setup

This section covers setting up Muraena for local development and testing, using a custom domain (e.g., github.local) that points to your local machine.

### Prerequisites for Local Development

- Docker and Docker Compose
- WSL2 (for Windows users)
- mkcert (for generating local SSL certificates)
- Administrative access (for modifying hosts files)

### Local SSL Certificates

1. Install mkcert:
   ```bash
   # On Ubuntu/Debian
   sudo apt install libnss3-tools
   sudo wget -O /usr/local/bin/mkcert https://github.com/FiloSottile/mkcert/releases/download/v1.4.3/mkcert-v1.4.3-linux-amd64
   sudo chmod +x /usr/local/bin/mkcert

   # On Windows (using Chocolatey)
   choco install mkcert
   ```

2. Generate local certificates:
   ```bash
   # Create certificates directory
   mkdir -p certs
   cd certs

   # Generate certificates for your local domain
   mkcert github.local "*.github.local"
   ```

### Local Configuration

1. Create a local SSL configuration file:
   ```bash
   cp config/config.toml config/local-ssl.toml
   ```

2. Configure local-ssl.toml:
   ```toml
   [proxy]
   phishing = "github.local"
   destination = "github.com"
   IP = "0.0.0.0"
   port = 443
   listener = "tcp4"

   [tls]
   enable = true
   expand = false
   certificate = "/app/certs/github.local+1.pem"
   key = "/app/certs/github.local+1-key.pem"
   root = "/app/certs/github.local+1.pem"
   minVersion = "TLS1.2"
   maxVersion = "TLS1.3"
   insecureSkipVerify = true
   ```

### Host File Configuration

#### Linux/WSL Setup

1. Create a setup script (setup-hosts.sh):
   ```bash
   #!/bin/bash
   DOMAIN="github.local"
   IP="0.0.0.0"

   # Remove existing entry if present
   sudo sed -i "/$DOMAIN/d" /etc/hosts

   # Add new entry
   echo "$IP    $DOMAIN" | sudo tee -a /etc/hosts

   # Display current entries
   echo "Current entries in /etc/hosts for $DOMAIN:"
   grep "$DOMAIN" /etc/hosts
   ```

2. Make it executable and run:
   ```bash
   chmod +x setup-hosts.sh
   sudo ./setup-hosts.sh
   ```

#### Windows Setup

1. Create a PowerShell script (setup-hosts-windows.ps1):
   ```powershell
   $domain = "github.local"
   $ip = "127.0.0.1"
   $hostsFile = "C:\Windows\System32\drivers\etc\hosts"

   # Remove existing entry
   $content = Get-Content $hostsFile | Where-Object { $_ -notmatch $domain }
   $content | Set-Content $hostsFile

   # Add new entry
   Add-Content -Path $hostsFile -Value "$ip`t$domain" -Force

   # Display current entries
   Write-Host "Current entries in hosts file for $domain:"
   Get-Content $hostsFile | Select-String $domain
   ```

2. Run as Administrator:
   ```powershell
   # Open PowerShell as Administrator and run:
   Set-ExecutionPolicy Bypass -Scope Process
   .\setup-hosts-windows.ps1
   ```

### Running the Local Environment

1. Start the services:
   ```bash
   docker-compose down
   CONFIG_FILE=/app/config/local-ssl.toml docker-compose up -d --build
   ```

2. Verify the setup:
   ```bash
   # Check service status
   docker-compose ps

   # View logs
   docker-compose logs -f muraena
   ```

3. Test the connection:
   ```bash
   # Test HTTPS connection (ignore certificate warnings)
   curl -k https://github.local
   ```

### Debugging

For debugging purposes, you can use the debug entrypoint:

1. Modify docker-compose.yml:
   ```yaml
   services:
     muraena:
       entrypoint: ["/app/debug-entrypoint.sh"]
   ```

2. The debug entrypoint will:
   - Show environment variables
   - List certificate files
   - Display configuration
   - Keep the container running for inspection

### Common Local Development Issues

1. Certificate Issues:
   - Ensure certificates are properly generated and placed in the certs directory
   - Check certificate paths in local-ssl.toml
   - Accept self-signed certificate in browser (first visit)

2. DNS Resolution:
   - Verify hosts file entries in both Windows and WSL
   - Clear DNS cache if needed:
     ```bash
     # Windows (PowerShell as Admin)
     ipconfig /flushdns

     # Linux/WSL
     sudo systemd-resolve --flush-caches
     ```

3. Port Conflicts:
   - Check if ports 80/443 are available
   - Stop any existing web servers
   - Modify port mappings in docker-compose.yml if needed

4. WSL/Windows Integration:
   - Ensure WSL2 is properly configured
   - Check network connectivity between Windows and WSL
   - Verify both hosts files are properly configured

For more detailed troubleshooting, refer to the [Troubleshooting](#troubleshooting) section above. 