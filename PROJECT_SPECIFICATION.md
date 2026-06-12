# CyberGuard: Security Incident Management & Response Platform

## Project Overview

**CyberGuard** is a cybersecurity-focused incident management platform designed to help Security Operations Center (SOC) teams monitor, track, investigate, and manage cybersecurity incidents throughout their lifecycle.

The system serves as a centralized platform where security analysts can record incidents such as brute-force attacks, malware infections, phishing attempts, suspicious logins, unauthorized access attempts, and network anomalies. It provides role-based access control, audit logging, investigation tracking, incident prioritization, and reporting capabilities.

The project simulates a real-world enterprise SOC environment and demonstrates backend engineering, secure API development, database management, authentication, authorization, and cybersecurity best practices.

---

# Problem Statement

Organizations often struggle with managing security incidents efficiently. Security events are frequently tracked using spreadsheets, emails, or disconnected tools, making investigation and response difficult.

CyberGuard solves this problem by providing:

* Centralized incident management
* Secure user authentication
* Role-based authorization
* Investigation documentation
* Audit trails
* Security reporting
* Incident lifecycle tracking

---

# Project Objectives

The primary objectives of CyberGuard are:

## Security Management

Provide a platform to manage cybersecurity incidents from detection to resolution.

## User Management

Allow different users to access the system according to their responsibilities.

## Investigation Tracking

Maintain detailed records of investigation progress and findings.

## Accountability

Track every user action through audit logging.

## Reporting

Generate incident statistics and security insights.

---

# Real-World Use Case

Imagine a company detects:

```text
150 failed login attempts
from IP: 185.76.43.22
within 5 minutes.
```

A security analyst logs into CyberGuard and creates:

```text
Incident Title:
Possible Brute Force Attack

Severity:
High

Status:
Open
```

The analyst then:

* Assigns the incident
* Adds investigation notes
* Uploads evidence
* Updates incident status
* Resolves the issue

Managers can later review:

* Timeline
* Actions taken
* Responsible personnel
* Final resolution

---

# Core Features

## 1. User Authentication System

Users can:

* Register
* Login
* Logout

Authentication uses JWT (JSON Web Tokens).

Every API request requires a valid token.

Example:

```http
Authorization: Bearer <jwt-token>
```

---

## 2. Role-Based Access Control (RBAC)

### Admin

Can:

* Create users
* Delete users
* Manage incidents
* View audit logs
* Access all reports

### Security Analyst

Can:

* Create incidents
* Update incidents
* Add investigation notes
* Assign incidents

Cannot:

* Manage users

### Viewer

Can:

* View incidents
* View reports

Cannot:

* Modify data

---

## 3. Incident Management

The core module.

Users can create incidents containing:

* Title
* Description
* Category
* Severity
* Assigned Analyst

### Incident Categories

* Malware
* Phishing
* DDoS
* Unauthorized Access
* Data Breach
* Insider Threat
* Brute Force
* Network Attack

### Severity Levels

* Critical
* High
* Medium
* Low
* Informational

### Status Workflow

```text
Open
↓
Investigating
↓
Resolved
↓
Closed
```

---

## 4. Investigation Notes

Each incident can contain multiple investigation updates.

Example:

```text
Analyst:
Mazin

Note:
Blocked attacking IP from firewall.
```

Each note stores:

* Author
* Timestamp
* Content

---

## 5. Audit Logging System

Every action performed is logged.

Examples:

* User Login
* Incident Created
* Severity Modified
* Status Changed
* User Deleted

Audit record example:

```text
2026-06-15 10:22

User:
Mazin

Action:
Changed Incident #24
Severity: High → Critical
```

---

## 6. Security Dashboard

Provides statistics such as:

### Incident Count

* Total Incidents
* Open Incidents
* Resolved Incidents

### Severity Distribution

```text
Critical: 5
High: 12
Medium: 18
Low: 25
```

### Incident Trends

* Incidents Per Day
* Incidents Per Week
* Incidents Per Month

### Analyst Activity

* Incidents Handled
* Notes Added
* Investigations Completed

---

## 7. Search and Filtering

Search by:

* Incident ID
* Title
* Analyst
* Category
* Severity
* Status

Example:

```text
Show all Critical incidents assigned to Mazin.
```

---

## 8. Reporting System

Generate reports such as:

### Monthly Security Report

Includes:

* Total Incidents
* Average Resolution Time
* Most Common Threat
* Critical Incidents
* Analyst Performance

---

# Database Architecture

## Users Table

| Field         | Type      |
| ------------- | --------- |
| id            | UUID      |
| name          | VARCHAR   |
| email         | VARCHAR   |
| password_hash | TEXT      |
| role          | VARCHAR   |
| created_at    | TIMESTAMP |
| updated_at    | TIMESTAMP |

---

## Incidents Table

| Field       | Type      |
| ----------- | --------- |
| id          | UUID      |
| title       | VARCHAR   |
| description | TEXT      |
| category    | VARCHAR   |
| severity    | VARCHAR   |
| status      | VARCHAR   |
| created_by  | UUID      |
| assigned_to | UUID      |
| created_at  | TIMESTAMP |
| updated_at  | TIMESTAMP |

---

## Investigation Notes Table

| Field       | Type      |
| ----------- | --------- |
| id          | UUID      |
| incident_id | UUID      |
| author_id   | UUID      |
| content     | TEXT      |
| created_at  | TIMESTAMP |

---

## Audit Logs Table

| Field         | Type      |
| ------------- | --------- |
| id            | UUID      |
| user_id       | UUID      |
| action        | TEXT      |
| resource_type | VARCHAR   |
| resource_id   | UUID      |
| timestamp     | TIMESTAMP |

---

# API Architecture

## Authentication APIs

```http
POST /api/auth/register
POST /api/auth/login
GET /api/auth/profile
```

## Incident APIs

```http
POST /api/incidents
GET /api/incidents
GET /api/incidents/:id
PUT /api/incidents/:id
DELETE /api/incidents/:id
```

## Notes APIs

```http
POST /api/incidents/:id/notes
GET /api/incidents/:id/notes
```

## Dashboard APIs

```http
GET /api/dashboard/stats
GET /api/dashboard/severity
GET /api/dashboard/trends
```

## Audit APIs

```http
GET /api/audit-logs
```

(Admin only)

---

# Security Features

## Password Hashing

Uses bcrypt.

Passwords are never stored in plain text.

## JWT Authentication

Secure token-based authentication.

## Authorization Middleware

Protects sensitive routes based on user roles.

## Input Validation

Prevents malformed requests and invalid data.

## SQL Injection Protection

Uses parameterized queries and ORM safety mechanisms.

---

# Technology Stack

## Backend

* Golang
* Gin Framework

## Database

* PostgreSQL

## Authentication

* JWT
* bcrypt

## Documentation

* Swagger

## Containerization

* Docker
* Docker Compose

## Version Control

* Git
* GitHub

---

# Future Enhancements

## Threat Intelligence Integration

Check suspicious IP addresses against threat intelligence feeds.

## SIEM Integration

Integrate with:

* Wazuh
* Splunk
* ELK Stack

## Email Notifications

Notify analysts when incidents are assigned.

## AI Incident Classification

Use machine learning to predict:

* Severity
* Threat Type
* Priority

---

# Resume Description

Developed CyberGuard, a cybersecurity incident management platform using Golang, Gin, PostgreSQL, JWT authentication, Docker, and REST APIs. Implemented role-based access control, audit logging, incident tracking, investigation workflows, dashboard analytics, and secure API architecture following enterprise cybersecurity best practices.
