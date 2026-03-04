# miroxy

**Micro Relay Proxy** — A tiny, stateless proxy relay built in Go for resource-constrained devices like Raspberry Pi Zero 2 W.

```
Client → Cloudflare Tunnel → miroxy → 3rd Party API
```

## Install

```bash
curl -fsSL https://pavunato.github.io/miroxy/install.sh | sh
```

## Usage

Send a POST to `/proxy` with the target request details:

```bash
curl -X POST http://localhost:8080/proxy \
  -H 'Content-Type: application/json' \
  -d '{
    "method": "GET",
    "url": "https://api.example.com/data",
    "headers": {
      "Authorization": "Bearer sk-xxx"
    },
    "timeout": 30
  }'
```

The upstream response (status code, headers, body) is streamed back directly.

### Health Check

```bash
curl http://localhost:8080/health
```

## Options

| Flag / Env | Description | Default |
|---|---|---|
| `--port` | Port to listen on | `8080` |
| `--token` | Bearer token for auth | — |
| `MIROXY_ADDR` | Listen address (env) | `:8080` |
| `MIROXY_TOKEN` | Bearer token (env) | — |

```bash
# Run on custom port
miroxy --port 3000

# Run with auth
miroxy --port 8080 --token my-secret
```

## Request Schema

| Field | Type | Required | Default | Description |
|---|---|---|---|---|
| `url` | string | yes | — | Upstream URL |
| `method` | string | no | `GET` | HTTP method |
| `headers` | object | no | — | Headers to send upstream |
| `body` | string | no | — | Request body |
| `timeout` | int | no | `30` | Timeout in seconds |

## Build from Source

```bash
# Local build
make build

# Cross-compile for Pi Zero 2 W
make build-pi2

# All platforms
make build-all
```

## Deploy with Cloudflare Tunnel

```bash
# 1. Install miroxy on your Pi
curl -fsSL https://pavunato.github.io/miroxy/install.sh | sh

# 2. Install cloudflared
curl -L https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-arm -o /usr/local/bin/cloudflared
chmod +x /usr/local/bin/cloudflared

# 3. Create tunnel pointing to miroxy
cloudflared tunnel login
cloudflared tunnel create miroxy
cloudflared tunnel route dns miroxy proxy.yourdomain.com
cloudflared tunnel run --url http://localhost:8080 miroxy
```

## License

MIT
