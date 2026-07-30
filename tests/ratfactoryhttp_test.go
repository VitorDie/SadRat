package tests

import (
	"testing"

	"github.com/VitorDie/SadRat/internal/domain"
	// O compilador vai chiar aqui porque o pacote factory ainda não existe!
	"github.com/VitorDie/SadRat/internal/factory" 
	"github.com/VitorDie/SadRat/internal/repository" 
)

func TestRatHttpFactory_CreateHandler(t *testing.T) {
	// 1. Especificação: Instanciamos a nossa fábrica concreta de HTTP
	repo := repository.NewServerDBPlainDataInMemory()
	var ratFactory domain.RatFactory = factory.NewRatFactoryHttp(repo, "http://localhost:8080")

	// 2. Ação: Pedimos para a fábrica nos dar o Handler pronto (com o Server já injetado dentro dele)
	handler := ratFactory.CreateHandler()

	// 3. Validação
	if handler == nil {
		t.Fatalf("Esperava receber um Handler instanciado, mas retornou nil")
	}

	// O Go garante em tempo de compilação que, se não for nil, 
	// o objeto retornado implementa a função Start(address string) error da interface domain.Handler!
}

func TestRatHttpFactory_CreateServer(t *testing.T) {
	repo := repository.NewServerDBPlainDataInMemory()
	var ratFactory domain.RatFactory = factory.NewRatFactoryHttp(repo, "http://localhost:8080")
	server := ratFactory.CreateServer()

	if server == nil {
		t.Fatalf("Esperava receber um Server instanciado, mas retornou nil")
	}
}

func TestRatHttpFactory_CreateClient(t *testing.T) {
	repo := repository.NewServerDBPlainDataInMemory()
	var ratFactory domain.RatFactory = factory.NewRatFactoryHttp(repo, "http://localhost:8080")
	client := ratFactory.CreateClient()

	if client == nil {
		t.Fatalf("Esperava receber um Client instanciado, mas retornou nil")
	}
}

func TestRatHttpFactory_CreateAgent(t *testing.T) {
	repo := repository.NewServerDBPlainDataInMemory()
	var ratFactory domain.RatFactory = factory.NewRatFactoryHttp(repo, "http://localhost:8080")
	agent := ratFactory.CreateAgent()

	if agent == nil {
		t.Fatalf("Esperava receber um Agent instanciado, mas retornou nil")
	}
}

func TestRatHttpFactory_CreateWorm(t *testing.T) {
	repo := repository.NewServerDBPlainDataInMemory()
	var ratFactory domain.RatFactory = factory.NewRatFactoryHttp(repo, "http://localhost:8080")
	worm := ratFactory.CreateWorm()

	if worm == nil {
		t.Fatalf("Esperava receber um Worm instanciado, mas retornou nil")
	}
}