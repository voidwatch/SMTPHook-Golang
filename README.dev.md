# SMTPHook (Development Setup)

This setup is for development and testing of the full SMTPHook stack, including optional services like:

- webhook-server
- webhook
- mailpit

It provides a full local environment using Podman Compose or Docker Compose.

---

## Features

- Modular Go microservices
- Parses emails and forwards to webhook receivers
- Includes test SMTP server (Mailpit)
- Structured logging
- `.env` config for each service
- One-command setup and reset scripts

---

## Getting Started (Development)

### 1. Clone the repository

```bash
git clone https://github.com/your-org/SMTPHook-Golang.git
cd SMTPHook-Golang
```

### 2. Run setup

```bash
chmod +x scripts/setup.sh
./scripts/setup.sh
```

---

## Environment Configuration

Each service has a `.env.example`. Copy and customize them:

```bash
cp parser/.env.example parser/.env
cp webhook/.env.example webhook/.env
cp webhook-server/.env.example webhook-server/.env
```

---

## Running the Stack

### Podman Compose (development)

```bash
podman-compose -f podman-compose.yml up --build -d
```

### Docker Compose (alternative)

```bash
docker-compose -f podman-compose.yml up --build -d
```

---

## Services Included

| Service         | Description                          |
|-----------------|--------------------------------------|
| mailpit         | Local SMTP mail server + web UI      |
| parser          | Converts email to JSON webhook POSTs |
| webhook         | Simple test receiver (logs to stdout)|
| webhook-server  | Advanced receiver with retry logic   |

---

## Email Testing

Send emails using swaks or other SMTP clients:

```bash
swaks --to test@example.com --server localhost:1025
```

View Mailpit UI at http://localhost:8025

---

## Webhook Format (from parser)

```json
{
  "from": "alerts@example.com",
  "to": "team@example.com",
  "subject": "Disk usage critical",
  "text": "90% used on server01"
}
```

---

## Folder Structure

```
SMTPHook-Golang/
├── Makefile
├── README.dev.md
├── sample-email.json
├── parser/
├── webhook/
├── webhook-server/
├── mailpit/
├── etc/quadlet/
├── scripts/
│   ├── setup.sh
│   ├── reset.sh
│   ├── uninstall.sh
│   ├── diagnose.sh
│   ├── run.sh
│   └── ...
├── podman-compose.yml
```

---

## Notes

- This is intended for local testing and development.
- Use `scripts/setup-production.sh` for production deployment with only the `parser` service.

---

## License

MIT
