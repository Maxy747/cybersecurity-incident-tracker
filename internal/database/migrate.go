package database

import (
	"log"
)

// Migrate runs the schema DDL statements in order.
// Uses CREATE TABLE IF NOT EXISTS so it is idempotent.
func Migrate() {
	statements := []string{
		// ── Users ────────────────────────────────────────────────────────
		`CREATE TABLE IF NOT EXISTS users (
			id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name          VARCHAR(150) NOT NULL,
			email         VARCHAR(255) NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			role          VARCHAR(20)  NOT NULL DEFAULT 'viewer',
			created_at    TIMESTAMP    NOT NULL DEFAULT NOW(),
			updated_at    TIMESTAMP    NOT NULL DEFAULT NOW()
		);`,

		// ── Incidents ────────────────────────────────────────────────────
		`CREATE TABLE IF NOT EXISTS incidents (
			id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			title       VARCHAR(255) NOT NULL,
			description TEXT         NOT NULL,
			category    VARCHAR(50)  NOT NULL,
			severity    VARCHAR(20)  NOT NULL,
			status      VARCHAR(20)  NOT NULL DEFAULT 'Open',
			created_by  UUID         NOT NULL REFERENCES users(id),
			assigned_to UUID         REFERENCES users(id),
			created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
			updated_at  TIMESTAMP    NOT NULL DEFAULT NOW()
		);`,

		// ── Investigation Notes ──────────────────────────────────────────
		`CREATE TABLE IF NOT EXISTS investigation_notes (
			id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			incident_id UUID      NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
			author_id   UUID      NOT NULL REFERENCES users(id),
			content     TEXT      NOT NULL,
			created_at  TIMESTAMP NOT NULL DEFAULT NOW()
		);`,

		// ── Audit Logs ───────────────────────────────────────────────────
		`CREATE TABLE IF NOT EXISTS audit_logs (
			id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id       UUID         NOT NULL REFERENCES users(id),
			action        TEXT         NOT NULL,
			resource_type VARCHAR(50),
			resource_id   UUID,
			timestamp     TIMESTAMP    NOT NULL DEFAULT NOW()
		);`,
	}

	for _, stmt := range statements {
		if _, err := DB.Exec(stmt); err != nil {
			log.Fatalf("[migrate] error running migration: %v\nStatement: %s", err, stmt)
		}
	}
	log.Println("[migrate] schema is up to date")
}
