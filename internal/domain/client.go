package domain

import (	
	"github.com/VitorDie/SadRat/internal/dto" 
)

type Client interface {
	RequestAvailableAgents() ([]dto.AgentUsedByClientDTO, error)
	SendCommand(command string, args []string, agentID string) (string, error)
	RequestJobResult(jobID string) (string, error)
}