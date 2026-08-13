# Cache Proxy

A simple caching proxy server written in Go. It sits in front of another server, forwards requests to it, and caches the responses so repeat requests don't hit the origin again.

Built as my submission for the [Caching Server project on roadmap.sh](https://roadmap.sh/projects/caching-server).

## How it works

Start the proxy with a port and the origin server you want to cache:

```bash
go run main.go --port 3000 --origin http://dummyjson.com
```

Now requests to `http://localhost:3000/products` get forwarded to `http://dummyjson.com/products`. The response gets saved to disk (in `./cache/`) and returned with a header telling you where it came from:

```
X-Cache: HIT   # served from the cache
X-Cache: MISS  # served from the origin server (and just got cached)
```

Make the same request again and you'll get `X-Cache: HIT`.

## Clearing the cache

```bash
go run main.go --clear-cache
```

This wipes everything in `./cache/`. No server starts when this flag is used.

## Flags

| Flag | Description |
|---|---|
| `--port` | Port the proxy server listens on |
| `--origin` | The server to proxy and cache requests for |
| `--clear-cache` | Clears the cache and exits (no server starts) |

## Running it

```bash
git clone https://github.com/coshi-muhammad/cache-proxy.git
cd cache-proxy
go run main.go --port 3000 --origin http://dummyjson.com
```

Or build a binary:

```bash
go build -o caching-proxy
./caching-proxy --port 3000 --origin http://dummyjson.com
```

## Notes

- Cached responses are stored as plain files in `./cache/`, one per URL.
- This is a learning project, so it's only been tested against simple GET requests.

---
Part of the [roadmap.sh backend projects](https://roadmap.sh/backend/projects).
