package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	// Ajuste este import para bater com a estrutura do seu projeto quando formos para a Fase Verde
	"github.com/VitorDie/SadRat/internal/agent" 
	"github.com/VitorDie/SadRat/internal/domain"
)

func TestAgentHTTP_RequestAnUUIDToMe(t *testing.T) {
	// 1. Mock do C&C: Criamos um servidor falso para fingir ser a nossa API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificamos se o Agente está chamando a rota certa com o método certo
		if r.URL.Path != "/api/agents" || r.Method != http.MethodPost {
			t.Fatalf("O Agente tentou acessar %s %s, mas deveria ser POST /api/agents", r.Method, r.URL.Path)
		}
		
		// O falso C&C devolve um UUID de mentira
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"uuid": "zumbi-mock-123"})
	}))
	// Garante que o servidor falso será desligado ao fim do teste
	defer mockServer.Close()

	// 2. Especificação: Instanciamos o AgenteHTTP passando a URL do Servidor Falso
	// (Ainda vamos criar o pacote agenthttp e a struct AgentHTTP)
	bot := agent.NewAgentHTTP(mockServer.URL)

	// 3. Ação: O Agente bate na porta do C&C para pedir sua identidade
	uuid, err := bot.RequestAnUUIDToMe()
	if err != nil {
		t.Fatalf("O Agente falhou ao pedir o UUID: %v", err)
	}

	// 4. Validação: O Agente conseguiu ler o JSON do C&C?
	if uuid != "zumbi-mock-123" {
		t.Errorf("Esperava que o Agente capturasse o UUID 'zumbi-mock-123', mas capturou '%s'", uuid)
	}
}

func TestAgentHTTP_RequestJob(t *testing.T) {
	// 1. Mock do C&C: Criamos um servidor falso para fingir ser a nossa API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// O Agente deve usar a técnica de Long Polling na rota correta
		if r.URL.Path != "/api/agents/agent-123/job" || r.Method != http.MethodGet {
			t.Fatalf("O Agente tentou acessar %s %s, mas deveria ser GET /api/agents/agent-123/job", r.Method, r.URL.Path)
		}
		
		// O falso C&C devolve um comando agendado pelo Operador em formato JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "job-999", "command": "whoami", "args": []}`))
	}))
	defer mockServer.Close()

	// 2. Especificação: Instanciamos o AgenteHTTP
	var bot = agent.NewAgentHTTP(mockServer.URL)

	// 3. Ação: O Agente pede o seu próximo comando passando sua identidade
	job, err := bot.RequestJob("agent-123")
	if err != nil {
		t.Fatalf("O Agente falhou ao pedir comandos: %v", err)
	}

	// 4. Validação: O Agente conseguiu decodificar o DTO do comando?
	if job.ID != "job-999" || job.Command != "whoami" {
		t.Errorf("Esperava o comando 'whoami' com ID 'job-999', mas recebeu: %+v", job)
	}
}

func TestAgentHTTP_ExecuteJob(t *testing.T) {
	// 1. Especificação: Instanciamos o Agente 
	// (Não precisamos de um C&C falso desta vez, pois não faremos requisições de rede)
	var bot = agent.NewAgentHTTP("")

	// 2. Preparamos uma Entidade de comando simulando a regra de negócio pura (Não é mais DTO!)
	job := domain.JobEntity{
		ID:      "job-teste-exec",
		Command: "echo",
		Args:    []string{"zumbi"},
	}

	// 3. Ação: O Agente invoca o Sistema Operacional para executar o comando
	output, err := bot.ExecuteJob(job)
	if err != nil {
		t.Fatalf("O Agente falhou ao tentar executar o comando no SO: %v", err)
	}

	// 4. Validação: A saída padrão do terminal (stdout) tem que ser exatamente o que o "echo" cospe
	if output != "zumbi\n" {
		t.Errorf("Esperava que o terminal devolvesse 'zumbi\\n', mas retornou '%s'", output)
	}
}

func TestAgentHTTP_SendJobResult(t *testing.T) {
	// 1. Mock do C&C: O Servidor Falso espera receber o relatório
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verifica se o Agente está batendo na rota certa
		if r.URL.Path != "/api/jobs/result" || r.Method != http.MethodPost {
			t.Fatalf("O Agente tentou acessar %s %s, mas deveria ser POST /api/jobs/result", r.Method, r.URL.Path)
		}

		// Valida se o DTO enviado pelo Agente está correto
		var req map[string]string
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("O Agente não enviou um JSON válido: %v", err)
		}

		if req["job_id"] != "job-teste-123" || req["output"] != "root\n" {
			t.Errorf("O Agente enviou os dados errados para o C&C: %+v", req)
		}

		w.WriteHeader(http.StatusOK) // O Mestre agradece
	}))
	defer mockServer.Close()

	// 2. Especificação: Instanciamos o AgenteHTTP
	var bot = agent.NewAgentHTTP(mockServer.URL)

	// 3. Ação: O Agente envia o resultado da execução para o servidor
	err := bot.SendJobResult("job-teste-123", "root\n")
	
	// 4. Validação
	if err != nil {
		t.Fatalf("O Agente falhou ao tentar relatar o resultado: %v", err)
	}
}
