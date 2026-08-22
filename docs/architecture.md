# Architecture

## Overview

The Url Mentalist is a real-time URL analysis engine. When a URL is submitted â€” via the web UI, Chrome extension, or REST API â€” the Go backend runs **18 concurrent analyzers** across 7 signal categories, aggregates a trust/risk score, assigns a verdict, and returns a fully explainable report. Results are cached in Valkey so repeat lookups are instant.

---

## System Architecture

![The Url Mentalist Architecture](../assets/architecture.png)

Four containerized services run on a shared Docker bridge network (`urlvet-net`). The Go backend is the **only** service that makes outbound calls â€” the frontend, Chrome, and cache are strictly internal.

| Service | Container | Role | Port |
|---|---|---|---|
| `urlvet-web` | SvelteKit UI | Renders the web interface; proxies API calls to the backend | `:3000` prod Â· `:5173` dev |
| `urlvet-backend` | Go REST API | Validates URLs, runs analyzers, aggregates scores, manages cache | `:8080` |
| `urlvet-chrome` | Headless Chrome | Takes page screenshots and serves content via WebSocket (chromedp) | `:9222` |
| `urlvet-valkey` | Valkey (Redis-compatible) | LRU result cache, volume-persisted across restarts | `:6379` |

**External services** (reached only by the backend over HTTPS/TCP):
- **PhishTank** â€” confirmed phishing database lookups
- **DNS resolvers** â€” NS, MX, IP resolution checks
- **WHOIS servers** â€” domain age and registration data

Clients (browser, Chrome extension, API consumers) communicate directly with the Go backend. The SvelteKit frontend forwards all `/api/v1/` requests server-side to avoid CORS complexity.

---

## Request Lifecycle

![The Url Mentalist Analyzer Pipeline](../assets/pipeline.png)

```
Client
  â”‚
  â”‚  GET /api/v1/analyze?url=...
  â–¼
Go Backend
  â”œâ”€ 1. Validate & normalize URL (add scheme if missing, reject private IPs)
  â”œâ”€ 2. Check Valkey cache
  â”‚      â””â”€ HIT  â†’ return full cached result immediately (sub-millisecond)
  â”‚      â””â”€ MISS â†’ continue
  â”œâ”€ 3. Launch 18 goroutines via sync.WaitGroup
  â”‚      â”œâ”€ Each task runs independently; panics are recovered per-task
  â”‚      â”œâ”€ Tasks share a read-only Input struct and write to a mutex-guarded Output
  â”‚      â””â”€ All 18 complete (or timeout) before proceeding
  â”œâ”€ 4. Aggregate scores â†’ apply formula â†’ assign verdict
  â”œâ”€ 5. Store result in Valkey (24 h TTL)
  â””â”€ 6. Return: trust score Â· verdict Â· per-signal reasons Â· redirect chain Â·
              screenshot Â· per-task timings
```

---

## Detection Engine

18 goroutines run across **7 signal categories**, producing **33 individual signals**. Every check emits a labeled reason string â€” good, bad, or neutral â€” so the final score is always fully explainable.

### Scoring Formula

```
finalScore = clamp(50 + (trustScore âˆ’ riskScore) Ã— 0.5, 0, 100)
```

- **50** is the neutral baseline â€” an unknown URL with no signals scores exactly 50
- Trust signals pull the score up; risk signals pull it down, each weighted at 0.5Ã— so neither dominates
- Both `trustScore` and `riskScore` are individually clamped to 0â€“100 before the formula runs

| Range | Verdict |
|---|---|
| â‰¥ 65 | Safe |
| 30 â€“ 64 | Suspicious |
| < 30 | Risky |

### Signal Categories

**URL Signals** â€” 8 checks, purely structural, no network call

1. Raw IP address as hostname
2. Punycode / IDN encoding (lookalike domain spoofing)
3. URL shortener (hides true destination)
4. Excessive URL length
5. Excessive URL path depth
6. Phishing keywords in URL path (`login`, `verify`, `secure`, `update` â€¦)
7. Excessive subdomain count
8. Non-ASCII Unicode characters in hostname (IDN homograph attack)

**HTTP / Network** â€” 4 checks, single HTTP request via `httpCombinedTask`

9. Redirect chain hop count
10. Cross-domain redirect (final destination differs from source)
11. HSTS support
12. HTTP status code

**DNS** â€” 3 checks

13. NS record validity
14. MX record validity
15. IP resolution

**TLS / SSL** â€” 2 checks, single TLS handshake

16. TLS presence and hostname mismatch
17. Certificate chain â€” validity, expiry, issuer, CT log status, known-bad fingerprints

**Domain Intelligence** â€” 6 checks

18. Domain rank (position in top-1M global popularity list)
19. TLD trust / risk / ICANN status
20. Domain age via WHOIS (newly registered = high risk)
21. DNSSEC (cryptographic DNS response integrity)
22. Shannon entropy score (flags algorithmically generated domains)
23. Typosquatting & combo-squatting across 500+ known brands

**Content Analysis** â€” 8 checks, one HTTP GET to fetch page HTML

24. Login form on unranked or newly registered domain
25. Payment form (credit card, CVV fields)
26. Personal information form
27. Hidden `<iframe>` (credential theft / clickjacking vector)
28. Tracking pixels (1Ã—1 hidden images)
29. Brand name in page content vs. hosting domain
30. Form submitting to an external domain
31. Password field over unencrypted HTTP

**Threat Intelligence** â€” 2 checks

32. PhishTank confirmed phishing (community-verified)
33. PhishTank reported phishing (awaiting verification, 3 h cache)

---

## Code Layout

```
server/
â”œâ”€â”€ cmd/urlvet/           entry point â€” init, router setup, graceful shutdown
â”œâ”€â”€ internal/
â”‚   â”œâ”€â”€ analyzer/
â”‚   â”‚   â”œâ”€â”€ analyze.go      task registration, cache integration
â”‚   â”‚   â”œâ”€â”€ runner.go       goroutine runner with panic recovery
â”‚   â”‚   â”œâ”€â”€ tasks.go        18 task implementations
â”‚   â”‚   â””â”€â”€ result.go       score aggregation, verdict assignment
â”‚   â”œâ”€â”€ handler/
â”‚   â”‚   â”œâ”€â”€ router.go       Gin router, middleware wiring
â”‚   â”‚   â”œâ”€â”€ analyze.go      /api/v1/analyze handler
â”‚   â”‚   â””â”€â”€ middleware/     rate limiter, auth, Prometheus, request logger
â”‚   â”œâ”€â”€ service/
â”‚   â”‚   â”œâ”€â”€ checks/         18 individual analyzer implementations
â”‚   â”‚   â”œâ”€â”€ screenshot/     headless Chrome integration (chromedp)
â”‚   â”‚   â”œâ”€â”€ cache/          Valkey client wrapper
â”‚   â”‚   â”œâ”€â”€ threatfeeds/    PhishTank client
â”‚   â”‚   â””â”€â”€ typosquat/      brand similarity engine
â”‚   â”œâ”€â”€ logger/             centralized slog-based logger (colors in DEV, JSON in prod)
â”‚   â””â”€â”€ admintoken/         admin JWT issuance and verification
web/
â”œâ”€â”€ website/                SvelteKit UI
â””â”€â”€ chrome-extension/       browser extension
docker/
â”œâ”€â”€ dev/                    dev Compose (hot reload, exposed ports)
â””â”€â”€ prod/                   prod Compose (optimized builds, restart policies)
docs/                       API reference, setup guide, architecture, security
```

---

## Deployment

Two fully separated Docker Compose stacks share the same image definitions but differ in configuration:

| | Dev | Prod |
|---|---|---|
| Backend | Air hot-reload, source mounted as volume | Compiled binary in distroless image |
| Frontend | Vite dev server `:5173` | Static build served by Nginx `:3000` |
| Chrome | Same `chromedp/headless-shell` image | Same |
| Valkey | Port exposed (`:6379`) for local inspection | Port not exposed; internal only |
| ENV | `ENV=DEV` â€” colored logs, debug endpoints | `ENV=PROD` â€” JSON logs, info level only |

Start everything with one command:

```bash
make start       # production stack
make dev         # development stack (hot reload)
```

See [docs/setup.md](setup.md) for full setup instructions including `.env` configuration.
