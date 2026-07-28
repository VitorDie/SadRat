package repository

import (
	"errors"
	"sync"
	"fmt"

	"github.com/VitorDie/SadRat/internal/domain"
)

// ServerDBPlainDataInMemory é o nosso banco de dados falso para testes ultrarrápidos.
type ServerDBPlainDataInMemory struct {
	mu     sync.RWMutex // Protege os mapas contra acesso concorrente (evita Data Races)
	agents map[string]domain.AgentRow // Simula a tabela 'agents'
	jobs   map[string]domain.JobRow   // Simula a tabela 'jobs'
}

// NewServerDBPlainDataInMemory cria uma nova instância do banco em memória.
func NewServerDBPlainDataInMemory() *ServerDBPlainDataInMemory {
	return &ServerDBPlainDataInMemory{
		agents: make(map[string]domain.AgentRow),
		jobs:   make(map[string]domain.JobRow),
	}
}

// SaveAgent insere uma nova linha na tabela de agentes.
func (r *ServerDBPlainDataInMemory) SaveAgent(agent domain.AgentRow) error {
	r.mu.Lock()         // Tranca o cofre para escrita
	defer r.mu.Unlock() // Garante que será destrancado ao final
	
	r.agents[agent.UUID] = agent
	return nil
}

// SaveJob insere uma nova linha na tabela de jobs.
func (r *ServerDBPlainDataInMemory) SaveJob(job domain.JobRow) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.jobs[job.ID] = job
	return nil
}

// GetAgent busca um agente pelo seu UUID.
func (r *ServerDBPlainDataInMemory) GetAgent(uuid string) (domain.AgentRow, error) {
	r.mu.RLock()         // Tranca o cofre apenas para leitura (mais rápido)
	defer r.mu.RUnlock()
	
	agent, exists := r.agents[uuid]
	if !exists {
		return domain.AgentRow{}, errors.New("agente não encontrado")
	}
	return agent, nil
}

// GetJob busca um comando pelo seu ID.
func (r *ServerDBPlainDataInMemory) GetJob(id string) (domain.JobRow, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	job, exists := r.jobs[id]
	if !exists {
		return domain.JobRow{}, errors.New("job não encontrado")
	}
	return job, nil
}

func (r *ServerDBPlainDataInMemory) GetJobForAgent(agentID string) (domain.JobRow, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	fmt.Printf("[DEBUG DB] Buscando job para o Agente: '%s'\n", agentID)
	fmt.Printf("[DEBUG DB] Total de jobs no banco: %d\n", len(r.jobs))

	for _, job := range r.jobs {
		fmt.Printf("[DEBUG DB] Analisando Job ID: %s | AgentID do Job: '%s' | ExecutedAt: %v\n", job.ID, job.AgentID, job.ExecutedAt)
		
		if job.AgentID == agentID && job.ExecutedAt == nil {
			return job, nil
		}
	}

	return domain.JobRow{}, errors.New("nenhum comando pendente encontrado")
}