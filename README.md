# Bloc Bootstrap Server

A **QUIC-based bootstrap server** for the Bloc blockchain network.  
Its sole responsibility is **peer discovery and coordination** — not consensus, validation, or block storage.

This server allows nodes behind NAT/firewalls to:
- register themselves
- discover other peers
- prepare for UDP hole punching

---

## 🚀 Features

- QUIC (UDP + TLS) based transport
- Secure peer registration
- Stateless protocol design (no blockchain state)
- JSON-based message envelope
- Graceful session handling
- Designed for NAT traversal workflows

---

## 🧠 Architecture Overview

Node ──QUIC──▶ Bootstrap Server
│ │
│ register │
│────────────────▶│
│ │ stores observed addr
│ │
│◀──────── peer list / responses


- **One QUIC session per peer**
- **One stream per request**
- Server does not retain long-term state (yet)

---

## 📡 Network Details

- **Protocol:** QUIC
- **Transport:** UDP
- **Default Port:** `4242`
- **TLS:** Self-signed certificates (for now)

---

## 📂 Project Structure

bootstrap/
├── certificate/
│ ├── server.crt
│ └── server.key
├── types/
│ ├── message.go
│ └── peer.go
├── utils/
│ └── logger/
├── main.go
├── peers.json
└── README.md


---

## 🧾 Message Protocol

All messages follow a **common envelope**:

```json
{
  "version": 1,
  "header": {
    "application/json": "true"
  },
  "type": "register",
  "length": 17,
  "payload": { ... }
}

| Type       | Description                      |
| ---------- | -------------------------------- |
| `register` | Register node with bootstrap     |
| `ping`     | Health check                     |
| `punch`    | Hole punching coordination (WIP) |

🧩 Example: Register Message

{
  "version": 1,
  "header": {
    "application/json": "true"
  },
  "type": "register",
  "length": 17,
  "payload": {
    "id": "node-123"
  }
}
