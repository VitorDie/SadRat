package domain

// RatFactory é a nossa Fábrica Abstrata
type RatFactory interface {
	CreateClient() Client
	CreateAgent() Agent
	// CreateWorm() Worm

	CreateServer() Server
	CreateHandler() Handler
}
