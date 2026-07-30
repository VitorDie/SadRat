package domain

type Agent interface {
	RequestAnUUIDToMe() (string, error)
	RequestJob(agentID string) (JobEntity, error)
	ExecuteJob(job JobEntity) (string, error)
	SendJobResult(jobID string, output string) error
	Run() error
}
