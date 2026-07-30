package domain

type Server interface {
	SendAvailableAgents() ([]AgentRow, error)
	ReceiveCommand(command string, args []string, agentID string) (string, error) // Retorna o ID do Job gerado
	SendJobs(agentID string) (JobRow, error)
	ReceiveJobResult(jobID string, output string) error
	SendJobResult(jobID string) (string, error) // Retorna o output do Job
	SendAllJobResults() ([]JobRow, error)
	GiveAnUUIDForAgentRequest() (string, error) // Retorna o UUID gerado
}