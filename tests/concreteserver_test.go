package tests

import (
	"testing"

	"github.com/VitorDie/SadRat/internal/domain"
	"github.com/VitorDie/SadRat/internal/service"
)

// Mock do nosso Banco de Dados para que a Camada de Serviço possa ser testada
// sem precisar de uma conexão real com PostgreSQL ou SQLite.
type MockServerDB struct {
	// Adicionamos essa variável para o Mock "lembrar" qual job foi atualizado
	UpdatedJob domain.JobRow
}

// Os métodos que realmente importam para o nosso primeiro teste:
func (m *MockServerDB) SaveAgent(agent domain.AgentRow) error { return nil }

// Os STUBS (funções vazias) para satisfazer o resto da interface domain.ServerDB:
func (m *MockServerDB) SaveJob(job domain.JobRow) error { return nil }
func (m *MockServerDB) GetJobForAgent(agentID string) (domain.JobRow, error) {
	// Simulamos que o banco encontrou um comando "whoami" pendente para este zumbi
	return domain.JobRow{ID: "job-mock-1", AgentID: agentID, Command: "whoami"}, nil
}

// Simulamos que o banco encontrou o Job e ele já possui um resultado (Output)
func (m *MockServerDB) GetJob(id string) (domain.JobRow, error) {
	saida := "root\n"
	return domain.JobRow{ID: id, Command: "whoami", Output: &saida}, nil
}

// O Mock agora "grava" a entidade que o Serviço mandou atualizar
func (m *MockServerDB) UpdateJob(job domain.JobRow) error {
	m.UpdatedJob = job
	return nil
}
// Retorna uma lista com 1 Job Falso simulando o histórico do banco
func (m *MockServerDB) GetAllJobs() ([]domain.JobRow, error) { 
	return []domain.JobRow{{ID: "job-mock-1", Command: "whoami"}}, nil 
}
func (m *MockServerDB) GetAllAgents() ([]domain.AgentRow, error) {
	// Simulamos que o banco de dados tem 1 zumbi infectado
	return []domain.AgentRow{{UUID: "zumbi-mock-1"}}, nil
}
func TestConcreteServer_GiveAnUUIDForAgentRequest(t *testing.T) {
	// 1. Instanciamos o Cofre falso (Mock)
	mockDB := &MockServerDB{}

	// 2. Especificação: Instanciamos a nossa Camada de Serviço (ConcreteServer)
	// injetando a dependência do banco de dados falso.
	// Como ainda não construímos o pacote, o compilador vai chiar.
	var regraDeNegocio domain.Server = service.NewConcreteServer(mockDB)

	// 3. Ação: Solicitamos a geração e registro de um UUID para um novo Agente
	uuidStr, err := regraDeNegocio.GiveAnUUIDForAgentRequest()

	// 4. Validação
	if err != nil {
		t.Fatalf("Esperava sucesso ao registrar agente na regra de negócio, mas falhou: %v", err)
	}

	// Um UUID padrão tem 36 caracteres (com hifens). Vamos garantir que a string não veio vazia.
	if len(uuidStr) < 32 {
		t.Errorf("Esperava um UUID válido, mas retornou algo muito curto: '%s'", uuidStr)
	}
}

func TestConcreteServer_ReceiveCommand(t *testing.T) {
	// 1. Instanciamos o Cofre falso (Mock)
	mockDB := &MockServerDB{}

	// 2. Especificação: Instanciamos a nossa Camada de Serviço
	var regraDeNegocio domain.Server = service.NewConcreteServer(mockDB)

	// Simula um Operador tentando atacar um Agente específico
	agentID := "zumbi-alfa-123"
	comando := "whoami"
	argumentos := []string{}

	// 3. Ação: Solicitamos o despacho do ataque
	jobID, err := regraDeNegocio.ReceiveCommand(comando, argumentos, agentID)

	// 4. Validação
	if err != nil {
		t.Fatalf("Esperava sucesso ao despachar comando, mas falhou: %v", err)
	}

	// O JobID retornado precisa ser um UUID válido (36 caracteres)
	if len(jobID) < 32 {
		t.Errorf("Esperava um JobID válido, mas retornou algo inválido ou vazio: '%s'", jobID)
	}
}

func TestConcreteServer_SendAvailableAgents(t *testing.T) {
	// 1. Instanciamos o Cofre falso (Mock)
	mockDB := &MockServerDB{}

	// 2. Especificação: Instanciamos a nossa Camada de Serviço
	var regraDeNegocio domain.Server = service.NewConcreteServer(mockDB)

	// 3. Ação: Solicitamos a lista de agentes disponíveis
	agentes, err := regraDeNegocio.SendAvailableAgents()

	// 4. Validação
	if err != nil {
		t.Fatalf("Esperava sucesso ao buscar agentes, mas falhou: %v", err)
	}

	if len(agentes) == 0 {
		t.Errorf("Esperava receber a lista de agentes do banco, mas retornou vazia")
	}

	if len(agentes) > 0 && agentes[0].UUID != "zumbi-mock-1" {
		t.Errorf("Esperava o agente 'zumbi-mock-1', mas recebeu: '%s'", agentes[0].UUID)
	}
}

func TestConcreteServer_SendJobs(t *testing.T) {
	// 1. Instanciamos o Cofre falso (Mock)
	mockDB := &MockServerDB{}

	// 2. Especificação: Instanciamos a nossa Camada de Serviço
	var regraDeNegocio domain.Server = service.NewConcreteServer(mockDB)

	// O Zumbi informa o seu ID
	agentID := "zumbi-mock-1"

	// 3. Ação: Solicitamos o trabalho pendente para aquele zumbi
	job, err := regraDeNegocio.SendJobs(agentID)

	// 4. Validação
	if err != nil {
		t.Fatalf("Esperava sucesso ao buscar jobs, mas falhou: %v", err)
	}

	// Verificamos se o serviço retornou o Job correto do banco
	// (Lembre-se de trocar .UUID pelo nome exato do campo na sua struct JobRow)
	if job.ID != "job-mock-1" {
		t.Errorf("Esperava repassar o Job 'job-mock-1', mas retornou: '%s'", job.ID)
	}
}

func TestConcreteServer_ReceiveJobResult(t *testing.T) {
	// 1. Instanciamos o Cofre falso (Mock)
	mockDB := &MockServerDB{}
	var regraDeNegocio domain.Server = service.NewConcreteServer(mockDB)

	jobID := "job-mock-1"
	resultadoDoZumbi := "root\n"

	// 2. Ação: O Serviço recebe o resultado e deve atualizar no banco
	err := regraDeNegocio.ReceiveJobResult(jobID, resultadoDoZumbi)

	// 3. Validação de Erro
	if err != nil {
		t.Fatalf("Esperava sucesso ao processar resultado do job, mas falhou: %v", err)
	}

	// 4. Validação Lógica (O coração do TDD)
	// Primeiro, verificamos se o ponteiro é nulo (o que DEVE acontecer agora, já que o stub do serviço está vazio)
	if mockDB.UpdatedJob.Output == nil {
		t.Errorf("A regra de negócio falhou em atualizar o banco. O campo Output continua nulo (nil)")
	} else if *mockDB.UpdatedJob.Output != resultadoDoZumbi {
		// Se não for nulo, usamos o asterisco '*' para desreferenciar o ponteiro e acessar a string real
		t.Errorf("A regra de negócio falhou. Esperava o output '%s', mas recebeu '%s'", resultadoDoZumbi, *mockDB.UpdatedJob.Output)
	}
}

func TestConcreteServer_SendJobResult(t *testing.T) {
	// 1. Instanciamos o Cofre falso (Mock)
	mockDB := &MockServerDB{}
	var regraDeNegocio domain.Server = service.NewConcreteServer(mockDB)

	jobID := "job-mock-1"
	resultadoEsperado := "root\n"

	// 2. Ação: O Operador pede para ver o resultado do ataque
	resultado, err := regraDeNegocio.SendJobResult(jobID)

	// 3. Validação de Erro
	if err != nil {
		t.Fatalf("Esperava sucesso ao buscar o resultado, mas falhou: %v", err)
	}

	// 4. Validação Lógica
	// Como o nosso stub em concreteserver.go ainda retorna apenas "", este teste deve falhar!
	if resultado != resultadoEsperado {
		t.Errorf("A regra de negócio falhou. Esperava repassar o output '%s', mas recebeu '%s'", resultadoEsperado, resultado)
	}
}

func TestConcreteServer_SendAllJobResults(t *testing.T) {
	// 1. Instanciamos o Cofre falso (Mock)
	mockDB := &MockServerDB{}
	var regraDeNegocio domain.Server = service.NewConcreteServer(mockDB)

	// 2. Ação: O Operador pede o histórico completo de ataques
	jobs, err := regraDeNegocio.SendAllJobResults()

	// 3. Validação de Erro
	if err != nil {
		t.Fatalf("Esperava sucesso ao buscar o histórico, mas falhou: %v", err)
	}

	// 4. Validação Lógica
	// Como o nosso stub em concreteserver.go ainda retorna nil, a lista virá vazia e o teste vai falhar!
	if len(jobs) == 0 {
		t.Errorf("A regra de negócio falhou. Esperava receber o histórico de jobs do banco, mas a lista retornou vazia")
	} else if jobs[0].ID != "job-mock-1" {
		t.Errorf("Esperava o job 'job-mock-1' na lista, mas recebeu '%s'", jobs[0].ID)
	}
}