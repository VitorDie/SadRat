package domain

type Client interface {
	RequestAvailableAgents() ([]Agent, error)
	SendCommand(command string, args []string, agentID string) (string, error)
	RequestJobResult(jobID string) (string, error)
}