package domain

import (
	"time"

	"github.com/google/uuid"
)

// JobRow representa exatamente uma linha na tabela de jobs do banco de dados do C&C.
// Utilizamos tags 'db' para que a biblioteca sqlx faça o mapeamento automático.
type JobRow struct {
	ID         string     `db:"id"`
	AgentID    string     `db:"agent_id"` // Faz referência direta ao UUID do Agente
	CreatedAt  time.Time  `db:"created_at"`
	ExecutedAt *time.Time `db:"executed_at"` // Ponteiro para permitir valor nulo no SQL
	Command    string     `db:"command"`
	Args       []string   `db:"args"` // Será salvo como JSONB no Postgres/MySQL
	Output     *string    `db:"output"` // Ponteiro para permitir valor nulo no SQL
}

// NewJobRow inicializa uma nova linha de comando pronta para ser salva no banco.
func NewJobRow(agentID, command string, args []string) JobRow {
	return JobRow{
		ID:         uuid.NewString(), // Gera o ID automaticamente, como exigido pelo seu teste
		AgentID:    agentID,
		CreatedAt:  time.Now(),
		ExecutedAt: nil, // Nasce sem data de execução
		Command:    command,
		Args:       args,
		Output:     nil, // Nasce sem saída do terminal
	}
}