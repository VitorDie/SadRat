package tests

import (
	"testing"
	"github.com/VitorDie/SadRat/internal/domain"
)

func TestNewJobInfoCreation(t *testing.T) {
	// 1. A nossa especificação (os dados que vamos enviar)
	agentID := "agent-123"
	command := "ls"
	args := []string{"-la", "/tmp"}

	// 2. A ação: chamar a função construtora que cria um JobInfo
	job := domain.NewJobInfo(agentID, command, args)

	// 3. As validações (checando se o JobInfo nasceu com as regras corretas)
	
	if job.ID == "" {
		t.Error("Esperava que o ID (UUID) fosse gerado automaticamente, mas veio vazio")
	}

	if job.AgentID != agentID {
		t.Errorf("Esperava AgentID %s, recebeu %s", agentID, job.AgentID)
	}

	if job.Command != command {
		t.Errorf("Esperava Command %s, recebeu %s", command, job.Command)
	}

	if len(job.Args) != 2 {
		t.Errorf("Esperava 2 argumentos no array, recebeu %d", len(job.Args))
	}

	// Um Job recém-criado pelo Operador nunca pode ter sido executado ainda
	if job.ExecutedAt != nil {
		t.Errorf("Esperava ExecutedAt como nil (vazio), recebeu %v", job.ExecutedAt)
	}
	
	// Um Job recém-criado não pode ter resposta do agente ainda
	if job.Output != nil {
		t.Errorf("Esperava Output como nil (vazio), recebeu %v", job.Output)
	}
}
