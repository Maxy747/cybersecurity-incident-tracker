# CyberGuard 🛡️

> A cybersecurity incident management and response platform for Security Operations Center (SOC) teams.

---

## Overview

CyberGuard is a full-featured backend platform built with **Golang + Gin + PostgreSQL**. It enables SOC teams to monitor, track, investigate, and resolve cybersecurity incidents through a secure REST API.

### Key capabilities

| Feature | Details |
|---|---|
| **Incident Management** | Full CRUD with categories, severity levels, and status workflow |
| **Role-Based Access** | Admin / Analyst / Viewer with enforced permissions per endpoint |
| **JWT Authentication** | HS256 tokens, 24-hour expiry, bcrypt password hashing |
| **Investigation Notes** | Timestamped analyst updates per incident |
| **Audit Logging** | Every mutating action recorded with user identity |
| **Dashboard Analytics** | Status counts, severity distribution, time-series trends |
| **Search & Filter** | Filter incidents by title, category, severity, status, assignee |
| **Swagger UI** | Interactive API docs at `/swagger/index.html` |

---

## Tech Stack

- **Language:** Go 1.22
- **Framework:** Gin
- **Database:** PostgreSQL 16
- **Auth:** JWT (golang-jwt/jwt v5) + bcrypt
- **Docs:** Swagger (swaggo)
- **Container:** Docker + Docker Compose

---

## Quick Start

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) & Docker Compose
- Git

### 1. Clone the repository

```bash
git clone https://github.com/Maxy747/cybersecurity-incident-tracker.git
cd cybersecurity-incident-tracker
```

### 2. Configure environment

```bash
cp .env.example .env
# Edit .env — at minimum change JWT_SECRET
openssl rand -hex 32   # use this output as your JWT_SECRET
```

### 3. Start the stack

```bash
docker compose up --build -d
```

The API will be available at **http://localhost:8080** once PostgreSQL is healthy and the schema has migrated.

### 4. Open Swagger UI

```
http://localhost:8080/swagger/index.html
```

---

## API Reference

### Authentication

| Method | Endpoint | Access | Description |
|---|---|---|---|
| POST | `/api/auth/register` | Public | Register a new user |
| POST | `/api/auth/login` | Public | Login and receive JWT |
| GET | `/api/auth/profile` | Authenticated | Get your profile |

> **First registered user is automatically made Admin.**

### Incidents

| Method | Endpoint | Access | Description |
|---|---|---|---|
| GET | `/api/incidents` | All | List / search incidents |
| POST | `/api/incidents` | Admin, Analyst | Create incident |
| GET | `/api/incidents/:id` | All | Get single incident |
| PUT | `/api/incidents/:id` | Admin, Analyst | Update incident |
| DELETE | `/api/incidents/:id` | Admin | Delete incident |

**Filter query params:** `title`, `category`, `severity`, `status`, `assigned_to`

### Investigation Notes

| Method | Endpoint | Access | Description |
|---|---|---|---|
| POST | `/api/incidents/:id/notes` | Admin, Analyst | Add a note |
| GET | `/api/incidents/:id/notes` | All | List notes |

### Dashboard

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/dashboard/stats` | Incident counts by status |
| GET | `/api/dashboard/severity` | Counts by severity level |
| GET | `/api/dashboard/trends` | Incidents per day (`?days=30`) |

### Audit Logs

| Method | Endpoint | Access |
|---|---|---|
| GET | `/api/audit-logs` | Admin only |

---

## Incident Model

### Categories
`Malware` · `Phishing` · `DDoS` · `Unauthorized Access` · `Data Breach` · `Insider Threat` · `Brute Force` · `Network Attack`

### Severity Levels
`Critical` · `High` · `Medium` · `Low` · `Informational`

### Status Workflow
```
Open → Investigating → Resolved → Closed
```

---

## Roles & Permissions

| Action | Admin | Analyst | Viewer |
|---|---|---|---|
| Register users | ✅ | ❌ | ❌ |
| Create incident | ✅ | ✅ | ❌ |
| Update any incident | ✅ | own only | ❌ |
| Delete incident | ✅ | ❌ | ❌ |
| Add notes | ✅ | ✅ | ❌ |
| View incidents/notes/dashboard | ✅ | ✅ | ✅ |
| View audit logs | ✅ | ❌ | ❌ |

---

## Development Commands

```bash
make up        # Start all services
make down      # Stop all services
make logs      # Tail API logs
make swag      # Regenerate Swagger docs
make tidy      # go mod tidy
make build     # Compile binary locally
make jwt-secret # Generate a secure JWT secret
```

---

## Project Structure

```
.
├── cmd/main.go                  # Application entry point + Swagger metadata
├── internal/
│   ├── config/                  # Environment config loader
│   ├── database/                # Connection pool + schema migrations
│   ├── handlers/                # HTTP request handlers (auth, incidents, notes, dashboard, audit)
│   ├── middleware/              # JWT auth + RBAC middleware
│   ├── models/                  # Data structs, request/response types, constants
│   ├── repository/              # Parameterized SQL queries (user, incident, note, audit)
│   └── routes/                  # Gin route registration
├── docs/                        # Auto-generated Swagger JSON/YAML
├── Dockerfile                   # Multi-stage build
├── docker-compose.yml           # PostgreSQL + API services
├── .env.example                 # Environment variable template
└── Makefile                     # Developer shortcuts
```

---

## Security

- Passwords hashed with **bcrypt** (cost 12) — never stored in plain text
- JWT signed with **HS256** — validated on every protected request
- **Parameterized queries** throughout — SQL injection safe
- **RBAC middleware** enforces role requirements per route
- `.env` is `.gitignore`d — secrets never committed

---

## Resume Description

> Developed **CyberGuard**, a cybersecurity incident management platform using Golang, Gin, PostgreSQL, JWT authentication, Docker, and REST APIs. Implemented role-based access control (Admin/Analyst/Viewer), bcrypt password hashing, audit logging, incident lifecycle tracking, investigation workflows, dashboard analytics, and Swagger documentation following enterprise cybersecurity best practices.
