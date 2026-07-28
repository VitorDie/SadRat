package tests

import (
	"testing"

	"github.com/VitorDie/SadRat/internal/domain"
	"github.com/VitorDie/SadRat/internal/repository" // O nosso novo pacote!
)

func TestServerDBPlainDataInMemory_CRUD(t *testing.T) {
	// 1. Especificação: Instanciamos o nosso banco de dados falso (em memória)
	// Ele vai guardar as informações em listas (slices/maps) internamente.
	repo := repository.NewServerDBPlainDataInMemory()

	// 2. Criamos as linhas que queremos salvar
	agent := domain.NewAgentRow("b6502209-1065-4cf5-b1a9-bfa15baef250")
	job := domain.NewJobRow(agent.UUID, "whoami", []string{})

	// 3. Ação: Testamos a operação CREATE (Salvar)
	errAgent := repo.SaveAgent(agent)
	if errAgent != nil {
		t.Fatalf("Não esperava erro ao salvar o agente: %v", errAgent)
	}

	errJob := repo.SaveJob(job)
	if errJob != nil {
		t.Fatalf("Não esperava erro ao salvar o job: %v", errJob)
	}

	// 4. Validação: Testamos a operação READ (Buscar)
	savedAgent, err := repo.GetAgent(agent.UUID)
	if err != nil {
		t.Fatalf("Não esperava erro ao buscar o agente: %v", err)
	}

	if savedAgent.UUID != agent.UUID {
		t.Errorf("Esperava o Agente %s, mas o banco retornou %s", agent.UUID, savedAgent.UUID)
	}

	savedJob, err := repo.GetJob(job.ID)
	if err != nil {
		t.Fatalf("Não esperava erro ao buscar o job: %v", err)
	}

	if savedJob.Command != "whoami" {
		t.Errorf("Esperava o comando whoami, mas o banco retornou %s", savedJob.Command)
	}
}