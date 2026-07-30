package tests

import (
	"testing"

	"github.com/VitorDie/SadRat/internal/domain"
	// O pacote que vamos criar na Fase Verde:
	// "github.com/VitorDie/SadRat/internal/service"
)

// Mock do nosso Banco de Dados para que a Camada de Serviço possa ser testada 
// sem precisar de uma conexão real com PostgreSQL, SQLite, etc.
type MockServerDB struct{}

// O Mock apenas finge que salvou o agente com sucesso
func (m *MockServerDB) SaveAgent(agent domain.AgentRow) error {
	return nil
}

// Simulamos os outros métodos do DB para satisfazer a interface (deixamos vazios por enquanto)
func (m *MockServerDB) SaveJob(job domain.JobRow) error { return nil }
// (Adicione aqui outros métodos vazios que a sua domain.ServerDB exige)


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
