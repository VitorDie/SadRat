package domain

// Repository exige os métodos que o Servidor precisa para funcionar
type ServerDB interface {
	SaveAgent(agent AgentRow) error
	SaveJob(job JobRow) error // Adicionamos a exigência de salvar o comando
	GetJobForAgent(agentID string) (JobRow, error) // Novo método
	GetJob(id string) (JobRow, error) 
	UpdateJob(job JobRow) error
	GetAllJobs() ([]JobRow, error) 
	GetAllAgents() ([]AgentRow, error)
}