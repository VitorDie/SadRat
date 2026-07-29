package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	// Ajuste este import para bater com a estrutura do seu projeto quando formos para a Fase Verde
	"github.com/VitorDie/SadRat/internal/agenthttp" 
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
	bot := agenthttp.NewAgentHTTP(mockServer.URL)

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