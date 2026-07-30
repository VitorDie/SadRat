package domain

// ServerDBPlainData é a interface específica para a fábrica HTTP.
// Ao embutir a ServerDB, ela herda automaticamente todos os métodos (SaveAgent, GetJob, etc)
// mantendo a assinatura idêntica sem duplicar linhas de código!
type ServerDBPlainData interface {
	ServerDB
}