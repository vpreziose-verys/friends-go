package handler

type Handler struct {
	Base
}

type Base interface {
	Close()
}
