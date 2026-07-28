package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitorDie/SadRat/internal/presentation"
	"github.com/VitorDie/SadRat/internal/repository"
)

func TestServerHTTP_GiveAnUUIDForAgentRequest(t *testing.T) {
	// 1. Especificação: Instanciamos o cofre e injetamos no Servidor HTTP
	repo := repository.NewServerDBPlainDataInMemory()
	server := presentation.NewServerHTTP(repo)

	// 2. Ação: Simulamos um Agente recém-infectado chamando o C&C via rede
	req, err := http.NewRequest(http.MethodPost, "/api/agents", nil)
	if err != nil {
		t.Fatalf("Erro ao criar a requisição: %v", err)
	}

	// O Recorder atua como se fosse o "navegador/cliente" para capturar a resposta
	rr := httptest.NewRecorder()
	
	// A função Router() deve devolver o roteador que contém as URLs mapeadas
	server.Router().ServeHTTP(rr, req)

	// 3. Validação de Rota e Status HTTP
	if status := rr.Code; status != http.StatusCreated { // Esperamos HTTP 201 (Created)
		t.Errorf("Esperava status %v, recebeu %v", http.StatusCreated, status)
	}

	// 4. Validação do DTO de Resposta (Encoding responses)
	var response map[string]string
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Erro ao decodificar JSON da API: %v", err)
	}

	if response["uuid"] == "" {
		t.Error("A API deveria ter gerado e devolvido um UUID em JSON, mas veio vazio")
	}
}
