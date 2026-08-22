# Glossary

Domain-specific terms, acronyms, and The Url Mentalist-internal concepts used throughout the codebase and documentation.

---

## A

**Analyzer**
The The Url Mentalist component that orchestrates all checks for a single URL. It launches 18 goroutines via `sync.WaitGroup`, collects their outputs, and feeds results to the scorer. See `server/internal/analyzer/`.

**AGPL-3.0**
GNU Affero General Public License v3. The open-source license The Url Mentalist is released under. Key clause: any modified version run over a network must make its source code available to users.

---

## C

**Cache hit / Cache miss**
When a URL is submitted, Valkey is checked first. A *hit* means a full result is already stored and is returned immediately (sub-millisecond). A *miss* means no cached result exists and the full analysis pipeline runs.

**chromedp**
A Go library that controls a headless Chrome browser over the Chrome DevTools Protocol (CDP) WebSocket. The Url Mentalist uses it to take page screenshots and for content fetching. Container: `urlvet-chrome` on port `:9222`.

**Combo-squatting**
A form of domain abuse where a brand name is combined with extra words to make the domain look legitimate, e.g. `paypal-login.com` or `apple-support.net`. The Url Mentalist checks the target domain against 500+ known brands.

**Content Analysis**
A category of 8 checks that fetch and parse the target page's HTML to look for phishing indicators â€” login forms, payment fields, hidden iframes, external form targets, and brand-domain mismatches.

**CT log (Certificate Transparency log)**
A public, append-only ledger of TLS certificates. Certificates without an embedded SCT (Signed Certificate Timestamp) may have been issued without public logging, which is a weak indicator of suspicious infrastructure.

---

## D

**DNS (Domain Name System)**
The internet's naming system. The Url Mentalist checks NS records (authoritative name servers), MX records (mail exchange), and IP resolution as part of domain infrastructure validation.

**DNSSEC (DNS Security Extensions)**
Cryptographic signatures on DNS responses that allow resolvers to verify authenticity. Absence of DNSSEC is a mild risk signal; many legitimate domains also lack it.

**Domain age**
How long ago a domain was first registered, retrieved via WHOIS/RDAP. Newly registered domains (days to weeks old) are a strong phishing indicator because attackers spin up fresh domains for each campaign.

**Domain Intelligence**
A category of 6 checks covering domain rank, TLD classification, domain age, DNSSEC, Shannon entropy, and typosquatting detection.

---

## E

**eTLD (effective Top-Level Domain)**
The public suffix of a domain as defined by the [Public Suffix List](https://publicsuffix.org/). For `mail.google.co.uk`, the eTLD is `co.uk` and the eTLD+1 (registrable domain) is `google.co.uk`. The Url Mentalist uses `golang.org/x/net/publicsuffix` for this.

**eTLD+1**
The registrable portion of a domain â€” the eTLD plus one additional label. Used to determine the true owner of a domain and to compare redirect chains for cross-domain hops.

**Entropy (Shannon entropy)**
A measure of randomness in a string, calculated as `âˆ’Î£ p(x) Â· logâ‚‚ p(x)`. High entropy in a domain name (e.g. `xq2k9f.net`) suggests it was algorithmically generated (DGA â€” Domain Generation Algorithm), which is common in malware C2 infrastructure.

---

## F

**False positive**
A legitimate URL that The Url Mentalist incorrectly scores as Suspicious or Risky. Common sources: newly registered legitimate domains, domains using IP-based hosting, or sites with unusual TLDs.

**Final URL**
The URL the browser lands on after following all redirects. The Url Mentalist records the full redirect chain and compares the final URL's domain to the original to detect cross-domain redirects.

---

## G

**Goroutine**
A lightweight concurrent execution unit in Go, cheaper than OS threads. The Url Mentalist launches one goroutine per analyzer task (18 total) via `sync.WaitGroup`, letting all checks run in parallel.

---

## H

**Homograph attack**
A phishing technique that uses visually similar Unicode characters to spoof domain names. Example: `Ð°pple.com` where `Ð°` is Cyrillic U+0430, not Latin `a`. Also called an IDN homograph attack.

**HSTS (HTTP Strict Transport Security)**
An HTTP response header (`Strict-Transport-Security`) that instructs browsers to always use HTTPS for a domain. Its absence on an HTTPS site is a weak risk signal.

**HTTP Combined Task** (`httpCombinedTask`)
An internal The Url Mentalist optimization. Instead of making separate HTTP requests for redirects, HSTS, and status code checks, a single HTTP request is shared across all three. Falls back to individual checks if the combined request fails.

---

## I

**IDN (Internationalized Domain Name)**
A domain name containing non-ASCII Unicode characters, encoded using Punycode for DNS transport (e.g. `mÃ¼nchen.de` â†’ `xn--mnchen-3ya.de`). The Url Mentalist flags IDN domains as they are frequently used in homograph attacks.

**iframe (hidden)**
An HTML `<iframe>` element with zero dimensions or `display:none`. Hidden iframes are used for clickjacking, credential theft, or silently loading malicious content.

**IP address as hostname**
Using a raw IPv4 or IPv6 address directly as the URL host (e.g. `http://192.168.1.1/login`). Legitimate services almost never use bare IP addresses; phishers use them to avoid domain registration and takedowns.

---

## L

**LRU (Least Recently Used)**
The cache eviction policy used by Valkey in The Url Mentalist. When the cache reaches its memory limit (`maxmemory`), the least recently accessed keys are evicted first. Configured as `allkeys-lru`.

---

## M

**MX record**
A DNS record that specifies the mail server for a domain. Missing or invalid MX records combined with other signals can indicate a disposable or fraudulent domain.

**Middleware**
In the Gin HTTP framework context, middleware is a function that runs before or after route handlers. The Url Mentalist middleware: request logger, rate limiter, CORS, URL length validator, Bearer auth, and Prometheus metrics.

---

## N

**NS record**
A DNS record that delegates a domain to its authoritative name servers. Missing or non-resolving NS records indicate a domain that is parked, expired, or misconfigured.

**Normalize (URL normalization)**
The process of bringing a URL to a canonical form before analysis: inferring a scheme (`https://` if missing), lowercasing the hostname, and percent-encoding special characters. Ensures consistent cache keys and check behavior.

---

## P

**PhishTank**
A community-driven database of verified phishing URLs maintained by Cisco Talos. The Url Mentalist queries the PhishTank API to check whether a URL has been confirmed or reported as phishing. Results are cached for 3 hours.

**Pipeline**
The end-to-end flow of a URL through The Url Mentalist: normalize â†’ cache check â†’ parallel analyzers â†’ score aggregation â†’ cache store â†’ response. Visualized in `assets/pipeline.png`.

**Punycode**
An ASCII-compatible encoding for Unicode domain labels, defined in RFC 3492. `xn--` is the ACE prefix that identifies a Punycode-encoded label. The Url Mentalist detects domains with Punycode labels as a potential homograph risk.

---

## R

**RDAP (Registration Data Access Protocol)**
The modern replacement for WHOIS (RFC 7483). Returns structured JSON for domain registration data including creation date, registrar, and name servers. The Url Mentalist uses RDAP when available, falls back to WHOIS.

**Reason**
A labeled explanation string emitted by each check â€” classified as `good`, `bad`, or `neutral`. Every signal that fires produces a reason, making the final verdict fully explainable. Stored in `GoodReasons`, `BadReasons`, `NeutralReasons` in the response.

**Redirect chain**
The ordered list of URLs traversed when following HTTP redirects. The Url Mentalist records each hop and flags chains that cross domain boundaries or are excessively long.

**Risk score**
A 0â€“100 accumulator for negative signals. Each bad check adds weighted points. Combined with the trust score in the final formula. Clamped to 100 before scoring.

---

## S

**SCT (Signed Certificate Timestamp)**
A cryptographic token embedded in a TLS certificate (or delivered via TLS extension/OCSP stapling) proving it was logged to a CT log before issuance. Required by major browsers for public certificates issued after 2018.

**Shannon entropy** â†’ see *Entropy*

**SLD (Second-Level Domain)**
The label directly to the left of the TLD. In `mail.google.com`, the SLD is `google`. The Url Mentalist's typosquatting check compares the SLD against known brand names.

**SSRF (Server-Side Request Forgery)**
An attack where a server is tricked into making HTTP requests to internal or private network addresses on behalf of an attacker. The Url Mentalist blocks SSRF by rejecting URLs that resolve to private IP ranges (RFC 1918, loopback, link-local) before fetching.

**Signal**
A single, binary or graded indicator produced by one check â€” e.g. "domain registered 3 days ago" or "HSTS present". The Url Mentalist produces 33 signals across 7 categories. Each signal maps to one or more reasons.

**Suspicious** â†’ see *Verdict*

---

## T

**Task**
The internal unit of work in the analyzer. Each task implements the `Task` interface (`Name() string`, `Run(*Input, *Output) error`) and corresponds to one analyzer goroutine. 18 tasks are registered in `analyze.go`.

**Threat Intelligence**
A category of 2 checks that query external databases of known-bad URLs. Currently: PhishTank confirmed and PhishTank reported phishing entries.

**TLD (Top-Level Domain)**
The rightmost label of a domain (`.com`, `.net`, `.gov`). The Url Mentalist classifies TLDs as trusted (`.gov`, `.edu`, `.mil`), risky (commonly abused ccTLDs and gTLDs), or neutral.

**TLS (Transport Layer Security)**
The cryptographic protocol that secures HTTPS connections. The Url Mentalist performs a full TLS handshake to inspect the certificate chain â€” validity, expiry, issuer, and CT log status.

**Trust score**
A 0â€“100 accumulator for positive signals. Each good check adds weighted points. Combined with the risk score in the final formula. Clamped to 100 before scoring.

**Typosquatting**
Registering domains that are intentional misspellings of well-known brands to capture mistyped traffic or trick users â€” e.g. `googie.com`, `paypa1.com`. The Url Mentalist checks the target domain against 500+ brand names using edit-distance and visual-similarity heuristics.

---

## U

**URL shortener**
A service that maps a short URL (e.g. `bit.ly/abc`) to a longer destination URL. Shorteners are used by phishers to obscure the true target. The Url Mentalist maintains a list of known shortener domains and flags their use as a risk signal.

---

## V

**Valkey**
A Redis-compatible, open-source key-value store. The Url Mentalist uses it as an LRU cache for full analysis results (24 h TTL), content analysis (configurable TTL), and threat feed lookups (3 h TTL). Data is persisted to a Docker volume.

**Verdict**
The human-readable classification assigned to a URL based on its final score:

| Score | Verdict |
|---|---|
| â‰¥ 65 | **Safe** |
| 30 â€“ 64 | **Suspicious** |
| < 30 | **Risky** |

---

## W

**WHOIS**
A text-based protocol (RFC 3912) for querying domain registration data. Returns domain age, registrar, name servers, and contact information. Being replaced by RDAP. The Url Mentalist uses WHOIS as a fallback when RDAP data is unavailable.

---

## Scoring Reference

```
finalScore = clamp(50 + (trustScore âˆ’ riskScore) Ã— 0.5, 0, 100)
```

| Symbol | Meaning |
|---|---|
| `50` | Neutral baseline â€” default for a URL with no signals |
| `trustScore` | Sum of positive signal weights, clamped 0â€“100 |
| `riskScore` | Sum of negative signal weights, clamped 0â€“100 |
| `Ã— 0.5` | Dampening factor so neither side dominates alone |
| `clamp(â€¦, 0, 100)` | Final score is always in [0, 100] |
