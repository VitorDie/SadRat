package tests

import (
	"testing"

	"github.com/VitorDie/SadRat/internal/domain"
)

func TestNewAgentRowCreation(t *testing.T) {
	// 1. A nossa especificação (O C&C gerou um UUID para o novo alvo)
	expectedUUID := "b6502209-1065-4cf5-b1a9-bfa15baef250"

	// 2. A ação: chamar a função construtora que cria a linha do Agente
	agent := domain.NewAgentRow(expectedUUID)

	// 3. As validações
	if agent.UUID != expectedUUID {
		t.Errorf("Esperava UUID %s, recebeu %s", expectedUUID, agent.UUID)
	}

	// Como a linha acabou de ser criada para ir pro banco, as datas não podem estar vazias
	if agent.CreatedAt.IsZero() {
		t.Error("Esperava que CreatedAt estivesse preenchido com a data e hora atual, mas veio zero")
	}

	if agent.LastSeenAt.IsZero() {
		t.Error("Esperava que LastSeenAt estivesse preenchido com a data e hora atual, mas veio zero")
	}
}