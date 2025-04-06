# Muraena Configuration Guide

This document provides detailed information about the `config.toml` configuration file structure and available options.

## Table of Contents
- [Proxy Configuration](#proxy-configuration)
- [Origins Configuration](#origins-configuration)
- [Transform Rules](#transform-rules)
- [Redirect Rules](#redirect-rules)
- [Logging](#logging)
- [Redis Configuration](#redis-configuration)
- [TLS Settings](#tls-settings)
- [Modules Configuration](#modules-configuration)

## Proxy Configuration

The `[proxy]` section controls how Muraena handles traffic routing between the phishing target and the legitimate destination.

```toml
[proxy]
phishing = "your-phishing-domain.com"    # Your phishing domain
destination = "target-site.com"          # Target legitimate site
IP = "0.0.0.0"                          # Listen interface (all interfaces)
port = 443                              # HTTPS port
listener = "tcp4"                       # Network protocol

# Optional port mapping
# portmapping = "1337:443"              # Format: <ListeningPort>:<TargetPort>

[proxy.HTTPtoHTTPS]
enable = true                           # Force HTTP to HTTPS redirect
port = 80                              # HTTP port to redirect from
```

## Origins Configuration

The `[origins]` section defines external origins that need to be handled by the proxy.

```toml
[origins]
externalOriginPrefix = "ext-"           # Prefix for external origins
externalOrigins = [                     # List of external domains to handle
    "github.com",
    "login.microsoft.com",
    # ... other domains
]
```

## Transform Rules

The `[transform]` section defines how Muraena transforms requests and responses.

### Base64 Transformation
```toml
[transform.base64]
enable = false                          # Enable base64 string transformation
padding = [ "=", "." ]                  # Base64 padding characters
```

### Request Transformation
```toml
[transform.request]
# userAgent = "MuraenaProxy"           # Custom User-Agent
# headers = ["Cookie", "Referer", "Origin"]  # Headers to transform

# Remove specific request headers
remove.headers = [
    "X-FORWARDED-FOR",
    "X-FORWARDED-PROTO"
]

# Add custom request headers
add.headers = [
    {name = "X-Phishing", value = "via Muraena"}
]
```

### Response Transformation
```toml
[transform.response]
skipContentType = [                     # Content types to skip transformation
    "image/jpeg",
    "font/*",
    "application/*"
]

headers = [                             # Headers to transform
    "Location",
    "Origin",
    "Set-Cookie",
    "Access-Control-Allow-Origin"
]

customContent = [                       # Custom content replacements
    ["integrity=", "integrify="]
]

# Security headers to remove
remove.headers = [
    "Content-Security-Policy",
    "X-Frame-Options",
    # ... other security headers
]
```

## Redirect Rules

Configure custom URL redirections:

```toml
[[redirect]]
hostname = "phishing.click"             # Domain to match
path = "/buyer/reset"                   # Path to match
redirectTo = "https://example.com"      # Redirect destination
httpStatusCode = 301                    # HTTP redirect code
```

## Logging

Configure logging settings:

```toml
[log]
enable = true                           # Enable logging
# filePath = "muraena.log"             # Custom log file path
```

## Redis Configuration

Optional Redis settings for session management:

```toml
[redis]
# host = "127.0.0.1"                   # Redis host
# port = 6379                          # Redis port
# password = ""                        # Redis password
```

## TLS Settings

Configure SSL/TLS settings:

```toml
[tls]
enable = true                           # Enable TLS
expand = false                          # Expand certificate validation
certificate = "/path/to/cert.pem"       # Certificate path
key = "/path/to/key.pem"               # Private key path
root = "/path/to/chain.pem"            # Certificate chain
minVersion = "TLS1.2"                  # Minimum TLS version
insecureSkipVerify = true              # Skip certificate verification
sessionTicketsDisabled = true          # Disable session tickets
renegotiationSupport = "Never"         # TLS renegotiation setting
```

## Modules Configuration

### Tracking Module

Configure user tracking and credential capturing:

```toml
[tracking]
enable = true                           # Enable tracking
trackRequestCookie = true               # Track request cookies

[tracking.trace]
identifier = "_gat"                     # Tracking identifier
validator = "[a-zA-Z0-9]{5}"           # Validation pattern
header = "X-Custom-Header"             # Custom tracking header

[tracking.secrets]
paths = [                              # Paths to monitor for credentials
    "/login",
    "/submit"
]

[[tracking.secrets.patterns]]           # Credential capture patterns
label = "Username"
start = "login="
end = "&"

[[tracking.secrets.patterns]]
label = "Password"
start = "passwd="
end = ""
```

### Static Server Module

Serve static content:

```toml
[staticServer]
enable = true                           # Enable static server
localPath = "./static/"                # Local path for static files
urlPath = "/evilpath/static"           # URL path for static content
```

### Optional Modules

#### NecroBrowser Module
```toml
[necrobrowser]
enable = false                          # Enable NecroBrowser
endpoint = "http://localhost:3000/instrument"
# Additional NecroBrowser settings...
```

#### Watchdog Module
```toml
[watchdog]
enable = false                          # Enable Watchdog
dynamic = true                          # Enable dynamic rules
rules = "./config/watchdog.rules"       # Rules file path
geoDB = "./config/geoDB.mmdb"          # GeoIP database
```

#### Telegram Module
```toml
[telegram]
enable = false                          # Enable Telegram notifications
botToken = "your-bot-token"            # Telegram bot token
chatIDs = ["chat-id-1"]                # Telegram chat IDs
```

## Best Practices

1. **Security**:
   - Always use HTTPS (TLS enabled)
   - Set appropriate security headers
   - Use strong TLS configuration
   - Regularly rotate certificates

2. **Performance**:
   - Configure appropriate content type skipping
   - Use Redis for better session management
   - Monitor log sizes

3. **Monitoring**:
   - Enable appropriate logging
   - Configure tracking for important paths
   - Set up notifications if needed

4. **Maintenance**:
   - Regularly update GeoIP database
   - Keep watchdog rules updated
   - Monitor Redis performance 