package client

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// Agent representa a máquina infectada (Zumbi) na visão do Operador
type Agent struct {
	ID string `json:"id"`
}

// Client é a Interface genérica para a ferramenta do Operador (preparada para ClientDNS no futuro!)
type Client interface {
	RequestAvailableAgents() ([]Agent, error)
}

// ClientHTTP é a implementação específica que fala com o C&C via HTTP
type ClientHTTP struct {
	serverURL  string
	httpClient *http.Client
}

// NewClientHTTP cria uma nova instância do painel do Operador
func NewClientHTTP(serverURL string) *ClientHTTP {
	return &ClientHTTP{
		serverURL: serverURL,
		// Regra de Ouro: Sempre usar timeout para evitar que a CLI trave!
		httpClient: &http.Client{Timeout: 10 * time.Second}, 
	}
}

// RequestAvailableAgents bate no C&C e busca todos os zumbis sob nosso controle
func (c *ClientHTTP) RequestAvailableAgents() ([]Agent, error) {
	// 1. Monta a URL da nossa API
	url := c.serverURL + "/api/agents"

	// 2. Prepara a requisição GET
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// 3. Dispara a requisição para o Servidor C&C
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 4. Valida se o servidor autorizou a requisição
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("falha ao buscar agentes no C&C, status: " + resp.Status)
	}

	// 5. Lemos o JSON com a lista de zumbis e convertemos para o nosso DTO (Slice de Agent)
	var agents []Agent
	if err := json.NewDecoder(resp.Body).Decode(&agents); err != nil {
		return nil, err
	}

	return agents, nil
}
