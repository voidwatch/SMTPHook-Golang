# SMTPHook

## Support

If you find this project useful, consider buying me a coffee:

[![Donate](https://img.shields.io/badge/Buy%20Me%20A%20Coffee-Donate-yellow)](https://coff.ee/voidwatch)

SMTPHook is a self-hosted email parsing pipeline written in Go. It converts emails into structured JSON and forwards them to an HTTP webhook.

This repository supports two separate modes:

1. **Parser-only production setup** (recommended)
2. **Full development environment** (includes webhook receiver and test SMTP)

---

## 1. Production Setup (Parser-Only)

If you're only interested in receiving and forwarding parsed email, use the `scripts/setup-parser.sh` script. This is the **simplest, containerized production setup**.

### Requirements:
- Linux with root access (Debian, Ubuntu, Fedora, Arch, etc.)
- Podman (Docker alternative)
- `pipx` to install `podman-compose`

### Steps:

```bash
chmod +x scripts/setup-parser.sh
./scripts/setup-parser.sh
```

This script will:
- Install `podman`, `podman-compose`, `pipx`
- Set up a `.env` file from `parser/.env.production.example`
- Create `mail/inbox` and `logs/` directories
- Prompt you to configure `.env`

Afterward, start the parser with:

```bash
podman-compose -f podman-compose-prod.yml up -d
```

### Example `.env`

```env
POLL_INTERVAL=5
WEBHOOK_URL=https://your.api/webhook
MAIL_DIR=/mail/inbox
```

### Volumes

Your `podman-compose-prod.yml` mounts:

```yaml
volumes:
  - ./logs:/logs
  - ./mail/inbox:/mail/inbox:Z
```

Place `.eml` test files in `mail/inbox/` to simulate incoming email.

---

### Important: Systemd Quadlet is Not Supported with Podman Compose

If you're using `scripts/setup-parser.sh` and `podman-compose`, **do not use Quadlet `.container` files**. They are not compatible with the 3.4.x `podman-compose` and may cause conflicts.

---

## 2. Full Development Setup

If you're contributing to SMTPHook or testing it locally, use `scripts/setup.sh` instead.

This will:
- Install all development dependencies
- Build `parser`, `webhook`, and `webhook-server`
- Enable local testing via `mailpit` and mock webhooks

```bash
chmod +x scripts/setup.sh
./scripts/setup.sh
```

Start all dev services with:

```bash
podman-compose -f podman-compose.yml up --build
```

---

## 3. Script Summary

| Script                       | Purpose                               |
|------------------------------|---------------------------------------|
| `scripts/setup-parser.sh`    | Minimal production setup (parser only)|
| `scripts/setup.sh`           | Full dev setup                        |
| `scripts/setup-production.sh`| Production setup with Go build        |
| `scripts/run.sh`             | Run dev stack using podman-compose    |
| `scripts/run-prod.sh`        | Run parser using podman-compose       |
| `scripts/reset.sh`           | Stop and clean dev environment        |
| `scripts/reset-prod.sh`      | Stop and clean parser environment     |
| `scripts/diagnose.sh`        | Dev environment diagnostics           |
| `scripts/diagnose-prod.sh`   | Production diagnostics                |
| `scripts/uninstall.sh`       | Remove systemd services and binaries  |
| `scripts/uninstall-prod.sh`  | Clean up production systemd units     |
| `scripts/start-production.sh`| Wrapper to manage prod lifecycle      |

---

## 4. Folder Structure

```
SMTPHook-Golang/
├── parser/                  # Main parser service
│   ├── main.go              # Parses mail and sends JSON to webhook
│   ├── .env.production.example
│   └── Dockerfile
├── scripts/                 # All shell scripts
│   ├── setup-parser.sh      # Parser-only setup (recommended)
│   ├── setup.sh             # Dev setup
│   ├── setup-production.sh  # Production setup with Go build
│   ├── run.sh               # Dev stack runner
│   ├── run-prod.sh          # Production runner
│   ├── reset.sh             # Dev reset
│   ├── reset-prod.sh        # Production reset
│   ├── diagnose.sh          # Dev diagnostics
│   ├── diagnose-prod.sh     # Production diagnostics
│   ├── uninstall.sh         # Dev uninstall
│   ├── uninstall-prod.sh    # Production uninstall
│   └── start-production.sh  # Production lifecycle wrapper
├── podman-compose-prod.yml  # For production use
├── podman-compose.yml       # For development stack
├── mail/inbox/              # Drop .eml files here
├── logs/                    # Log output from parser
```

---

## 5. License

MIT

---

## 6. Contributions

Pull requests and improvements are welcome. Please use `scripts/setup.sh` for dev testing and lint your Go code before submitting.
