package cli

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	
	"github.com/VitorDie/SadRat/internal/domain"
)

// O Model do Bubble Tea guarda o estado da interface e o nosso Cliente de Domínio
type RatModel struct {
	client      domain.Client // O nosso ClientHTTP injetado!
	
	leftBanner  string
	rightBanner string
	
	commandInput string
	logs         []string
}

// NewRatModel inicializa a UI do Operador carregando os banners (você pode usar go:embed para os .txt)
func NewRatModel(c domain.Client, leftBanner, rightBanner string) RatModel {
	return RatModel{
		client:      c,
		leftBanner:  leftBanner,
		rightBanner: rightBanner,
		logs:        []string{"[+] SadRat C&C Interface iniciada..."},
	}
}

// Init é disparado assim que a UI abre
func (m RatModel) Init() tea.Cmd {
	// Podemos retornar um comando aqui, por exemplo, buscar zumbis disponíveis logo de cara!
	return nil 
}

// Update é o cérebro da UI (onde capturamos o que você digita e atualizamos os logs)
func (m RatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	
	// Captura teclas do teclado
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			// 1. O Operador deu Enter! Pegamos o comando
			cmdStr := m.commandInput
			m.commandInput = "" 
			
			// 2. Adicionamos no log da direita
			m.logs = append(m.logs, fmt.Sprintf("[Comando Enviado]: %s", cmdStr))
			
			// 3. Aqui você usaria o m.client.SendCommand(...) em uma Goroutine ou tea.Cmd!
			
			return m, nil
		default:
			// Adiciona a letra digitada no painel da esquerda
			m.commandInput += msg.String()
		}
	}
	return m, nil
}

// View desenha os dois painéis na tela!
func (m RatModel) View() string {
	// Estilo para dividir a tela no meio usando Lip Gloss
	paneStyle := lipgloss.NewStyle().
		Width(50).
		Height(20).
		Border(lipgloss.RoundedBorder())

	// ==========================================
	// PAINEL ESQUERDO: Banner SadRat + Input
	// ==========================================
	leftContent := fmt.Sprintf("%s\n\n> %s_", m.leftBanner, m.commandInput)
	leftPane := paneStyle.Render(leftContent)

	// ==========================================
	// PAINEL DIREITO: Coração Partido + Logs
	// ==========================================
	logsStr := strings.Join(m.logs, "\n")
	rightContent := fmt.Sprintf("%s\n\n%s", m.rightBanner, logsStr)
	rightPane := paneStyle.Render(rightContent)

	// Junta os dois painéis lado a lado
	return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
}