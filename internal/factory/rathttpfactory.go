package factory

import (
	"github.com/VitorDie/SadRat/internal/domain"
	"github.com/VitorDie/SadRat/internal/presentation"
	"github.com/VitorDie/SadRat/internal/repository"
	"github.com/VitorDie/SadRat/internal/service"
)

// RatHttpFactory é a nossa fábrica concreta de HTTP em texto plano
type RatHttpFactory struct{}

// NewRatHttpFactory instancia a fábrica
func NewRatHttpFactory() *RatHttpFactory {
	return &RatHttpFactory{}
}

// CreateServer (protected na sua modelagem) constrói o Repositório e injeta na Regra de Negócio
func (f *RatHttpFactory) CreateServer() domain.Server {
	// A Fábrica HTTP usa o banco Plain Data
	repo := repository.NewServerDBPlainDataInMemory()
	
	// Retornamos a lógica (ConcreteServer) encapsulada na interface domain.Server
	return service.NewConcreteServer(repo)
}

// CreateHandler (public) constrói a cadeia inteira chamando o CreateServer internamente!
func (f *RatHttpFactory) CreateHandler() domain.Handler {
	// 1. A fábrica pede para si mesma construir o Server configurado
	server := f.CreateServer()
	
	// 2. Injeta o Server no Handler HTTP
	// Como HandlerHTTP possui o método Start(), ele automaticamente é um domain.Handler
	return presentation.NewHandlerHTTP(server)
}

// STUBS (Deixaremos preparados para quando criarmos os componentes)
// func (f *RatHttpFactory) CreateClient() domain.Client { return nil }
// func (f *RatHttpFactory) CreateAgent() domain.Agent { return nil }
// func (f *RatHttpFactory) CreateWorm() domain.Worm { return nil }