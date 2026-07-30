package domain

// Worm representa um Agente com capacidade de propagação lateral autônoma
type Worm interface {
	Agent // Embutimos a interface Agent! O Worm herda a exigência do método Run() automaticamente.
	Spread() error
}