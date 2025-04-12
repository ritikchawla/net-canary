# NetCanary 🕊️

A stealthy network threat tripwire that sets up decoy services to detect and monitor unauthorized network activity.

## Features 🛠️

- **Fake Service Honeypots**: Simulates various network services:
  - SSH (port 22)
  - HTTP (port 80)
  - FTP (port 21)
  - MySQL (port 3306)
  - Redis (port 6379)

- **Silent Logging**: Records all connection attempts and interactions
  - IP addresses and timestamps
  - Connection details and headers
  - Attack patterns and behaviors

- **Configurable Responses**: Each service can be customized with:
  - Custom port numbers
  - Service banners
  - Response behaviors

- **Security Alerts**: Optional alerting via:
  - Local log files
  - Webhooks
  - Slack notifications
  - Email alerts

## Installation 🚀

```bash
go install github.com/ritikchawla/net-canary@latest
```

Or build from source:

```bash
git clone https://github.com/ritikchawla/net-canary.git
cd net-canary
go build
```

## Configuration 📝

Configuration is done via a YAML file. Copy the default `config.yaml`:

```yaml
services:
  ssh:
    enabled: true
    port: 22
    host: "0.0.0.0"
    banner: "OpenSSH_8.2p1 Ubuntu-4ubuntu0.5"
  
  # Additional services...

logging:
  file: "netcanary.log"
  level: "info"

alerts:
  slack: "" # Optional Slack webhook URL
  webhook: "" # Optional webhook URL
  email: "" # Optional email address
  paranoia: false
```

## Usage 💡

1. Start NetCanary:
```bash
./net-canary -config config.yaml
```

2. Monitor the log file for activity:
```bash
tail -f netcanary.log
```

## Security Recommendations 🔒

1. Run NetCanary with limited privileges
2. Monitor log files regularly
3. Configure alerts for immediate notification
4. Use in conjunction with existing security measures