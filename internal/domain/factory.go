package domain

// RatFactory é a nossa Fábrica Abstrata
type RatFactory interface {
	// Deixaremos Client, Agent e Worm comentados por enquanto até criarmos essas entidades
	// CreateClient() Client
	// CreateAgent() Agent
	// CreateWorm() Worm
	
	CreateServer() Server
	CreateHandler() Handler
}
