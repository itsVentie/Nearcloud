# Nearcloud 

Local-first, zero-configuration peer-to-peer (P2P) file sharing platform bypassing the cloud entirely. Fast, secure, and works offline.

Nearcloud establishes direct, high-throughput TLS connections over local network interfaces (Wi-Fi or Ethernet). Discovery is entirely decentralized, removing the requirement for static IP allocation, internet access, or external orchestration.


```

```
                       [ LAN / Wi-Fi Network ]
                                  │
           ┌──────────────────────┴──────────────────────┐
           ▼ (mDNS Discovery)                            ▼ (mDNS Discovery)
 ┌───────────────────┐                         ┌───────────────────┐
 │   Nearcloud (A)   │                         │   Nearcloud (B)   │
 │ ┌───────────────┐ │                         │ ┌───────────────┐ │
 │ │   Wails GUI   │ │                         │ │   Wails GUI   │ │
 │ └───────┬───────┘ │                         │ └───────▲───────┘ │
 │         │ (File)  │                         │         │ (File)  │
 │ ┌───────▼───────┐ │   TLS Tunnel over TCP   │ ┌───────┴───────┐ │
 │ │  Go Engine    ├─┼────────────────────────►│  Go Engine    │ │
 │ └───────────────┘ │      (Direct Stream)    │ └───────────────┘ │
 └───────────────────┘                         └───────────────────┘

```

```

## Features & Architecture

### 1. Zero-Config Peer Discovery (mDNS)
Devices continuously advertise and browse for the specialized service `_nearcloud._tcp` in the `local.` multicast domain using Multicast DNS (RFC 6762). No central registry, no setup.

### 2. High-Speed Zero-Copy Transfers
Files are read sequentially and piped directly into the network socket stream using a pooled memory architecture. This guarantees near-wire gigabit speeds while maintaining a minimal and stable RAM footprint.

### 3. Mutual TLS Security
All traffic is encrypted in transit using TLS 1.3 with ephemeral keys. Self-signed certificates are generated dynamically on startup, protecting payloads from eavesdropping or tampering even on untrusted public Wi-Fi networks.

---

## Roadmap

### Phase 1: Core Network Engine (CLI MVP)
- [ ] Implement mDNS-based peer discovery and service registration.
- [ ] Develop a lightweight TCP server/client architecture.
- [ ] Implement TLS 1.3 with dynamic self-signed certificate generation.
- [ ] Establish a basic framing protocol for file metadata exchange.
- [ ] Add a CLI interactive CLI prompt for transfer confirmation.

### Phase 2: Desktop GUI & Integration
- [ ] Initialize Wails v2 project structure with TS + React.
- [ ] Connect Go backend handlers to frontend state.
- [ ] Implement a visual "radar" showing available local peers.
- [ ] Design drag-and-drop file staging area.
- [ ] Add native OS system tray integration and notifications.

### Phase 3: Performance & Robustness
- [ ] Support directory transfers (recursive compression/archiving on the fly).
- [ ] Implement transfer pause, resume, and cancellation mechanics.
- [ ] Optimize buffer pools (`sync.Pool`) for multi-gigabyte files.
- [ ] Automatic fallback to alternative network interfaces (e.g., ethernet over Wi-Fi).

---

## Tech Stack

* **Core Daemon:** Go (optimized networking, zero-allocation pooling, cross-platform)
* **Frontend UI:** TypeScript, React, TailwindCSS
* **Application Wrapper:** Wails v2 (native OS webviews, no Chromium memory overhead)

## Development

### Prerequisites
* Go 1.22+
* Node.js 20+
* OS-specific C compilers (GCC on Linux, MSVC on Windows)

### Build CLI (MVP)
```bash
go build -o nearcloud ./cmd/cli

```

### Run Desktop App

```bash
go install [github.com/wailsapp/wails/v2/cmd/wails@latest](https://github.com/wailsapp/wails/v2/cmd/wails@latest)
wails dev
```