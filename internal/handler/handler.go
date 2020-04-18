package handler

type Handler struct {
	base
}

type base interface {
	Close()
	LogTest()
}

func NewHandler(base base) *Handler {
	return &Handler{base}
}
