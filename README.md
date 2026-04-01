# Stockyard Chronicle

**Changelog and release notes — write entries, publish a public page, RSS and email subscribe**

Part of the [Stockyard](https://stockyard.dev) family of self-hosted developer tools.

## Quick Start

```bash
docker run -p 9290:9290 -v chronicle_data:/data ghcr.io/stockyard-dev/stockyard-chronicle
```

Or with docker-compose:

```bash
docker-compose up -d
```

Open `http://localhost:9290` in your browser.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `9290` | HTTP port |
| `DATA_DIR` | `./data` | SQLite database directory |
| `CHRONICLE_LICENSE_KEY` | *(empty)* | Pro license key |

## Free vs Pro

| | Free | Pro |
|-|------|-----|
| Limits | 1 project, 20 entries | Unlimited projects and entries |
| Price | Free | $1.99/mo |

Get a Pro license at [stockyard.dev/tools/](https://stockyard.dev/tools/).

## Category

Creator & Small Business

## License

Apache 2.0
