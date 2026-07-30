package domain

type Handler interface {
	Start(address string) error
}