package tests

import (
	"testing"
	"github.com/VitorDie/SadRat/internal/agent" 
)

func TestWormHTTP_Spread(t *testing.T) {
	// 1. Especificação: Instanciamos o nosso Worm apontando para um C&C fictício
	worm := agent.NewWormHTTP("http://localhost:8080")

	// 2. Ação: Mandamos o Worm se propagar
	err := worm.Spread()

	// 3. Validação: Como ainda é um stub (ou uma rotina que roda silenciosamente), 
	// esperamos que ele não retorne erros catastróficos.
	if err != nil {
		t.Fatalf("Esperava que o Spread executasse sem erros, mas retornou: %v", err)
	}
}