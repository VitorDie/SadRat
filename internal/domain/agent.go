package domain

import (	
	"github.com/VitorDie/SadRat/internal/dto" 
)


type Agent interface {
	RequestAnUUIDToMe() (string, error)
	RequestJob(agentID string) (dto.JobUsedByAgentDTO, error)
	ExecuteJob(job dto.JobUsedByAgentDTO) (string, error)
	SendJobResult(jobID string, output string) error
}
