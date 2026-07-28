package domain

import (
	"time"
)

// AgentRow representa estritamente uma linha na tabela de agentes no banco de dados do C&C.
type AgentRow struct {
	UUID       string    `db:"uuid"`
	CreatedAt  time.Time `db:"created_at"`
	LastSeenAt time.Time `db:"last_seen_at"`
}

// NewAgentRow inicializa uma nova linha de agente recém-registrado para ser salva no banco.
func NewAgentRow(uuid string) AgentRow {
	now := time.Now()
	return AgentRow{
		UUID:       uuid,
		CreatedAt:  now,
		LastSeenAt: now,
	}
}