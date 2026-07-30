package tests

import (
	"testing"

	"github.com/VitorDie/SadRat/internal/domain"
	// O compilador vai chiar aqui porque o pacote factory ainda não existe!
	"github.com/VitorDie/SadRat/internal/factory" 
)

func TestRatHttpFactory_CreateHandler(t *testing.T) {
	// 1. Especificação: Instanciamos a nossa fábrica concreta de HTTP
	var ratFactory domain.RatFactory = factory.NewRatHttpFactory()

	// 2. Ação: Pedimos para a fábrica nos dar o Handler pronto (com o Server já injetado dentro dele)
	handler := ratFactory.CreateHandler()

	// 3. Validação
	if handler == nil {
		t.Fatalf("Esperava receber um Handler instanciado, mas retornou nil")
	}

	// O Go garante em tempo de compilação que, se não for nil, 
	// o objeto retornado implementa a função Start(address string) error da interface domain.Handler!
}