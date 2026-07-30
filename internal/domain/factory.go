package domain

// RatFactory é a nossa Fábrica Abstrata
type RatFactory interface {
	// createClient() Client
	// createAgent() Agent
	// createWorm() Worm
	
	CreateServer() Server
	CreateHandler() Handler
}
